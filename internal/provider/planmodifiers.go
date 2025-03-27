package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SchemaSetPlanModifier struct {}

func(s SchemaSetPlanModifier) Description(_ context.Context) string{
	return ""
}

func(s SchemaSetPlanModifier) MarkdownDescription(_ context.Context) string {
    return ""
}

func (s SchemaSetPlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {

	containsSchemas := checkSchemas(ctx, req, resp)

	if resp.Diagnostics.HasError(){
		return
	}

    if req.ConfigValue.IsNull() || !containsSchemas{
        
    }

	x, _ := types.SetValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
	})

    resp.PlanValue = x
}

func SchemaModifier() planmodifier.Set{
	return SchemaSetPlanModifier{}
}

func checkSchemas(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) (bool) {

	var schemas []string
	diags := req.ConfigValue.ElementsAs(ctx, &schemas, true)

	if diags.HasError(){
		resp.Diagnostics.Append(diags...)
		return false
	}

	_ = []string{
		"urn:ietf:params:scim:schemas:core:2.0:Group",
		"urn:sap:cloud:scim:schemas:extension:custom:2.0:Group",
	}


	return true
}