package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var regexpUUID = utils.UuidRegexp

func TestResourceApplication(t *testing.T) {

	application := applications.Application{
		Name:        "basic-test-app",
		Description: "application for testing purposes",
		AuthenticationSchema: &applications.AuthenticationSchema{
			SubjectNameIdentifier:         "mail",
			SubjectNameIdentifierFunction: "lowerCase",
			AssertionAttributes: []applications.AssertionAttribute{
				{
					AssertionAttributeName: "param1",
					UserAttributeName:      "firstName",
				},
				{
					AssertionAttributeName: "param2",
					UserAttributeName:      "mail",
				},
			},
			AdvancedAssertionAttributes: []applications.AdvancedAssertionAttribute{
				{
					AttributeName:  "adv_param1",
					AttributeValue: "value1",
				},
				{
					AttributeName:  "adv_param2",
					AttributeValue: "value2",
				},
			},
			DefaultAuthenticatingIdpId: "c93f6b04-7a0f-42c1-b3c5-3b30d0ad8910",
			ConditionalAuthentication: []applications.AuthenicationRule{
				{
					UserType:           "employee",
					UserEmailDomain:    "gmail.com",
					IpNetworkRange:     "10.0.0.1/8",
					IdentityProviderId: "c93f6b04-7a0f-42c1-b3c5-3b30d0ad8910",
				},
			},
			// SapManagedAttributes: &applications.SapManagedAttributes{
			// 		ServiceInstanceId: "c93f6b04-7a0f-42c1-b3c5-3b30d0ad8910",
			// 		SourceAppId: 	 "c93f6b04-7a0f-42c1-b3c5-3b30d0ad8910",
			// 		SourceTenantId: "c93f6b04-7a0f-42c1-b3c5-3b30d0ad8910",
			// 		AppTenantId:    "c93f6b04-7a0f-42c1-b3c5-3b30d0ad8910",
			// 		Type:   	"xsuaa",
			// 		PlanName: "application",
			// 		BtpTenantType: "customer",
			// },
		},
	}
	oidcApplication := applications.Application{
		Name:        "OIDC-test-app",
		Description: "application for testing purposes",
		AuthenticationSchema: &applications.AuthenticationSchema{
			SsoType: "openIdConnect",
			OidcConfig: &applications.OidcConfig{
				RedirectUris: []string{
					"https://redirectUris.com",
				},
				PostLogoutRedirectUris: []string{
					"https://postLogoutRedirectUris.com",
				},
				FrontChannelLogoutUris: []string{
					"https://frontChannelLogoutUris.com",
				},
				BackChannelLogoutUris: []string{
					"https://backChannelLogoutUris.com",
				},
				TokenPolicy: &applications.TokenPolicy{
					JwtValidity:                  41000,
					RefreshValidity:              1500000,
					RefreshParallel:              5,
					MaxExchangePeriod:            "unlimited",
					RefreshTokenRotationScenario: "off",
					AccessTokenFormat:            "default",
				},
				RestrictedGrantTypes: []applications.GrantType{
					"clientCredentials",
					"authorizationCode",
				},
				ProxyConfig: &applications.OidcProxyConfig{
					Acrs: []string{
						"acrs",
					},
				},
			},
		},
	}

	saml2App := applications.Application{
		Name: "SAML2 - Test App",
		AuthenticationSchema: &applications.AuthenticationSchema{
			SsoType: "saml2",
			Saml2Configuration: &applications.SamlConfiguration{
				SamlMetadataUrl: "https://test.com",
				AcsEndpoints: []applications.Saml2AcsEndpoint{
					{
						BindingName: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
						Location:    "https://test.1.com",
						Index:       1,
						IsDefault:   true,
					},
				},
				SloEndpoints: []applications.Saml2SLOEndpoint{
					{
						BindingName:      "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
						Location:         "https://logout.1.com",
						ResponseLocation: "https://logout-response.1.com",
					},
				},
				CertificatesForSigning: []corporateidps.SigningCertificateData{
					{
						// Always replace with a valid certificate for recording of fixtures
						Base64Certificate: "-----BEGIN CERTIFICATE-----\\nredacted\\n-----END CERTIFICATE-----",
						Dn:                "CN=Test Cert",
						IsDefault:         true,
					},
				},
				CertificateForEncryption: &applications.EncryptionCertificateData{
					// Always replace with a valid certificate for recording of fixtures
					Base64Certificate: "-----BEGIN CERTIFICATE-----\\nredacted\\n-----END CERTIFICATE-----",
				},
				ResponseElementsToEncrypt: "wholeAssertion",
				DefaultNameIdFormat:       "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
				SignSLOMessages:           true,
				RequireSignedSLOMessages:  true,
				RequireSignedAuthnRequest: true,
				SignAssertions:            true,
				SignAuthnResponses:        true,
				DigestAlgorithm:           "sha1",
			},
		},
	}

	t.Parallel()

	t.Run("happy path - application", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_application")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceApplication("testApp", application),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", application.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", application.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier.value", application.AuthenticationSchema.SubjectNameIdentifier),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier_function", application.AuthenticationSchema.SubjectNameIdentifierFunction),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_name", application.AuthenticationSchema.AssertionAttributes[0].AssertionAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", application.AuthenticationSchema.AssertionAttributes[0].UserAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.1.attribute_name", application.AuthenticationSchema.AssertionAttributes[1].AssertionAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.1.attribute_value", application.AuthenticationSchema.AssertionAttributes[1].UserAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_name", application.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_value", application.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_name", application.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_value", application.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.default_authenticating_idp", application.AuthenticationSchema.DefaultAuthenticatingIdpId),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.identity_provider_id", application.AuthenticationSchema.ConditionalAuthentication[0].IdentityProviderId),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.user_type", application.AuthenticationSchema.ConditionalAuthentication[0].UserType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.user_email_domain", application.AuthenticationSchema.ConditionalAuthentication[0].UserEmailDomain),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.ip_network_range", application.AuthenticationSchema.ConditionalAuthentication[0].IpNetworkRange),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.service_instance_id", application.AuthenticationSchema.SapManagedAttributes.ServiceInstanceId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.source_app_id", application.AuthenticationSchema.SapManagedAttributes.SourceAppId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.source_tenant_id", application.AuthenticationSchema.SapManagedAttributes.SourceTenantId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.app_tenant_id", application.AuthenticationSchema.SapManagedAttributes.AppTenantId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.type", application.AuthenticationSchema.SapManagedAttributes.Type),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.plan_name", application.AuthenticationSchema.SapManagedAttributes.PlanName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.btp_tenant_type", application.AuthenticationSchema.SapManagedAttributes.BtpTenantType),
					),
				},
				{
					ResourceName:      "sci_application.testApp",
					ImportState:       true,
					ImportStateVerify: true,
					// Given that the API always returns the internal ID of the IdP, the state verificiation of the attribute can be ignored in this test since it is configured with the UUID as seen above
					// The mismatch of IDs is expected behaviour and does not indicate an error, as the parameter can be configured with both the UUID and the internal ID of the IdP
					ImportStateVerifyIgnore: []string{"authentication_schema.default_authenticating_idp", "authentication_schema.sap_managed_attributes"},
				},
			},
		})
	})
	t.Run("happy path - bundled application1", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_bundled_application1")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceApplicationWithBundledApp("testBundledApp", "name"),
					ResourceName: 	"sci_application.testBundledApp",
					ImportState: 	true,
					ImportStateId: "73afa691-5946-4bb1-9c39-b404e4b21594",
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "id", "73afa691-5946-4bb1-9c39-b404e4b21594"),
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "name", "XSUAA_b75a605d-151c-4485-83f4-64604378e4ec"),
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.type", "xsuaa"),
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.app_tenant_id", "b75a605d-151c-4485-83f4-64604378e4ec"),
					),
				},
			},
		})
	})

	t.Run("happy path - bundled application2", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_bundled_application2")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceApplicationWithBundledApp("testBundledApp", "name"),
					ResourceName: 	"sci_application.testBundledApp",
					ImportState: 	true,
					ImportStateId: "31e38d9c-ca48-4227-963d-32e7dfcb5007",
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "id", "31e38d9c-ca48-4227-963d-32e7dfcb5007"),
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "name", "identity-subscription-c6c390f4-c9a2-4a6c-9cc7-01675a31e4f6-in-b75a605d-151c-4485-83f4-64604378e4ec"),
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.type", "xsuaa"),
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.app_tenant_id", "b75a605d-151c-4485-83f4-64604378e4ec"),
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.source_app_id", "3cc4b385-fe8b-423a-a8c0-34e15c9970c"),
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.source_tenant_id", "sapdas"),
						resource.TestCheckResourceAttr("sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.service_instance_id", "c6c390f4-c9a2-4a6c-9cc7-01675a31e4f6"),

					),
				},
			},
		})
	})

	t.Run("happy path - application update", func(t *testing.T) {

		updatedApplication := applications.Application{
			Name:        "test-app-updated",
			Description: "application for testing purposes",
			AuthenticationSchema: &applications.AuthenticationSchema{
				SubjectNameIdentifier:         "userUuid",
				SubjectNameIdentifierFunction: "upperCase",
				AssertionAttributes: []applications.AssertionAttribute{
					{
						AssertionAttributeName: "param1",
						UserAttributeName:      "lastName",
					},
				},
				AdvancedAssertionAttributes: []applications.AdvancedAssertionAttribute{
					{
						AttributeName:  "adv_param1",
						AttributeValue: "updated_value1",
					},
					{
						AttributeName:  "adv_param2",
						AttributeValue: "value2",
					},
					{
						AttributeName:  "adv_param3",
						AttributeValue: "value3",
					},
				},
				DefaultAuthenticatingIdpId: "7b56ab2b-dfc1-4a56-a8c3-830c2697a4d1",
				ConditionalAuthentication: []applications.AuthenicationRule{
					{
						UserType:           "customer",
						UserEmailDomain:    "sap.com",
						IpNetworkRange:     "192.168.1.1/24",
						IdentityProviderId: "7b56ab2b-dfc1-4a56-a8c3-830c2697a4d1",
					},
				},
				// SapManagedAttributes: &applications.SapManagedAttributes{
				// 	ServiceInstanceId: "7b56ab2b-dfc1-4a56-a8c3-830c2697a4d1",
				// 	SourceAppId: 	 "7b56ab2b-dfc1-4a56-a8c3-830c2697a4d1",
				// 	SourceTenantId: "7b56ab2b-dfc1-4a56-a8c3-830c2697a4d1",
				// 	AppTenantId:    "7b56ab2b-dfc1-4a56-a8c3-830c2697a4d1",
				// 	Type:   	"subscription",
				// 	PlanName: "application",
				// 	BtpTenantType: "customer",
				// },
			},
		}

		rec, user := setupVCR(t, "fixtures/resource_application_updated")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceApplication("testApp", application),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", application.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", application.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier.value", application.AuthenticationSchema.SubjectNameIdentifier),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier_function", application.AuthenticationSchema.SubjectNameIdentifierFunction),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_name", application.AuthenticationSchema.AssertionAttributes[0].AssertionAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", application.AuthenticationSchema.AssertionAttributes[0].UserAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.1.attribute_name", application.AuthenticationSchema.AssertionAttributes[1].AssertionAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.1.attribute_value", application.AuthenticationSchema.AssertionAttributes[1].UserAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_name", application.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_value", application.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_name", application.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_value", application.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.default_authenticating_idp", application.AuthenticationSchema.DefaultAuthenticatingIdpId),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.user_type", application.AuthenticationSchema.ConditionalAuthentication[0].UserType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.user_email_domain", application.AuthenticationSchema.ConditionalAuthentication[0].UserEmailDomain),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.ip_network_range", application.AuthenticationSchema.ConditionalAuthentication[0].IpNetworkRange),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.service_instance_id", application.AuthenticationSchema.SapManagedAttributes.ServiceInstanceId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.source_app_id", application.AuthenticationSchema.SapManagedAttributes.SourceAppId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.source_tenant_id", application.AuthenticationSchema.SapManagedAttributes.SourceTenantId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.app_tenant_id", application.AuthenticationSchema.SapManagedAttributes.AppTenantId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.type", application.AuthenticationSchema.SapManagedAttributes.Type),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.plan_name", application.AuthenticationSchema.SapManagedAttributes.PlanName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.btp_tenant_type", application.AuthenticationSchema.SapManagedAttributes.BtpTenantType),
					),
				},
				{
					Config: providerConfig("", user) + ResourceApplication("testApp", updatedApplication),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", updatedApplication.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", updatedApplication.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier.value", updatedApplication.AuthenticationSchema.SubjectNameIdentifier),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier_function", updatedApplication.AuthenticationSchema.SubjectNameIdentifierFunction),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_name", updatedApplication.AuthenticationSchema.AssertionAttributes[0].AssertionAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", updatedApplication.AuthenticationSchema.AssertionAttributes[0].UserAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_name", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_value", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_name", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_value", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.2.attribute_name", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[2].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.2.attribute_value", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[2].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.default_authenticating_idp", updatedApplication.AuthenticationSchema.DefaultAuthenticatingIdpId),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.user_type", updatedApplication.AuthenticationSchema.ConditionalAuthentication[0].UserType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.user_email_domain", updatedApplication.AuthenticationSchema.ConditionalAuthentication[0].UserEmailDomain),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.conditional_authentication.0.ip_network_range", updatedApplication.AuthenticationSchema.ConditionalAuthentication[0].IpNetworkRange),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.service_instance_id", updatedApplication.AuthenticationSchema.SapManagedAttributes.ServiceInstanceId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.source_app_id", updatedApplication.AuthenticationSchema.SapManagedAttributes.SourceAppId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.source_tenant_id", updatedApplication.AuthenticationSchema.SapManagedAttributes.SourceTenantId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.app_tenant_id", updatedApplication.AuthenticationSchema.SapManagedAttributes.AppTenantId),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.type", updatedApplication.AuthenticationSchema.SapManagedAttributes.Type),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.plan_name", updatedApplication.AuthenticationSchema.SapManagedAttributes.PlanName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sap_managed_attributes.btp_tenant_type", updatedApplication.AuthenticationSchema.SapManagedAttributes.BtpTenantType),
					),
				},
			},
		})
	})

	t.Run("happy path - oidc application", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_oidc_application")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + OidcResourceApplication("testApp", oidcApplication),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", oidcApplication.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", oidcApplication.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", oidcApplication.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.redirect_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.RedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.post_logout_redirect_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.PostLogoutRedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.front_channel_logout_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.FrontChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.back_channel_logout_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.BackChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.jwt_validity", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.JwtValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_validity", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_parallel", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshParallel)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.max_exchange_period", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.MaxExchangePeriod),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_token_rotation_scenario", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshTokenRotationScenario),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.access_token_format", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.AccessTokenFormat),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "authorizationCode"),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "clientCredentials"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.proxy_config.acrs.0", oidcApplication.AuthenticationSchema.OidcConfig.ProxyConfig.Acrs[0]),
					),
					ImportStateVerifyIgnore: []string{"authentication_schema.default_authenticating_idp", "authentication_schema.sap_managed_attributes"},
				},
			},
		})
	})

	t.Run("happy path - oidc application update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_oidc_application_updated")
		defer stopQuietly(rec)

		oidcUpdatedApplication := applications.Application{
			Name:        "OIDC-test-app-update",
			Description: "application for testing purposes",
			AuthenticationSchema: &applications.AuthenticationSchema{
				SsoType: "openIdConnect",
				OidcConfig: &applications.OidcConfig{
					RedirectUris: []string{
						"https://redirectUris2.com",
					},
					PostLogoutRedirectUris: []string{
						"https://postLogoutRedirectUris2.com",
					},
					FrontChannelLogoutUris: []string{
						"https://frontChannelLogoutUris2.com",
					},
					BackChannelLogoutUris: []string{
						"https://backChannelLogoutUris2.com",
					},
					TokenPolicy: &applications.TokenPolicy{
						JwtValidity:                  42000,
						RefreshValidity:              1600000,
						RefreshParallel:              6,
						MaxExchangePeriod:            "maxSessionValidity",
						RefreshTokenRotationScenario: "online",
						AccessTokenFormat:            "jwt",
					},
					RestrictedGrantTypes: []applications.GrantType{
						"refreshToken",
						"password",
					},
					ProxyConfig: &applications.OidcProxyConfig{
						Acrs: []string{
							"acrs",
						},
					},
				},
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + OidcResourceApplication("testApp", oidcApplication),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", oidcApplication.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", oidcApplication.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", oidcApplication.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.redirect_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.RedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.post_logout_redirect_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.PostLogoutRedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.front_channel_logout_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.FrontChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.back_channel_logout_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.BackChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.jwt_validity", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.JwtValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_validity", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_parallel", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshParallel)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.max_exchange_period", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.MaxExchangePeriod),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_token_rotation_scenario", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshTokenRotationScenario),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.access_token_format", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.AccessTokenFormat),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "authorizationCode"),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "clientCredentials"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.proxy_config.acrs.0", oidcApplication.AuthenticationSchema.OidcConfig.ProxyConfig.Acrs[0]),
					),
				},
				{
					Config: providerConfig("", user) + OidcResourceApplication("testApp", oidcUpdatedApplication),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", oidcUpdatedApplication.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", oidcUpdatedApplication.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", oidcUpdatedApplication.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.redirect_uris.0", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.RedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.post_logout_redirect_uris.0", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.PostLogoutRedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.front_channel_logout_uris.0", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.FrontChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.back_channel_logout_uris.0", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.BackChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.jwt_validity", fmt.Sprintf("%d", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.TokenPolicy.JwtValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_validity", fmt.Sprintf("%d", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_parallel", fmt.Sprintf("%d", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshParallel)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.max_exchange_period", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.TokenPolicy.MaxExchangePeriod),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_token_rotation_scenario", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshTokenRotationScenario),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.access_token_format", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.TokenPolicy.AccessTokenFormat),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "refreshToken"),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "password"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.proxy_config.acrs.0", oidcUpdatedApplication.AuthenticationSchema.OidcConfig.ProxyConfig.Acrs[0]),
					),
				},
			},
		})
	})

	t.Run("happy path - saml2 application", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_application_saml2")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceSaml2Application("testApp", saml2App),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", saml2App.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", saml2App.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.saml_metadata_url", saml2App.AuthenticationSchema.Saml2Configuration.SamlMetadataUrl),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.binding_name", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.location", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.index", fmt.Sprintf("%d", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Index)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.default", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.binding_name", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.location", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.response_location", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].ResponseLocation),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.signing_certificates.0.default", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.CertificatesForSigning[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.response_elements_to_encrypt", saml2App.AuthenticationSchema.Saml2Configuration.ResponseElementsToEncrypt),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.default_name_id_format", saml2App.AuthenticationSchema.Saml2Configuration.DefaultNameIdFormat),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_slo_messages", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_slo_messages", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.RequireSignedSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_auth_requests", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.RequireSignedAuthnRequest)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_assertions", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignAssertions)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_auth_responses", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignAuthnResponses)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.digest_algorithm", saml2App.AuthenticationSchema.Saml2Configuration.DigestAlgorithm),
					),
				},
				{
					ResourceName:      "sci_application.testApp",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - saml2 application update", func(t *testing.T) {

		updatedSaml2App := applications.Application{
			Name: "SAML2 - Updated Test App",
			AuthenticationSchema: &applications.AuthenticationSchema{
				SsoType: "saml2",
				Saml2Configuration: &applications.SamlConfiguration{
					SamlMetadataUrl: "https://updated.test.com",
					AcsEndpoints: []applications.Saml2AcsEndpoint{
						{
							BindingName: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
							Location:    "https://test.1.com",
							Index:       1,
							IsDefault:   false,
						},
						{
							BindingName: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
							Location:    "https://test.2.com",
							Index:       2,
							IsDefault:   true,
						},
					},
					SloEndpoints: []applications.Saml2SLOEndpoint{
						{
							BindingName:      "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
							Location:         "https://logout.1.com",
							ResponseLocation: "https://logout-response.1.com",
						},
						{
							BindingName:      "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
							Location:         "https://logout.2.com",
							ResponseLocation: "https://logout-response.2.com",
						},
					},
					CertificatesForSigning: []corporateidps.SigningCertificateData{
						{
							Base64Certificate: "-----BEGIN CERTIFICATE-----\\nredacted\\n-----END CERTIFICATE-----",
							IsDefault:         false,
						},
						{
							Base64Certificate: "-----BEGIN CERTIFICATE-----\\nredacted\\n-----END CERTIFICATE-----",
							IsDefault:         true,
						},
					},
					CertificateForEncryption: &applications.EncryptionCertificateData{
						Base64Certificate: "-----BEGIN CERTIFICATE-----\\nredacted\\n-----END CERTIFICATE-----",
					},
					ResponseElementsToEncrypt: "subjectNameId",
					DefaultNameIdFormat:       "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
					SignSLOMessages:           false,
					RequireSignedSLOMessages:  false,
					RequireSignedAuthnRequest: false,
					SignAssertions:            false,
					SignAuthnResponses:        false,
					DigestAlgorithm:           "sha256",
				},
			},
		}

		rec, user := setupVCR(t, "fixtures/resource_application_saml2_updated")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceSaml2Application("testApp", saml2App),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", saml2App.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", saml2App.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.saml_metadata_url", saml2App.AuthenticationSchema.Saml2Configuration.SamlMetadataUrl),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.binding_name", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.location", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.index", fmt.Sprintf("%d", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Index)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.default", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.binding_name", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.location", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.response_location", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].ResponseLocation),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.signing_certificates.0.default", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.CertificatesForSigning[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.response_elements_to_encrypt", saml2App.AuthenticationSchema.Saml2Configuration.ResponseElementsToEncrypt),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.default_name_id_format", saml2App.AuthenticationSchema.Saml2Configuration.DefaultNameIdFormat),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_slo_messages", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_slo_messages", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.RequireSignedSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_auth_requests", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.RequireSignedAuthnRequest)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_assertions", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignAssertions)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_auth_responses", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignAuthnResponses)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.digest_algorithm", saml2App.AuthenticationSchema.Saml2Configuration.DigestAlgorithm),
					),
				},
				{
					Config: providerConfig("", user) + ResourceSaml2Application("testApp", updatedSaml2App),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", updatedSaml2App.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", saml2App.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.saml_metadata_url", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SamlMetadataUrl),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.binding_name", updatedSaml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.location", updatedSaml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.index", fmt.Sprintf("%d", updatedSaml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Index)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.default", fmt.Sprintf("%t", updatedSaml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.1.binding_name", updatedSaml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[1].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.1.location", updatedSaml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[1].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.1.index", fmt.Sprintf("%d", updatedSaml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[1].Index)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.1.default", fmt.Sprintf("%t", updatedSaml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[1].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.binding_name", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.location", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.response_location", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].ResponseLocation),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.1.binding_name", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[1].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.1.location", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[1].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.1.response_location", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[1].ResponseLocation),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.signing_certificates.0.default", fmt.Sprintf("%t", updatedSaml2App.AuthenticationSchema.Saml2Configuration.CertificatesForSigning[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.signing_certificates.1.default", fmt.Sprintf("%t", updatedSaml2App.AuthenticationSchema.Saml2Configuration.CertificatesForSigning[1].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.response_elements_to_encrypt", updatedSaml2App.AuthenticationSchema.Saml2Configuration.ResponseElementsToEncrypt),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.default_name_id_format", updatedSaml2App.AuthenticationSchema.Saml2Configuration.DefaultNameIdFormat),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_slo_messages", fmt.Sprintf("%t", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SignSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_slo_messages", fmt.Sprintf("%t", updatedSaml2App.AuthenticationSchema.Saml2Configuration.RequireSignedSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_auth_requests", fmt.Sprintf("%t", updatedSaml2App.AuthenticationSchema.Saml2Configuration.RequireSignedAuthnRequest)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_assertions", fmt.Sprintf("%t", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SignAssertions)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_auth_responses", fmt.Sprintf("%t", updatedSaml2App.AuthenticationSchema.Saml2Configuration.SignAuthnResponses)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.digest_algorithm", updatedSaml2App.AuthenticationSchema.Saml2Configuration.DigestAlgorithm),
					),
				},
			},
		})
	})

	t.Run("happy path - oidc to saml2 application update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_oidcToSaml_application_updated")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + OidcResourceApplication("testApp", oidcApplication),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", oidcApplication.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", oidcApplication.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", oidcApplication.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.redirect_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.RedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.post_logout_redirect_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.PostLogoutRedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.front_channel_logout_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.FrontChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.back_channel_logout_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.BackChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.jwt_validity", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.JwtValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_validity", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_parallel", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshParallel)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.max_exchange_period", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.MaxExchangePeriod),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_token_rotation_scenario", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshTokenRotationScenario),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.access_token_format", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.AccessTokenFormat),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "authorizationCode"),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "clientCredentials"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.proxy_config.acrs.0", oidcApplication.AuthenticationSchema.OidcConfig.ProxyConfig.Acrs[0]),
					),
				},
				{
					Config: providerConfig("", user) + ResourceSaml2Application("testApp", saml2App),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", saml2App.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.saml_metadata_url", saml2App.AuthenticationSchema.Saml2Configuration.SamlMetadataUrl),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.binding_name", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.location", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.index", fmt.Sprintf("%d", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Index)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.default", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.binding_name", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.location", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.response_location", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].ResponseLocation),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.signing_certificates.0.default", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.CertificatesForSigning[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.response_elements_to_encrypt", saml2App.AuthenticationSchema.Saml2Configuration.ResponseElementsToEncrypt),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.default_name_id_format", saml2App.AuthenticationSchema.Saml2Configuration.DefaultNameIdFormat),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_slo_messages", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_slo_messages", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.RequireSignedSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_auth_requests", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.RequireSignedAuthnRequest)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_assertions", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignAssertions)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_auth_responses", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignAuthnResponses)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.digest_algorithm", saml2App.AuthenticationSchema.Saml2Configuration.DigestAlgorithm),
					),
				},
			},
		})
	})

	t.Run("happy path - saml2 to oidc application update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_application_saml2_to_oidc")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceSaml2Application("testApp", saml2App),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", saml2App.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.saml_metadata_url", saml2App.AuthenticationSchema.Saml2Configuration.SamlMetadataUrl),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.binding_name", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.location", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.index", fmt.Sprintf("%d", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].Index)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.default", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.AcsEndpoints[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.binding_name", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.location", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.response_location", saml2App.AuthenticationSchema.Saml2Configuration.SloEndpoints[0].ResponseLocation),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.signing_certificates.0.default", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.CertificatesForSigning[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.response_elements_to_encrypt", saml2App.AuthenticationSchema.Saml2Configuration.ResponseElementsToEncrypt),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.default_name_id_format", saml2App.AuthenticationSchema.Saml2Configuration.DefaultNameIdFormat),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_slo_messages", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_slo_messages", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.RequireSignedSLOMessages)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.require_signed_auth_requests", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.RequireSignedAuthnRequest)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_assertions", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignAssertions)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.sign_auth_responses", fmt.Sprintf("%t", saml2App.AuthenticationSchema.Saml2Configuration.SignAuthnResponses)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.saml2_config.digest_algorithm", saml2App.AuthenticationSchema.Saml2Configuration.DigestAlgorithm),
					),
				},
				{
					Config: providerConfig("", user) + OidcResourceApplication("testApp", oidcApplication),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", oidcApplication.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", oidcApplication.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", oidcApplication.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.redirect_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.RedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.post_logout_redirect_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.PostLogoutRedirectUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.front_channel_logout_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.FrontChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.back_channel_logout_uris.0", oidcApplication.AuthenticationSchema.OidcConfig.BackChannelLogoutUris[0]),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.jwt_validity", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.JwtValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_validity", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshValidity)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_parallel", fmt.Sprintf("%d", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshParallel)),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.max_exchange_period", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.MaxExchangePeriod),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_token_rotation_scenario", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshTokenRotationScenario),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.token_policy.access_token_format", oidcApplication.AuthenticationSchema.OidcConfig.TokenPolicy.AccessTokenFormat),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "authorizationCode"),
						resource.TestCheckTypeSetElemAttr("sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "clientCredentials"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.oidc_config.proxy_config.acrs.0", oidcApplication.AuthenticationSchema.OidcConfig.ProxyConfig.Acrs[0]),
					),
				},
			},
		})
	})

	t.Run("error path - name is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithoutAppName("testApp"),
					ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - parent_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithParent("testApp", "test-app", "application for testing purposes", "this-is-not-uuid"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute parent_application_id value must be a valid UUID, got:\n%s", "this-is-not-uuid")),
				},
			},
		})
	})

	t.Run("error path - sso_type needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSsoType("testApp", "test-app", "application for testing purposes", "this-is-not-a-valid-sso_type", ""),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.sso_type value must be one of:\n\\[\"openIdConnect\" \"saml2\"], got: \"%s\"", "this-is-not-a-valid-sso_type")),
				},
			},
		})
	})

	t.Run("error path - subject_name_identifier requires sub-attributes: source, value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSubjectNameIdentifier("testApp", "test-app", "application for testing purposes", "source = \"source\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.subject_name_identifier.value\" must be\nspecified when \"authentication_schema.subject_name_identifier\" is specified"),
				},
				{
					Config:      ResourceApplicationWithSubjectNameIdentifier("testApp", "test-app", "application for testing purposes", "value = \"value\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.subject_name_identifier.source\" must be\nspecified when \"authentication_schema.subject_name_identifier\" is specified"),
				},
			},
		})
	})

	t.Run("error path - subject_name_identifier.source needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSubjectNameIdentifier("testApp", "test-app", "application for testing purposes", "source = \"this-is-not-a-valid-source\", value = \"value\""),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.subject_name_identifier.source value must be\none of: \\[\"Identity Directory\" \"Corporate Identity Provider\" \"Expression\"],\ngot: \"%s\"", "this-is-not-a-valid-source")),
				},
			},
		})
	})

	t.Run("error path - subject_name_identifier_function needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSubjectNameIdentifierFunction("testApp", "test-app", "application for testing purposes", "invalid-function-name"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.subject_name_identifier_function value must\nbe one of: \\[\"none\" \"upperCase\" \"lowerCase\"], got: \"%s\"", "invalid-function-name")),
				},
			},
		})
	})

	t.Run("error path - assertion_attributes requires sub-attributes: attribute_name, attribute_value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAssertionAttributes("testApp", "test-app", "application for testing purposes", "attribute_value = \"value\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.assertion_attributes\\[0].attribute_name\" must\nbe specified when \"authentication_schema.assertion_attributes\" is specified"),
				},
				{
					Config:      ResourceApplicationWithAssertionAttributes("testApp", "test-app", "application for testing purposes", "attribute_name = \"name\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.assertion_attributes\\[0].attribute_value\"\nmust be specified when \"authentication_schema.assertion_attributes\" is\nspecified"),
				},
			},
		})
	})

	t.Run("error path - advanced_assertion_attributes requires sub-attributes: source, attribute_name, attribute_value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAdvancedAssertionAttributes("testApp", "test-app", "application for testing purposes", "source = \"source\", attribute_value = \"value\""),
					ExpectError: regexp.MustCompile("Attribute\n\"authentication_schema.advanced_assertion_attributes\\[0].attribute_name\" must\nbe specified when \"authentication_schema.advanced_assertion_attributes\" is\nspecified"),
				},
				{
					Config:      ResourceApplicationWithAdvancedAssertionAttributes("testApp", "test-app", "application for testing purposes", "attribute_name = \"name\", attribute_value = \"value\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.advanced_assertion_attributes\\[0].source\"\nmust be specified when \"authentication_schema.advanced_assertion_attributes\"\nis specified"),
				},
				{
					Config:      ResourceApplicationWithAdvancedAssertionAttributes("testApp", "test-app", "application for testing purposes", "source = \"source\", attribute_name = \"name\""),
					ExpectError: regexp.MustCompile("Attribute\n\"authentication_schema.advanced_assertion_attributes\\[0].attribute_value\" must\nbe specified when \"authentication_schema.advanced_assertion_attributes\" is\nspecified"),
				},
			},
		})
	})

	t.Run("error path - advanced_assertion_attributes.source needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAdvancedAssertionAttributes("testApp", "test-app", "application for testing purposes", "source = \"this-is-not-a-valid-source\", attribute_name = \"name\", attribute_value = \"value\""),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.advanced_assertion_attributes\\[0].source value\nmust be one of: \\[\"Corporate Identity Provider\" \"Expression\"], got:\n\"%s\"", "this-is-not-a-valid-source")),
				},
			},
		})
	})

	t.Run("error path - conditional_authentication requires sub-attribute: identity_provider_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "user_type = \"user\""),
					ExpectError: regexp.MustCompile("Attribute\n\"authentication_schema.conditional_authentication\\[0].identity_provider_id\"\nmust be specified when \"authentication_schema.conditional_authentication\" is\nspecified"),
				},
			},
		})
	})

	t.Run("error path - conditional_authentication requires atleast one of the following sub-attribute: user_type, user_group, user_email_domain, ip_network_range", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "identity_provider_id = \"664c660e25cff252c5c202dc\""),
					ExpectError: regexp.MustCompile("At least one attribute out of\n\\[authentication_schema.conditional_authentication\\[\\*].user_group,authentication_schema.conditional_authentication\\[\\*].user_email_domain,authentication_schema.conditional_authentication\\[\\*].ip_network_range"),
				},
			},
		})
	})

	t.Run("error path - email not a valid domain", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "identity_provider_id = \"664c660e25cff252c5c202dc\", user_email_domain=\"this-is-not-email-domain\""),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute\nauthentication_schema.conditional_authentication\\[0].user_email_domain value\nmust be a valid Email Domain, got: %s", "this-is-not-email-domain")),
				},
			},
		})
	})

	t.Run("error path - ip address is invalid", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "identity_provider_id = \"664c660e25cff252c5c202dc\", ip_network_range = \"this-is-not-ip-address\""),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute\nauthentication_schema.conditional_authentication\\[0].ip_network_range value\nmust be a valid IP Address with a valid CIDR notation, got:\n%s", "this-is-not-ip-address")),
				},
			},
		})
	})

	t.Run("error path - authentication_schema.sso_type must be a valid value when saml2 is configured", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSsoType("testApp", "test-app", "application for testing purposes", "openIdConnect", "saml2_config = {}"),
					ExpectError: regexp.MustCompile("Attribute authentication_schema.saml2_config : value of attribute\n\"authentication_schema.sso_type\" must be modified to match the IDP\nconfiguration provided. Acceptable values are : `saml2`, got: openIdConnect"),
				},
			},
		})
	})

	t.Run("error path - saml2_config.acs_endpoints requires sub-attributes: binding_name, location, index", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					// Missing binding_name
					Config:      ResourceApplicationWithSaml2AcsEndpoints("testApp", "test-app", "location = \"https://test.1.com\", index = 1"),
					ExpectError: regexp.MustCompile(`Attribute "authentication_schema.saml2_config.acs_endpoints\[0].binding_name"\nmust be specified when "authentication_schema.saml2_config.acs_endpoints" is\nspecified`),
				},
				{
					// Missing location
					Config:      ResourceApplicationWithSaml2AcsEndpoints("testApp", "test-app", "binding_name = \"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\", index = 1"),
					ExpectError: regexp.MustCompile(`Attribute "authentication_schema.saml2_config.acs_endpoints\[0].location" must\nbe specified when "authentication_schema.saml2_config.acs_endpoints" is\nspecified`),
				},
				{
					// Missing index
					Config:      ResourceApplicationWithSaml2AcsEndpoints("testApp", "test-app", "binding_name = \"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\", location = \"https://test.1.com\""),
					ExpectError: regexp.MustCompile(`Attribute "authentication_schema.saml2_config.acs_endpoints\[0].index" must be\nspecified when "authentication_schema.saml2_config.acs_endpoints" is\nspecified`),
				},
			},
		})
	})

	t.Run("error path - saml2_config.acs_endpoints.binding_name must be a valid value", func(t *testing.T) {

		bindingName := "invalid-binding"
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSaml2AcsEndpoints("testApp", "test-app", fmt.Sprintf(`binding_name = "%s", location = "https://test.1.com", index = 1`, bindingName)),
					ExpectError: regexp.MustCompile(fmt.Sprintf(`Attribute authentication_schema.saml2_config.acs_endpoints\[0].binding_name\nvalue must be one of: \["urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"\n"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"\n"urn:oasis:names:tc:SAML:2.0:bindings:SOAP"\n"urn:oasis:names:tc:SAML:2.0:bindings:URI"], got: "%s"`, bindingName)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.slo_endpoints requires sub-attributes: binding_name, location", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					// Missing binding_name
					Config:      ResourceApplicationWithSaml2SloEndpoints("testApp", "test-app", "location = \"https://test.1.com\""),
					ExpectError: regexp.MustCompile(`Attribute "authentication_schema.saml2_config.slo_endpoints\[0].binding_name"\nmust be specified when "authentication_schema.saml2_config.slo_endpoints" is\nspecified`),
				},
				{
					// Missing location
					Config:      ResourceApplicationWithSaml2SloEndpoints("testApp", "test-app", "binding_name = \"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\""),
					ExpectError: regexp.MustCompile(`Attribute "authentication_schema.saml2_config.slo_endpoints\[0].location" must\nbe specified when "authentication_schema.saml2_config.slo_endpoints" is\nspecified`),
				},
			},
		})
	})

	t.Run("error path - saml2_config.slo_endpoints.binding_name must be a valid value", func(t *testing.T) {

		bindingName := "invalid-binding"
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSaml2SloEndpoints("testApp", "test-app", fmt.Sprintf(`binding_name = "%s", location = "https://test.1.com"`, bindingName)),
					ExpectError: regexp.MustCompile(fmt.Sprintf(`Attribute authentication_schema.saml2_config.slo_endpoints\[0].binding_name\nvalue must be one of: \["urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"\n"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"\n"urn:oasis:names:tc:SAML:2.0:bindings:SOAP"\n"urn:oasis:names:tc:SAML:2.0:bindings:URI"], got: "%s"`, bindingName)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.signing_certificates requires sub-attribute: base64_certificate and default", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					// Missing base64_certificate
					Config:      ResourceApplicationWithSaml2SigningCertificates("testApp", "test-app", "default = true"),
					ExpectError: regexp.MustCompile(`Attribute\n"authentication_schema.saml2_config.signing_certificates\[0].base64_certificate"\nmust be specified when\n"authentication_schema.saml2_config.signing_certificates" is specified`),
				},
				{
					// Missing default
					Config:      ResourceApplicationWithSaml2SigningCertificates("testApp", "test-app", `base64_certificate = "-----BEGIN CERTIFICATE-----\\nredacted\\n-----END CERTIFICATE-----"`),
					ExpectError: regexp.MustCompile(`Attribute\n"authentication_schema.saml2_config.signing_certificates\[0].default" must be\nspecified when "authentication_schema.saml2_config.signing_certificates" is\nspecified`),
				},
			},
		})
	})

	t.Run("error path - saml2_config.signing_certificates.base64_certificate must be a valid value", func(t *testing.T) {

		certificate := "invalid-base64-certificate"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSaml2SigningCertificates("testApp", "test-app", fmt.Sprintf(`base64_certificate = "%s", default = true`, certificate)),
					ExpectError: regexp.MustCompile(fmt.Sprintf(`Attribute\nauthentication_schema\.saml2_config\.signing_certificates\[0\]\.base64_certificate\nvalue must be a valid PEM string in the format -----BEGIN\nCERTIFICATE-----\\n<certificate-content>\\n-----END CERTIFICATE-----, got:\n%s`, certificate)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.encryption_certificates requires sub-attribute: base64_certificate", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					// Missing base64_certificate
					Config:      ResourceApplicationWithSaml2EncryptionCertificates("testApp", "test-app", " "),
					ExpectError: regexp.MustCompile(`Attribute\n"authentication_schema.saml2_config.encryption_certificate.base64_certificate"\nmust be specified when\n"authentication_schema.saml2_config.encryption_certificate" is specified`),
				},
			},
		})
	})

	t.Run("error path - saml2_config.signing_certificates.base64_certificate must be a valid value", func(t *testing.T) {

		certificate := "invalid-base64-certificate"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSaml2EncryptionCertificates("testApp", "test-app", fmt.Sprintf(`base64_certificate = "%s"`, certificate)),
					ExpectError: regexp.MustCompile(fmt.Sprintf(`Attribute\nauthentication_schema.saml2_config.encryption_certificate.base64_certificate\nvalue must be a valid PEM string in the format -----BEGIN\nCERTIFICATE-----\\n<certificate-content>\\n-----END CERTIFICATE-----, got:\n%s`, certificate)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.default_name_id_format must be a valid value", func(t *testing.T) {
		defaultNameIdFormat := "invalid-name-id-format"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSaml2DefaultNameIdFormat("testApp", "test-app", defaultNameIdFormat),
					ExpectError: regexp.MustCompile(fmt.Sprintf(`Attribute authentication_schema.saml2_config.default_name_id_format value\nmust be one of: \["urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"\n"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"\n"urn:oasis:names:tc:SAML:2.0:nameid-format:persistent"\n"urn:oasis:names:tc:SAML:2.0:nameid-format:transient"], got:\n"%s"`, defaultNameIdFormat)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.response_elements_to_encrypt must be a valid value", func(t *testing.T) {
		responseElementsToEncrypt := "invalid-response-element"
		encryptionCertificate := "-----BEGIN CERTIFICATE-----\\nredacted\\n-----END CERTIFICATE-----"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSaml2ResponseElementsToEncrypt("testApp", "test-app", fmt.Sprintf("encryption_certificate = {base64_certificate = \"%s\"}", encryptionCertificate), responseElementsToEncrypt),
					ExpectError: regexp.MustCompile(fmt.Sprintf(`Attribute authentication_schema.saml2_config.response_elements_to_encrypt\nvalue must be one of: \["none" "wholeAssertion" "subjectNameId" "attributes"\n"subjectNameIdAndAttributes"], got: "%s"`, responseElementsToEncrypt)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.response_elements_to_encrypt requires attribute encryption_certificate", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSaml2ResponseElementsToEncrypt("testApp", "test-app", "", "none"),
					ExpectError: regexp.MustCompile(`Attribute "authentication_schema.saml2_config.encryption_certificate" must be\nspecified when\n"authentication_schema.saml2_config.response_elements_to_encrypt" is\nspecified`),
				},
			},
		})
	})

	t.Run("error path - saml2_config.digest_algorithm must be a valid value", func(t *testing.T) {
		digestAlgorithm := "invalid-digest-algorithm"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSaml2DigestAlgorithm("testApp", "test-app", digestAlgorithm),
					ExpectError: regexp.MustCompile(fmt.Sprintf(`Attribute authentication_schema.saml2_config.digest_algorithm value must be\none of: \["sha1" "sha256" "sha512"], got: "%s"`, digestAlgorithm)),
				},
			},
		})
	})
	t.Run("error path - oidc_config.redirect_uris is mandatory with oidc configuration", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithoutRedirectUris("testApp", "test-app", "application for testing purposes"),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.oidc_config.redirect_uris\" must be specified\nwhen \"authentication_schema.oidc_config\" is specified"),
				},
			},
		})

	})

	t.Run("error path - oidc_config.front_channel_logout invalid URI", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithFrontChannelLogoutUris("testApp", "test-app", "application for testing purposes", []string{"https://validUri.com", "this-is-not-a uri"}),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute\nauthentication_schema.oidc_config.front_channel_logout_uris\\[Value\\(\"this-is-not-a\nuri\"\\)\\] value must be a valid URL, got: %s", "this-is-not-a uri")),
				},
			},
		})
	})

	t.Run("error path - oidc_config.back_channel_logout invalid URI", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithBackChannelLogoutUris("testApp", "test-app", "application for testing purposes", []string{"https://validUri.com", "this-is-not-a uri"}),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute\nauthentication_schema.oidc_config.back_channel_logout_uris\\[Value\\(\"this-is-not-a\nuri\"\\)\\] value must be a valid URL, got: %s", "this-is-not-a uri")),
				},
			},
		})
	})

	t.Run("error path - oidc_config.token_policy.max_exchange_period needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithMaxExchangePeriod("testApp", "test-app", "application for testing purposes", "this-is-not-a-valid-max_exchange_period"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.oidc_config.token_policy.max_exchange_period\nvalue must be one of: \\[\"unlimited\" \"maxSessionValidity\"\n\"initialRefreshTokenValidity\"\\], got:\n\"%s\"", "this-is-not-a-valid-max_exchange_period")),
				},
			},
		})
	})

	t.Run("error path - oidc_config.token_policy.refresh_token_rotation_scenario needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithRefreshTokenRotationScenario("testApp", "test-app", "application for testing purposes", "this-is-not-a-valid-refresh_token_rotation_scenario"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute\nauthentication_schema.oidc_config.token_policy.refresh_token_rotation_scenario\nvalue must be one of: \\[\"off\" \"online\" \"mobile\"\\], got:\n\"%s\"", "this-is-not-a-valid-refresh_token_rotation_scenario")),
				},
			},
		})
	})

	t.Run("error path - oidc_config.token_policy.access_token_format needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAccessTokenFormat("testApp", "test-app", "application for testing purposes", "this-is-not-a-valid-access_token_format"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.oidc_config.token_policy.access_token_format\nvalue must be one of: \\[\"default\" \"jwt\" \"opaque\"\\], got:\n\"%s\"", "this-is-not-a-valid-access_token_format")),
				},
			},
		})
	})

	t.Run("error path - oidc_config.restricted_grant_types needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithRestrictedGrantTypes("testApp", "test-app", "application for testing purposes", []string{"this-is-not-a-valid-restricted_grant_type"}),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute\nauthentication_schema.oidc_config.restricted_grant_types\\[Value\\(\"this-is-not-a-valid-restricted_grant_type\"\\)\\]\nvalue must be one of: \\[\"clientCredentials\" \"authorizationCode\" \"refreshToken\"\n\"password\" \"implicit\" \"jwtBearer\" \"authorizationCodePkceS256\"\n\"tokenExchange\"\\], got: \"%s\"", "this-is-not-a-valid-restricted_grant_type")),
				},
			},
		})
	})

	// t.Run("error path - sap_managed_attributes.source_app_id needs to be a valid UUID", func(t *testing.T) {
	// 	resource.Test(t, resource.TestCase{
	// 		IsUnitTest:               true,
	// 		ProtoV6ProviderFactories: getTestProviders(nil),
	// 		Steps: []resource.TestStep{
	// 			{
	// 				Config:      ResourceApplicationWithSapManagedAttributesSourceAppID("testApp", "test-app", "application for testing purposes"),
	// 				//ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.sap_managed_attributes.source_app_id value must be a valid UUID, got:\n%s", "this-is-not-uuid")),
	// 				Check: resource.ComposeTestCheckFunc(
	// 					resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.0.sap_managed_attributes.0.source_app_id", ""),
	// 				),
	// 			},
	// 		},
	// 	})
	// })

	// t.Run("error path - sap_managed_attributes.type needs to be a valid value", func(t *testing.T) {
	// 	resource.Test(t, resource.TestCase{
	// 		IsUnitTest:               true,
	// 		ProtoV6ProviderFactories: getTestProviders(nil),
	// 		Steps: []resource.TestStep{
	// 			{
	// 				Config:      ResourceApplicationWithSapManagedAttributesType("testApp", "test-app", "application for testing purposes", "this-is-not-a-valid-type"),
	// 				ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.sap_managed_attributes.type value must be one of:\n\\[\"identityInstance\" \"subscription\" \"reuseInstance\" \"xsuaa\"\\], got:\n\"%s\"", "this-is-not-a-valid-type")),
	// 			},
	// 		},
	// 	})
	// })
}

