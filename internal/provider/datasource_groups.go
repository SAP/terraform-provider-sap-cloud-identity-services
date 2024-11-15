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

func newGroupsDataSource() datasource.DataSource{
	return &groupsDataSource{}
}

var membersObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"value": types.StringType,
		"type": types.StringType,
	},
}

var groupObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id": types.StringType,
		"schemas": types.SetType{
			ElemType: types.StringType,
		},
		"display_name": types.StringType,
		"group_members": types.ListType{
			ElemType: membersObjType,
		},
		"external_id": types.StringType,
		"description": types.StringType,
	},	
}

type groupsDataSource struct {
	cli *cli.IasClient
}

func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) { 
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (d *groupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) { 
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { 
	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{
			"values" : schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							// MarkdownDescription: ,
							Validators: []validator.String{
								ValidUUID(),
							},
						},
						"schemas": schema.SetAttribute{
							ElementType: types.StringType,
							Computed: true,
							//MarkdownDescription:
						},
						"display_name": schema.StringAttribute{
							Computed: true,
							// MarkdownDescription: ,
						},
						"group_members": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										Computed: true,
										// MarkdownDescription: ,
									},
									"type": schema.StringAttribute{
										Computed: true,
										// MarkdownDescription: ,
									},
								},
							},
							Computed: true,
						},
						"external_id": schema.StringAttribute{
							Computed: true,
							// MarkdownDescription: ,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},

	}
}

func (d *groupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { 
	
	var config groupsData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, err := d.cli.Group.Get(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving groups", fmt.Sprintf("%s",err))
		return
	}

	resGroups, diags := groupsValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	config.Values, diags = types.ListValueFrom(ctx, groupObjType, resGroups)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
