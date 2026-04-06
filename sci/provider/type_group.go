package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type groupExtensionData struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type memberData struct {
	Value types.String `tfsdk:"value"`
	Type  types.String `tfsdk:"type"`
}

type groupData struct {
	Id             types.String `tfsdk:"id"`
	Schemas        types.Set    `tfsdk:"schemas" json:"schemas"`
	DisplayName    types.String `tfsdk:"display_name" json:"displayName"`
	GroupMembers   types.Set    `tfsdk:"group_members" json:"members"`
	GroupExtension types.Object `tfsdk:"group_extension" json:"urn:sap:cloud:scim:schemas:extension:custom:2.0:Group"`
}

type groupsData struct {
	Values types.List `tfsdk:"values"`
}

func groupValueFrom(ctx context.Context, g groups.Group) (groupData, diag.Diagnostics) {

	var diagnostics, diags diag.Diagnostics

	group := groupData{
		Id:          types.StringValue(g.Id),
		DisplayName: types.StringValue(g.DisplayName),
	}

	//Schemas
	schemas, diags := types.SetValueFrom(ctx, types.StringType, g.Schemas)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return group, diagnostics
	}

	group.Schemas = schemas

	// Group Extension
	groupExtension := groupExtensionData{
		Name: types.StringValue(g.GroupExtension.Name),
	}

	// mapping is done manually to handle null values
	if len(g.GroupExtension.Description) > 0 {
		groupExtension.Description = types.StringValue(g.GroupExtension.Description)
	}

	groupExtensionData, diags := types.ObjectValueFrom(ctx, groupExtensionObjType, groupExtension)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return group, diagnostics
	}

	group.GroupExtension = groupExtensionData

	// Group Members
	if len(g.GroupMembers) > 0 {

		groupMembers, diags := types.SetValueFrom(ctx, membersObjType, g.GroupMembers)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return group, diagnostics
		}

		group.GroupMembers = groupMembers

	} else {
		group.GroupMembers = types.SetNull(membersObjType)
	}
	diagnostics.Append(diags...)

	return group, diagnostics
}

func groupsValueFrom(ctx context.Context, g groups.GroupsResponse) ([]groupData, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	groups := []groupData{}

	for _, groupRes := range g.Resources {

		group, diags := groupValueFrom(ctx, groupRes)
		groups = append(groups, group)
		diagnostics.Append(diags...)

	}

	return groups, diagnostics
}

func (r *groupResource) GetGroupRequest(ctx context.Context, plan groupData) (*groups.Group, diag.Diagnostics) {

	var diagnostics diag.Diagnostics

	var schemas []string
	diags := plan.Schemas.ElementsAs(ctx, &schemas, true)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	args := &groups.Group{
		Schemas:     schemas,
		DisplayName: plan.DisplayName.ValueString(),
	}

	if !plan.GroupMembers.IsNull() {

		var members []memberData
		diags = plan.GroupMembers.ElementsAs(ctx, &members, true)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		// the mapping is done manually in order to carry out the member validation
		for _, member := range members {

			// validate the member as a valid user or group as the API does not handle this
			err := validateMembers(ctx, r.cli, member.Value.ValueString())
			if err != nil {
				diagnostics.AddError(
					fmt.Sprintf("%s", err),
					"please provide a valid member UUID",
				)
				return nil, diagnostics
			}

			groupMember := groups.GroupMember{
				Value: member.Value.ValueString(),
			}

			if !member.Type.IsNull() {
				groupMember.Type = member.Type.ValueString()
			}

			args.GroupMembers = append(args.GroupMembers, groupMember)
		}
	}

	if !plan.GroupExtension.IsNull() && !plan.GroupExtension.IsUnknown() {

		var groupExtension groups.GroupExtension
		diags = plan.GroupExtension.As(ctx, &groupExtension, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		args.GroupExtension = &groupExtension
	}

	return args, diagnostics
}

func getGroupUpdateRequest(ctx context.Context, plan groupData, state groupData) ([]generic.PatchRequest, diag.Diagnostics) {

	var diags diag.Diagnostics
	reqs := []generic.PatchRequest{}

	argsType := reflect.TypeFor[groupData]()

	if !plan.DisplayName.Equal(state.DisplayName) {
		var displayName string
		if !plan.DisplayName.IsNull() && !plan.DisplayName.IsUnknown() {
			displayName = plan.DisplayName.ValueString()
		}

		patchReq, diags := utils.GetScimPatchRequest("DisplayName", "", displayName, argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.Schemas.Equal(state.Schemas) {
		var schemas []string

		if !plan.Schemas.IsNull() {
			diags = plan.Schemas.ElementsAs(ctx, &schemas, false)
			if diags.HasError() {
				return reqs, diags
			}
		}

		patchReq, diags := utils.GetScimPatchRequest("Schemas", "", schemas, argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.GroupMembers.Equal(state.GroupMembers) {
		members := []memberData{}

		if !plan.GroupMembers.IsNull() {
			diags = plan.GroupMembers.ElementsAs(ctx, &members, true)
			if diags.HasError() {
				return reqs, diags
			}
		}
		scimMembers := []map[string]interface{}{}

		for _, m := range members {
			if m.Value.IsNull() || m.Value.IsUnknown() {
				continue
			}

			scimMembers = append(scimMembers, map[string]interface{}{
				"value": m.Value.ValueString(),
				"type":  m.Type.ValueString(),
			})
		}

		patchReq, diags := utils.GetScimPatchRequest("GroupMembers", "", scimMembers, argsType)
		if diags.HasError() {
			return reqs, diags
		}
		reqs = append(reqs, patchReq)
	}

	if !plan.GroupExtension.Equal(state.GroupExtension) {

		groupExtensionPath, diags := utils.GetAttributeTag("GroupExtension", argsType)
		if diags.HasError() {
			return reqs, diags
		}

		groupExtensionArgsType := reflect.TypeFor[groupExtensionData]()

		var planExt, stateExt groupExtensionData

		diags = plan.GroupExtension.As(ctx, &planExt, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		if diags.HasError() {
			return reqs, diags
		}
		diags = state.GroupExtension.As(ctx, &stateExt, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		if diags.HasError() {
			return reqs, diags
		}
		if !planExt.Name.Equal(stateExt.Name) {
			var name string
			if !planExt.Name.IsNull() && !planExt.Name.IsUnknown() {
				name = planExt.Name.ValueString()
			}

			patchReq, diags := utils.GetScimPatchRequest("Name", groupExtensionPath, name, groupExtensionArgsType)
			if diags.HasError() {
				return reqs, diags
			}
			reqs = append(reqs, patchReq)
		}
		if !planExt.Description.Equal(stateExt.Description) {
			var description string
			if !planExt.Description.IsNull() && !planExt.Description.IsUnknown() {
				description = planExt.Description.ValueString()
			}

			patchReq, diags := utils.GetScimPatchRequest("Description", groupExtensionPath, description, groupExtensionArgsType)
			if diags.HasError() {
				return reqs, diags
			}
			reqs = append(reqs, patchReq)
		}
	}

	return reqs, diags
}
