package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	cli *cli.IasClient
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique ID of the resource.",
				Required: true,
				Validators: []validator.String{
					ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				// MarkdownDescription: "",
				ElementType: types.StringType,
				Computed: true,
			},
			"user_name": schema.StringAttribute{
				MarkdownDescription: "Unique user name of the user.",
				Computed: true,
			},
			"name": schema.SingleNestedAttribute{
				MarkdownDescription: "Name of the user",
				Attributes: map[string]schema.Attribute{
					"family_name": schema.StringAttribute{
						MarkdownDescription: "The following characters: <, >, : are not allowed.",
						Computed: true,
					},
					"given_name": schema.StringAttribute{
						MarkdownDescription: "The following characters: <, >, : are not allowed.",
						Computed: true,
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
							Computed: true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the email of the user.",
							Computed: true,
						},
						"display": schema.StringAttribute{
							Computed: true,
						},
						"primary": schema.BoolAttribute{
							MarkdownDescription: "Set the email to be primary",
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to be set for the user.",
				Computed: true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The name to be displayed for the user.",
				Computed: true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The title to be given for the user.",
				Computed: true,
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "Specifies the type of the user.The default type is \"public\".",
				Computed: true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Determines whether the user is active or not.The default value for the attribute is false.",
				Computed: true,
			},
			"send_mail": schema.BoolAttribute{
				MarkdownDescription: "Specifies if an activation mail should be sent. The value of the attribute only matters when creating the user.",
				Computed: true,
			},
			"mail_verified": schema.BoolAttribute{
				MarkdownDescription: "The attribute specifies if the e-mail of the newly created user is verified or not. So if the values of the \"mail_verified\" and \"send_mail\" attributes are true, the user will receive e-mail and they will be able to log on. On the other hand, if the \"send_mail\" is true, but the \"mail_verified\" is false, the user will receive e-mail and they have to click the verification link in the e-mail. If the attribute \"verified\" is not passed in the request body, the default value of \"mail_erified\" is false.",
				Computed: true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Specifies if the user is created as active, inactive or new. If the attribute \"active\" is not passed in the request body, the default value of the attribute \"status\" is inactive.",
				Optional: true,
				Computed: true,
			},
		},	

	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config userData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("User ID is missing", "Please provide a valid ID")
		return
	}

	res, err := d.cli.User.GetByUserId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", fmt.Sprintf("%s", err))
		return
	}

	state, diags := userValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
