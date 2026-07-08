package provider

import (
	"context"
	"reflect"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func groupBaseValueFrom(ctx context.Context, g groups.Group) (groupBaseData, diag.Diagnostics) {

	var diagnostics, diags diag.Diagnostics

	group := groupBaseData{
		Id:          types.StringValue(g.Id),
		DisplayName: types.StringValue(g.DisplayName),
	}

	schemas, diags := types.SetValueFrom(ctx, types.StringType, g.Schemas)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return group, diagnostics
	}

	group.Schemas = schemas

	groupExt := groupExtensionData{
		Name: types.StringValue(g.GroupExtension.Name),
	}

	if len(g.GroupExtension.Description) > 0 {
		groupExt.Description = types.StringValue(g.GroupExtension.Description)
	}

	groupExtObj, diags := types.ObjectValueFrom(ctx, groupExtensionObjType, groupExt)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return group, diagnostics
	}

	group.GroupExtension = groupExtObj

	return group, diagnostics
}

func groupBasesValueFrom(ctx context.Context, g groups.GroupsResponse) ([]groupBaseData, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	groups := []groupBaseData{}

	for _, groupRes := range g.Resources {

		group, diags := groupBaseValueFrom(ctx, groupRes)
		groups = append(groups, group)
		diagnostics.Append(diags...)

	}

	return groups, diagnostics
}

func getGroupBaseRequest(ctx context.Context, plan groupBaseData) (*groups.Group, diag.Diagnostics) {

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

func getGroupBaseUpdateRequest(ctx context.Context, plan groupBaseData, state groupBaseData) ([]generic.PatchRequest, diag.Diagnostics) {

	var diags diag.Diagnostics
	reqs := []generic.PatchRequest{}

	argsType := reflect.TypeFor[groupBaseData]()

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
