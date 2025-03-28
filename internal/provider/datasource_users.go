package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"
	"terraform-provider-ias/internal/utils"

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
	cli *cli.IasClient
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
	"family_name":    types.StringType,
	"given_name":     types.StringType,
	"formatted":      types.StringType,
	"middle_name":    types.StringType,
	"honoric_prefix": types.StringType,
	"honoric_suffix": types.StringType,
}

var emailObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"value":   types.StringType,
		"type":    types.StringType,
		"display": types.StringType,
		"primary": types.BoolType,
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
		"password":     types.StringType,
		"display_name": types.StringType,
		"title":        types.StringType,
		"user_type":    types.StringType,
		"active":       types.BoolType,
		"sap_extension_user": types.ObjectType{
			AttrTypes: sapExtensionUserObjType,
		},
		"custom_schemas": types.StringType,
	},
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
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
							Computed: true,
							Validators: []validator.String{
								utils.ValidUUID(),
							},
						},
						"schemas": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"user_name": schema.StringAttribute{
							MarkdownDescription: "Unique user name of the user.",
							Computed:            true,
						},
						"name": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"family_name": schema.StringAttribute{
									MarkdownDescription: "The following characters: <, >, : are not allowed.",
									Computed:            true,
								},
								"given_name": schema.StringAttribute{
									MarkdownDescription: "The following characters: <, >, : are not allowed.",
									Computed:            true,
								},
								"formatted": schema.StringAttribute{
									Computed: true,
								},
								"middle_name": schema.StringAttribute{
									Computed: true,
								},
								"honoric_prefix": schema.StringAttribute{
									Computed: true,
								},
								"honoric_suffix": schema.StringAttribute{
									Computed: true,
								},
							},
							Computed: true,
						},
						"emails": schema.SetNestedAttribute{
							MarkdownDescription: "Email of the user.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										MarkdownDescription: "Value of the email of the user.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										MarkdownDescription: "Type of the email of the user.",
										Computed:            true,
									},
									"display": schema.StringAttribute{
										Computed: true,
									},
									"primary": schema.BoolAttribute{
										MarkdownDescription: "Set the email to be primary or not.",
										Computed:            true,
									},
								},
							},
							Computed: true,
						},
						"password": schema.StringAttribute{
							MarkdownDescription: "The password to be set for the user.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The name to be displayed for the user.",
							Computed:            true,
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "The title to be given for the user.",
							Computed:            true,
						},
						"user_type": schema.StringAttribute{
							MarkdownDescription: "Specifies the type of the user.The default type is \"public\".",
							Computed:            true,
						},
						"active": schema.BoolAttribute{
							MarkdownDescription: "Determines whether the user is active or not.The default value for the attribute is false.",
							Computed:            true,
						},
						"sap_extension_user": schema.SingleNestedAttribute{
							// MarkdownDescription
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"send_mail": schema.BoolAttribute{
									MarkdownDescription: "Specifies if an activation mail should be sent. The value of the attribute only matters when creating the user.",
									Computed:            true,
								},
								"mail_verified": schema.BoolAttribute{
									MarkdownDescription: "The attribute specifies if the e-mail of the newly created user is verified or not. So if the values of the \"mail_verified\" and \"send_mail\" attributes are true, the user will receive e-mail and they will be able to log on. On the other hand, if the \"send_mail\" is true, but the \"mail_verified\" is false, the user will receive e-mail and they have to click the verification link in the e-mail. If the attribute \"verified\" is not passed in the request body, the default value of \"mail_erified\" is false.",
									Computed:            true,
								},
								"status": schema.StringAttribute{
									MarkdownDescription: "Specifies if the user is created as active, inactive or new. If the attribute \"active\" is not passed in the request body, the default value of the attribute \"status\" is inactive.",
									Computed:            true,
								},
							},
						},
						"custom_schemas": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Furthur enhance the user created with custom schemas.",
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

	res, customSchemasRes, err := d.cli.User.Get(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving users", fmt.Sprintf("%s", err))
		return
	}

	resUsers := usersValueFrom(ctx, res, customSchemasRes)

	config.Values, diags = types.ListValueFrom(ctx, userObjType, resUsers)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
