package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newCorporateIdPsDataSource() datasource.DataSource {
	return &corporateIdPsDataSource{}
}

type corporateIdPsDataSource struct {
	cli *cli.SciClient
}

type corporateIdPsData struct {
	Values types.List `tfsdk:"values"`
}

var OidcCAdditionalConfigObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"enforce_nonce":                types.BoolType,
		"enforce_issuer_check":         types.BoolType,
		"disable_logout_id_token_hint": types.BoolType,
	},
}

var oidcConfigObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"discovery_url":              types.StringType,
		"client_id":                  types.StringType,
		"client_secret":              types.StringType,
		"subject_name_identifier":    types.StringType,
		"token_endpoint_auth_method": types.StringType,
		"scopes": types.SetType{
			ElemType: types.StringType,
		},
		"enable_pkce":                 types.BoolType,
		"additional_config":           OidcCAdditionalConfigObjType,
		"issuer":                      types.StringType,
		"jwks_uri":                    types.StringType,
		"jwks":                        types.StringType,
		"token_endpoint":              types.StringType,
		"authorization_endpoint":      types.StringType,
		"logout_endpoint":             types.StringType,
		"user_info_endpoint":          types.StringType,
		"is_client_secret_configured": types.BoolType,
	},
}

var saml2SloEndpointObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"binding_name":      types.StringType,
		"location":          types.StringType,
		"response_location": types.StringType,
		"default":           types.BoolType,
	},
}

var saml2SsoEndpointObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"binding_name": types.StringType,
		"location":     types.StringType,
		"default":      types.BoolType,
	},
}

var saml2SigningCertificateObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"base64_certificate": types.StringType,
		"dn":                 types.StringType,
		"default":            types.BoolType,
		"valid_from":         types.StringType,
		"valid_to":           types.StringType,
	},
}

var saml2AssertionAttributeObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	},
}

var saml2ConfigObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"saml_metadata_url": types.StringType,
		"assertion_attributes": types.ListType{
			ElemType: saml2AssertionAttributeObjType,
		},
		"digest_algorithm": types.StringType,
		"include_scoping":  types.BoolType,
		"name_id_format":   types.StringType,
		"allow_create":     types.StringType,
		"signing_certificates": types.ListType{
			ElemType: saml2SigningCertificateObjType,
		},
		"sso_endpoints": types.ListType{
			ElemType: saml2SsoEndpointObjType,
		},
		"slo_endpoints": types.ListType{
			ElemType: saml2SloEndpointObjType,
		},
	},
}

var loginHintConfigObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"login_hint_type": types.StringType,
		"send_method":     types.StringType,
	},
}

var identityFederationObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"use_local_user_store":            types.BoolType,
		"allow_local_users_only":          types.BoolType,
		"apply_local_idp_auth_and_checks": types.BoolType,
		"required_groups": types.SetType{
			ElemType: types.StringType,
		},
	},
}

var corporateIdPObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                       types.StringType,
		"display_name":             types.StringType,
		"name":                     types.StringType,
		"type":                     types.StringType,
		"logout_url":               types.StringType,
		"forward_all_sso_requests": types.BoolType,
		"identity_federation":      identityFederationObjType,
		"login_hint_config":        loginHintConfigObjType,
		"saml2_config":             saml2ConfigObjType,
		"oidc_config":              oidcConfigObjType,
	},
}

func (d *corporateIdPsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *corporateIdPsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_corporate_idps"
}

