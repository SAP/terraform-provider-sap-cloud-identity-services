package provider

import (
	"context"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-ias/internal/cli/apiObjects/applications"
)

type assertionAttributesData struct {
	AssertionAttributeName 		types.String 		`tfsdk:"assertion_attribute_name"`
	UserAttributeName 			types.String 		`tfsdk:"user_attribute_name"`
	Inherited 					types.Bool 			`tfsdk:"inherited"`
}

type applicationData struct {
	//INPUT
	Id 							types.String 	 `tfsdk:"id"`
	//OUTPUT
	Name        				types.String	 `tfsdk:"name"`
	Description 				types.String	 `tfsdk:"description"`
	ParentApplicationId			types.String 	 `tfsdk:"parent_application_id"`
	MultiTenantApp				types.Bool	 	 `tfsdk:"multi_tenant_app"`
	GlobalAccount 				types.String	 `tfsdk:"global_account"`
	SsoType 					types.String 	 `tfsdk:"sso_type"`
	AssertionAttributes    		types.Set		 `tfsdk:"assertion_attributes"`
}

func applicationValueFrom(ctx context.Context, a applications.Application) (applicationData, diag.Diagnostics) {
	
	var diagnostics, diags diag.Diagnostics

	application := applicationData{
		Id:         			types.StringValue(a.Id),
		Name:       			types.StringValue(a.Name),
		Description: 			types.StringValue(a.Description),
		ParentApplicationId: 	types.StringValue(a.ParentApplicationId),
		MultiTenantApp: 		types.BoolValue(a.MultiTenantApp),
		GlobalAccount: 			types.StringValue(a.GlobalAccount),
		SsoType: 				types.StringValue(a.AuthenticationSchema.SsoType),
	}

	attributes := []assertionAttributesData{}
	for _, attributeRes := range a.AuthenticationSchema.AssertionAttributes{
		
		attribute := assertionAttributesData{
			AssertionAttributeName: types.StringValue(attributeRes.AssertionAttributeName),
			UserAttributeName: types.StringValue(attributeRes.UserAttributeName),
			Inherited: types.BoolValue(attributeRes.Inherited),
		}

		attributes = append(attributes, attribute)
	}
	application.AssertionAttributes, diags = types.SetValueFrom(ctx, attributesObjType, attributes)
	diagnostics.Append(diags...)

	return application, diagnostics
}

func applicationsValueFrom(ctx context.Context, a applications.ApplicationsResponse) []applicationData {
	apps := []applicationData{}

	for _, appRes := range a.Applications {

		app, _:= applicationValueFrom(ctx, appRes)
		apps = append(apps, app)

	}

	return apps
}
