package provider

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
)

type authenticationSchemaData struct {
	SsoType                       types.String `tfsdk:"sso_type" json:"ssoType"`
	SubjectNameIdentifier         types.Object `tfsdk:"subject_name_identifier" json:"subjectNameIdentifier"`
	SubjectNameIdentifierFunction types.String `tfsdk:"subject_name_identifier_function" json:"subjectNameIdentifierFunction"`
	AssertionAttributes           types.List   `tfsdk:"assertion_attributes" json:"assertionAttributes"`
	AdvancedAssertionAttributes   types.List   `tfsdk:"advanced_assertion_attributes" json:"advancedAssertionAttributes"`
	DefaultAuthenticatingIdpId    types.String `tfsdk:"default_authenticating_idp" json:"defaultAuthenticatingIdpId"`
	AuthenticationRules           types.List   `tfsdk:"conditional_authentication" json:"conditionalAuthentication"`
	OpenIdConnectConfiguration    types.Object `tfsdk:"oidc_config" json:"openIdConnectConfiguration"`
	Saml2Configuration            types.Object `tfsdk:"saml2_config" json:"saml2Configuration"`
	SapManagedAttributes          types.Object `tfsdk:"sap_managed_attributes"`
}

type AppSaml2ConfigData struct {
	SamlMetadataUrl           types.String `tfsdk:"saml_metadata_url" json:"samlMetadataUrl"`
	AcsEndpoints              types.List   `tfsdk:"acs_endpoints" json:"acsEndpoints"`
	SloEndpoints              types.List   `tfsdk:"slo_endpoints" json:"sloEndpoints"`
	CertificatesForSigning    types.List   `tfsdk:"signing_certificates" json:"certificatesForSigning"`
	CertificateForEncryption  types.Object `tfsdk:"encryption_certificate" json:"certificateForEncryption"`
	ResponseElementsToEncrypt types.String `tfsdk:"response_elements_to_encrypt" json:"responseElementsToEncrypt"`
	DefaultNameIdFormat       types.String `tfsdk:"default_name_id_format" json:"defaultNameIdFormat"`
	SignSloMessages           types.Bool   `tfsdk:"sign_slo_messages" json:"signSLOMessages"`
	RequireSignedSloMessages  types.Bool   `tfsdk:"require_signed_slo_messages" json:"requireSignedSLOMessages"`
	RequireSignedAuthnRequest types.Bool   `tfsdk:"require_signed_auth_requests" json:"requireSignedAuthnRequest"`
	SignAssertions            types.Bool   `tfsdk:"sign_assertions" json:"signAssertions"`
	SignAuthnResponses        types.Bool   `tfsdk:"sign_auth_responses" json:"signAuthnResponses"`
	DigestAlgorithm           types.String `tfsdk:"digest_algorithm" json:"digestAlgorithm"`
}

type AcsSsoEndpointData struct {
	BindingName types.String `tfsdk:"binding_name" json:"bindingName"`
	Location    types.String `tfsdk:"location" json:"location"`
	Index       types.Int32  `tfsdk:"index" json:"index"`
	IsDefault   types.Bool   `tfsdk:"default" json:"isDefault"`
}

type AppSloEndpointData struct {
	BindingName      types.String `tfsdk:"binding_name" json:"bindingName"`
	Location         types.String `tfsdk:"location" json:"location"`
	ResponseLocation types.String `tfsdk:"response_location" json:"responseLocation"`
}

type authenticationRulesData struct {
	UserType           types.String `tfsdk:"user_type" json:"userType"`
	UserGroup          types.String `tfsdk:"user_group" json:"userGroup"`
	UserEmailDomain    types.String `tfsdk:"user_email_domain" json:"userEmailDomain"`
	IdentityProviderId types.String `tfsdk:"identity_provider_id" json:"identityProviderId"`
	IpNetworkRange     types.String `tfsdk:"ip_network_range" json:"ipNetworkRange"`
}

type advancedAssertionAttributesData struct {
	Source         types.String `tfsdk:"source"`
	AttributeName  types.String `tfsdk:"attribute_name" json:"attributeName"`
	AttributeValue types.String `tfsdk:"attribute_value" json:"attributeValue"`
	Inherited      types.Bool   `tfsdk:"inherited"`
}

type subjectNameIdentifierData struct {
	Source types.String `tfsdk:"source"`
	Value  types.String `tfsdk:"value" json:"subjectNameIdentifier"`
}

type openIdConnectConfigurationData struct {
	RedirectUris           types.Set    `tfsdk:"redirect_uris" json:"redirectUris"`
	PostLogoutRedirectUris types.Set    `tfsdk:"post_logout_redirect_uris" json:"postLogoutRedirectUris"`
	FrontChannelLogoutUris types.Set    `tfsdk:"front_channel_logout_uris" json:"frontChannelLogoutUris"`
	BackChannelLogoutUris  types.Set    `tfsdk:"back_channel_logout_uris" json:"backChannelLogoutUris"`
	TokenPolicy            types.Object `tfsdk:"token_policy" json:"tokenPolicy"`
	RestrictedGrantTypes   types.Set    `tfsdk:"restricted_grant_types" json:"restrictedGrantTypes"`
	ProxyConfig            types.Object `tfsdk:"proxy_config" json:"proxyConfig"`
}

type tokenPolicyData struct {
	JwtValidity                  types.Int32  `tfsdk:"jwt_validity" json:"jwtValidity"`
	RefreshValidity              types.Int32  `tfsdk:"refresh_validity" json:"refreshValidity"`
	RefreshParallel              types.Int32  `tfsdk:"refresh_parallel" json:"refreshParallel"`
	MaxExchangePeriod            types.String `tfsdk:"max_exchange_period" json:"maxExchangePeriod"`
	RefreshTokenRotationScenario types.String `tfsdk:"refresh_token_rotation_scenario" json:"refreshTokenRotationScenario"`
	AccessTokenFormat            types.String `tfsdk:"access_token_format" json:"accessTokenFormat"`
}

