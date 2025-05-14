package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
)

type authenticationSchemaData struct {
	SsoType                       types.String               `tfsdk:"sso_type"`
	SubjectNameIdentifier         *subjectNameIdentifierData `tfsdk:"subject_name_identifier"`
	SubjectNameIdentifierFunction types.String               `tfsdk:"subject_name_identifier_function"`
	AssertionAttributes           types.List                 `tfsdk:"assertion_attributes"`
	AdvancedAssertionAttributes   types.List                 `tfsdk:"advanced_assertion_attributes"`
	DefaultAuthenticatingIdpId    types.String               `tfsdk:"default_authenticating_idp"`
	AuthenticationRules           types.List                 `tfsdk:"authentication_rules"`
}

type advancedAssertionAttributesData struct {
	Source         types.String `tfsdk:"source"`
	AttributeName  types.String `tfsdk:"attribute_name"`
	AttributeValue types.String `tfsdk:"attribute_value"`
	Inherited      types.Bool   `tfsdk:"inherited"`
}

type assertionAttributesData struct {
	AttributeName  types.String `tfsdk:"attribute_name"`
	AttributeValue types.String `tfsdk:"attribute_value"`
	Inherited      types.Bool   `tfsdk:"inherited"`
}

type subjectNameIdentifierData struct {
	Source types.String `tfsdk:"source"`
	Value  types.String `tfsdk:"value"`
}

type authenticationRulesData struct {
	UserType           types.String `tfsdk:"user_type"`
	UserGroup          types.String `tfsdk:"user_group"`
	UserEmailDomain    types.String `tfsdk:"user_email_domain"`
	IdentityProviderId types.String `tfsdk:"identity_provider_id"`
	IpNetworkRange     types.String `tfsdk:"ip_network_range"`
}

type applicationData struct {
	//INPUT
	Id types.String `tfsdk:"id"`
	//OUTPUT
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	ParentApplicationId  types.String `tfsdk:"parent_application_id"`
	MultiTenantApp       types.Bool   `tfsdk:"multi_tenant_app"`
	GlobalAccount        types.String `tfsdk:"global_account"`
	AuthenticationSchema types.Object `tfsdk:"authentication_schema"`
}

