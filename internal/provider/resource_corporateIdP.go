package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	idpTypeValues         = []string{"sapSSO", "microsoftADFS", "saml2", "openIdConnect"}
	loginHintTypeValues   = []string{"none", "userInput", "mail", "loginName"}
	sendMethodValues      = []string{"urlParam", "authRequest"}
	digestAlgorithmValues = []string{"sha1", "sha256", "sha512"}
	nameIdFormatValues    = []string{"default", "none", "unspecified", "email"}
	allowCreateValues     = []string{"default", "none", "true", "false"}
	endpointBindingValues = []string{"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST", "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect", "urn:oasis:names:tc:SAML:2.0:bindings:SOAP", "urn:oasis:names:tc:SAML:2.0:bindings:URI"}
)

func newCorporateIdPResource() resource.Resource {
	return &corporateIdPResource{}
}

type corporateIdPResource struct {
	cli *cli.SciClient
}

func (r *corporateIdPResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.cli = req.ProviderData.(*cli.SciClient)
}

func (r *corporateIdPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_corporate_idp"
}

func (r *corporateIdPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a Corporate Identity Provider in the SAP Cloud Identity Services.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the Corporate Identity Provider",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the Corporate Identity Provider",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the Corporate Identity Provider",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the Corporate Identity Provider. " + utils.ValidValuesString(idpTypeValues),
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(idpTypeValues...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"logout_url": schema.StringAttribute{
				MarkdownDescription: "URL to redirect users after successful logout.",
				Optional:            true,
				Validators: []validator.String{
					utils.ValidUrl(),
				},
			},
			"forward_all_sso_requests": schema.BoolAttribute{
				MarkdownDescription: "If set to true, all authentication requests will be sent to this corporate IdP when it is chosen as the default identity provider.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_federation": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure how the user and user attributes are handled when authenticating via the Corporate Identity Provider.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"use_local_user_store": schema.BoolAttribute{
						MarkdownDescription: `Configure if user attributes will be taken from the corporate IdP assertion or from Identity Authentication user store.
							By default, Identity Authentication takes all assertion attributes and Subject Name Identifier from the corporate IdP assertion and sends them to the application. 
							If set to true, data from Identity Authentication user store will be used. For users with no profile in Identity Authentication, the application receives the subject name identifier from the corporate IdP assertion and attributes according to the application configuration.`,
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"allow_local_users_only": schema.BoolAttribute{
						MarkdownDescription: `Restrict access based on user profile. By default, all users successfully authenticated to corporate IdP are allowed.
							If set to true, only users with profiles in Identity Authentication are allowed access.`,
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"apply_local_idp_auth_and_checks": schema.BoolAttribute{
						MarkdownDescription: "Configure if local authentication and access policies must be applied if users with profiles in Identity Authentication are authenticated via corporate IdP.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"required_groups": schema.SetAttribute{
						MarkdownDescription: "Restrict access to specific user groups. Only members of these groups will be allowed to access applications after successful authentication to the corporate IdP.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"login_hint_config": schema.SingleNestedAttribute{
				MarkdownDescription: `Configure the value of the login hint attribute and how it is sent to the corporate IdP.
					 This parameter helps the user when they are known to the service provider (SP) or relying party. Thus it prevents the user from re-typing the user identifier on the logon. 
					 If the corporate IdP supports the login hint parameter, then it requests only the user credentials.`,
				Optional: true,
				Computed: true,
				Validators: []validator.Object{
					objectvalidator.AlsoRequires(
						path.MatchRoot("login_hint_config").AtName("login_hint_type"),
					),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"login_hint_type": schema.StringAttribute{
						MarkdownDescription: "The value of the parameter sent. " + utils.ValidValuesString(loginHintTypeValues),
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(loginHintTypeValues...),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"send_method": schema.StringAttribute{
						MarkdownDescription: "Configure how this parameter is sent to the corporate IdP. " + utils.ValidValuesString(sendMethodValues),
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(sendMethodValues...),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"saml2_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure trust with an identity provider by providing the necessary metadata for web-based authentication.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Object{
					objectvalidator.AlsoRequires(
						path.MatchRoot("type"),
						path.MatchRoot("name"),
					),
					// The API does not validate the type of the corporate IdP depending on the configuration provided.
					// When the saml2 configuration is provided but the type of the IdP is set to "openIdConnect", the IdP listed on the Admin Console is of type OIDC.
					// Although the provided saml2 configurationd details are returned in the GET call, this validator ensures the consistency of the type and config provided
					objectvalidator.All(
						utils.ValidType(
							path.MatchRoot("type"),
							idpTypeValues[:3],
						),
					),
				},
				Attributes: map[string]schema.Attribute{
					"saml_metadata_url": schema.StringAttribute{
						MarkdownDescription: "The URL with identity provider metadata.",
						Optional:            true,
						Validators: []validator.String{
							utils.ValidUrl(),
						},
					},
					"assertion_attributes": schema.ListNestedAttribute{
						MarkdownDescription: "Enrich the assertion attributes coming from the corporate IdP.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("saml2_config").AtName("assertion_attributes").AtAnyListIndex().AtName("name"),
								path.MatchRoot("saml2_config").AtName("assertion_attributes").AtAnyListIndex().AtName("value"),
							),
							listvalidator.SizeAtLeast(1),
							listvalidator.SizeAtMost(30),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Set the name of the attribute.",
									Optional:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "Set the value of the attribute.",
									Optional:            true,
								},
							},
						},
					},
					"signing_certificates": schema.ListNestedAttribute{
						MarkdownDescription: "Base64-encoded certificates used by the service provider to sign digitally, SAML protocol messages sent to Identity Authentication. A maximum of 2 certificates are allowed.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.SizeAtMost(2),
							listvalidator.AlsoRequires(
								path.MatchRoot("saml2_config").AtName("signing_certificates").AtAnyListIndex().AtName("base64_certificate"),
								path.MatchRoot("saml2_config").AtName("signing_certificates").AtAnyListIndex().AtName("default"),
							),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"base64_certificate": schema.StringAttribute{
									MarkdownDescription: "The content of the Base64 certificate. The certificate must be in PEM format.",
									Optional:            true,
									// Validator enforces the presence of PEM boundary markers ("-----BEGIN CERTIFICATE-----" / "-----END CERTIFICATE-----").
									// The API accepts certificates without markers and returns the response by wrapping them
									// which can result in a difference between the planned value and the state.
									// Hence the validation maintains consistency of the parameter
									Validators: []validator.String{
										utils.ValidCertificate(),
									},
								},
								"dn": schema.StringAttribute{
									MarkdownDescription: "A unique identifier for the certificate.",
									Optional:            true,
								},
								"default": schema.BoolAttribute{
									MarkdownDescription: "Configure if the certificate is the default one to be used.",
									Optional:            true,
								},
								// Validator enforces that the date-time string is in UTC format with no offset.
								// The API accepts numeric values and offset-based timestamps, then normalizes them to UTC
								// without preserving the original format. This can cause a difference between the planned and actual state.
								// Validation ensures consistent formatting of the parameter.
								"valid_from": schema.StringAttribute{
									MarkdownDescription: "Set the date from which the certificate is valid.",
									Optional:            true,
									Validators: []validator.String{
										utils.ValidDateTime(),
									},
								},
								"valid_to": schema.StringAttribute{
									MarkdownDescription: "Set the date uptil which the certificate is valid.",
									Optional:            true,
									Validators: []validator.String{
										utils.ValidDateTime(),
									},
								},
							},
						},
					},
					"sso_endpoints": schema.ListNestedAttribute{
						MarkdownDescription: "Configure the URLs of the identity provider single sign-on endpoint that receive authentication requests.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("saml2_config").AtName("sso_endpoints").AtAnyListIndex().AtName("binding_name"),
								path.MatchRoot("saml2_config").AtName("sso_endpoints").AtAnyListIndex().AtName("location"),
							),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding_name": schema.StringAttribute{
									MarkdownDescription: "Specify the SAML binding for the endpoint. " + utils.ValidValuesString(endpointBindingValues),
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(endpointBindingValues...),
									},
								},
								"location": schema.StringAttribute{
									MarkdownDescription: "The value of the URL or endpoint to be called.",
									Optional:            true,
									Validators: []validator.String{
										utils.ValidUrl(),
									},
								},
								"default": schema.BoolAttribute{
									MarkdownDescription: "Configure if the endpoint is the default one to be used.",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"slo_endpoints": schema.ListNestedAttribute{
						MarkdownDescription: "Configure the URLs of the identity provider single logout endpoint that receive logout messages.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.AlsoRequires(
								path.MatchRoot("saml2_config").AtName("slo_endpoints").AtAnyListIndex().AtName("binding_name"),
								path.MatchRoot("saml2_config").AtName("slo_endpoints").AtAnyListIndex().AtName("location"),
							),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding_name": schema.StringAttribute{
									MarkdownDescription: "Specify the SAML binding for the endpoint. " + utils.ValidValuesString(endpointBindingValues),
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(endpointBindingValues...),
									},
								},
								"location": schema.StringAttribute{
									MarkdownDescription: "The value of the URL or endpoint to be called.",
									Optional:            true,
									Validators: []validator.String{
										utils.ValidUrl(),
									},
								},
								"response_location": schema.StringAttribute{
									MarkdownDescription: "The URL or endpoint to which logout response messages are sent.",
									Optional:            true,
									Validators: []validator.String{
										utils.ValidUrl(),
									},
								},
								"default": schema.BoolAttribute{
									MarkdownDescription: "Configure if the endpoint is the default one to be used",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"digest_algorithm": schema.StringAttribute{
						MarkdownDescription: "Configure the Signing Algorithm. " + utils.ValidValuesString(digestAlgorithmValues),
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf(digestAlgorithmValues...),
						},
					},
					"include_scoping": schema.BoolAttribute{
						MarkdownDescription: "Configure whether to include or exclude the Scoping element in the SAML 2.0 request.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"name_id_format": schema.StringAttribute{
						MarkdownDescription: "Configure preferred Name ID format. The attribute is sent to the corporate identity provider as name ID format to the Identity Provider. " + utils.ValidValuesString(nameIdFormatValues),
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf(nameIdFormatValues...),
						},
					},
					"allow_create": schema.StringAttribute{
						MarkdownDescription: "Configure if the `AllowCreate` attribute sent by the Service Provider is forwarded to the Corporate IdP or not. " + utils.ValidValuesString(allowCreateValues),
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf(allowCreateValues...),
						},
					},
				},
			},
		},
	}
}

