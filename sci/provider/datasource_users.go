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

func newUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	cli *cli.SciClient
}

type usersData struct {
	Values types.List `tfsdk:"values"`
}

var sapExtensionUserObjType = map[string]attr.Type{
	"send_mail":     types.BoolType,
	"mail_verified": types.BoolType,
	"status":        types.StringType,
}

var nameObjType = map[string]attr.Type{
	"family_name":      types.StringType,
	"given_name":       types.StringType,
	"honorific_prefix": types.StringType,
}

var emailObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"value":   types.StringType,
		"type":    types.StringType,
		"primary": types.BoolType,
	},
}

var groupListObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type":    types.StringType,
		"value":   types.StringType,
		"display": types.StringType,
	},
}

var userObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id": types.StringType,
		"schemas": types.SetType{
			ElemType: types.StringType,
		},
		"user_name": types.StringType,
		"name": types.ObjectType{
			AttrTypes: nameObjType,
		},
		"emails": types.SetType{
			ElemType: emailObjType,
		},
		"initial_password": types.StringType,
		"display_name":     types.StringType,
		"user_type":        types.StringType,
		"active":           types.BoolType,
		"sap_extension_user": types.ObjectType{
			AttrTypes: sapExtensionUserObjType,
		},
		"custom_schemas": types.StringType,
		"groups": types.ListType{
			ElemType: groupListObjType,
		},
	},
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a list of users from the SAP Cloud Identity services.`,
		Attributes: map[string]schema.Attribute{

			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "ID of the user.",
							Computed:            true,
							Validators: []validator.String{
								utils.ValidUUID(),
							},
						},
						"schemas": schema.SetAttribute{
							MarkdownDescription: "List of SCIM schemas to configure users. The attribute is configured with default values :\n" +
								utils.PrintDefaultSchemas(defaultUserSchemas),
							ElementType: types.StringType,
							Computed:    true,
						},
						"user_name": schema.StringAttribute{
							MarkdownDescription: "Unique user name of the user.",
							Computed:            true,
						},
						"emails": schema.SetNestedAttribute{
							MarkdownDescription: "Emails of the user.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										MarkdownDescription: "Value of the user's email.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										MarkdownDescription: "Type of the user's email.",
										Computed:            true,
									},
									"primary": schema.BoolAttribute{
										MarkdownDescription: "Set the email to be primary or not.",
										Computed:            true,
									},
								},
							},
						},
						"name": schema.SingleNestedAttribute{
							MarkdownDescription: "Name of the user",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"family_name": schema.StringAttribute{
									MarkdownDescription: "Last name of the user. The following characters: <, >, : are not allowed.",
									Computed:            true,
								},
								"given_name": schema.StringAttribute{
									MarkdownDescription: "First name of the user. The following characters: <, >, : are not allowed.",
									Computed:            true,
								},
								"honorific_prefix": schema.StringAttribute{
									MarkdownDescription: "HonorificPrefix is part of the Master Data attributes and have canonical values. The specific values for this attribute can be found on `<tenantUrl>/service/md/salutations`",
									Computed:            true,
								},
							},
						},
						"initial_password": schema.StringAttribute{
							MarkdownDescription: "The initial password to be configured for the user.",
							Computed:            true,
							Sensitive:           true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The name to be displayed for the user.",
							Computed:            true,
						},
						"user_type": schema.StringAttribute{
							MarkdownDescription: "Specifies the type of the user. The default type is \"public\".",
							Computed:            true,
						},
						"active": schema.BoolAttribute{
							MarkdownDescription: "Determines whether the user is active or not. The default value for the attribute is false.",
							Computed:            true,
						},
						"sap_extension_user": schema.SingleNestedAttribute{
							MarkdownDescription: "Configure attributes particular to the schema `" + defaultUserSchemas[1].String() + "`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"send_mail": schema.BoolAttribute{
									MarkdownDescription: "Specifies if an activation mail should be sent. The value of the attribute only matters when creating the user.",
									Computed:            true,
								},
								"mail_verified": schema.BoolAttribute{
									MarkdownDescription: "The attribute specifies if the e-mail of the newly created user is verified or not. So if the values of the \"mail_verified\" and \"send_mail\" attributes are true, the user will receive an e-mail and they will be able to log on. On the other hand, if the \"send_mail\" is true, but the \"mail_verified\" is false, the user will receive e-mail and they have to click the verification link in the e-mail. If the attribute \"mail_verified\" is not configured, the default value is false.",
									Computed:            true,
								},
								"status": schema.StringAttribute{
									MarkdownDescription: "Specifies if the user is created as active, inactive or new. If the attribute \"active\" is not configured, the default value is inactive.",
									Computed:            true,
								},
							},
						},
						"custom_schemas": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Furthur enhance your user with custom schemas.",
							Validators: []validator.String{
								utils.ValidJSON(),
							},
						},
						"groups": schema.ListNestedAttribute{
							MarkdownDescription: "The list of Groups that the user belongs to.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										MarkdownDescription: "The unique UUID of the Group.",
										Computed:            true,
									},
									"display": schema.StringAttribute{
										MarkdownDescription: "The display name of the Group.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										MarkdownDescription: "The type of the Group.",
										Computed:            true,
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

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config usersData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, customSchemasRes, err := d.cli.User.Get(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving users", fmt.Sprintf("%s", err))
		return
	}

	resUsers := usersValueFrom(ctx, res, customSchemasRes)

	config.Values, diags = types.ListValueFrom(ctx, userObjType, resUsers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
