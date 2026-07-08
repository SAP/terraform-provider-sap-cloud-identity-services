package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type groupAssignmentData struct {
	GroupId     types.String `tfsdk:"group_id"`
	GroupMember types.Object `tfsdk:"group_member"`
}

func groupAssignmentValueFrom(ctx context.Context, g groups.Group, groupId string, memberValue string) (groupAssignmentData, diag.Diagnostics) {

	var diagnostics diag.Diagnostics

	result := groupAssignmentData{
		GroupId: types.StringValue(groupId),
	}

	for _, m := range g.GroupMembers {
		if m.Value == memberValue {
			memberObj, diags := memberObjectFrom(ctx, m)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return result, diagnostics
			}
			result.GroupMember = memberObj
			return result, diagnostics
		}
	}

	diagnostics.AddError(
		"Member not found after operation",
		fmt.Sprintf("member %s was not found in group %s after the API call", memberValue, groupId),
	)
	return result, diagnostics
}

func memberObjectFrom(ctx context.Context, m groups.GroupMember) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(ctx, membersObjType.AttrTypes, memberData{
		Value: types.StringValue(m.Value),
		Type:  types.StringValue(m.Type),
	})
}
