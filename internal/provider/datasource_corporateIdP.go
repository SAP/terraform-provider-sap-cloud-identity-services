package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newCorporateIdPDataSource() datasource.DataSource {
	return &corporateIdPDataSource{}
}

type corporateIdPDataSource struct {
	cli *cli.SciClient
}

func (d *corporateIdPDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (r *corporateIdPDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_corporate_idp"
}

func (d *corporateIdPDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Corporate Identity Provider from the SAP Cloud Identity Services.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the Corporate Identity Provider",
				Required:            true,
				Validators: []validator.String{
					utils.ValidUUID(),
				},
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
	}
}

func (d *corporateIdPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config corporateIdPData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := d.cli.CorporateIdP.GetByIdPId(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Corporate Identity Provider", fmt.Sprintf("%s", err))
		return
	}

	state, diags := corporateIdPValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}
