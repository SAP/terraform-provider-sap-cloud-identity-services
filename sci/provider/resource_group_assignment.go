package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var memberAsOptions = basetypes.ObjectAsOptions{
	UnhandledNullAsEmpty:    true,
	UnhandledUnknownAsEmpty: true,
}

func newGroupAssignmentResource() resource.Resource {
	return &groupAssignmentResource{}
}

type groupAssignmentResource struct {
	cli *cli.SciClient
}

func (r *groupAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.cli = req.ProviderData.(*cli.SciClient)
}

func (r *groupAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_assignment"
}

func (r *groupAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assign and manage a single member assigned to a group in the SAP Cloud Identity Services tenant.`,
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique ID of the group.",
				Validators: []validator.String{
					utils.ValidUUID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_member": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The member assigned to the group.",
				Attributes: map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "SCIM ID of the user or group assigned.",
						Validators: []validator.String{
							utils.ValidUUID(),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Type of the member. " + utils.ValidValuesString(memberTypeValues),
						Validators: []validator.String{
							stringvalidator.OneOf(memberTypeValues...),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

func (r *groupAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan groupAssignmentData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var member memberData
	diags = plan.GroupMember.As(ctx, &member, memberAsOptions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateMembers(ctx, r.cli, member.Value.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("%s", err), "please provide a valid member UUID")
		return
	}

	groupMember := groups.GroupMember{Value: member.Value.ValueString()}
	if !member.Type.IsNull() && !member.Type.IsUnknown() {
		groupMember.Type = member.Type.ValueString()
	}

	patchOp := utils.GenerateAddPatchRequest("members", []groups.GroupMember{groupMember})

	res, _, err := r.cli.Group.Update(ctx, []generic.PatchRequest{patchOp}, plan.GroupId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error adding group member", fmt.Sprintf("%s", err))
		return
	}

	state, diags := groupAssignmentValueFrom(ctx, res, plan.GroupId.ValueString(), member.Value.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var config groupAssignmentData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Group.GetByGroupId(ctx, config.GroupId.ValueString())
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

	resp.Diagnostics.AddError("Unable to read group member", "Member not found")

}

func (r *groupAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// the resource has to be destroyed and re-created if any parameter is changed
	resp.Diagnostics.AddError("Resource sci_group_assignment cannot be updated", "Modify either the group_id or the group_member.value to re-create the assignment")
}

func (r *groupAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var config groupAssignmentData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var member memberData
	diags = config.GroupMember.As(ctx, &member, memberAsOptions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	removeOp := utils.GenerateDeletePatchRequest(
		fmt.Sprintf(`members[value eq "%s"]`, member.Value.ValueString()),
	)

	_, _, err := r.cli.Group.Update(ctx, []generic.PatchRequest{removeOp}, config.GroupId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error removing group member", fmt.Sprintf("%s", err))
		return
	}
}

func (r *groupAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: group_id,group_member.value Got: %q", req.ID),
		)
		return
	}

	memberObj, diags := memberObjectFrom(ctx, groups.GroupMember{Value: idParts[1]})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_member"), memberObj)...)
}
