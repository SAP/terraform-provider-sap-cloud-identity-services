package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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
	OpenIdConnectConfiguration    types.Object `tfsdk:"oidc_config"`
	Saml2Configuration            types.Object `tfsdk:"saml2_config"`
	SapManagedAttributes          types.Object `tfsdk:"sap_managed_attributes"`
}

type AppSaml2ConfigData struct {
	SamlMetadataUrl           types.String `tfsdk:"saml_metadata_url"`
	AcsEndpoints              types.List   `tfsdk:"acs_endpoints"`
	SloEndpoints              types.List   `tfsdk:"slo_endpoints"`
	CertificatesForSigning    types.List   `tfsdk:"signing_certificates"`
	CertificateForEncryption  types.Object `tfsdk:"encryption_certificate"`
	ResponseElementsToEncrypt types.String `tfsdk:"response_elements_to_encrypt"`
	DefaultNameIdFormat       types.String `tfsdk:"default_name_id_format"`
	SignSloMessages           types.Bool   `tfsdk:"sign_slo_messages"`
	RequireSignedSloMessages  types.Bool   `tfsdk:"require_signed_slo_messages"`
	RequireSignedAuthnRequest types.Bool   `tfsdk:"require_signed_auth_requests"`
	SignAssertions            types.Bool   `tfsdk:"sign_assertions"`
	SignAuthnResponses        types.Bool   `tfsdk:"sign_auth_responses"`
	DigestAlgorithm           types.String `tfsdk:"digest_algorithm"`
}

type AppSloEndpointData struct {
	BindingName      types.String `tfsdk:"binding_name"`
	Location         types.String `tfsdk:"location"`
	ResponseLocation types.String `tfsdk:"response_location"`
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

type openIdConnectConfigurationData struct {
	RedirectUris           types.Set    `tfsdk:"redirect_uris"`
	PostLogoutRedirectUris types.Set    `tfsdk:"post_logout_redirect_uris"`
	FrontChannelLogoutUris types.Set    `tfsdk:"front_channel_logout_uris"`
	BackChannelLogoutUris  types.Set    `tfsdk:"back_channel_logout_uris"`
	TokenPolicy            types.Object `tfsdk:"token_policy"`
	RestrictedGrantTypes   types.Set    `tfsdk:"restricted_grant_types"`
	ProxyConfig            types.Object `tfsdk:"proxy_config"`
}

type tokenPolicyData struct {
	JwtValidity                  types.Int32  `tfsdk:"jwt_validity"`
	RefreshValidity              types.Int32  `tfsdk:"refresh_validity"`
	RefreshParallel              types.Int32  `tfsdk:"refresh_parallel"`
	MaxExchangePeriod            types.String `tfsdk:"max_exchange_period"`
	RefreshTokenRotationScenario types.String `tfsdk:"refresh_token_rotation_scenario"`
	AccessTokenFormat            types.String `tfsdk:"access_token_format"`
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
	if a.AuthenticationSchema.OidcConfig != nil {
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
	} else {
		authenticationSchema.OpenIdConnectConfiguration = types.ObjectNull(openIdConnectConfigurationObjType)
	}

	// Authentication Schema SAML2
	// the mapping is done manually in order to handle the null values
	if a.AuthenticationSchema.Saml2Configuration != nil {
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
	} else {
		authenticationSchema.Saml2Configuration = types.ObjectNull(appSaml2ConfigObjType.AttrTypes)
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
		authenticationSchema.SapManagedAttributes = types.ObjectNull(sapManagedAttributesObjType)
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

		//OPEN_ID_CONNECT_CONFIGURATION
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

		if !authenticationSchema.SapManagedAttributes.IsNull() && !authenticationSchema.SapManagedAttributes.IsUnknown() {
			var sapManagedAttributes applications.SapManagedAttributes
			diags := authenticationSchema.SapManagedAttributes.As(ctx, &sapManagedAttributes, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return nil, diagnostics
			}

			args.AuthenticationSchema.SapManagedAttributes = &sapManagedAttributes
		}

	}
	return args, diagnostics
}