func applicationValueFrom(ctx context.Context, a applications.Application) (applicationData, diag.Diagnostics) {

	var diagnostics, diags diag.Diagnostics

	// check for expressions?
	re := regexp.MustCompile(`\$\{corporateIdP\.([^\}]+)\}`)

	application := applicationData{
		Id:                  types.StringValue(a.Id),
		Name:                types.StringValue(a.Name),
		ParentApplicationId: types.StringValue(a.ParentApplicationId),
		MultiTenantApp:      types.BoolValue(a.MultiTenantApp),
		GlobalAccount:       types.StringValue(a.GlobalAccount),
	}

	if len(a.Description) > 0 {
		application.Description = types.StringValue(a.Description)
	}

	authenticationSchema := authenticationSchemaData{}

	authenticationSchema.SsoType = types.StringValue(a.AuthenticationSchema.SsoType)
	authenticationSchema.DefaultAuthenticatingIdpId = types.StringValue(a.AuthenticationSchema.DefaultAuthenticatingIdpId)

	authenticationSchema.SubjectNameIdentifier = &subjectNameIdentifierData{}

	//regex for expressions needed
	if re.MatchString(a.AuthenticationSchema.SubjectNameIdentifier) {
		match := re.FindStringSubmatch(a.AuthenticationSchema.SubjectNameIdentifier)
		authenticationSchema.SubjectNameIdentifier.Value = types.StringValue(match[1])
		authenticationSchema.SubjectNameIdentifier.Source = types.StringValue(sourceValues[1])
	} else {
		authenticationSchema.SubjectNameIdentifier.Value = types.StringValue(a.AuthenticationSchema.SubjectNameIdentifier)
		authenticationSchema.SubjectNameIdentifier.Source = types.StringValue(sourceValues[0])
	}

	if len(a.AuthenticationSchema.SubjectNameIdentifierFunction) > 0 {
		authenticationSchema.SubjectNameIdentifierFunction = types.StringValue((a.AuthenticationSchema.SubjectNameIdentifierFunction))
	}

	attributes := []assertionAttributesData{}
	for _, attributeRes := range a.AuthenticationSchema.AssertionAttributes {

		attribute := assertionAttributesData{
			AttributeName:  types.StringValue(attributeRes.AssertionAttributeName),
			AttributeValue: types.StringValue(attributeRes.UserAttributeName),
			Inherited:      types.BoolValue(attributeRes.Inherited),
		}

		attributes = append(attributes, attribute)
	}
	authenticationSchema.AssertionAttributes, diags = types.ListValueFrom(ctx, assertionAttributesObjType, attributes)
	diagnostics.Append(diags...)

	advancedAttributes := []advancedAssertionAttributesData{}
	for _, attributeRes := range a.AuthenticationSchema.AdvancedAssertionAttributes {

		attribute := advancedAssertionAttributesData{
			AttributeName: types.StringValue(attributeRes.AttributeName),
			Inherited:     types.BoolValue(attributeRes.Inherited),
		}

		if re.MatchString(attributeRes.AttributeValue) {
			attribute.Source = types.StringValue(sourceValues[1])
			match := re.FindStringSubmatch(attributeRes.AttributeValue)
			attribute.AttributeValue = types.StringValue(match[1])

		} else {
			attribute.Source = types.StringValue(sourceValues[2])
			attribute.AttributeValue = types.StringValue(attributeRes.AttributeValue)
		}

		advancedAttributes = append(advancedAttributes, attribute)
	}

	if len(advancedAttributes) > 0 {
		authenticationSchema.AdvancedAssertionAttributes, diags = types.ListValueFrom(ctx, advancedAssertionAttributesObjType, advancedAttributes)
	} else {
		authenticationSchema.AdvancedAssertionAttributes = types.ListNull(advancedAssertionAttributesObjType)
	}

	diagnostics.Append(diags...)

	authRules := []authenticationRulesData{}
	for _, authRulesRes := range a.AuthenticationSchema.ConditionalAuthentication {

		rule := authenticationRulesData{}

		if len(authRulesRes.UserType) > 0 {
			rule.UserType = types.StringValue(authRulesRes.UserType)
		}
		if len(authRulesRes.UserGroup) > 0 {
			rule.UserGroup = types.StringValue(authRulesRes.UserGroup)
		}
		if len(authRulesRes.UserEmailDomain) > 0 {
			rule.UserEmailDomain = types.StringValue(authRulesRes.UserEmailDomain)
		}
		if len(authRulesRes.IdentityProviderId) > 0 {
			rule.IdentityProviderId = types.StringValue(authRulesRes.IdentityProviderId)
		}
		if len(authRulesRes.IpNetworkRange) > 0 {
			rule.IpNetworkRange = types.StringValue(authRulesRes.IpNetworkRange)
		}

		authRules = append(authRules, rule)
	}

	if len(authRules) > 0 {
		authenticationSchema.AuthenticationRules, diags = types.ListValueFrom(ctx, authenticationRulesObjType, authRules)
	} else {
		authenticationSchema.AuthenticationRules = types.ListNull(authenticationRulesObjType)
	}

	diagnostics.Append(diags...)

	application.AuthenticationSchema, diags = types.ObjectValueFrom(ctx, authenticationSchemaObjType, authenticationSchema)
	diagnostics.Append(diags...)

	return application, diagnostics
}

func applicationsValueFrom(ctx context.Context, a applications.ApplicationsResponse) []applicationData {
	apps := []applicationData{}

	for _, appRes := range a.Applications {

		app, _ := applicationValueFrom(ctx, appRes)
		apps = append(apps, app)

	}

	return apps
}