func ResourceApplication(resourceName string, app applications.Application) string {

	authSchema := app.AuthenticationSchema

	var assertionAttributes string
	for _, attribute := range app.AuthenticationSchema.AssertionAttributes {
		assertionAttributes += fmt.Sprintf(`
				{
					attribute_name = "%s"
					attribute_value = "%s"
				},`, attribute.AssertionAttributeName, attribute.UserAttributeName)
	}

	var advancedAssertionAttributes string
	for _, attribute := range app.AuthenticationSchema.AdvancedAssertionAttributes {
		advancedAssertionAttributes += fmt.Sprintf(`
				{
					source = "Corporate Identity Provider",
					attribute_name = "%s",
					attribute_value = "%s",

				},`, attribute.AttributeName, attribute.AttributeValue)
	}

	var authenticationRules string
	for _, rule := range app.AuthenticationSchema.ConditionalAuthentication {
		authenticationRules += fmt.Sprintf(`{
					user_type = "%s"
					user_email_domain = "%s"
					ip_network_range = "%s"
					identity_provider_id = "%s"
					},`, rule.UserType, rule.UserEmailDomain, rule.IpNetworkRange, rule.IdentityProviderId)
	}

	authSchemaConfig := fmt.Sprintf(`
		subject_name_identifier = {
			source = "Identity Directory"
			value = "%s"
		}
		subject_name_identifier_function = "%s"
		assertion_attributes = [%s]
		advanced_assertion_attributes = [%s]
		default_authenticating_idp = "%s"
		conditional_authentication = [%s]
	`, authSchema.SubjectNameIdentifier, authSchema.SubjectNameIdentifierFunction, assertionAttributes, advancedAssertionAttributes, authSchema.DefaultAuthenticatingIdpId, authenticationRules)

	application := fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			%s
		}
	}`, resourceName, app.Name, app.Description, authSchemaConfig)

	return application
}

func ResourceApplicationWithBundledApp(resourceName string, appName string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
	}
	`, resourceName, appName)
}

