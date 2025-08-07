package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	sourceValues                        = []string{"Identity Directory", "Corporate Identity Provider", "Expression"}
	ssoValues                           = []string{"openIdConnect", "saml2"}
	usersTypeValues                     = []string{"public", "employee", "customer", "partner", "external", "onboardee"}
	subjectNameIdentifierFunctionValues = []string{"none", "upperCase", "lowerCase"}
	maxExchangePeriodValues             = []string{"unlimited", "maxSessionValidity", "initialRefreshTokenValidity"}
	refreshTokenRotationScenarioValues  = []string{"off", "online", "mobile"}
	accessTokenFormatValues             = []string{"default", "jwt", "opaque"}
	restrictedGrantTypesValues          = []string{"clientCredentials", "authorizationCode", "refreshToken", "password", "implicit", "jwtBearer", "authorizationCodePkceS256", "tokenExchange"}
)

func newApplicationResource() resource.Resource {
	return &applicationResource{}
}

type applicationResource struct {
	cli *cli.SciClient
}

func (d *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates an application in the SAP Cloud Identity Services.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the application",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					utils.ValidUUID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the application",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Free text description of the Application",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"parent_application_id": schema.StringAttribute{
				MarkdownDescription: "ID of the parent, from which the application will inherit its configurations",
				Optional:            true,
				Validators: []validator.String{
					utils.ValidUUID(),
				},
			},
			"multi_tenant_app": schema.BoolAttribute{
				MarkdownDescription: "Only for Internal Use",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"authentication_schema": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure attributes particular to the schema \"urn:sap:identity:application:schemas:extension:sci:1.0:Authentication\"",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"sso_type": schema.StringAttribute{
						MarkdownDescription: "The preferred protocol for the application. " + utils.ValidValuesString(ssoValues),
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(ssoValues...),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"subject_name_identifier": schema.SingleNestedAttribute{
						MarkdownDescription: "The attribute by which the application uses to identify the users. Used by the application to uniquely identify the user during logon.\n" +
							fmt.Sprintln("Identity Authentication sends the attribute to the application as :") +
							fmt.Sprintln("\t - subject in OpenID Connect tokens") +
							fmt.Sprintln("\t - name ID in SAML 2.0 assertions"),
						Optional: true,
						Computed: true,
						Validators: []validator.Object{
							objectvalidator.AlsoRequires(
								path.MatchRoot("authentication_schema").AtName("subject_name_identifier").AtName("source"),
								path.MatchRoot("authentication_schema").AtName("subject_name_identifier").AtName("value"),
							),
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"source": schema.StringAttribute{
								MarkdownDescription: utils.ValidValuesString(sourceValues),
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(sourceValues...),
								},
							},
							"value": schema.StringAttribute{
								MarkdownDescription: "If the source is Identity Directory, the only acceptable values are `none`, `uid`, `mail`, `loginName`, `displayName`, `personnelNumber`, `userUuid`",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 255),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"subject_name_identifier_function": schema.StringAttribute{
						MarkdownDescription: "Convert the subject name identifier to uppercase or lowercase. " + utils.ValidValuesString(subjectNameIdentifierFunctionValues),
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(subjectNameIdentifierFunctionValues...),
						},
					},
					"assertion_attributes": schema.ListNestedAttribute{
						MarkdownDescription: "User attributes to be sent to the application. The Source of these attributes is always the Identity Directory, thus only valid attribute values will be accepted.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("authentication_schema").AtName("assertion_attributes").AtAnyListIndex().AtName("attribute_name"),
								path.MatchRoot("authentication_schema").AtName("assertion_attributes").AtAnyListIndex().AtName("attribute_value"),
							),
						},
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_name": schema.StringAttribute{
									MarkdownDescription: "Name of the attribute",
									Optional:            true,
									Computed:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 255),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"attribute_value": schema.StringAttribute{
									MarkdownDescription: "Value of the attribute.",
									Optional:            true,
									Computed:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 255),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"inherited": schema.BoolAttribute{
									MarkdownDescription: "Indicates whether the attribute has been inherited from a parent application.",
									Computed:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"advanced_assertion_attributes": schema.ListNestedAttribute{
						MarkdownDescription: "Identical to the assertion attributes, except that the assertion attributes can come from other Sources.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("authentication_schema").AtName("advanced_assertion_attributes").AtAnyListIndex().AtName("source"),
								path.MatchRoot("authentication_schema").AtName("advanced_assertion_attributes").AtAnyListIndex().AtName("attribute_name"),
								path.MatchRoot("authentication_schema").AtName("advanced_assertion_attributes").AtAnyListIndex().AtName("attribute_value"),
							),
							listvalidator.SizeAtLeast(1),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"source": schema.StringAttribute{
									MarkdownDescription: utils.ValidValuesString(sourceValues[1:]),
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(sourceValues[1:]...),
									},
								},
								"attribute_name": schema.StringAttribute{
									MarkdownDescription: "Name of the attribute",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 512),
									},
								},
								"attribute_value": schema.StringAttribute{
									MarkdownDescription: "Value of the attribute",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 512),
									},
								},
								"inherited": schema.BoolAttribute{
									MarkdownDescription: "Indicates whether the attribute has been inherited from a parent application.",
									Computed:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"default_authenticating_idp": schema.StringAttribute{
						MarkdownDescription: "A default identity provider can be used for users with any user domain, group and type. This identity provider is used when none of the defined authentication rules meets the criteria.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"conditional_authentication": schema.ListNestedAttribute{
						MarkdownDescription: "Define rules for authenticating identity provider according to email domain, user type, user group, and IP range. Each rule is evaluated by priority until the criteria of a rule are fulfilled.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("authentication_schema").AtName("conditional_authentication").AtAnyListIndex().AtName("identity_provider_id"),
							),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identity_provider_id": schema.StringAttribute{
									MarkdownDescription: "The identity provider to delegate authentication to when all the defined conditions are met.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 255),
									},
								},
								"user_type": schema.StringAttribute{
									MarkdownDescription: "The type of user to be authenticated." + utils.ValidValuesString(usersTypeValues),
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(usersTypeValues...),
										stringvalidator.AtLeastOneOf(
											path.MatchRoot("authentication_schema").AtName("conditional_authentication").AtAnyListIndex().AtName("user_group"),
											path.MatchRoot("authentication_schema").AtName("conditional_authentication").AtAnyListIndex().AtName("user_email_domain"),
											path.MatchRoot("authentication_schema").AtName("conditional_authentication").AtAnyListIndex().AtName("ip_network_range"),
										),
									},
								},
								"user_group": schema.StringAttribute{
									MarkdownDescription: "The user group to be authenticated.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 255),
									},
								},
								"user_email_domain": schema.StringAttribute{
									MarkdownDescription: "Valid email domain to be authenticated.",
									Optional:            true,
									Validators: []validator.String{
										utils.ValidEmailDomain(),
									},
								},
								"ip_network_range": schema.StringAttribute{
									MarkdownDescription: "Valid IP range to be authenticated.",
									Optional:            true,
									Validators: []validator.String{
										utils.ValidIPAddress(),
									},
								},
							},
						},
					},
					"openid_connect_configuration": schema.SingleNestedAttribute{
						MarkdownDescription: "OpenID Connect (OIDC) configuration options for this application.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
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
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
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
										MarkdownDescription: "The format of the access token issued.",
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
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
							"proxy_config": schema.SingleNestedAttribute{
								MarkdownDescription: "Optional proxy configuration including accepted ACR values.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
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
	}
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan applicationData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	args, diags := getApplicationRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Application.Create(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error creating application", fmt.Sprintf("%s", err))
		return
	}

	state, diags := applicationValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var config applicationData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Application.GetByAppId(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application", fmt.Sprintf("%s", err))
		return
	}

	state, diags := applicationValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan applicationData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the current state to get the existing application ID
	var state applicationData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() {
		resp.Diagnostics.AddError("Application ID is missing", "Please provide a valid ID")
		return
	}

	// Update the application details
	args, diags := getApplicationRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	args.Id = state.Id.ValueString()

	res, _, err := r.cli.Application.Update(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := applicationValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var config applicationData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Application ID is missing", "Please provide a valid ID")
		return
	}

	err := r.cli.Application.Delete(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error deleting application", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// retrieve the API Request body from the plan data
