package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

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
	"saml2_config": types.ObjectType{
		AttrTypes: appSaml2ConfigObjType.AttrTypes,
	},
}

var appSaml2ConfigObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"saml_metadata_url": types.StringType,
		"acs_endpoints": types.ListType{
			ElemType: acsEndpointsObjType,
		},
		"slo_endpoints": types.ListType{
			ElemType: appSaml2SloEndpointObjType,
		},
		"signing_certificates": types.ListType{
			ElemType: saml2SigningCertificateObjType,
		},
		"encryption_certificate":       saml2EncryptionCertificateObjType,
		"response_elements_to_encrypt": types.StringType,
		"default_name_id_format":       types.StringType,
		"sign_slo_messages":            types.BoolType,
		"require_signed_slo_messages":  types.BoolType,
		"require_signed_auth_requests": types.BoolType,
		"sign_assertions":              types.BoolType,
		"sign_auth_responses":          types.BoolType,
		"digest_algorithm":             types.StringType,
	},
}

var acsEndpointsObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"binding_name": types.StringType,
		"location":     types.StringType,
		"index":        types.Int32Type,
		"default":      types.BoolType,
	},
}

var appSaml2SloEndpointObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"binding_name":      types.StringType,
		"location":          types.StringType,
		"response_location": types.StringType,
	},
}

var saml2EncryptionCertificateObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"base64_certificate": types.StringType,
		"dn":                 types.StringType,
		"valid_from":         types.StringType,
		"valid_to":           types.StringType,
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
								"saml2_config": schema.SingleNestedAttribute{
									MarkdownDescription: "Configure a SAML 2.0 service provider by providing the necessary metadata.",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"saml_metadata_url": schema.StringAttribute{
											MarkdownDescription: "The URL with service provider metadata. The metadata URL must not contain a query parameter.",
											Computed:            true,
										},
										"acs_endpoints": schema.ListNestedAttribute{
											MarkdownDescription: "Configure the allowed domains for browser flows.",
											Computed:            true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"binding_name": schema.StringAttribute{
														MarkdownDescription: "Specify the SAML binding for the endpoint.",
														Computed:            true,
													},
													"location": schema.StringAttribute{
														MarkdownDescription: "The value of the URL or endpoint to be called.",
														Computed:            true,
													},
													"index": schema.Int32Attribute{
														MarkdownDescription: "A unique index for the endpoint.",
														Computed:            true,
													},
													"default": schema.BoolAttribute{
														MarkdownDescription: "Configure if the endpoint is the default one to be used.",
														Computed:            true,
													},
												},
											},
										},
										"slo_endpoints": schema.ListNestedAttribute{
											MarkdownDescription: "Configure the URLs of the service provider's single logout endpoints that will receive the logout response or request from Identity Authentication.",
											Computed:            true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"binding_name": schema.StringAttribute{
														MarkdownDescription: "Specify the SAML binding for the endpoint.",
														Computed:            true,
													},
													"location": schema.StringAttribute{
														MarkdownDescription: "The value of the URL or endpoint to be called.",
														Computed:            true,
													},
													"response_location": schema.StringAttribute{
														MarkdownDescription: "The URL or endpoint to which logout response messages are sent.",
														Computed:            true,
													},
												},
											},
										},
										"signing_certificates": schema.ListNestedAttribute{
											MarkdownDescription: "Base64-encoded certificates used by the service provider to sign digitally, SAML protocol messages sent to Identity Authentication. A maximum of 2 certificates are allowed.",
											Computed:            true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"base64_certificate": schema.StringAttribute{
														MarkdownDescription: "The content of the Base64 certificate. The certificate must be in PEM format.",
														Computed:            true,
													},
													"default": schema.BoolAttribute{
														MarkdownDescription: "Configure if the certificate is the default one to be used.",
														Computed:            true,
													},
													"dn": schema.StringAttribute{
														MarkdownDescription: "A unique identifier for the certificate.",
														Computed:            true,
													},
													"valid_from": schema.StringAttribute{
														MarkdownDescription: "Set the date from which the certificate is valid.",
														Computed:            true,
													},
													"valid_to": schema.StringAttribute{
														MarkdownDescription: "Set the date uptil which the certificate is valid.",
														Computed:            true,
													},
												},
											},
										},
										"encryption_certificate": schema.SingleNestedAttribute{
											MarkdownDescription: "The certificate used for encryption of SAML2 requests and responses.",
											Computed:            true,
											Attributes: map[string]schema.Attribute{
												"base64_certificate": schema.StringAttribute{
													MarkdownDescription: "The content of the Base64 certificate. The certificate must be in PEM format.",
													Computed:            true,
												},
												"dn": schema.StringAttribute{
													MarkdownDescription: "A unique identifier for the certificate.",
													Computed:            true,
												},
												"valid_from": schema.StringAttribute{
													MarkdownDescription: "Set the date from which the certificate is valid.",
													Computed:            true,
												},
												"valid_to": schema.StringAttribute{
													MarkdownDescription: "Set the date uptil which the certificate is valid.",
													Computed:            true,
												},
											},
										},
										"response_elements_to_encrypt": schema.StringAttribute{
											MarkdownDescription: "Specify which SAML response elements should be encrypted. " + utils.ValidValuesString(responseElementsToEncrypt),
											Computed:            true,
										},
										"default_name_id_format": schema.StringAttribute{
											MarkdownDescription: "Configure the default Name ID format. The attribute is sent as name ID format in SAML 2.0 authentication requests to Identity Provider.",
											Computed:            true,
										},
										"sign_slo_messages": schema.BoolAttribute{
											MarkdownDescription: "Enable if the single logout messages must be signed or not.",
											Computed:            true,
										},
										"require_signed_slo_messages": schema.BoolAttribute{
											MarkdownDescription: "Enable if the single logout messages must be signed or not.",
											Computed:            true,
										},
										"require_signed_auth_requests": schema.BoolAttribute{
											MarkdownDescription: "Enable if the authentication request must be signed or not.",
											Computed:            true,
										},
										"sign_assertions": schema.BoolAttribute{
											MarkdownDescription: "Enable if the SAML assertions must be signed or not.",
											Computed:            true,
										},
										"sign_auth_responses": schema.BoolAttribute{
											MarkdownDescription: "Enable if the SAML authentication responses must be signed or not.",
											Computed:            true,
										},
										"digest_algorithm": schema.StringAttribute{
											MarkdownDescription: "Configure the algorithm for signing outgoing messages. " + utils.ValidValuesString(digestAlgorithmValues),
											Computed:            true,
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
