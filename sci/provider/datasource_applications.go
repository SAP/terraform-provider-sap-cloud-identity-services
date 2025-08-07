package provider

import (
	"context"
	"fmt"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newApplicationsDataSource() datasource.DataSource {
	return &applicationsDataSource{}
}

type applicationsDataSource struct {
	cli *cli.SciClient
}

type applicationsData struct {
	Values types.List `tfsdk:"values"`
}

var authenticationSchemaObjType = map[string]attr.Type{
	"sso_type": types.StringType,
	"subject_name_identifier": types.ObjectType{
		AttrTypes: subjectNameIdentitfierObjType,
	},
	"subject_name_identifier_function": types.StringType,
	"assertion_attributes": types.ListType{
		ElemType: assertionAttributesObjType,
	},
	"advanced_assertion_attributes": types.ListType{
		ElemType: advancedAssertionAttributesObjType,
	},
	"default_authenticating_idp": types.StringType,
	"conditional_authentication": types.ListType{
		ElemType: authenticationRulesObjType,
	},
	"openid_connect_configuration": types.ObjectType{
		AttrTypes: openIdConnectConfigurationObjType,
	},
}

var subjectNameIdentitfierObjType = map[string]attr.Type{
	"source": types.StringType,
	"value":  types.StringType,
}

var advancedAssertionAttributesObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"source":          types.StringType,
		"attribute_name":  types.StringType,
		"attribute_value": types.StringType,
		"inherited":       types.BoolType,
	},
}

var assertionAttributesObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"attribute_name":  types.StringType,
		"attribute_value": types.StringType,
		"inherited":       types.BoolType,
	},
}

var authenticationRulesObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"user_type":            types.StringType,
		"user_group":           types.StringType,
		"user_email_domain":    types.StringType,
		"identity_provider_id": types.StringType,
		"ip_network_range":     types.StringType,
	},
}

var openIdConnectConfigurationObjType = map[string]attr.Type{
	"redirect_uris": types.SetType{
		ElemType: types.StringType,
	},
	"post_logout_redirect_uris": types.SetType{
		ElemType: types.StringType,
	},
	"front_channel_logout_uris": types.SetType{
		ElemType: types.StringType,
	},
	"back_channel_logout_uris": types.SetType{
		ElemType: types.StringType,
	},
	"token_policy": types.ObjectType{
		AttrTypes: tokenPolicyObjType,
	},
	"restricted_grant_types": types.SetType{
		ElemType: types.StringType,
	},
	"proxy_config": types.ObjectType{
		AttrTypes: proxyConfigObjType,
	},
}

var tokenPolicyObjType = map[string]attr.Type{
	"jwt_validity":                    types.Int32Type,
	"refresh_validity":                types.Int32Type,
	"refresh_parallel":                types.Int32Type,
	"max_exchange_period":             types.StringType,
	"refresh_token_rotation_scenario": types.StringType,
	"access_token_format":             types.StringType,
}

var proxyConfigObjType = map[string]attr.Type{
	"acrs": types.SetType{
		ElemType: types.StringType,
	},
}

var appObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                    types.StringType,
		"name":                  types.StringType,
		"description":           types.StringType,
		"parent_application_id": types.StringType,
		"multi_tenant_app":      types.BoolType,
		"authentication_schema": types.ObjectType{
			AttrTypes: authenticationSchemaObjType,
		},
	},
}