type proxyConfigData struct {
	Acrs types.Set `tfsdk:"acrs"`
}

type sapManagedAttributesData struct {
	ServiceInstanceId types.String `tfsdk:"service_instance_id"`
	SourceAppId       types.String `tfsdk:"source_app_id"`
	SourceTenantId    types.String `tfsdk:"source_tenant_id"`
	AppTenantId       types.String `tfsdk:"app_tenant_id"`
	Type              types.String `tfsdk:"type"`
	PlanName          types.String `tfsdk:"plan_name"`
	BtpTenantType     types.String `tfsdk:"btp_tenant_type"`
}

type applicationData struct {
	//INPUT
	Id types.String `tfsdk:"id" json:"id"`
	//OUTPUT
	Name                 types.String `tfsdk:"name" json:"name"`
	Description          types.String `tfsdk:"description" json:"description"`
	ParentApplicationId  types.String `tfsdk:"parent_application_id" json:"parentApplicationId"`
	MultiTenantApp       types.Bool   `tfsdk:"multi_tenant_app" json:"multiTenantApp"`
	AuthenticationSchema types.Object `tfsdk:"authentication_schema" json:"urn:sap:identity:application:schemas:extension:sci:1.0:Authentication"`
}

func applicationValueFrom(ctx context.Context, a applications.Application) (applicationData, diag.Diagnostics) {

	var diagnostics, diags diag.Diagnostics

	// regex for attribute values whose sources are the corporate IDP
	re := regexp.MustCompile(`\$\{corporateIdP\.([^\}]+)\}`)

	// Id, Name, Multi Tenant App and Global Account
	application := applicationData{
		Id:             types.StringValue(a.Id),
		Name:           types.StringValue(a.Name),
		MultiTenantApp: types.BoolValue(a.MultiTenantApp),
	}

	// Description and Parent Application Id
	if len(a.Description) > 0 {
		application.Description = types.StringValue(a.Description)
	}

	if len(a.ParentApplicationId) > 0 {
		application.ParentApplicationId = types.StringValue(a.ParentApplicationId)
	}

	// Authentication Schema Sso Type & Default Authenticating Idp
	authenticationSchema := authenticationSchemaData{
		SsoType:                    types.StringValue(a.AuthenticationSchema.SsoType),
		DefaultAuthenticatingIdpId: types.StringValue(a.AuthenticationSchema.DefaultAuthenticatingIdpId),
	}

	// Authentication Schema Subject Name Identifier
	// mapping is done manually to handle the inconsistency between the structure of the API response body and the schema
	// the schema defines the parameter subject_name_identitifier as an object whereas the response body returns it as a string
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

	// Authentication Schema Subject Name Identifier Function
	if len(a.AuthenticationSchema.SubjectNameIdentifierFunction) > 0 {
		authenticationSchema.SubjectNameIdentifierFunction = types.StringValue((a.AuthenticationSchema.SubjectNameIdentifierFunction))
	}

	// Authentication Schema Assertion Attributes
	attributes, diags := types.ListValueFrom(ctx, assertionAttributesObjType, a.AuthenticationSchema.AssertionAttributes)
	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return application, diagnostics
	}

	authenticationSchema.AssertionAttributes = attributes

	// Authentication Schema Advanced Assertion Attributes
	// mapping is done manually to handle the inconsistency between the structure of the API response body and the schema
	// the schema defines the parameter advanced_assertion_attributes with an additional parameter source which is not returned by the response
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

	// Authentication Schema Conditional Authentication
	// the mapping is done manually in order to handle the null values
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

	// Authentication Schema OIDC
	// the mapping is done manually in order to handle the null values

	oidc := openIdConnectConfigurationData{}

	oidc.RedirectUris, diags = types.SetValueFrom(ctx, types.StringType, a.AuthenticationSchema.OidcConfig.RedirectUris)
	diagnostics.Append(diags...)

	oidc.PostLogoutRedirectUris, diags = types.SetValueFrom(ctx, types.StringType, a.AuthenticationSchema.OidcConfig.PostLogoutRedirectUris)
	diagnostics.Append(diags...)

	oidc.FrontChannelLogoutUris, diags = types.SetValueFrom(ctx, types.StringType, a.AuthenticationSchema.OidcConfig.FrontChannelLogoutUris)
	diagnostics.Append(diags...)

	oidc.BackChannelLogoutUris, diags = types.SetValueFrom(ctx, types.StringType, a.AuthenticationSchema.OidcConfig.BackChannelLogoutUris)
	diagnostics.Append(diags...)

	if a.AuthenticationSchema.OidcConfig.TokenPolicy != nil {
		tokenpolicy := tokenPolicyData{
			JwtValidity:                  types.Int32Value(a.AuthenticationSchema.OidcConfig.TokenPolicy.JwtValidity),
			RefreshValidity:              types.Int32Value(a.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshValidity),
			RefreshParallel:              types.Int32Value(a.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshParallel),
			MaxExchangePeriod:            types.StringValue(a.AuthenticationSchema.OidcConfig.TokenPolicy.MaxExchangePeriod),
			RefreshTokenRotationScenario: types.StringValue(a.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshTokenRotationScenario),
			AccessTokenFormat:            types.StringValue(a.AuthenticationSchema.OidcConfig.TokenPolicy.AccessTokenFormat),
		}
		oidc.TokenPolicy, diags = types.ObjectValueFrom(ctx, tokenPolicyObjType, tokenpolicy)
		diagnostics.Append(diags...)
	} else {
		oidc.TokenPolicy = types.ObjectNull(tokenPolicyObjType)
	}
	var restrictedGrants []string
	for _, g := range a.AuthenticationSchema.OidcConfig.RestrictedGrantTypes {
		restrictedGrants = append(restrictedGrants, string(g))
	}
	oidc.RestrictedGrantTypes, diags = types.SetValueFrom(ctx, types.StringType, restrictedGrants)
	diagnostics.Append(diags...)

	// Proxy Config
	if a.AuthenticationSchema.OidcConfig.ProxyConfig != nil {
		proxyConfig := proxyConfigData{}
		proxyConfig.Acrs, diags = types.SetValueFrom(ctx, types.StringType, a.AuthenticationSchema.OidcConfig.ProxyConfig.Acrs)
		diagnostics.Append(diags...)

		oidc.ProxyConfig, diags = types.ObjectValueFrom(ctx, proxyConfigObjType, proxyConfig)
		diagnostics.Append(diags...)
	} else {
		oidc.ProxyConfig = types.ObjectNull(proxyConfigObjType)
	}

	authenticationSchema.OpenIdConnectConfiguration, diags = types.ObjectValueFrom(ctx, openIdConnectConfigurationObjType, oidc)
	diagnostics.Append(diags...)

	// Authentication Schema SAML2
	// the mapping is done manually in order to handle the null values

	saml2Res := a.AuthenticationSchema.Saml2Configuration
	saml2Config := AppSaml2ConfigData{
		ResponseElementsToEncrypt: types.StringValue(saml2Res.ResponseElementsToEncrypt),
		DefaultNameIdFormat:       types.StringValue(saml2Res.DefaultNameIdFormat),
		SignSloMessages:           types.BoolValue(saml2Res.SignSLOMessages),
		RequireSignedSloMessages:  types.BoolValue(saml2Res.RequireSignedSLOMessages),
		RequireSignedAuthnRequest: types.BoolValue(saml2Res.RequireSignedAuthnRequest),
		SignAssertions:            types.BoolValue(saml2Res.SignAssertions),
		SignAuthnResponses:        types.BoolValue(saml2Res.SignAuthnResponses),
		DigestAlgorithm:           types.StringValue(saml2Res.DigestAlgorithm),
	}

	// SAML2
	// Saml Metadata URL
	if len(saml2Res.SamlMetadataUrl) > 0 {
		saml2Config.SamlMetadataUrl = types.StringValue(saml2Res.SamlMetadataUrl)
	}

	// SAML2 ACS Endpoints
	if len(a.AuthenticationSchema.Saml2Configuration.AcsEndpoints) > 0 {
		saml2Config.AcsEndpoints, diags = types.ListValueFrom(ctx, acsEndpointsObjType, a.AuthenticationSchema.Saml2Configuration.AcsEndpoints)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return application, diagnostics
		}
	} else {
		saml2Config.AcsEndpoints = types.ListNull(acsEndpointsObjType)
	}

	// SAML2 SLO Endpoints
	if len(a.AuthenticationSchema.Saml2Configuration.SloEndpoints) > 0 {

		endpointsData := []AppSloEndpointData{}

		for _, endpoint := range a.AuthenticationSchema.Saml2Configuration.SloEndpoints {
			endpointData := AppSloEndpointData{
				BindingName: types.StringValue(endpoint.BindingName),
				Location:    types.StringValue(endpoint.Location),
			}

			if len(endpoint.ResponseLocation) > 0 {
				endpointData.ResponseLocation = types.StringValue(endpoint.ResponseLocation)
			}

			endpointsData = append(endpointsData, endpointData)
		}

		saml2Endpoints, diags := types.ListValueFrom(ctx, appSaml2SloEndpointObjType, endpointsData)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return application, diagnostics
		}

		saml2Config.SloEndpoints = saml2Endpoints

	} else {
		saml2Config.SloEndpoints = types.ListNull(appSaml2SloEndpointObjType)
	}

	// SAML2 Signing Certificates
	if len(saml2Res.CertificatesForSigning) > 0 {
		certificates, diags := types.ListValueFrom(ctx, saml2SigningCertificateObjType, saml2Res.CertificatesForSigning)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return application, diagnostics
		}

		saml2Config.CertificatesForSigning = certificates
	} else {
		saml2Config.CertificatesForSigning = types.ListNull(saml2SigningCertificateObjType)
	}

	//SAML2 Encryption Certificate
	if saml2Res.CertificateForEncryption != nil {
		encryptionCertificate, diags := types.ObjectValueFrom(ctx, saml2EncryptionCertificateObjType.AttrTypes, saml2Res.CertificateForEncryption)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return application, diagnostics
		}

		saml2Config.CertificateForEncryption = encryptionCertificate
	} else {
		saml2Config.CertificateForEncryption = types.ObjectNull(saml2EncryptionCertificateObjType.AttrTypes)
	}

	if diagnostics.HasError() {
		return application, diagnostics
	}

	authenticationSchema.Saml2Configuration, diags = types.ObjectValueFrom(ctx, appSaml2ConfigObjType.AttrTypes, saml2Config)
	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return application, diagnostics
	}

	if a.AuthenticationSchema.SapManagedAttributes != nil {
		sapManagedAttributes := sapManagedAttributesData{
			ServiceInstanceId: types.StringValue(a.AuthenticationSchema.SapManagedAttributes.ServiceInstanceId),
			SourceAppId:       types.StringValue(a.AuthenticationSchema.SapManagedAttributes.SourceAppId),
			SourceTenantId:    types.StringValue(a.AuthenticationSchema.SapManagedAttributes.SourceTenantId),
			AppTenantId:       types.StringValue(a.AuthenticationSchema.SapManagedAttributes.AppTenantId),
			Type:              types.StringValue(a.AuthenticationSchema.SapManagedAttributes.Type),
			PlanName:          types.StringValue(a.AuthenticationSchema.SapManagedAttributes.PlanName),
			BtpTenantType:     types.StringValue(a.AuthenticationSchema.SapManagedAttributes.BtpTenantType),
		}
		authenticationSchema.SapManagedAttributes, diags = types.ObjectValueFrom(ctx, sapManagedAttributesObjType, sapManagedAttributes)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return application, diagnostics
		}
	} else {
		// each attribute is set to null as setting the whole object to null causes in place updates
		sapManagedAttributes := sapManagedAttributesData{
			ServiceInstanceId: types.StringNull(),
			SourceAppId:       types.StringNull(),
			SourceTenantId:    types.StringNull(),
			AppTenantId:       types.StringNull(),
			Type:              types.StringNull(),
			PlanName:          types.StringNull(),
			BtpTenantType:     types.StringNull(),
		}
		authenticationSchema.SapManagedAttributes, diags = types.ObjectValueFrom(ctx, sapManagedAttributesObjType, sapManagedAttributes)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return application, diagnostics
		}
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

