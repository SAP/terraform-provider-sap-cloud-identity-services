package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newGroupBasesDataSource() datasource.DataSource {
	return &groupBasesDataSource{}
}

var groupBaseObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id": types.StringType,
		"schemas": types.SetType{
			ElemType: types.StringType,
		},
		"display_name": types.StringType,
		"group_extension": types.ObjectType{
			AttrTypes: groupExtensionObjType,
		},
	},
}

type groupBasesData struct {
	Values types.List `tfsdk:"values"`
}

type groupBasesDataSource struct {
	cli *cli.SciClient
}

func (d *groupBasesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *groupBasesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_bases"
}

func (d *groupBasesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a list of groups without its member assignments from the SAP Cloud Identity Services tenant.`,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				Computed: true,
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
			},
		},
	}
}

func (d *groupBasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config groupBasesData
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

	groupBases := make([]groupBaseData, 0, len(res.Resources))
	var diagnostics diag.Diagnostics
	for _, g := range res.Resources {
		gb, diags := groupBaseValueFrom(ctx, g)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			resp.Diagnostics.Append(diagnostics...)
			return
		}
		groupBases = append(groupBases, gb)
	}

	config.Values, diags = types.ListValueFrom(ctx, groupBaseObjType, groupBases)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