func (r *corporateIdPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan corporateIdPData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	args, diags := r.getCorporateIdPRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.CorporateIdP.Create(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Corporate Identity Provider", fmt.Sprintf("%s", err))
		return
	}

	state, diags := corporateIdPValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

// TODO implement the Update operation once the API is available
func (r *corporateIdPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Error updating Corporate Identity Provider.", "This resource does not support updates")
}

func (r *corporateIdPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var config corporateIdPData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.CorporateIdP.GetByIdPId(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Corporate Identity Provider", fmt.Sprintf("%s", err))
		return
	}

	state, diags := corporateIdPValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

func (r *corporateIdPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var config corporateIdPData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.cli.CorporateIdP.Delete(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Corporate Identity Provider", fmt.Sprintf("%s", err))
		return
	}
}

func (r *corporateIdPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *corporateIdPResource) getCorporateIdPRequest(ctx context.Context, plan corporateIdPData) (*corporateidps.IdentityProvider, diag.Diagnostics) {
	var diags, diagnostics diag.Diagnostics

	corporateIdP := &corporateidps.IdentityProvider{
		DisplayName: plan.DisplayName.ValueString(),
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		corporateIdP.Name = plan.Name.ValueString()
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		corporateIdP.Type = plan.Type.ValueString()
	}

	if !plan.LogoutUrl.IsNull() && !plan.LogoutUrl.IsUnknown() {
		corporateIdP.LogoutUrl = plan.LogoutUrl.ValueString()
	}

	if !plan.ForwardAllSsoRequests.IsNull() && !plan.ForwardAllSsoRequests.IsUnknown() {
		corporateIdP.ForwardAllSsoRequests = plan.ForwardAllSsoRequests.ValueBool()
	}

	if !plan.IdentityFederation.IsNull() && !plan.IdentityFederation.IsUnknown() {

		var idF corporateidps.IdentityFederation
		diags = plan.IdentityFederation.As(ctx, &idF, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		corporateIdP.IdentityFederation = &idF
	}

	if !plan.LoginHintConfig.IsNull() && !plan.LoginHintConfig.IsUnknown() {

		var loginHintConfig corporateidps.LoginHintConfiguration
		diags = plan.LoginHintConfig.As(ctx, &loginHintConfig, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return &corporateidps.IdentityProvider{}, diagnostics
		}

		corporateIdP.LoginHintConfiguration = &loginHintConfig

	}

	if !plan.Saml2Config.IsNull() && !plan.Saml2Config.IsUnknown() {

		var saml2Config corporateidps.SAML2Configuration
		diags = plan.Saml2Config.As(ctx, &saml2Config, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		corporateIdP.Saml2Configuration = &saml2Config
	}

	return corporateIdP, diagnostics
}