func ResourceSaml2Application(resourceName string, app applications.Application) string {

	saml2Config := app.AuthenticationSchema.Saml2Configuration

	var acsEndpoints string
	for _, endpoint := range saml2Config.AcsEndpoints {
		acsEndpoints += fmt.Sprintf(`
				{
					binding_name = "%s"
					location = "%s"
					index = %d
					default = %t
				},
			`, endpoint.BindingName, endpoint.Location, endpoint.Index, endpoint.IsDefault)
	}

	var sloEndpoints string
	for _, endpoint := range saml2Config.SloEndpoints {
		sloEndpoints += fmt.Sprintf(`
                {
                    binding_name = "%s"
                    location = "%s"
                    response_location = "%s"
                },
            `, endpoint.BindingName, endpoint.Location, endpoint.ResponseLocation)
	}

	var signingCertificates string
	for _, cert := range saml2Config.CertificatesForSigning {
		signingCertificates += fmt.Sprintf(`
                {
                    base64_certificate = "%s"
                    default = %t
                },
            `, cert.Base64Certificate, cert.IsDefault)
	}

	encryptionCertificate := fmt.Sprintf(`{ base64_certificate = "%s"}`, saml2Config.CertificateForEncryption.Base64Certificate)

	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		authentication_schema = {
			sso_type = "saml2"
			saml2_config = {
				saml_metadata_url = "%s"
				acs_endpoints = [%s]
				slo_endpoints = [%s]
				signing_certificates = [%s]
				encryption_certificate = %s
				response_elements_to_encrypt = "%s"
				default_name_id_format = "%s"
				sign_slo_messages = %t
				require_signed_slo_messages = %t
				require_signed_auth_requests = %t
				sign_assertions = %t
				sign_auth_responses = %t
				digest_algorithm = "%s"
			}
		}
	}`, resourceName, app.Name, saml2Config.SamlMetadataUrl, acsEndpoints, sloEndpoints, signingCertificates, encryptionCertificate, saml2Config.ResponseElementsToEncrypt, saml2Config.DefaultNameIdFormat, saml2Config.SignSLOMessages, saml2Config.RequireSignedSLOMessages, saml2Config.RequireSignedAuthnRequest, saml2Config.SignAssertions, saml2Config.SignAuthnResponses, saml2Config.DigestAlgorithm)
}

func OidcResourceApplication(resourceName string, app applications.Application) string {

	var redirectUris string
	for _, uri := range app.AuthenticationSchema.OidcConfig.RedirectUris {
		redirectUris += fmt.Sprintf(`"%s",`, uri)
	}
	var postLogoutRedirectUris string
	for _, uri := range app.AuthenticationSchema.OidcConfig.PostLogoutRedirectUris {
		postLogoutRedirectUris += fmt.Sprintf(`"%s",`, uri)
	}
	var frontChannelLogoutUris string
	for _, uri := range app.AuthenticationSchema.OidcConfig.FrontChannelLogoutUris {
		frontChannelLogoutUris += fmt.Sprintf(`"%s",`, uri)
	}
	var backChannelLogoutUris string
	for _, uri := range app.AuthenticationSchema.OidcConfig.BackChannelLogoutUris {
		backChannelLogoutUris += fmt.Sprintf(`"%s",`, uri)
	}
	var restrictedGrantTypes string
	for _, grantType := range app.AuthenticationSchema.OidcConfig.RestrictedGrantTypes {
		restrictedGrantTypes += fmt.Sprintf(`"%s",`, grantType)
	}
	var acrs string
	for _, acr := range app.AuthenticationSchema.OidcConfig.ProxyConfig.Acrs {
		acrs += fmt.Sprintf(`"%s",`, acr)
	}

	return fmt.Sprintf(
		`resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			sso_type = "%s"
			oidc_config = {
				redirect_uris=[%s]
				post_logout_redirect_uris=[%s]
				front_channel_logout_uris=[%s]
				back_channel_logout_uris=[%s]
				token_policy = {
					jwt_validity = %d
					refresh_validity = %d		
					refresh_parallel = %d
					max_exchange_period = "%s"
					refresh_token_rotation_scenario="%s"
					access_token_format="%s"
				}
				restricted_grant_types=[%s]
				proxy_config = {
					acrs = [%s]
				}
			}
			
		}
	}`, resourceName, app.Name, app.Description, app.AuthenticationSchema.SsoType, redirectUris, postLogoutRedirectUris, frontChannelLogoutUris, backChannelLogoutUris, app.AuthenticationSchema.OidcConfig.TokenPolicy.JwtValidity, app.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshValidity, app.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshParallel, app.AuthenticationSchema.OidcConfig.TokenPolicy.MaxExchangePeriod, app.AuthenticationSchema.OidcConfig.TokenPolicy.RefreshTokenRotationScenario, app.AuthenticationSchema.OidcConfig.TokenPolicy.AccessTokenFormat, restrictedGrantTypes, acrs)
}

func ResourceApplicationWithParent(resourceName string, appName string, description string, parentAppId string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		parent_application_id = "%s"
	}
	`, resourceName, appName, description, parentAppId)
}

