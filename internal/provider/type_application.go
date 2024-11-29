package provider

import (
	"context"
	"regexp"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-ias/internal/cli/apiObjects/applications"
)

type subjectNameIdentifierData struct {
	Source 		types.String 	`tfsdk:"source"`
	Value 		types.String 	`tfsdk:"value"`
}

type assertionAttributesData struct {
	Source 						types.String  		`tfsdk:"source"`
	AssertionAttributeName 		types.String 		`tfsdk:"assertion_attribute_name"`
	UserAttributeName 			types.String 		`tfsdk:"user_attribute_name"`
	Inherited 					types.Bool 			`tfsdk:"inherited"`
}

type applicationData struct {
	//INPUT
	Id 							types.String 				 `tfsdk:"id"`
	//OUTPUT
	Name        				types.String				 `tfsdk:"name"`
	Description 				types.String				 `tfsdk:"description"`
	ParentApplicationId			types.String 				 `tfsdk:"parent_application_id"`
	MultiTenantApp				types.Bool	 				 `tfsdk:"multi_tenant_app"`
	GlobalAccount 				types.String				 `tfsdk:"global_account"`
	SsoType 					types.String 	 			 `tfsdk:"sso_type"`
	SubjectNameIdentifier 		*subjectNameIdentifierData   `tfsdk:"subject_name_identifier"`
	AssertionAttributes    		*[]assertionAttributesData   `tfsdk:"assertion_attributes"`
}

func applicationValueFrom(_ context.Context, a applications.Application) (applicationData, diag.Diagnostics) {
	
	var diagnostics, diags diag.Diagnostics

	application := applicationData{
		Id:         			types.StringValue(a.Id),
		Name:       			types.StringValue(a.Name),
		Description: 			types.StringValue(a.Description),
		ParentApplicationId: 	types.StringValue(a.ParentApplicationId),
		MultiTenantApp: 		types.BoolValue(a.MultiTenantApp),
		GlobalAccount: 			types.StringValue(a.GlobalAccount),
		SsoType: 				types.StringValue(a.AuthenticationSchema.SsoType),
		SubjectNameIdentifier:  &subjectNameIdentifierData{
			Value: types.StringValue(a.AuthenticationSchema.SubjectNameIdentifier),
		},
	}

	attributes := []assertionAttributesData{}
	for _, attributeRes := range a.AuthenticationSchema.AssertionAttributes{
		
		attribute := assertionAttributesData{
			Source: types.StringValue("Identity Directory"),
			AssertionAttributeName: types.StringValue(attributeRes.AssertionAttributeName),
			UserAttributeName: types.StringValue(attributeRes.UserAttributeName),
			Inherited: types.BoolValue(attributeRes.Inherited),
		}

		attributes = append(attributes, attribute)
	}

	re := regexp.MustCompile(`\$\{corporateIdP\.([^\}]+)\}`)
	for _, attributeRes := range a.AuthenticationSchema.AdvancedAssertionAttributes{
		
		attribute := assertionAttributesData{
			AssertionAttributeName: types.StringValue(attributeRes.AttributeName),
			Inherited: types.BoolValue(attributeRes.Inherited),
		}


		if re.MatchString(attributeRes.AttributeValue) {
			attribute.Source = types.StringValue("Corporate Identity Provider")
			match := re.FindStringSubmatch(attributeRes.AttributeValue)
			attribute.UserAttributeName = types.StringValue(match[1])

		} else {
			attribute.Source = types.StringValue("Expression")
			attribute.UserAttributeName = types.StringValue(attributeRes.AttributeValue)
		}

		attributes = append(attributes, attribute)
	}
	application.AssertionAttributes = &attributes

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
