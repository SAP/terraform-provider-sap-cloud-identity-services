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
				Required: true,
				Validators: []validator.String{
					ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				ElementType: types.StringType,
				Computed: true,
			},
			"user_name": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"family_name": schema.StringAttribute{
						Computed: true,
					},
					"given_name": schema.StringAttribute{
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
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"display": schema.StringAttribute{
							Computed: true,
						},
						"primary": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"password": schema.StringAttribute{
				Computed: true,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
			},
			"title": schema.StringAttribute{
				Computed: true,
			},
			"user_type": schema.StringAttribute{
				Computed: true,
			},
			"active": schema.BoolAttribute{
				Computed: true,
			},
			"send_mail": schema.BoolAttribute{
				Computed: true,
			},
			"mail_verified": schema.BoolAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
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
