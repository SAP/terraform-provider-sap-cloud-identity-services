package provider

import (
	"context"
	"terraform-provider-sci/internal/cli/apiObjects/groups"

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
	GroupMembers   types.List   `tfsdk:"group_members"`
	ExternalId     types.String `tfsdk:"external_id"`
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

	group.Schemas, diags = types.SetValueFrom(ctx, types.StringType, g.Schemas)
	diagnostics.Append(diags...)

	if len(g.ExternalId) > 0 {
		group.ExternalId = types.StringValue(g.ExternalId)
	}

	groupExtension := groupExtensionData{
		Name:        types.StringValue(g.GroupExtension.Name),
		Description: types.StringValue(g.GroupExtension.Description),
	}

	group.GroupExtension, diags = types.ObjectValueFrom(ctx, groupExtensionObjType, groupExtension)
	diagnostics.Append(diags...)

	groupMembers := []memberData{}
	for _, memberRes := range g.GroupMembers {

		member := memberData{
			Value: types.StringValue(memberRes.Value),
			Type:  types.StringValue(memberRes.Type),
		}

		groupMembers = append(groupMembers, member)
	}

	if len(groupMembers) > 0 {
		group.GroupMembers, diags = types.ListValueFrom(ctx, membersObjType, groupMembers)
	} else {
		group.GroupMembers = types.ListNull(membersObjType)
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
