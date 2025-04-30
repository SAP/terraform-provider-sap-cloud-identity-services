package utils

import (
	"context"

	// "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/types"
	// "github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateAction() planmodifier.Object {
	return updateAction{}
}

// useStateForUnknownModifier implements the plan modifier.
type updateAction struct{}

// Description returns a human-readable description of the plan modifier.
func (m updateAction) Description(_ context.Context) string {
	return "If plan is unknown and state doesn't have the default 5 values, update the state"
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m updateAction) MarkdownDescription(_ context.Context) string {
	return "If plan is unknown and state doesn't have the default 5 values, update the state"
}

// PlanModifyList implements the plan modification logic.
func (m updateAction) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {

	if req.StateValue.IsNull() {
		return
	}

	// if !req.PlanValue.IsUnknown() { 

	// 	attrVals := req.PlanValue.Attributes()

	// 	defaultAction := attrVals["default_action"].(types.List)
	// 	rules := attrVals["rules"].(types.List)

	// 	if (defaultAction.IsNull() || defaultAction.IsUnknown())
	// }

	// attrVals := req.StateValue.Attributes()

	// defaultAction := attrVals["default_action"].(types.List)
	// rules := attrVals["rules"].(types.List)

	// if 

	// if (defaultAction.IsNull() || defaultAction.IsUnknown()) && (rules.IsNull() || len(rules.Elements()) == 0)  {
	// 	defaultAction = types.ListNull(types.StringType)



	// }

	// req.PlanValue = types.ListNull(types.StringType)
	resp.PlanValue = req.PlanValue

	
}
