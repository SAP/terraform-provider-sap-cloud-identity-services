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
