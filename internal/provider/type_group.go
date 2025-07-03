package provider

import (
	"context"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Schemas        types.Set    `tfsdk:"schemas"`
	DisplayName    types.String `tfsdk:"display_name"`
	GroupMembers   types.Set    `tfsdk:"group_members"`
	GroupExtension types.Object `tfsdk:"group_extension"`
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
