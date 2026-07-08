package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func newGroupAssignmentDataSource() datasource.DataSource {
	return &groupAssignmentDataSource{}
}

type groupAssignmentDataSource struct {
	cli *cli.SciClient
}

func (d *groupAssignmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *groupAssignmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_assignment"
}

func (d *groupAssignmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a single member assignment for a group from the SAP Cloud Identity Services tenant.`,
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique ID of the group.",
				Validators: []validator.String{
					utils.ValidUUID(),
				},
			},
			"group_member": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The member assigned to the group.",
				Attributes: map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "SCIM ID of the user or group assigned",
						Validators: []validator.String{
							utils.ValidUUID(),
						},
					},
					"type": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Type of the member. " + utils.ValidValuesString(memberTypeValues),
					},
				},
			},
		},
	}
}

func (d *groupAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config groupAssignmentData
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

	var currentMember memberData
	diags = config.GroupMember.As(ctx, &currentMember, memberAsOptions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, m := range res.GroupMembers {
		if m.Value == currentMember.Value.ValueString() {
			memberObj, diags := memberObjectFrom(ctx, m)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			state := groupAssignmentData{
				GroupId:     config.GroupId,
				GroupMember: memberObj,
			}

			diags = resp.State.Set(ctx, &state)
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	resp.Diagnostics.AddError(
		"Member not found",
		fmt.Sprintf("member %s was not found in group %s", currentMember.Value.ValueString(), config.GroupId.ValueString()),
	)
}