// retrieve the API Request body from the plan data
func getApplicationRequest(ctx context.Context, plan applicationData) (*applications.Application, diag.Diagnostics) {

	var diagnostics, diags diag.Diagnostics

	args := &applications.Application{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		MultiTenantApp: plan.MultiTenantApp.ValueBool(),
	}

	if !plan.ParentApplicationId.IsNull() {
		args.ParentApplicationId = plan.ParentApplicationId.ValueString()
	}

	// mapping of the plan data to the API request body must be done manually
	// this is to ensure the proper handling of attributes where the API request body structure differs from that of the schema
	// these attributes are : subject_name_identifier and advanced_assertion_attributes
	// the schema defines them as objects with sub-attributes source and value, but the API request body expects them to be strings
	if !plan.AuthenticationSchema.IsNull() && !plan.AuthenticationSchema.IsUnknown() {

		args.AuthenticationSchema = &applications.AuthenticationSchema{}

		var authenticationSchema authenticationSchemaData
		diags = plan.AuthenticationSchema.As(ctx, &authenticationSchema, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})

		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		//SSO_TYPE
		if !authenticationSchema.SsoType.IsNull() && !authenticationSchema.SsoType.IsUnknown() {
			args.AuthenticationSchema.SsoType = authenticationSchema.SsoType.ValueString()
		}

		//DEFAULT_AUTHENTICATING_IDP
		if !authenticationSchema.DefaultAuthenticatingIdpId.IsNull() && !authenticationSchema.DefaultAuthenticatingIdpId.IsUnknown() {
			args.AuthenticationSchema.DefaultAuthenticatingIdpId = authenticationSchema.DefaultAuthenticatingIdpId.ValueString()
		}

		//SUBJECT_NAME_IDENTIFIER
		if !authenticationSchema.SubjectNameIdentifier.IsNull() && !authenticationSchema.SubjectNameIdentifier.IsUnknown() {

			var subjectNameIdentifier subjectNameIdentifierData
			diags = authenticationSchema.SubjectNameIdentifier.As(ctx, &subjectNameIdentifier, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})

			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return nil, diagnostics
			}

			// the mapping is done manually, in order to handle the parameter value when the source is set to "Corporate Identity Provider"
			if subjectNameIdentifier.Source.ValueString() == sourceValues[0] || subjectNameIdentifier.Source.ValueString() == sourceValues[2] {
				args.AuthenticationSchema.SubjectNameIdentifier = subjectNameIdentifier.Value.ValueString()
			} else {
				args.AuthenticationSchema.SubjectNameIdentifier = "${corporateIdP." + subjectNameIdentifier.Value.ValueString() + "}"
			}
		}

		//SUBJECT_NAME_IDENTIFIER_FUNCTION
		if !authenticationSchema.SubjectNameIdentifierFunction.IsNull() && !authenticationSchema.SubjectNameIdentifierFunction.IsUnknown() {
			args.AuthenticationSchema.SubjectNameIdentifierFunction = authenticationSchema.SubjectNameIdentifierFunction.ValueString()
		}

		//ASSERTION_ATTRIBUTES
		if !authenticationSchema.AssertionAttributes.IsNull() && !authenticationSchema.AssertionAttributes.IsUnknown() {

			var attributes []applications.AssertionAttribute
			diags := authenticationSchema.AssertionAttributes.ElementsAs(ctx, &attributes, true)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return nil, diagnostics
			}

			args.AuthenticationSchema.AssertionAttributes = attributes
		}

		//ADVANCED_ASSERTION_ATTRIBUTES
		if !authenticationSchema.AdvancedAssertionAttributes.IsNull() && !authenticationSchema.AdvancedAssertionAttributes.IsUnknown() {

			var advancedAssertionAttributes []advancedAssertionAttributesData
			diags := authenticationSchema.AdvancedAssertionAttributes.ElementsAs(ctx, &advancedAssertionAttributes, true)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return nil, diagnostics
			}

			for _, attribute := range advancedAssertionAttributes {

				assertionAttribute := applications.AdvancedAssertionAttribute{
					AttributeName: attribute.AttributeName.ValueString(),
				}

				// the mapping is done manually, in order to handle the parameter attribute_value when the source is set to "Corporate Identity Provider"
				if attribute.Source == types.StringValue(sourceValues[1]) {
					assertionAttribute.AttributeValue = "${corporateIdP." + attribute.AttributeValue.ValueString() + "}"
				} else {
					assertionAttribute.AttributeValue = attribute.AttributeValue.ValueString()
				}

				args.AuthenticationSchema.AdvancedAssertionAttributes = append(args.AuthenticationSchema.AdvancedAssertionAttributes, assertionAttribute)
			}
		}

		//AUTHENTICATION_RULES
		if !authenticationSchema.AuthenticationRules.IsNull() {

			var authrules []applications.AuthenicationRule
			diags = authenticationSchema.AuthenticationRules.ElementsAs(ctx, &authrules, true)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return nil, diagnostics
			}

			args.AuthenticationSchema.ConditionalAuthentication = authrules
		}

		//SAML2 CONFIGURATION
		if authenticationSchema.SsoType.ValueString() != "openIdConnect" {
			if !authenticationSchema.Saml2Configuration.IsNull() && !authenticationSchema.Saml2Configuration.IsUnknown() {

				var saml2config applications.SamlConfiguration
				diags := authenticationSchema.Saml2Configuration.As(ctx, &saml2config, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    true,
					UnhandledUnknownAsEmpty: true,
				})
				diagnostics.Append(diags...)

				if diagnostics.HasError() {
					return nil, diagnostics
				}

				args.AuthenticationSchema.Saml2Configuration = &saml2config
			}
		}

		//OPEN_ID_CONNECT_CONFIGURATION
		if authenticationSchema.SsoType.ValueString() != "saml2" {
			if !authenticationSchema.OpenIdConnectConfiguration.IsNull() && !authenticationSchema.OpenIdConnectConfiguration.IsUnknown() {

				var openIdConnectConfiguration openIdConnectConfigurationData
				diags := authenticationSchema.OpenIdConnectConfiguration.As(ctx, &openIdConnectConfiguration, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    true,
					UnhandledUnknownAsEmpty: true,
				})
				diagnostics.Append(diags...)

				oidc := &applications.OidcConfig{}

				// Redirect URIs
				if !openIdConnectConfiguration.RedirectUris.IsNull() {
					diags := openIdConnectConfiguration.RedirectUris.ElementsAs(ctx, &oidc.RedirectUris, true)
					diagnostics.Append(diags...)
				}

				// Post Logout Redirect URIs
				if !openIdConnectConfiguration.PostLogoutRedirectUris.IsNull() {
					diags := openIdConnectConfiguration.PostLogoutRedirectUris.ElementsAs(ctx, &oidc.PostLogoutRedirectUris, true)
					diagnostics.Append(diags...)
				}

				// Front Channel Logout URIs
				if !openIdConnectConfiguration.FrontChannelLogoutUris.IsNull() {
					diags := openIdConnectConfiguration.FrontChannelLogoutUris.ElementsAs(ctx, &oidc.FrontChannelLogoutUris, true)
					diagnostics.Append(diags...)
				}

				// Back Channel Logout URIs
				if !openIdConnectConfiguration.BackChannelLogoutUris.IsNull() {
					diags := openIdConnectConfiguration.BackChannelLogoutUris.ElementsAs(ctx, &oidc.BackChannelLogoutUris, true)
					diagnostics.Append(diags...)
				}

				// Restricted Grant Types
				if !openIdConnectConfiguration.RestrictedGrantTypes.IsNull() {
					diags := openIdConnectConfiguration.RestrictedGrantTypes.ElementsAs(ctx, &oidc.RestrictedGrantTypes, true)
					diagnostics.Append(diags...)
				}
				if diagnostics.HasError() {
					return nil, diagnostics
				}

				// Token Policy
				if !openIdConnectConfiguration.TokenPolicy.IsNull() {
					var token tokenPolicyData
					diags := openIdConnectConfiguration.TokenPolicy.As(ctx, &token, basetypes.ObjectAsOptions{
						UnhandledNullAsEmpty:    true,
						UnhandledUnknownAsEmpty: true,
					})
					diagnostics.Append(diags...)
					if diagnostics.HasError() {
						return nil, diagnostics
					}

					oidc.TokenPolicy = &applications.TokenPolicy{
						JwtValidity:                  token.JwtValidity.ValueInt32(),
						RefreshValidity:              token.RefreshValidity.ValueInt32(),
						RefreshParallel:              token.RefreshParallel.ValueInt32(),
						MaxExchangePeriod:            token.MaxExchangePeriod.ValueString(),
						RefreshTokenRotationScenario: token.RefreshTokenRotationScenario.ValueString(),
						AccessTokenFormat:            token.AccessTokenFormat.ValueString(),
					}
				}

				// Proxy Config
				if !openIdConnectConfiguration.ProxyConfig.IsNull() {
					var proxy proxyConfigData
					diags := openIdConnectConfiguration.ProxyConfig.As(ctx, &proxy, basetypes.ObjectAsOptions{
						UnhandledNullAsEmpty:    true,
						UnhandledUnknownAsEmpty: true,
					})
					diagnostics.Append(diags...)
					if diagnostics.HasError() {
						return nil, diagnostics
					}

					var acrs []string
					if !proxy.Acrs.IsNull() {
						diags := proxy.Acrs.ElementsAs(ctx, &acrs, true)
						diagnostics.Append(diags...)
					}
					if diagnostics.HasError() {
						return nil, diagnostics
					}
					oidc.ProxyConfig = &applications.OidcProxyConfig{
						Acrs: acrs,
					}
				}

				if diagnostics.HasError() {
					return nil, diagnostics
				}
				args.AuthenticationSchema.OidcConfig = oidc
			}
		}

		// SAP MANAGED ATTRIBUTES
		if !authenticationSchema.SapManagedAttributes.IsNull() && !authenticationSchema.SapManagedAttributes.IsUnknown() {

			var sapManagedAttributes sapManagedAttributesData
			diags := authenticationSchema.SapManagedAttributes.As(ctx, &sapManagedAttributes, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return nil, diagnostics
			}

			// reflect over the sapManagedAttributesData to set the attributes in the API request body
			// this avoids multiple if statements for each attribute
			// also helps to determine if at least one attribute is set, to decide whether to set the object in the request body or leave it as nil

			var attributes applications.SapManagedAttributes
			attributesVal := reflect.ValueOf(&attributes)

			setAttributes := false

			sapManagedAttributesVal := reflect.ValueOf(sapManagedAttributes)

			for i := 0; i < sapManagedAttributesVal.NumField(); i++ {

				fieldName := sapManagedAttributesVal.Type().Field(i)
				fieldValue := sapManagedAttributesVal.Field(i).Interface().(types.String)

				elem := attributesVal.Elem()
				field := elem.FieldByName(fieldName.Name)

				if len(fieldValue.ValueString()) > 0 {
					field.SetString(fieldValue.ValueString())
					setAttributes = true
				} else {
					field.SetString(types.StringNull().ValueString())
				}
			}

			if setAttributes {
				attributes = attributesVal.Elem().Interface().(applications.SapManagedAttributes)
				args.AuthenticationSchema.SapManagedAttributes = &attributes
			} else {
				args.AuthenticationSchema.SapManagedAttributes = nil
			}

		}

	}
	return args, diagnostics
}

