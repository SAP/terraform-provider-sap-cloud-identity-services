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

func newGroupsDataSource() datasource.DataSource {
	return &groupsDataSource{}
}

var groupExtensionObjType = map[string]attr.Type{
	"name":        types.StringType,
	"description": types.StringType,
}

var membersObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"value": types.StringType,
		"type":  types.StringType,
	},
}

var groupObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id": types.StringType,
		"schemas": types.SetType{
			ElemType: types.StringType,
		},
		"display_name": types.StringType,
		"group_members": types.SetType{
			ElemType: membersObjType,
		},
		"group_extension": types.ObjectType{
			AttrTypes: groupExtensionObjType,
		},
	},
}

type groupsDataSource struct {
	cli *cli.SciClient
}

func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *groupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets an list of groups from the SAP Cloud Identity services.`,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique ID of the group.",
							Validators: []validator.String{
								utils.ValidUUID(),
							},
						},
						"schemas": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							MarkdownDescription: "List of SCIM schemas to configure groups. The attribute is configured with default values :\n" +
								utils.PrintDefaultSchemas(defaultGroupSchemas),
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Display Name of the group.",
						},
						"group_members": schema.SetNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Specify the members to be part of the group.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "SCIM ID of the user or the group",
									},
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Type of the member added to the group.",
									},
								},
							},
						},
						"group_extension": schema.SingleNestedAttribute{
							MarkdownDescription: "Configure attributes particular to the schema `" + defaultGroupSchemas[1].String() + "`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Provide a unique name for the group.",
									Computed:            true,
								},
								"description": schema.StringAttribute{
									Computed:            true,
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
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := d.cli.Group.Get(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving groups", fmt.Sprintf("%s", err))
		return
	}

	resGroups, diags := groupsValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Values, diags = types.ListValueFrom(ctx, groupObjType, resGroups)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
