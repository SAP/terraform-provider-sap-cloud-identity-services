package provider

import (
	"context"
	"fmt"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func newApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

type applicationDataSource struct {
	cli *cli.SciClient
}

func (d *applicationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *applicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *applicationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets an application from the SAP Cloud Identity services.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the application",
				Required:            true,
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
							fmt.Sprintf(`Identity Authentication sends the attribute to the application as : 
							                     - subject in OpenID Connect tokens
							                     - name ID in SAML 2.0 assertions`),
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
				},
			},
		},
	}
}

func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config applicationData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Application ID is missing", "Please provide a valid ID")
		return
	}

	res, _, err := d.cli.Application.GetByAppId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application", fmt.Sprintf("%s", err))
		return
	}

	state, _ := applicationValueFrom(ctx, res)
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
}
