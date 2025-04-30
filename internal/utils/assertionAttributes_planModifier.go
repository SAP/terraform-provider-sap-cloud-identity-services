package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateUnknown() planmodifier.List {
	return updateUnknown{}
}

// useStateForUnknownModifier implements the plan modifier.
type updateUnknown struct{}

// Description returns a human-readable description of the plan modifier.
func (m updateUnknown) Description(_ context.Context) string {
	return "If plan is unknown and state doesn't have the default 5 values, update the state"
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m updateUnknown) MarkdownDescription(_ context.Context) string {
	return "If plan is unknown and state doesn't have the default 5 values, update the state"
}

// PlanModifyList implements the plan modification logic.
func (m updateUnknown) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	if req.PlanValue.IsUnknown() {
		resp.PlanValue = req.StateValue
	}

	if (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) && len(req.StateValue.Elements()) < 5 {
		req.PlanValue = types.ListUnknown(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"attribute_name":  types.StringType,
				"attribute_value": types.StringType,
				"inherited":       types.BoolType,
			},
		})
		resp.PlanValue = req.PlanValue
	}

	
}
