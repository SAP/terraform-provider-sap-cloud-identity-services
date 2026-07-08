package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newGroupAssignmentsDataSource() datasource.DataSource {
	return &groupAssignmentsDataSource{}
}

type groupAssignmentsData struct {
	GroupId types.String `tfsdk:"group_id"`
	Values  types.List   `tfsdk:"values"`
}

type groupAssignmentsDataSource struct {
	cli *cli.SciClient
}

func (d *groupAssignmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *groupAssignmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_assignments"
}

func (d *groupAssignmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all member assignments for a group from the SAP Cloud Identity Services tenant.`,
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique ID of the group.",
				Validators: []validator.String{
					utils.ValidUUID(),
				},
			},
			"values": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of members assigned to the group.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "SCIM ID of the user or group assigned.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Type of the member.",
						},
					},
				},
			},
		},
	}
}

func (d *groupAssignmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config groupAssignmentsData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := d.cli.Group.GetByGroupId(ctx, config.GroupId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving group", fmt.Sprintf("%s", err))
		return
	}

	config.Values, diags = types.ListValueFrom(ctx, membersObjType, res.GroupMembers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