func (d *corporateIdPsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a list of Corporate Identity Providers from the SAP Cloud Identity services.`,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Id of the Corporate Identity Provider",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name of the Corporate Identity Provider",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Unique name of the Corporate Identity Provider",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the Corporate Identity Provider.",
							Computed:            true,
						},
						"logout_url": schema.StringAttribute{
							MarkdownDescription: "URL to redirect users after successful logout.",
							Computed:            true,
						},
						"forward_all_sso_requests": schema.BoolAttribute{
							MarkdownDescription: "If set to true, all authentication requests will be sent to this corporate IdP when it is chosen as the default identity provider.",
							Computed:            true,
						},
						"identity_federation": schema.SingleNestedAttribute{
							MarkdownDescription: "Configure how the user and user attributes are handled when authenticating via the Corporate Identity Provider.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"use_local_user_store": schema.BoolAttribute{
									MarkdownDescription: `Configure if user attributes will be taken from the corporate IdP assertion or from Identity Authentication user store.
										By default, Identity Authentication takes all assertion attributes and Subject Name Identifier from the corporate IdP assertion and sends them to the application. 
										If set to true, data from Identity Authentication user store will be used. For users with no profile in Identity Authentication, the application receives the subject name identifier from the corporate IdP assertion and attributes according to the application configuration.`,
									Computed: true,
								},
								"allow_local_users_only": schema.BoolAttribute{
									MarkdownDescription: `Restrict access based on user profile. By default, all users successfully authenticated to corporate IdP are allowed.
										If set to true, only users with profiles in Identity Authentication are allowed access.`,
									Computed: true,
								},
								"apply_local_idp_auth_and_checks": schema.BoolAttribute{
									MarkdownDescription: "Configure if local authentication and access policies must be applied if users with profiles in Identity Authentication are authenticated via corporate IdP.",
									Computed:            true,
								},
								"required_groups": schema.SetAttribute{
									MarkdownDescription: "Restrict access to specific user groups. Only members of these groups will be allowed to access applications after successful authentication to the corporate IdP.",
									Computed:            true,
									ElementType:         types.StringType,
								},
							},
						},
						"login_hint_config": schema.SingleNestedAttribute{
							MarkdownDescription: `Configure the value of the login hint attribute and how it is sent to the corporate IdP.
								This parameter helps the user when they are known to the service provider (SP) or relying party. Thus it prevents the user from re-typing the user identifier on the logon. 
								If the corporate IdP supports the login hint parameter, then it requests only the user credentials.`,
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"login_hint_type": schema.StringAttribute{
									MarkdownDescription: "The value of the parameter sent.",
									Computed:            true,
								},
								"send_method": schema.StringAttribute{
									MarkdownDescription: "Configure how this parameter is sent to the corporate IdP.",
									Computed:            true,
								},
							},
						},
						"saml2_config": schema.SingleNestedAttribute{
							MarkdownDescription: "Configure trust with an identity provider by providing the necessary metadata for web-based authentication.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"saml_metadata_url": schema.StringAttribute{
									MarkdownDescription: "The URL with identity provider metadata.",
									Computed:            true,
								},
								"assertion_attributes": schema.ListNestedAttribute{
									MarkdownDescription: "Enrich the assertion attributes coming from the corporate IdP.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												MarkdownDescription: "Set the name of the attribute.",
												Computed:            true,
											},
											"value": schema.StringAttribute{
												MarkdownDescription: "Set the value of the attribute.",
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
											"dn": schema.StringAttribute{
												MarkdownDescription: "A unique identifier for the certificate.",
												Computed:            true,
											},
											"default": schema.BoolAttribute{
												MarkdownDescription: "Configure if the certificate is the default one to be used.",
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
								"sso_endpoints": schema.ListNestedAttribute{
									MarkdownDescription: "Configure the URLs of the identity provider single sign-on endpoint that receive authentication requests.",
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
											"default": schema.BoolAttribute{
												MarkdownDescription: "Configure if the endpoint is the default one to be used.",
												Computed:            true,
											},
										},
									},
								},
								"slo_endpoints": schema.ListNestedAttribute{
									MarkdownDescription: "Configure the URLs of the identity provider single logout endpoint that receive logout messages.",
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
											"default": schema.BoolAttribute{
												MarkdownDescription: "Configure if the endpoint is the default one to be used",
												Computed:            true,
											},
										},
									},
								},
								"digest_algorithm": schema.StringAttribute{
									MarkdownDescription: "Configure the Signing Algorithm.",
									Computed:            true,
								},
								"include_scoping": schema.BoolAttribute{
									MarkdownDescription: "Configure whether to include or exclude the Scoping element in the SAML 2.0 request.",
									Computed:            true,
								},
								"name_id_format": schema.StringAttribute{
									MarkdownDescription: "Configure preferred Name ID format. The attribute is sent to the corporate identity provider as name ID format to the Identity Provider.",
									Computed:            true,
								},
								"allow_create": schema.StringAttribute{
									MarkdownDescription: "Configure if the `AllowCreate` attribute sent by the Service Provider is forwarded to the Corporate IdP or not.",
									Computed:            true,
								},
							},
						},
						"oidc_config": schema.SingleNestedAttribute{
							MarkdownDescription: "Configure trust with an identity provider by providing the necessary metadata for web-based authentication.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"discovery_url": schema.StringAttribute{
									MarkdownDescription: "Specify the Issuer or Metadata URL",
									Computed:            true,
								},
								"client_id": schema.StringAttribute{
									MarkdownDescription: "Configure the Client ID for Client Authentication.",
									Computed:            true,
								},
								"client_secret": schema.StringAttribute{
									MarkdownDescription: "Configure the Client Secret for Client Authentication.",
									Computed:            true,
								},
								"token_endpoint_auth_method": schema.StringAttribute{
									MarkdownDescription: "Configure the Client Authentication Method.",
									Computed:            true,
								},
								"subject_name_identifier": schema.StringAttribute{
									MarkdownDescription: "Define the claim which is used as subject name identifier. The Subject Name Identifier configuration defines with which value the identity provider user will be searched in the Identity Authentication user store.",
									Computed:            true,
								},
								"scopes": schema.SetAttribute{
									MarkdownDescription: "Configure additional scopes required by the Identity Provider. By default, the \"openid\" scope is added.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"enable_pkce": schema.BoolAttribute{
									MarkdownDescription: "Configure Proof Key for Code Exchange (PKCE) for the corporate IdP. This is an enhancement of the authorization code flow to prevent the interception of authorization code. This feature is recommended only if the corporate IdP supports PKCE and you have public applications that aren't capable of keeping client secrets.",
									Computed:            true,
								},
								"additional_config": schema.SingleNestedAttribute{
									MarkdownDescription: "Configure additional settings of the corporate IdP.",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"enforce_nonce": schema.BoolAttribute{
											MarkdownDescription: "Configure if the authenticating application is required to send nonces to the corporate IdP. A nonce is a string associated with a client session and is used to mitigate replay attacks. If supplied by an application, Identity Authentication forwards the nonce with requests to the corporate IdP.",
											Computed:            true,
										},
										"enforce_issuer_check": schema.BoolAttribute{
											MarkdownDescription: "Configure if Identity Authentication should enforce Issuer Validation. If set to true, responses from the corporate IdP which don't provide the iss attribute are rejected.",
											Computed:            true,
										},
										"disable_logout_id_token_hint": schema.BoolAttribute{
											MarkdownDescription: "Configure if the Identity Authentication should not include the ID token in the id_token_hint URL parameter.",
											Computed:            true,
										},
									},
								},
								"issuer": schema.StringAttribute{
									MarkdownDescription: "The unique field that identifies the IdP.",
									Computed:            true,
								},
								"jwks_uri": schema.StringAttribute{
									MarkdownDescription: "The endpoint called to request JSON Web Keys for JWT validation.",
									Computed:            true,
								},
								"jwks": schema.StringAttribute{
									MarkdownDescription: "The JSON Web Keys used for the JSON Web Token Validation.",
									Computed:            true,
								},
								"token_endpoint": schema.StringAttribute{
									MarkdownDescription: "The endpoint called to request the ID token for SSO.",
									Computed:            true,
								},
								"authorization_endpoint": schema.StringAttribute{
									MarkdownDescription: "The endpoint to which SSO requests are forwarded to, in order to retrieve an authorization code.",
									Computed:            true,
								},
								"logout_endpoint": schema.StringAttribute{
									MarkdownDescription: "The endpoint called to log out the current user session.",
									Computed:            true,
								},
								"user_info_endpoint": schema.StringAttribute{
									MarkdownDescription: "The endpoint called to get information about a user.",
									Computed:            true,
								},
								"is_client_secret_configured": schema.BoolAttribute{
									MarkdownDescription: "Indicates if a client secret is configured or not.",
									Computed:            true,
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

func (d *corporateIdPsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config corporateIdPsData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, _, err := d.cli.CorporateIdP.Get(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving corporate idps", fmt.Sprintf("%s", err))
		return
	}

	resIdPs := corporateIdPsValueFrom(ctx, res)

	config.Values, diags = types.ListValueFrom(ctx, corporateIdPObjType, resIdPs)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
