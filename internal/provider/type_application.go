package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
)

type authenticationSchemaData struct {
	SsoType                       types.String `tfsdk:"sso_type"`
	SubjectNameIdentifier         types.Object `tfsdk:"subject_name_identifier"`
	SubjectNameIdentifierFunction types.String `tfsdk:"subject_name_identifier_function"`
	AssertionAttributes           types.List   `tfsdk:"assertion_attributes"`
	AdvancedAssertionAttributes   types.List   `tfsdk:"advanced_assertion_attributes"`
	DefaultAuthenticatingIdpId    types.String `tfsdk:"default_authenticating_idp"`
	AuthenticationRules           types.List   `tfsdk:"conditional_authentication"`
}

type authenticationRulesData struct {
	UserType           types.String `tfsdk:"user_type"`
	UserGroup          types.String `tfsdk:"user_group"`
	UserEmailDomain    types.String `tfsdk:"user_email_domain"`
	IdentityProviderId types.String `tfsdk:"identity_provider_id"`
	IpNetworkRange     types.String `tfsdk:"ip_network_range"`
}

type advancedAssertionAttributesData struct {
	Source         types.String `tfsdk:"source"`
	AttributeName  types.String `tfsdk:"attribute_name"`
	AttributeValue types.String `tfsdk:"attribute_value"`
	Inherited      types.Bool   `tfsdk:"inherited"`
}

type subjectNameIdentifierData struct {
	Source types.String `tfsdk:"source"`
	Value  types.String `tfsdk:"value"`
}

type applicationData struct {
	//INPUT
	Id types.String `tfsdk:"id"`
	//OUTPUT
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	ParentApplicationId  types.String `tfsdk:"parent_application_id"`
	MultiTenantApp       types.Bool   `tfsdk:"multi_tenant_app"`
	AuthenticationSchema types.Object `tfsdk:"authentication_schema"`
}

func applicationValueFrom(ctx context.Context, a applications.Application) (applicationData, diag.Diagnostics) {

	var diagnostics, diags diag.Diagnostics

	// regex for attribute values whose sources are the corporate IDP
	re := regexp.MustCompile(`\$\{corporateIdP\.([^\}]+)\}`)

	// reading attributes : id, name, multi_tenant_app and global_account
	application := applicationData{
		Id:             types.StringValue(a.Id),
		Name:           types.StringValue(a.Name),
		MultiTenantApp: types.BoolValue(a.MultiTenantApp),
	}

	// reading attributes : description and parent_application_id
	if len(a.Description) > 0 {
		application.Description = types.StringValue(a.Description)
	}

	if len(a.ParentApplicationId) > 0 {
		application.ParentApplicationId = types.StringValue(a.ParentApplicationId)
	}

	// reading attributes of the Authentication Schema : sso_type, default_authenticating_idp
	authenticationSchema := authenticationSchemaData{
		SsoType:                    types.StringValue(a.AuthenticationSchema.SsoType),
		DefaultAuthenticatingIdpId: types.StringValue(a.AuthenticationSchema.DefaultAuthenticatingIdpId),
	}

	// reading attribute of the Authentication Schema : subject_name_identifier
	subjectNameIdentifier := subjectNameIdentifierData{}

	if re.MatchString(a.AuthenticationSchema.SubjectNameIdentifier) {
		match := re.FindStringSubmatch(a.AuthenticationSchema.SubjectNameIdentifier)
		subjectNameIdentifier.Value = types.StringValue(match[1])
		subjectNameIdentifier.Source = types.StringValue(sourceValues[1])
	} else {
		subjectNameIdentifier.Value = types.StringValue(a.AuthenticationSchema.SubjectNameIdentifier)
		subjectNameIdentifier.Source = types.StringValue(sourceValues[0])
	}

	subjectNameIdentifierData, diags := types.ObjectValueFrom(ctx, subjectNameIdentitfierObjType, subjectNameIdentifier)
	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return application, diagnostics
	}

	authenticationSchema.SubjectNameIdentifier = subjectNameIdentifierData

	// reading attributes of the Authentication Schema : subject_name_identifier_function
	if len(a.AuthenticationSchema.SubjectNameIdentifierFunction) > 0 {
		authenticationSchema.SubjectNameIdentifierFunction = types.StringValue((a.AuthenticationSchema.SubjectNameIdentifierFunction))
	}

	// reading attributes of the Authentication Schema : assertion_attributes
	attributes, diags := types.ListValueFrom(ctx, assertionAttributesObjType, a.AuthenticationSchema.AssertionAttributes)
	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return application, diagnostics
	}

	authenticationSchema.AssertionAttributes = attributes

	// reading attributes of the Authentication Schema : advanced_assertion_attributes
	if len(a.AuthenticationSchema.AdvancedAssertionAttributes) > 0 {

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

		advancedAttributesData, diags := types.ListValueFrom(ctx, advancedAssertionAttributesObjType, advancedAttributes)

		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return application, diagnostics
		}

		authenticationSchema.AdvancedAssertionAttributes = advancedAttributesData
	} else {
		authenticationSchema.AdvancedAssertionAttributes = types.ListNull(advancedAssertionAttributesObjType)
	}

	// reading attributes of the Authentication Schema : conditional_authentication
	if len(a.AuthenticationSchema.ConditionalAuthentication) > 0 {

		authRules := []authenticationRulesData{}

		for _, rule := range a.AuthenticationSchema.ConditionalAuthentication {

			authRule := authenticationRulesData{}

			if len(rule.UserType) > 0 {
				authRule.UserType = types.StringValue(rule.UserType)
			}

			if len(rule.UserGroup) > 0 {
				authRule.UserGroup = types.StringValue(rule.UserGroup)
			}

			if len(rule.UserEmailDomain) > 0 {
				authRule.UserEmailDomain = types.StringValue(rule.UserEmailDomain)
			}

			if len(rule.IdentityProviderId) > 0 {
				authRule.IdentityProviderId = types.StringValue(rule.IdentityProviderId)
			}

			if len(rule.IpNetworkRange) > 0 {
				authRule.IpNetworkRange = types.StringValue(rule.IpNetworkRange)
			}

			authRules = append(authRules, authRule)
		}

		authRulesData, diags := types.ListValueFrom(ctx, authenticationRulesObjType, authRules)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return application, diagnostics
		}

		authenticationSchema.AuthenticationRules = authRulesData
	} else {
		authenticationSchema.AuthenticationRules = types.ListNull(authenticationRulesObjType)
	}

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
