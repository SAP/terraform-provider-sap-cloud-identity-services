package provider

import (
	"context"
	"fmt"
	"strings"
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

var groupExtensionObjType = map[string]attr.Type{
	"name": types.StringType,
	"description": types.StringType,
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
		"group_extension": types.ObjectType{
			AttrTypes: groupExtensionObjType,
		},
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
							MarkdownDescription: "Unique ID of the group.",
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
							MarkdownDescription: "Display Name of the group.",
						},
						"group_members": schema.ListNestedAttribute{
							Computed: true,
							MarkdownDescription: "Specify the members to be part of the group.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										Computed: true,
										MarkdownDescription: "SCIM ID of the user or the group",
									},
									"type": schema.StringAttribute{
										Computed: true,
										MarkdownDescription: fmt.Sprintf("Type of the member added to the group. Valid Values can be one of the following : %s",strings.Join(memberTypeValues, ",")),
									},
								},
							},
						},
						"external_id": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "Unique and global identifier for the given group",
						},
						"group_extension": schema.SingleNestedAttribute{
							// MarkdownDescription: ,
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Provide a unique name for the group.",
									Computed: true,
								},
								"description": schema.StringAttribute{
									Computed: true,
									MarkdownDescription: "Briefly describe the nature of the group.",
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
