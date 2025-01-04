package provider

import (
	"context"
	"regexp"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-ias/internal/cli/apiObjects/applications"
)

type advancedAssertionAttributesData struct {
	Source 				types.String 			`tfsdk:"source"`
	AttributeName		types.String 			`tfsdk:"attribute_name"`
	AttributeValue 		types.String 			`tfsdk:"attribute_value"`
	Inherited 			types.Bool 				`tfsdk:"inherited"`
}

type assertionAttributesData struct {
	AttributeName 				types.String 		`tfsdk:"attribute_name"`
	AttributeValue 			    types.String 		`tfsdk:"attribute_value"`
	Inherited 					types.Bool 			`tfsdk:"inherited"`
}

type subjectNameIdentifierData struct {
	Source 		types.String 	`tfsdk:"source"`
	Value 		types.String 	`tfsdk:"value"`
}

type authenticationRulesData struct {
	UserType				types.String		`tfsdk:"user_type"`
	UserGroup 				types.String 		`tfsdk:"user_group"`
	UserEmailDomain			types.String 		`tfsdk:"user_email_domain"`
	IdentityProviderId		types.String 		`tfsdk:"identity_provider_id"`
	IpNetworkRange 			types.String 		`tfsdk:"ip_network_range"`
}

type applicationData struct {
	//INPUT
	Id 									types.String 							 `tfsdk:"id"`
	//OUTPUT
	Name        						types.String							 `tfsdk:"name"`
	Description 						types.String							 `tfsdk:"description"`
	ParentApplicationId					types.String 							 `tfsdk:"parent_application_id"`
	MultiTenantApp						types.Bool	 							 `tfsdk:"multi_tenant_app"`
	GlobalAccount 						types.String							 `tfsdk:"global_account"`
	SsoType 							types.String 	 						 `tfsdk:"sso_type"`
	SubjectNameIdentifier 				*subjectNameIdentifierData  			 `tfsdk:"subject_name_identifier"`
	AssertionAttributes    				types.List 			 					 `tfsdk:"assertion_attributes"`
	AdvancedAssertionAttributes			types.List								 `tfsdk:"advanced_assertion_attributes"`
	DefaultAuthenticatingIdpId  		types.String 				 			 `tfsdk:"default_authenticating_idp"`
	AuthenticationRules 				types.List 					   			 `tfsdk:"authentication_rules"`
}

func applicationValueFrom(ctx context.Context, a applications.Application) (applicationData, diag.Diagnostics) {
	
	var diagnostics, diags diag.Diagnostics

	// check for expressions?
	re := regexp.MustCompile(`\$\{corporateIdP\.([^\}]+)\}`)

	application := applicationData{
		Id:         			types.StringValue(a.Id),
		Name:       			types.StringValue(a.Name),
		ParentApplicationId: 	types.StringValue(a.ParentApplicationId),
		MultiTenantApp: 		types.BoolValue(a.MultiTenantApp),
		GlobalAccount: 			types.StringValue(a.GlobalAccount),
		SsoType: 				types.StringValue(a.AuthenticationSchema.SsoType),
		DefaultAuthenticatingIdpId: types.StringValue(a.AuthenticationSchema.DefaultAuthenticatingIdpId),
	}

	if len(a.Description) > 0 {
		application.Description = types.StringValue(a.Description)
	}

	application.SubjectNameIdentifier = &subjectNameIdentifierData{}
	if re.MatchString(a.AuthenticationSchema.SubjectNameIdentifier) {
		match := re.FindStringSubmatch(a.AuthenticationSchema.SubjectNameIdentifier)
		application.SubjectNameIdentifier.Value = types.StringValue(match[1])
		application.SubjectNameIdentifier.Source = types.StringValue("Corporate Identity Provider")
	} else {
		application.SubjectNameIdentifier.Value = types.StringValue(a.AuthenticationSchema.SubjectNameIdentifier)
		application.SubjectNameIdentifier.Source = types.StringValue("Identity Directory")
	}

	attributes := []assertionAttributesData{}
	for _, attributeRes := range a.AuthenticationSchema.AssertionAttributes{
		
		attribute := assertionAttributesData{
			AttributeName: types.StringValue(attributeRes.AssertionAttributeName),
			AttributeValue: types.StringValue(attributeRes.UserAttributeName),
			Inherited: types.BoolValue(attributeRes.Inherited),
		}

		attributes = append(attributes, attribute)
	}
	application.AssertionAttributes, diags = types.ListValueFrom(ctx, assertionAttributesObjType, attributes)
	diagnostics.Append(diags...)

	advancedAttributes := []advancedAssertionAttributesData{}
	for _, attributeRes := range a.AuthenticationSchema.AdvancedAssertionAttributes{
		
		attribute := advancedAssertionAttributesData{
			AttributeName: types.StringValue(attributeRes.AttributeName),
			Inherited: types.BoolValue(attributeRes.Inherited),
		}

		//generalise these strings
		if re.MatchString(attributeRes.AttributeValue) {
			attribute.Source = types.StringValue("Corporate Identity Provider")
			match := re.FindStringSubmatch(attributeRes.AttributeValue)
			attribute.AttributeValue = types.StringValue(match[1])

		} else {
			attribute.Source = types.StringValue("Expression")
			attribute.AttributeValue = types.StringValue(attributeRes.AttributeValue)
		}

		advancedAttributes = append(advancedAttributes, attribute)
	}
	application.AdvancedAssertionAttributes, diags = types.ListValueFrom(ctx, advancedAssertionAttributesObjType, advancedAttributes)
	diagnostics.Append(diags...) 

	authRules := []authenticationRulesData{}
	for _, authRulesRes := range a.AuthenticationSchema.ConditionalAuthentication{

		//alt logic?
		rule := authenticationRulesData{}

		if len(authRulesRes.UserType)>0 		  	{ rule.UserType = types.StringValue(authRulesRes.UserType) } 
		if len(authRulesRes.UserGroup)>0		  	{ rule.UserGroup = types.StringValue(authRulesRes.UserGroup) }
		if len(authRulesRes.UserEmailDomain)>0 	  	{ rule.UserEmailDomain = types.StringValue(authRulesRes.UserEmailDomain) }
		if len(authRulesRes.IdentityProviderId)>0 	{ rule.IdentityProviderId = types.StringValue(authRulesRes.IdentityProviderId) }	
		if len(authRulesRes.IpNetworkRange)>0 		{ rule.IpNetworkRange = types.StringValue(authRulesRes.IpNetworkRange) }

		authRules = append(authRules, rule)
	}

	if len(authRules) > 0 {
		application.AuthenticationRules, diags = types.ListValueFrom(ctx, authenticationRulesObjType, authRules)
	} else {
		application.AuthenticationRules = types.ListNull(authenticationRulesObjType)
	}

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