func ResourceApplicationWithAppId(resourceName string, appID string, appName string, description string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		id = "%s"
		name = "%s"
		description = "%s"
	}
	`, resourceName, appID, appName, description)
}

func ResourceApplicationWithoutAppName(resourceName string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
	}
	`, resourceName)
}

func ResourceApplicationWithSubjectNameIdentifier(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			subject_name_identifier = {
				%s
			}
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithSubjectNameIdentifierFunction(resourceName string, appName string, description string, functionName string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			subject_name_identifier_function = "%s"
		}
	}
	`, resourceName, appName, description, functionName)
}

func ResourceApplicationWithAssertionAttributes(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			assertion_attributes = [
				{	%s	}
			]
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithAdvancedAssertionAttributes(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			advanced_assertion_attributes = [
				{	%s	}
			]
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithSsoType(resourceName string, appName string, description string, ssoType string, config string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			sso_type = "%s"
			%s
		}
	}
	`, resourceName, appName, description, ssoType, config)
}

func ResourceApplicationWithAuthenticationRules(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			conditional_authentication = [
				{
					%s
				}
			]
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithFrontChannelLogoutUris(resourceName string, appName string, description string, subAttribute []string) string {

	var builder string
	for _, uri := range subAttribute {
		builder += fmt.Sprintf(`"%s",`, uri)
	}

	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			oidc_config = {
				redirect_uris = ["https://redirecturi.com"]
				front_channel_logout_uris = [%s]
			}	
		}
	}
	`, resourceName, appName, description, builder)
}

func ResourceApplicationWithBackChannelLogoutUris(resourceName string, appName string, description string, subAttribute []string) string {
	var builder string
	for _, uri := range subAttribute {
		builder += fmt.Sprintf(`"%s",`, uri)
	}

	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			oidc_config = {
				redirect_uris = ["https://redirecturi.com"]
				back_channel_logout_uris = [%s]
			}	
		}
	}
	`, resourceName, appName, description, builder)
}

