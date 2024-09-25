package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"

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

type usersData struct{
	Values 	types.List 		`tfsdk:"values"`
}

var nameObjType =  map[string]attr.Type{
	"family_name": types.StringType,
	"given_name": types.StringType,
	"formatted": types.StringType,
	"middle_name": types.StringType,
	"honoric_prefix": types.StringType,
	"honoric_suffix": types.StringType,
}

var emailObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"value": types.StringType,
		"type": types.StringType,
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
		"password": types.StringType,
		"display_name": types.StringType,
		"title": types.StringType,
		"user_type": types.StringType,
		"active": types.BoolType,
		"send_mail": types.BoolType,
		"mail_verified": types.BoolType,
		"status": types.StringType,
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

		Attributes: map[string]schema.Attribute{
			
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
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

	res, err := d.cli.Directory.User.Get(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving users", fmt.Sprintf("%s", err))
		return
	}

	resUsers := usersValueFrom(ctx, res)

	config.Values, diags = types.ListValueFrom(ctx, userObjType, resUsers)
	resp.Diagnostics.Append(diags...)
	
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