func getUpdateRequest(ctx context.Context, plan applicationData, state applicationData) []generic.PatchRequest {

	reqs := []generic.PatchRequest{}

	argsType := reflect.TypeFor[applicationData]()

	if !plan.Name.Equal(state.Name) {
		reqs = append(reqs, getPatchRequest("replace", "Name", "", plan.Name.ValueString(), argsType))
	}

	if !plan.Description.Equal(state.Description) {
		reqs = append(reqs, getPatchRequest("replace", "Description", "", plan.Description.ValueString(), argsType))
	}

	if !plan.MultiTenantApp.Equal(state.MultiTenantApp) {
		reqs = append(reqs, getPatchRequest("replace", "MultiTenantApp", "", plan.MultiTenantApp.ValueBool(), argsType))
	}

	if !plan.ParentApplicationId.Equal(state.ParentApplicationId) {
		reqs = append(reqs, getPatchRequest("replace", "ParentApplicationId", "", plan.ParentApplicationId.ValueString(), argsType))
	}

	if !plan.AuthenticationSchema.Equal(state.AuthenticationSchema) {

		arg, _ := argsType.FieldByName("AuthenticationSchema")
		path := fmt.Sprintf("/%s", arg.Tag.Get("json"))

		argsType = reflect.TypeFor[authenticationSchemaData]()

		var planAuthSchema, stateAuthSchema authenticationSchemaData

		_ = plan.AuthenticationSchema.As(ctx, &planAuthSchema, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})

		_ = state.AuthenticationSchema.As(ctx, &stateAuthSchema, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})

		if !planAuthSchema.SsoType.Equal(stateAuthSchema.SsoType) {
			reqs = append(reqs, getPatchRequest("replace", "SsoType", path, planAuthSchema.SsoType.ValueString(), argsType))
		}

		if !planAuthSchema.SubjectNameIdentifier.Equal(stateAuthSchema.SubjectNameIdentifier) {
			var planSubjectNameIdentifier subjectNameIdentifierData

			_ = planAuthSchema.SubjectNameIdentifier.As(ctx, &planSubjectNameIdentifier, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})

			if planSubjectNameIdentifier.Source.ValueString() == sourceValues[0] || planSubjectNameIdentifier.Source.ValueString() == sourceValues[2] {
				reqs = append(reqs, getPatchRequest("replace", "SubjectNameIdentifier", path, planSubjectNameIdentifier.Value.ValueString(), argsType))
			} else {
				reqs = append(reqs, getPatchRequest("replace", "SubjectNameIdentifier", path, "${corporateIdP."+planSubjectNameIdentifier.Value.ValueString()+"}", argsType))
			}
		}

		if !planAuthSchema.SubjectNameIdentifierFunction.Equal(stateAuthSchema.SubjectNameIdentifierFunction) {
			reqs = append(reqs, getPatchRequest("replace", "SubjectNameIdentifierFunction", path, planAuthSchema.SubjectNameIdentifierFunction.ValueString(), argsType))
		}

		if !planAuthSchema.AssertionAttributes.Equal(stateAuthSchema.AssertionAttributes) {

			planAssertionAttributes := []applications.AssertionAttribute{}

			if !planAuthSchema.AssertionAttributes.IsNull() {
				_ = planAuthSchema.AssertionAttributes.ElementsAs(ctx, &planAssertionAttributes, true)
			}

			reqs = append(reqs, getPatchRequest("replace", "AssertionAttributes", path, planAssertionAttributes, argsType))
		}

		if !planAuthSchema.AdvancedAssertionAttributes.Equal(stateAuthSchema.AdvancedAssertionAttributes) {

			attributes := []applications.AdvancedAssertionAttribute{}

			if !planAuthSchema.AdvancedAssertionAttributes.IsNull() {
				var planAdvancedAssertionAttributes []advancedAssertionAttributesData
				_ = planAuthSchema.AdvancedAssertionAttributes.ElementsAs(ctx, &planAdvancedAssertionAttributes, true)

				for _, attribute := range planAdvancedAssertionAttributes {

					assertionAttribute := applications.AdvancedAssertionAttribute{
						AttributeName: attribute.AttributeName.ValueString(),
					}

					// the mapping is done manually, in order to handle the parameter attribute_value when the source is set to "Corporate Identity Provider"
					if attribute.Source == types.StringValue(sourceValues[1]) {
						assertionAttribute.AttributeValue = "${corporateIdP." + attribute.AttributeValue.ValueString() + "}"
					} else {
						assertionAttribute.AttributeValue = attribute.AttributeValue.ValueString()
					}

					attributes = append(attributes, assertionAttribute)
				}
			}

			reqs = append(reqs, getPatchRequest("replace", "AdvancedAssertionAttributes", path, attributes, argsType))
		}

		if !planAuthSchema.DefaultAuthenticatingIdpId.Equal(stateAuthSchema.DefaultAuthenticatingIdpId) {
			reqs = append(reqs, getPatchRequest("replace", "DefaultAuthenticatingIdpId", path, planAuthSchema.DefaultAuthenticatingIdpId.ValueString(), argsType))
		}

		if !planAuthSchema.AuthenticationRules.Equal(stateAuthSchema.AuthenticationRules) {

			rules := []applications.AuthenicationRule{}

			if !planAuthSchema.AuthenticationRules.IsNull() {
				_ = planAuthSchema.AuthenticationRules.ElementsAs(ctx, &rules, true)
			}

			reqs = append(reqs, getPatchRequest("replace", "AuthenticationRules", path, rules, argsType))
		}

		if !planAuthSchema.OpenIdConnectConfiguration.Equal(stateAuthSchema.OpenIdConnectConfiguration) {

			arg, _ := argsType.FieldByName("OidcConfig")
			path = fmt.Sprintf("%s/%s", path, arg.Tag.Get("json"))

			argsType = reflect.TypeFor[applications.OidcConfig]()

			var planOidcSchema, stateOidcSchema openIdConnectConfigurationData

			_ = planAuthSchema.OpenIdConnectConfiguration.As(ctx, &planOidcSchema, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})

			_ = stateAuthSchema.OpenIdConnectConfiguration.As(ctx, &stateOidcSchema, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})

			if !planOidcSchema.RedirectUris.Equal(stateOidcSchema.RedirectUris) {
				val := []string{}

				if !planOidcSchema.RedirectUris.IsNull() {
					_ = planOidcSchema.RedirectUris.ElementsAs(ctx, &val, true)
				}

				reqs = append(reqs, getPatchRequest("replace", "RedirectUris", path, val, argsType))
			}

			if !planOidcSchema.PostLogoutRedirectUris.Equal(stateOidcSchema.PostLogoutRedirectUris) {
				val := []string{}

				if !planOidcSchema.PostLogoutRedirectUris.IsNull() {
					_ = planOidcSchema.PostLogoutRedirectUris.ElementsAs(ctx, &val, true)
				}

				reqs = append(reqs, getPatchRequest("replace", "PostLogoutRedirectUris", path, val, argsType))
			}

			if !planOidcSchema.FrontChannelLogoutUris.Equal(stateOidcSchema.FrontChannelLogoutUris) {
				val := []string{}

				if !planOidcSchema.FrontChannelLogoutUris.IsNull() {
					_ = planOidcSchema.FrontChannelLogoutUris.ElementsAs(ctx, &val, true)
				}

				reqs = append(reqs, getPatchRequest("replace", "FrontChannelLogoutUris", path, val, argsType))
			}

			if !planOidcSchema.BackChannelLogoutUris.Equal(stateOidcSchema.BackChannelLogoutUris) {
				val := []string{}

				if !planOidcSchema.BackChannelLogoutUris.IsNull() {
					_ = planOidcSchema.BackChannelLogoutUris.ElementsAs(ctx, &val, true)
				}

				reqs = append(reqs, getPatchRequest("replace", "BackChannelLogoutUris", path, val, argsType))
			}

			if !planOidcSchema.TokenPolicy.Equal(stateOidcSchema.TokenPolicy) {

				val := applications.TokenPolicy{}

				if !planOidcSchema.TokenPolicy.IsNull() {
					_ = planOidcSchema.TokenPolicy.As(ctx, &val, basetypes.ObjectAsOptions{
						UnhandledNullAsEmpty:    true,
						UnhandledUnknownAsEmpty: true,
					})
				}

				reqs = append(reqs, getPatchRequest("replace", "TokenPolicy", path, val, argsType))
			}

			if !planOidcSchema.RestrictedGrantTypes.Equal(stateOidcSchema.RestrictedGrantTypes) {
				val := []string{}

				if !planOidcSchema.RestrictedGrantTypes.IsNull() {
					_ = planOidcSchema.RestrictedGrantTypes.ElementsAs(ctx, &val, true)
				}

				reqs = append(reqs, getPatchRequest("replace", "RestrictedGrantTypes", path, val, argsType))
			}

			if !planOidcSchema.ProxyConfig.Equal(stateOidcSchema.ProxyConfig) {
				val := applications.OidcProxyConfig{}

				_ = planOidcSchema.ProxyConfig.As(ctx, &val, basetypes.ObjectAsOptions{
					UnhandledNullAsEmpty:    true,
					UnhandledUnknownAsEmpty: true,
				})

				reqs = append(reqs, getPatchRequest("replace", "ProxyConfig", path, val, argsType))
			}
		}

		if !planAuthSchema.Saml2Configuration.Equal(stateAuthSchema.Saml2Configuration) {

			arg, _ := argsType.FieldByName("Saml2Configuration")
			path = fmt.Sprintf("%s/%s", path, arg.Tag.Get("json"))

			argsType = reflect.TypeFor[applications.SamlConfiguration]()

			var planSaml2Schema, stateSaml2Schema AppSaml2ConfigData

			_ = planAuthSchema.Saml2Configuration.As(ctx, &planSaml2Schema, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})

			_ = stateAuthSchema.Saml2Configuration.As(ctx, &stateSaml2Schema, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})

			if !planSaml2Schema.SamlMetadataUrl.Equal(stateSaml2Schema.SamlMetadataUrl) {
				reqs = append(reqs, getPatchRequest("replace", "SamlMetadataUrl", path, planSaml2Schema.SamlMetadataUrl.ValueString(), argsType))
			}

			if !planSaml2Schema.AcsEndpoints.Equal(stateSaml2Schema.AcsEndpoints) {
				val := []applications.Saml2AcsEndpoint{}

				if !planSaml2Schema.AcsEndpoints.IsNull() {
					_ = planSaml2Schema.AcsEndpoints.ElementsAs(ctx, &val, true)
				}

				reqs = append(reqs, getPatchRequest("replace", "AcsEndpoints", path, val, argsType))
			}

			if !planSaml2Schema.SloEndpoints.Equal(stateSaml2Schema.SloEndpoints) {
				val := []applications.Saml2SLOEndpoint{}

				if !planSaml2Schema.SloEndpoints.IsNull() {
					_ = planSaml2Schema.SloEndpoints.ElementsAs(ctx, &val, true)
				}

				reqs = append(reqs, getPatchRequest("replace", "SloEndpoints", path, val, argsType))
			}

			if !planSaml2Schema.CertificatesForSigning.Equal(stateSaml2Schema.CertificatesForSigning) {
				val := []corporateidps.SigningCertificateData{}

				if !planSaml2Schema.CertificatesForSigning.IsNull() {
					_ = planSaml2Schema.CertificatesForSigning.ElementsAs(ctx, &val, true)
				}

				reqs = append(reqs, getPatchRequest("replace", "CertificatesForSigning", path, val, argsType))
			}

			if !planSaml2Schema.CertificateForEncryption.Equal(stateSaml2Schema.CertificateForEncryption) {
				val := applications.EncryptionCertificateData{}

				if !planSaml2Schema.CertificatesForSigning.IsNull() {
					_ = planSaml2Schema.CertificateForEncryption.As(ctx, &val, basetypes.ObjectAsOptions{
						UnhandledNullAsEmpty:    true,
						UnhandledUnknownAsEmpty: true,
					})
				}

				reqs = append(reqs, getPatchRequest("replace", "CertificateForEncryption", path, val, argsType))
			}

			if !planSaml2Schema.ResponseElementsToEncrypt.Equal(stateSaml2Schema.ResponseElementsToEncrypt) {
				reqs = append(reqs, getPatchRequest("replace", "ResponseElementsToEncrypt", path, planSaml2Schema.ResponseElementsToEncrypt.ValueString(), argsType))
			}

			if !planSaml2Schema.DefaultNameIdFormat.Equal(stateSaml2Schema.DefaultNameIdFormat) {
				reqs = append(reqs, getPatchRequest("replace", "DefaultNameIdFormat", path, planSaml2Schema.DefaultNameIdFormat.ValueString(), argsType))
			}

			if !planSaml2Schema.SignSloMessages.Equal(stateSaml2Schema.SignSloMessages) {
				reqs = append(reqs, getPatchRequest("replace", "SignSLOMessages", path, planSaml2Schema.SignSloMessages.ValueBool(), argsType))
			}

			if !planSaml2Schema.RequireSignedSloMessages.Equal(stateSaml2Schema.RequireSignedSloMessages) {
				reqs = append(reqs, getPatchRequest("replace", "RequireSignedSLOMessages", path, planSaml2Schema.RequireSignedSloMessages.ValueBool(), argsType))
			}

			if !planSaml2Schema.RequireSignedAuthnRequest.Equal(stateSaml2Schema.RequireSignedAuthnRequest) {
				reqs = append(reqs, getPatchRequest("replace", "RequireSignedAuthnRequest", path, planSaml2Schema.RequireSignedAuthnRequest.ValueBool(), argsType))
			}

			if !planSaml2Schema.SignAssertions.Equal(stateSaml2Schema.SignAssertions) {
				reqs = append(reqs, getPatchRequest("replace", "SignAssertions", path, planSaml2Schema.SignAssertions.ValueBool(), argsType))
			}

			if !planSaml2Schema.SignAuthnResponses.Equal(stateSaml2Schema.SignAuthnResponses) {
				reqs = append(reqs, getPatchRequest("replace", "SignAuthnResponses", path, planSaml2Schema.SignAuthnResponses.ValueBool(), argsType))
			}

			if !planSaml2Schema.DigestAlgorithm.Equal(stateSaml2Schema.DigestAlgorithm) {
				reqs = append(reqs, getPatchRequest("replace", "DigestAlgorithm", path, planSaml2Schema.DigestAlgorithm.ValueString(), argsType))
			}
		}
	}

	return reqs
}

func getPatchRequest(operation string, attrName string, path string, value any, argsType reflect.Type) generic.PatchRequest {

	arg, _ := argsType.FieldByName(attrName)
	tag := fmt.Sprintf("/%s", arg.Tag.Get("json"))

	if path != "" {
		tag = fmt.Sprintf("%s%s", path, tag)
	}

	return generic.PatchRequest{
		Op:    operation,
		Path:  tag,
		Value: value,
	}

}