func ResourceApplicationWithMaxExchangePeriod(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s"{
		name = "%s"
		description = "%s"
		authentication_schema = {
			oidc_config = {
				redirect_uris = ["https://redirectUris.com"]
				token_policy = {
					max_exchange_period = "%s"
				}
			}
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithRefreshTokenRotationScenario(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s"{
		name = "%s"
		description = "%s"
		authentication_schema = {
			oidc_config = {
				redirect_uris = ["https://redirectUris.com"]
				token_policy = {
					refresh_token_rotation_scenario = "%s"
				}
			}
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithAccessTokenFormat(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s"{
		name = "%s"
		description = "%s"
		authentication_schema = {
			oidc_config = {
				redirect_uris = ["https://redirectUris.com"]
				token_policy = {
					access_token_format = "%s"
				}
			}
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithRestrictedGrantTypes(resourceName string, appName string, description string, subAttribute []string) string {

	var builder string
	for _, uri := range subAttribute {
		builder += fmt.Sprintf(`"%s",`, uri)
	}

	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			oidc_config = {
				redirect_uris = ["https://redirecturi.com"]
				restricted_grant_types = [%s]
			}	
		}
	}
	`, resourceName, appName, description, builder)
}

func ResourceApplicationWithoutRedirectUris(resourceName string, appName string, description string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s"{
		name = "%s"
		description = "%s"
		authentication_schema = {
			oidc_config = {
			}
		}
	}
	`, resourceName, appName, description)
}
func ResourceApplicationWithSaml2AcsEndpoints(resourceName string, appName string, acsEndpoints string) string {
	return fmt.Sprintf(`
    resource "sci_application" "%s" {
        name = "%s"
        authentication_schema = {
            sso_type = "saml2"
            saml2_config = {
                acs_endpoints = [
                    {
						%s
					}
                ]
            }
        }
    }
    `, resourceName, appName, acsEndpoints)
}

func ResourceApplicationWithSaml2SloEndpoints(resourceName string, appName string, sloEndpoints string) string {
	return fmt.Sprintf(`
    resource "sci_application" "%s" {
        name = "%s"
        authentication_schema = {
            sso_type = "saml2"
            saml2_config = {
                slo_endpoints = [
                    {
						%s
					}
                ]
            }
        }
    }
    `, resourceName, appName, sloEndpoints)
}

func ResourceApplicationWithSaml2SigningCertificates(resourceName string, appName string, signingCertificate string) string {
	return fmt.Sprintf(`
    resource "sci_application" "%s" {
        name = "%s"
        authentication_schema = {
            sso_type = "saml2"
            saml2_config = {
                signing_certificates = [
                    {
                        %s
                    }
                ]
            }
        }
    }
    `, resourceName, appName, signingCertificate)
}

func ResourceApplicationWithSaml2EncryptionCertificates(resourceName string, appName string, encryptionCertificate string) string {
	return fmt.Sprintf(`
    resource "sci_application" "%s" {
        name = "%s"
        authentication_schema = {
            sso_type = "saml2"
            saml2_config = {
                encryption_certificate = {
					%s
				}
            }
        }
    }
    `, resourceName, appName, encryptionCertificate)
}

func ResourceApplicationWithSaml2DefaultNameIdFormat(resourceName string, appName string, nameIdFormat string) string {
	return fmt.Sprintf(`
    resource "sci_application" "%s" {
        name = "%s"
        authentication_schema = {
            sso_type = "saml2"
            saml2_config = {
                default_name_id_format = "%s"
            }
        }
    }
    `, resourceName, appName, nameIdFormat)
}

func ResourceApplicationWithSaml2ResponseElementsToEncrypt(resourceName string, appName string, encryptionCertificate string, element string) string {
	return fmt.Sprintf(`
    resource "sci_application" "%s" {
        name = "%s"
        authentication_schema = {
            sso_type = "saml2"
            saml2_config = {
				%s
                response_elements_to_encrypt = "%s"
            }
        }
    }
    `, resourceName, appName, encryptionCertificate, element)
}

func ResourceApplicationWithSaml2DigestAlgorithm(resourceName string, appName string, digestAlgorithm string) string {
	return fmt.Sprintf(`
    resource "sci_application" "%s" {
        name = "%s"
        authentication_schema = {
            sso_type = "saml2"
            saml2_config = {
				digest_algorithm = "%s"
            }
        }
    }
    `, resourceName, appName, digestAlgorithm)
}

// func ResourceApplicationWithSapManagedAttributesType(resourceName string, appName string, description string, appType string) string{
// 	return fmt.Sprintf(`
// 	resource "sci_application" "%s" {
// 		name = "%s"
// 		description = "%s"
// 		authentication_schema = {
// 			sap_managed_attributes = {
// 				type = "%s"
// 			}
// 		}
// 	}`, resourceName, appName, description, appType)
// }

// func ResourceApplicationWithSapManagedAttributesSourceAppID(resourceName string, appName string, description string) string{
// 	return fmt.Sprintf(`
// 	resource "sci_application" "%s" {
// 		name = "%s"
// 		description = "%s"
// 	}`, resourceName, appName, description)
// }