func (d *applicationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *applicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

func (d *applicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a list of applications from the SAP Cloud Identity services.`,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Id of the application",
							Computed:            true,
							Validators: []validator.String{
								utils.ValidUUID(),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the application",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Free text description of the Application",
							Computed:            true,
						},
						"parent_application_id": schema.StringAttribute{
							MarkdownDescription: "ID of the parent, from which the application will inherit its configurations",
							Computed:            true,
							Validators: []validator.String{
								utils.ValidUUID(),
							},
						},
						"multi_tenant_app": schema.BoolAttribute{
							MarkdownDescription: "Only for Internal Use",
							Computed:            true,
						},
						"authentication_schema": schema.SingleNestedAttribute{
							MarkdownDescription: "Configure attributes particular to the schema \"urn:sap:identity:application:schemas:extension:sci:1.0:Authentication\"",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"sso_type": schema.StringAttribute{
									MarkdownDescription: "The preferred protocol for the application",
									Computed:            true,
								},
								"subject_name_identifier": schema.SingleNestedAttribute{
									MarkdownDescription: "The attribute by which the application uses to identify the users. Used by the application to uniquely identify the user during logon.\n" +
										fmt.Sprintln("Identity Authentication sends the attribute to the application as :") +
										fmt.Sprintln("\t - subject in OpenID Connect tokens") +
										fmt.Sprintln("\t - name ID in SAML 2.0 assertions"),
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"source": schema.StringAttribute{
											MarkdownDescription: utils.ValidValuesString(sourceValues),
											Computed:            true,
										},
										"value": schema.StringAttribute{
											MarkdownDescription: "If the source is Identity Directory, the only acceptable values are \" none, uid, mail, loginName, displayName, personnelNumber, userUuid\"",
											Computed:            true,
										},
									},
								},
								"subject_name_identifier_function": schema.StringAttribute{
									MarkdownDescription: "Convert the subject name identifier to uppercase or lowercase",
									Computed:            true,
								},
								"assertion_attributes": schema.ListNestedAttribute{
									MarkdownDescription: "User attributes to be sent to the application. The Source of these attributes is always the Identity Directory",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"attribute_name": schema.StringAttribute{
												MarkdownDescription: "Name of the attribute",
												Computed:            true,
											},
											"attribute_value": schema.StringAttribute{
												MarkdownDescription: "Value of the attribute.",
												Computed:            true,
											},
											"inherited": schema.BoolAttribute{
												MarkdownDescription: "Indicates whether the attribute has been inherited from a parent application.",
												Computed:            true,
											},
										},
									},
								},
								"advanced_assertion_attributes": schema.ListNestedAttribute{
									MarkdownDescription: "Identical to the assertion attributes, except that the assertion attributes can come from other Sources.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"source": schema.StringAttribute{
												MarkdownDescription: utils.ValidValuesString(sourceValues[1:]),
												Computed:            true,
											},
											"attribute_name": schema.StringAttribute{
												MarkdownDescription: "Name of the attribute",
												Computed:            true,
											},
											"attribute_value": schema.StringAttribute{
												MarkdownDescription: "Value of the attribute",
												Computed:            true,
											},
											"inherited": schema.BoolAttribute{
												MarkdownDescription: "Indicates whether the attribute has been inherited from a parent application.",
												Computed:            true,
											},
										},
									},
								},
								"default_authenticating_idp": schema.StringAttribute{
									MarkdownDescription: "A default identity provider can be used for users with any user domain, group and type. This identity provider is used when none of the defined authentication rules meets the criteria.",
									Computed:            true,
								},
								"conditional_authentication": schema.ListNestedAttribute{
									MarkdownDescription: "Define rules for authenticating identity provider according to email domain, user type, user group, and IP range. Each rule is evaluated by priority until the criteria of a rule are fulfilled.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"identity_provider_id": schema.StringAttribute{
												MarkdownDescription: "The identity provider to delegate authentication to when all the defined conditions are met.",
												Computed:            true,
											},
											"user_type": schema.StringAttribute{
												MarkdownDescription: "The type of user to be authenticated.",
												Computed:            true,
											},
											"user_group": schema.StringAttribute{
												MarkdownDescription: "The user group to be authenticated.",
												Computed:            true,
											},
											"user_email_domain": schema.StringAttribute{
												MarkdownDescription: "Valid email domain to be authenticated.",
												Computed:            true,
											},
											"ip_network_range": schema.StringAttribute{
												MarkdownDescription: "Valid IP range to be authenticated.",
												Computed:            true,
											},
										},
									},
								},
								"openid_connect_configuration": schema.SingleNestedAttribute{
									MarkdownDescription: "oidc",
									Optional:            true,
									Computed:            true,
									Validators: []validator.Object{
										objectvalidator.AlsoRequires(
											path.MatchRoot("authentication_schema").AtName("openid_connect_configuration").AtName("redirect_uris"),
										),
									},
									Attributes: map[string]schema.Attribute{
										"redirect_uris": schema.SetAttribute{
											MarkdownDescription: "A list of redirect URIs that the OpenID Provider is allowed to redirect to after authentication. Must contain 1 to 20 valid URIs.",
											ElementType:         types.StringType,
											Optional:            true,
											Validators: []validator.Set{
												setvalidator.SizeBetween(1, 20),
											},
										},
										"post_logout_redirect_uris": schema.SetAttribute{
											MarkdownDescription: "List of URIs to which the user will be redirected after logging out from the application. Can include up to 20 URIs.",
											ElementType:         types.StringType,
											Optional:            true,
											Validators: []validator.Set{
												setvalidator.SizeBetween(1, 20),
											},
										},
										"front_channel_logout_uris": schema.SetAttribute{
											MarkdownDescription: "List of front-channel logout URIs that support browser-based logout. Each must be a valid URL and up to 20 URIs are allowed.",
											ElementType:         types.StringType,
											Optional:            true,
											Validators: []validator.Set{
												setvalidator.SizeBetween(1, 20),
												setvalidator.ValueStringsAre(utils.ValidUrl()),
											},
										},
										"back_channel_logout_uris": schema.SetAttribute{
											MarkdownDescription: "List of back-channel logout URIs that support server-to-server logout notifications. Each must be a valid URL. Up to 20 URIs allowed.",
											ElementType:         types.StringType,
											Optional:            true,
											Validators: []validator.Set{
												setvalidator.SizeBetween(1, 20),
												setvalidator.ValueStringsAre(utils.ValidUrl()),
											},
										},
										"token_policy": schema.SingleNestedAttribute{
											MarkdownDescription: "Defines the token policy for the application.",
											Optional:            true,
											Computed:            true,
											Attributes: map[string]schema.Attribute{
												"jwt_validity": schema.Int32Attribute{
													MarkdownDescription: "JWT access token validity in seconds. Must be between 60 seconds (1 minute) and 43200 seconds (12 hours).",
													Optional:            true,
													Computed:            true,
													Validators: []validator.Int32{
														int32validator.Between(60, 43200),
													},
												},
												"refresh_validity": schema.Int32Attribute{
													MarkdownDescription: "Refresh token validity in seconds. Can range from 0 to 15552000 seconds (180 days).",
													Optional:            true,
													Computed:            true,
													Validators: []validator.Int32{
														int32validator.Between(0, 15552000),
													},
												},
												"refresh_parallel": schema.Int32Attribute{
													MarkdownDescription: "Maximum number of refresh tokens that can be used in parallel. Valid values range from 1 to 10.",
													Optional:            true,
													Computed:            true,
													Validators: []validator.Int32{
														int32validator.Between(1, 10),
													},
												},
												"max_exchange_period": schema.StringAttribute{
													MarkdownDescription: "Maximum token exchange period. Must be one of the allowed values.",
													Optional:            true,
													Computed:            true,
													Validators: []validator.String{
														stringvalidator.OneOf(maxExchangePeriodValues...),
													},
												},
												"refresh_token_rotation_scenario": schema.StringAttribute{
													MarkdownDescription: "Defines the scenario for refresh token rotation. Must be one of the allowed values.",
													Optional:            true,
													Computed:            true,
													Validators: []validator.String{
														stringvalidator.OneOf(refreshTokenRotationScenarioValues...),
													},
												},
												"access_token_format": schema.StringAttribute{
													MarkdownDescription: "The format of the access token issued. Must be one of the allowed values.",
													Optional:            true,
													Computed:            true,
													Validators: []validator.String{
														stringvalidator.OneOf(accessTokenFormatValues...),
													},
												},
											},
										},
										"restricted_grant_types": schema.SetAttribute{
											MarkdownDescription: "Set of OAuth 2.0 grant types that are restricted for the application. Must match one of the supported grant types.",
											Optional:            true,
											Computed:            true,
											ElementType:         types.StringType,
											Validators: []validator.Set{
												setvalidator.ValueStringsAre(stringvalidator.OneOf(restrictedGrantTypesValues...)),
											},
										},
										"proxy_config": schema.SingleNestedAttribute{
											MarkdownDescription: "Optional proxy configuration including accepted ACR values.",
											Optional:            true,
											Computed:            true,
											Attributes: map[string]schema.Attribute{
												"acrs": schema.SetAttribute{
													MarkdownDescription: "Set of accepted ACR (Authentication Context Class Reference) values. Up to 20 values allowed.",
													Optional:            true,
													ElementType:         types.StringType,
													Validators: []validator.Set{
														setvalidator.SizeAtMost(20),
														setvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 99)),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *applicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config applicationsData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, _, err := d.cli.Application.Get(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application", fmt.Sprintf("%s", err))
		return
	}

	resApps := applicationsValueFrom(ctx, res)

	config.Values, diags = types.ListValueFrom(ctx, appObjType, resApps)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
