package provider

import (
	"fmt"
	"regexp"
	"testing"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceCorporateIdP(t *testing.T) {

	corporateIdP := corporateidps.IdentityProvider{
		DisplayName:           "Test Corporate IdP",
		Name:                  "Test IdP",
		ForwardAllSsoRequests: true,
		IdentityFederation: &corporateidps.IdentityFederation{
			UseLocalUserStore:        true,
			AllowLocalUsersOnly:      true,
			ApplyLocalIdPAuthnChecks: true,
			RequiredGroups: []string{
				"Test Group",
			},
		},
		LoginHintConfiguration: &corporateidps.LoginHintConfiguration{
			LoginHintType: "mail",
			SendMethod:    "urlParam",
		},
		LogoutUrl: "https://example.com/logout",
	}

	saml2Config := corporateidps.SAML2Configuration{
		SamlMetadataUrl: "https://example.com/saml2-metadata",
		AssertionAttributes: []corporateidps.AssertionAttribute{
			{
				Name:  "attr_name",
				Value: "attr_value",
			},
		},
		DigestAlgorithm:     "sha1",
		IncludeScoping:      true,
		DefaultNameIdFormat: "email",
		AllowCreate:         "true",
		CertificatesForSigning: []corporateidps.SigningCertificateData{
			{
				// Always replace with a valid certificate for recording of fixtures
				Base64Certificate: "-----BEGIN CERTIFICATE-----\\nredacted\\n-----END CERTIFICATE-----",
				IsDefault:         true,
				Dn:                "Test",
				ValidFrom:         "1999-01-01T00:00:00Z",
				ValidTo:           "9999-12-31T23:59:59Z",
			},
		},
		SsoEndpoints: []corporateidps.SAML2SSOEndpoint{
			{
				BindingName: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
				Location:    "https://test.com",
				IsDefault:   true,
			},
		},
		SloEndpoints: []corporateidps.SAML2SLOEndpoint{
			{
				BindingName:      "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
				Location:         "https://test.com",
				ResponseLocation: "https://test.com",
				IsDefault:        true,
			},
		},
	}

	oidcCoporateIdP := corporateidps.OIDCConfiguration{
		DiscoveryUrl:            "https://accounts.sap.com",
		ClientId:                "test-client-id",
		ClientSecret:            "test-client-secret",
		SubjectNameIdentifier:   "email",
		TokenEndpointAuthMethod: "clientSecretBasic",
		Scopes: []string{
			"test-value-1",
			"openid",
		},
		PkceEnabled: true,
		AdditionalConfig: &corporateidps.OIDCAdditionalConfig{
			OmitIDTokenHintForLogout: true,
			EnforceIssuerCheck:       true,
			EnforceNonce:             true,
		},
	}

	t.Parallel()

	t.Run("happy path - saml2 corporate idp", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config

		rec, user := setupVCR(t, "fixtures/resource_corporateIdP_saml2")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceCorporateIdP("testIdP", saml2IdP),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_corporate_idp.testIdP", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "display_name", saml2IdP.DisplayName),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "name", saml2IdP.Name),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "type", saml2IdP.Type),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "logout_url", saml2IdP.LogoutUrl),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "forward_all_sso_requests", fmt.Sprintf("%t", saml2IdP.ForwardAllSsoRequests)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.use_local_user_store", fmt.Sprintf("%t", saml2IdP.IdentityFederation.UseLocalUserStore)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.allow_local_users_only", fmt.Sprintf("%t", saml2IdP.IdentityFederation.AllowLocalUsersOnly)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.apply_local_idp_auth_and_checks", fmt.Sprintf("%t", saml2IdP.IdentityFederation.ApplyLocalIdPAuthnChecks)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.required_groups.0", saml2IdP.IdentityFederation.RequiredGroups[0]),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "login_hint_config.login_hint_type", saml2IdP.LoginHintConfiguration.LoginHintType),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "login_hint_config.send_method", saml2IdP.LoginHintConfiguration.SendMethod),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.saml_metadata_url", saml2IdP.Saml2Configuration.SamlMetadataUrl),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.assertion_attributes.0.name", saml2IdP.Saml2Configuration.AssertionAttributes[0].Name),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.assertion_attributes.0.value", saml2IdP.Saml2Configuration.AssertionAttributes[0].Value),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.digest_algorithm", saml2IdP.Saml2Configuration.DigestAlgorithm),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.include_scoping", fmt.Sprintf("%t", saml2IdP.Saml2Configuration.IncludeScoping)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.name_id_format", saml2IdP.Saml2Configuration.DefaultNameIdFormat),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.allow_create", saml2IdP.Saml2Configuration.AllowCreate),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.signing_certificates.0.dn", saml2IdP.Saml2Configuration.CertificatesForSigning[0].Dn),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.signing_certificates.0.default", fmt.Sprintf("%t", saml2IdP.Saml2Configuration.CertificatesForSigning[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.signing_certificates.0.valid_from", saml2IdP.Saml2Configuration.CertificatesForSigning[0].ValidFrom),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.signing_certificates.0.valid_to", saml2IdP.Saml2Configuration.CertificatesForSigning[0].ValidTo),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.sso_endpoints.0.binding_name", saml2IdP.Saml2Configuration.SsoEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.sso_endpoints.0.location", saml2IdP.Saml2Configuration.SsoEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.sso_endpoints.0.default", fmt.Sprintf("%t", saml2IdP.Saml2Configuration.SsoEndpoints[0].IsDefault)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.slo_endpoints.0.binding_name", saml2IdP.Saml2Configuration.SloEndpoints[0].BindingName),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.slo_endpoints.0.location", saml2IdP.Saml2Configuration.SloEndpoints[0].Location),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.slo_endpoints.0.response_location", saml2IdP.Saml2Configuration.SloEndpoints[0].ResponseLocation),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.slo_endpoints.0.default", fmt.Sprintf("%t", saml2IdP.Saml2Configuration.SloEndpoints[0].IsDefault)),
					),
				},
				{
					ResourceName:      "sci_corporate_idp.testIdP",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - oidc corporate idp", func(t *testing.T) {

		oidcIdP := corporateIdP
		oidcIdP.DisplayName = "OIDC - Test Corporate IdP"
		oidcIdP.Name = "OIDC - Test IdP"
		oidcIdP.Type = "openIdConnect"
		oidcIdP.OidcConfiguration = &oidcCoporateIdP

		rec, user := setupVCR(t, "fixtures/resource_corporateIdP_oidc")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceCorporateIdP("testIdP", oidcIdP),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_corporate_idp.testIdP", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "display_name", oidcIdP.DisplayName),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "name", oidcIdP.Name),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "type", oidcIdP.Type),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "logout_url", oidcIdP.LogoutUrl),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "forward_all_sso_requests", fmt.Sprintf("%t", oidcIdP.ForwardAllSsoRequests)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.use_local_user_store", fmt.Sprintf("%t", oidcIdP.IdentityFederation.UseLocalUserStore)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.allow_local_users_only", fmt.Sprintf("%t", oidcIdP.IdentityFederation.AllowLocalUsersOnly)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.apply_local_idp_auth_and_checks", fmt.Sprintf("%t", oidcIdP.IdentityFederation.ApplyLocalIdPAuthnChecks)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.required_groups.0", oidcIdP.IdentityFederation.RequiredGroups[0]),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "login_hint_config.login_hint_type", oidcIdP.LoginHintConfiguration.LoginHintType),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "login_hint_config.send_method", oidcIdP.LoginHintConfiguration.SendMethod),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.discovery_url", oidcIdP.OidcConfiguration.DiscoveryUrl),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.client_id", oidcIdP.OidcConfiguration.ClientId),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.client_secret", oidcIdP.OidcConfiguration.ClientSecret),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.subject_name_identifier", oidcIdP.OidcConfiguration.SubjectNameIdentifier),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.token_endpoint_auth_method", oidcIdP.OidcConfiguration.TokenEndpointAuthMethod),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.scopes.#", "2"),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.enable_pkce", fmt.Sprintf("%t", oidcCoporateIdP.PkceEnabled)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.additional_config.enforce_nonce", fmt.Sprintf("%t", oidcCoporateIdP.AdditionalConfig.EnforceNonce)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.additional_config.enforce_issuer_check", fmt.Sprintf("%t", oidcCoporateIdP.AdditionalConfig.EnforceIssuerCheck)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.additional_config.disable_logout_id_token_hint", fmt.Sprintf("%t", oidcCoporateIdP.AdditionalConfig.EnforceIssuerCheck)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.issuer", oidcCoporateIdP.DiscoveryUrl),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.jwks_uri", oidcCoporateIdP.DiscoveryUrl+"/oauth2/certs"),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.logout_endpoint", oidcCoporateIdP.DiscoveryUrl+"/oauth2/logout"),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.token_endpoint", oidcCoporateIdP.DiscoveryUrl+"/oauth2/token"),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.user_info_endpoint", oidcCoporateIdP.DiscoveryUrl+"/oauth2/userinfo"),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.authorization_endpoint", oidcCoporateIdP.DiscoveryUrl+"/oauth2/authorize"),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "oidc_config.is_client_secret_configured", fmt.Sprintf("%t", len(oidcIdP.OidcConfiguration.ClientSecret) != 0)),
					),
				},
				{
					ResourceName:      "sci_corporate_idp.testIdP",
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"oidc_config.client_secret", // The client secret is not returned by the GET call, so it cannot be verified
					},
				},
			},
		})
	})

	t.Run("error path - display_name is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdPWithoutDisplayName("testIdP"),
					ExpectError: regexp.MustCompile("The argument \"display_name\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - logout_url is not a valid URL", func(t *testing.T) {

		corporateIdP.LogoutUrl = "invalid-logout-url"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", corporateIdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute logout_url value must be a valid URL, got: %s", corporateIdP.LogoutUrl)),
				},
			},
		})
	})

	t.Run("error path - login_hint_config requires sub-attribute login_hint_type", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdPWithEmptyLoginHintConfig("testIdP", corporateIdP.Name),
					ExpectError: regexp.MustCompile("Attribute \"login_hint_config.login_hint_type\" must be specified when\n\"login_hint_config\" is specified"),
				},
			},
		})
	})

	t.Run("error path - login_hint_config.login_hint_type needs to be a valid value", func(t *testing.T) {

		corporateIdP.LoginHintConfiguration.LoginHintType = "invalid-login-hint-type"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", corporateIdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute login_hint_config.login_hint_type value must be one of: \\[\"none\"\n\"userInput\" \"mail\" \"loginName\"\\], got: \"%s\"", corporateIdP.LoginHintConfiguration.LoginHintType)),
				},
			},
		})
	})

	t.Run("error path - login_hint_config.send_method needs to be a valid value", func(t *testing.T) {

		corporateIdP.LoginHintConfiguration.SendMethod = "invalid-send-method"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", corporateIdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute login_hint_config.send_method value must be one of: \\[\"urlParam\"\n\"authRequest\"\\], got: \"%s\"", corporateIdP.LoginHintConfiguration.SendMethod)),
				},
			},
		})
	})

	t.Run("error path - saml2_config requires root attribute name & type", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdPWithoutNameType("testIdP", "Test SAML2 IDP", "saml2_config = {}"),
					ExpectError: regexp.MustCompile("Attribute \"name\" must be specified when \"saml2_config\" is specified"),
				},
				{
					Config:      ResourceCorporateIdPWithoutNameType("testIdP", "Test SAML2 IDP", "saml2_config = {}"),
					ExpectError: regexp.MustCompile("Attribute \"type\" must be specified when \"saml2_config\" is specified"),
				},
			},
		})
	})

	t.Run("error path - type needs to be a valid value when saml2 is configured", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdPConfigTypeMismtach("testIdP", "Test IDP", "openIdConnect", "saml2_config = {}"),
					ExpectError: regexp.MustCompile("Attribute saml2_config : value of attribute \"type\" must be modified to match\nthe IDP configuration provided. Acceptable values are : `sapSSO`,\n`microsoftADFS`, `saml2`, got: openIdConnect"),
				},
			},
		})
	})

	t.Run("error path - saml2_config.saml_metadata_url is not a valid URL", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config
		saml2IdP.Saml2Configuration.SamlMetadataUrl = "invalid-url"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", saml2IdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute saml2_config.saml_metadata_url value must be a valid URL, got:\n%s", saml2IdP.Saml2Configuration.SamlMetadataUrl)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.assertion_attributes requires sub-attributes name and value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSaml2CorporateIdPWithAssertionAttributes("testIdP", "Test SAML2 IDP", "value=\"value\""),
					ExpectError: regexp.MustCompile("Attribute \"saml2_config.assertion_attributes\\[0].name\" must be specified when\n\"saml2_config.assertion_attributes\" is specified"),
				},
				{
					Config:      ResourceSaml2CorporateIdPWithAssertionAttributes("testIdP", "Test SAML2 IDP", "name=\"name\""),
					ExpectError: regexp.MustCompile("Attribute \"saml2_config.assertion_attributes\\[0].value\" must be specified when\n\"saml2_config.assertion_attributes\" is specified"),
				},
			},
		})
	})

	t.Run("error path - saml2_config.sso_endpoints requires sub-attributes binding_name and location", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSaml2CorporateIdPWithSsoEndpoints("testIdP", "Test SAML2 IDP", "location=\"https://test.com\""),
					ExpectError: regexp.MustCompile("Attribute \"saml2_config.sso_endpoints\\[0].binding_name\" must be specified when\n\"saml2_config.sso_endpoints\" is specified"),
				},
				{
					Config:      ResourceSaml2CorporateIdPWithSsoEndpoints("testIdP", "Test SAML2 IDP", "binding_name=\"binding\""),
					ExpectError: regexp.MustCompile("Attribute \"saml2_config.sso_endpoints\\[0].location\" must be specified when\n\"saml2_config.sso_endpoints\" is specified"),
				},
			},
		})
	})

	t.Run("error path - saml2_config.sso_endpoints.binding_name needs to be a valid value", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config
		saml2IdP.Saml2Configuration.SsoEndpoints[0].BindingName = "invalid-binding"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", saml2IdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute saml2_config.sso_endpoints\\[0].binding_name value must be one of:\n\\[\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\"\n\"urn:oasis:names:tc:SAML:2.0:bindings:SOAP\"\n\"urn:oasis:names:tc:SAML:2.0:bindings:URI\"], got: \"%s\"", saml2IdP.Saml2Configuration.SsoEndpoints[0].BindingName)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.sso_endpoints.location is not a valid URL", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config
		saml2IdP.Saml2Configuration.SsoEndpoints[0].Location = "invalid-url"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", saml2IdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute saml2_config.sso_endpoints\\[0].location value must be a valid URL,\ngot: %s", saml2IdP.Saml2Configuration.SsoEndpoints[0].Location)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.slo_endpoints requires sub-attributes binding_name and location", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSaml2CorporateIdPWithSloEndpoints("testIdP", "Test SAML2 IDP", "location=\"https://test.com\""),
					ExpectError: regexp.MustCompile("Attribute \"saml2_config.slo_endpoints\\[0].binding_name\" must be specified when\n\"saml2_config.slo_endpoints\" is specified"),
				},
				{
					Config:      ResourceSaml2CorporateIdPWithSloEndpoints("testIdP", "Test SAML2 IDP", "binding_name=\"binding\""),
					ExpectError: regexp.MustCompile("Attribute \"saml2_config.slo_endpoints\\[0].location\" must be specified when\n\"saml2_config.slo_endpoints\" is specified"),
				},
			},
		})
	})

	t.Run("error path - saml2_config.slo_endpoints.binding_name needs to be a valid value", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config
		saml2IdP.Saml2Configuration.SloEndpoints[0].BindingName = "invalid-binding"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", saml2IdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute saml2_config.slo_endpoints\\[0].binding_name value must be one of:\n\\[\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\"\n\"urn:oasis:names:tc:SAML:2.0:bindings:SOAP\"\n\"urn:oasis:names:tc:SAML:2.0:bindings:URI\"], got: \"%s\"", saml2IdP.Saml2Configuration.SloEndpoints[0].BindingName)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.slo_endpoints.location is not a valid URL", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config
		saml2IdP.Saml2Configuration.SloEndpoints[0].Location = "invalid-url"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", saml2IdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute saml2_config.slo_endpoints\\[0].location value must be a valid URL,\ngot: %s", saml2IdP.Saml2Configuration.SloEndpoints[0].Location)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.slo_endpoints.response_location is not a valid URL", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config
		saml2IdP.Saml2Configuration.SloEndpoints[0].ResponseLocation = "invalid-url"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", saml2IdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute saml2_config.slo_endpoints\\[0].response_location value must be a\nvalid URL, got: %s", saml2IdP.Saml2Configuration.SloEndpoints[0].ResponseLocation)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.digest_algorithm needs to be a valid value", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config
		saml2IdP.Saml2Configuration.DigestAlgorithm = "invalid-algorithm"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", saml2IdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute saml2_config.digest_algorithm value must be one of: \\[\"sha1\"\n\"sha256\" \"sha512\"], got: \"%s\"", saml2IdP.Saml2Configuration.DigestAlgorithm)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.name_id_format needs to be a valid value", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config
		saml2IdP.Saml2Configuration.DefaultNameIdFormat = "invalid-name-id-format"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", saml2IdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute saml2_config.name_id_format value must be one of: \\[\"default\" \"none\"\n\"unspecified\" \"email\"\\], got: \"%s\"", saml2IdP.Saml2Configuration.DefaultNameIdFormat)),
				},
			},
		})
	})

	t.Run("error path - saml2_config.allow_create needs to be a valid value", func(t *testing.T) {

		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = &saml2Config
		saml2IdP.Saml2Configuration.AllowCreate = "invalid-value"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", saml2IdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute saml2_config.allow_create value must be one of: \\[\"default\" \"none\"\n\"true\" \"false\"\\], got: \"%s\"", saml2IdP.Saml2Configuration.AllowCreate)),
				},
			},
		})
	})

	t.Run("error path - oidc_config requires root attributes name & type", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdPWithoutNameType("testIdP", "Test OIDC IDP", "oidc_config = {}"),
					ExpectError: regexp.MustCompile("Attribute \"name\" must be specified when \"oidc_config\" is specified"),
				},
				{
					Config:      ResourceCorporateIdPWithoutNameType("testIdP", "Test OIDC IDP", "oidc_config = {}"),
					ExpectError: regexp.MustCompile("Attribute \"type\" must be specified when \"oidc_config\" is specified"),
				},
			},
		})
	})

	t.Run("error path - type needs to be a valid value when oidc is configured", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdPConfigTypeMismtach("testIdP", "Test IDP", "saml2", "oidc_config = { discovery_url = \"https://test.com\"\nclient_id = \"test-client-id\" }"),
					ExpectError: regexp.MustCompile("Attribute oidc_config : value of attribute \"type\" must be modified to match\nthe IDP configuration provided. Acceptable values are : `openIdConnect`, got:\nsaml2"),
				},
			},
		})
	})

	t.Run("error path - oidc_config requires sub-attribute discovery_url and client_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceOidcCorporateIdP("testIdP", "Test IdP", "Test OIDC IDP", "openIdConnect", "client_id=\"test-client-id\""),
					ExpectError: regexp.MustCompile("Attribute \"oidc_config.discovery_url\" must be specified when \"oidc_config\" is\nspecified"),
				},
				{
					Config:      ResourceOidcCorporateIdP("testIdP", "Test IdP", "Test OIDC IDP", "openIdConnect", "discovery_url=\"https://test.com\""),
					ExpectError: regexp.MustCompile("Attribute \"oidc_config.client_id\" must be specified when \"oidc_config\" is\nspecified"),
				},
			},
		})

	})

	t.Run("error path - oidc_config.discovery_url is not a valid URL", func(t *testing.T) {

		oidcIdP := corporateIdP
		oidcIdP.DisplayName = "OIDC - Test Corporate IdP"
		oidcIdP.Name = "OIDC - Test IdP"
		oidcIdP.Type = "openIdConnect"
		oidcIdP.OidcConfiguration = &oidcCoporateIdP
		oidcIdP.OidcConfiguration.DiscoveryUrl = "invalid-url"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", oidcIdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute oidc_config.discovery_url value must be a valid URL, got:\n%s", oidcIdP.OidcConfiguration.DiscoveryUrl)),
				},
			},
		})
	})

	t.Run("error path - oidc_config.subject_name_identifier needs to be a valid value", func(t *testing.T) {

		oidcIdP := corporateIdP
		oidcIdP.DisplayName = "OIDC - Test Corporate IdP"
		oidcIdP.Name = "OIDC - Test IdP"
		oidcIdP.Type = "openIdConnect"
		oidcIdP.OidcConfiguration = &oidcCoporateIdP
		oidcIdP.OidcConfiguration.SubjectNameIdentifier = "invalid-value"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", oidcIdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute oidc_config.subject_name_identifier value must be one of: \\[\"none\"\n\"email\"], got: \"%s\"", oidcIdP.OidcConfiguration.SubjectNameIdentifier)),
				},
			},
		})
	})

	t.Run("error path - oidc_config.token_endpoint_auth_method needs to be a valid value", func(t *testing.T) {

		oidcIdP := corporateIdP
		oidcIdP.DisplayName = "OIDC - Test Corporate IdP"
		oidcIdP.Name = "OIDC - Test IdP"
		oidcIdP.Type = "openIdConnect"
		oidcIdP.OidcConfiguration = &oidcCoporateIdP
		oidcIdP.OidcConfiguration.TokenEndpointAuthMethod = "invalid-value"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", oidcIdP),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute oidc_config.token_endpoint_auth_method value must be one of:\n\\[\"clientSecretPost\" \"clientSecretBasic\" \"privateKeyJwt\"\n\"privateKeyJwtRfc7523\"], got: \"%s\"", oidcIdP.OidcConfiguration.TokenEndpointAuthMethod)),
				},
			},
		})
	})

	t.Run("error path - oidc_config.client_secret needs to be specified depending on the value of token_endpoint_auth_method", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					// when token_endpoint_auth_method is not specified, it defaults to "clientSecretPost" for which client_secret is required
					Config:      ResourceOidcCorporateIdP("testIdP", "Test IdP", "Test OIDC Idp", "openIdConnect", "discovery_url = \"https://test.com\"\nclient_id = \"test-client-id\""),
					ExpectError: regexp.MustCompile("Attribute oidc_config.client_secret attribute \"oidc_config.client_secret\"\nmust be specified, got: <null>"),
				},
				{
					// when token_endpoint_auth_method is set to "clientSecretPost" or "clientSecretBasic", client_secret is required
					Config:      ResourceOidcCorporateIdP("testIdP", "Test IdP", "Test OIDC Idp", "openIdConnect", "discovery_url = \"https://test.com\"\nclient_id = \"test-client-id\"\ntoken_endpoint_auth_method = \"clientSecretBasic\""),
					ExpectError: regexp.MustCompile("Attribute oidc_config.client_secret attribute \"oidc_config.client_secret\"\nmust be specified when oidc_config.token_endpoint_auth_method is set to one\nof: clientSecretPost, clientSecretBasic, got: <null>"),
				},
			},
		})
	})

	t.Run("error path - oidc_config.scopes needs to be configured with default values", func(t *testing.T) {

		oidcIdP := corporateIdP
		oidcIdP.DisplayName = "OIDC - Test Corporate IdP"
		oidcIdP.Name = "OIDC - Test IdP"
		oidcIdP.Type = "openIdConnect"
		oidcIdP.OidcConfiguration = &oidcCoporateIdP
		oidcIdP.OidcConfiguration.Scopes = []string{
			"test-value-1",
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceCorporateIdP("testIdP", oidcIdP),
					ExpectError: regexp.MustCompile("Please add the values :\n\\[openid]"),
				},
			},
		})
	})
}

func ResourceCorporateIdP(resourceName string, idp corporateidps.IdentityProvider) string {

	var groups string
	for _, group := range idp.IdentityFederation.RequiredGroups {
		groups += fmt.Sprintf(`"%s",`, group)
	}

	resourceIdP := fmt.Sprintf(`
	resource "sci_corporate_idp" "%s" {
		display_name = "%s"
		name = "%s"
		logout_url = "%s"
		forward_all_sso_requests = %t
		identity_federation = {
			use_local_user_store = %t
			allow_local_users_only = %t
			apply_local_idp_auth_and_checks = %t
			required_groups = [%s]
		}
		login_hint_config = {
			login_hint_type = "%s"
			send_method = "%s"
		}
	`, resourceName, idp.DisplayName, idp.Name, idp.LogoutUrl, idp.ForwardAllSsoRequests, idp.IdentityFederation.UseLocalUserStore, idp.IdentityFederation.AllowLocalUsersOnly, idp.IdentityFederation.ApplyLocalIdPAuthnChecks, groups, idp.LoginHintConfiguration.LoginHintType, idp.LoginHintConfiguration.SendMethod)

	switch idp.Type {
	case "openIdConnect":
		oidcConfig := idp.OidcConfiguration

		var scopes string
		for _, scope := range oidcConfig.Scopes {
			scopes += fmt.Sprintf(`
					"%s",
				`, scope)
		}

		additionalConfig := fmt.Sprintf(`
				additional_config = {
					enforce_nonce = %t
					enforce_issuer_check = %t
					disable_logout_id_token_hint = %t
				}
			`, oidcConfig.AdditionalConfig.EnforceNonce, oidcConfig.AdditionalConfig.EnforceIssuerCheck, oidcConfig.AdditionalConfig.OmitIDTokenHintForLogout)

		resourceIdP += fmt.Sprintf(`
				type = "%s"
				oidc_config = {
					discovery_url = "%s"
					client_id = "%s"
					client_secret = "%s"
					subject_name_identifier = "%s"
					token_endpoint_auth_method = "%s"
					scopes = [%s]
					enable_pkce = %t
					%s
				}
			`, idp.Type, oidcConfig.DiscoveryUrl, oidcConfig.ClientId, oidcConfig.ClientSecret, oidcConfig.SubjectNameIdentifier, oidcConfig.TokenEndpointAuthMethod, scopes, oidcConfig.PkceEnabled, additionalConfig)

	case "saml2", "sapSSO", "microsoftADFS":
		saml2Config := idp.Saml2Configuration

		var assertionAttributes string
		for _, attribute := range saml2Config.AssertionAttributes {
			assertionAttributes += fmt.Sprintf(`
					{
						name = "%s"
						value = "%s"
					},
				`, attribute.Name, attribute.Value)
		}

		var certificates string
		for _, certificate := range saml2Config.CertificatesForSigning {
			certificates += fmt.Sprintf(`
					{
						base64_certificate = "%s"
						dn = "%s"
						default = %t
						valid_from = "%s"
						valid_to = "%s"
					}
				`, certificate.Base64Certificate, certificate.Dn, certificate.IsDefault, certificate.ValidFrom, certificate.ValidTo)
		}

		var ssoEndpoints string
		for _, endpoint := range saml2Config.SsoEndpoints {
			ssoEndpoints += fmt.Sprintf(`
					{
						binding_name = "%s"
						location = "%s"
						default = %t
					}
				`, endpoint.BindingName, endpoint.Location, endpoint.IsDefault)
		}

		var sloEndpoints string
		for _, endpoint := range saml2Config.SloEndpoints {
			sloEndpoints += fmt.Sprintf(`
					{
						binding_name = "%s"
						location = "%s"
						response_location = "%s"
						default = %t
					}
				`, endpoint.BindingName, endpoint.Location, endpoint.ResponseLocation, endpoint.IsDefault)
		}

		resourceIdP += fmt.Sprintf(`
				type = "%s"
				saml2_config = {
					saml_metadata_url = "%s"
					assertion_attributes = [%s]
					digest_algorithm = "%s"
					include_scoping = %t
					name_id_format = "%s"
					allow_create = "%s"
					signing_certificates = [%s]
					sso_endpoints = [%s]
					slo_endpoints = [%s]
				}
			`, idp.Type, saml2Config.SamlMetadataUrl, assertionAttributes, saml2Config.DigestAlgorithm, saml2Config.IncludeScoping, saml2Config.DefaultNameIdFormat, saml2Config.AllowCreate, certificates, ssoEndpoints, sloEndpoints)
	}

	resourceIdP += `}`

	return resourceIdP
}

func ResourceCorporateIdPWithoutDisplayName(resourceName string) string {
	return fmt.Sprintf(`
	resource "sci_corporate_idp" "%s" {}
	`, resourceName)
}

func ResourceCorporateIdPWithEmptyLoginHintConfig(resourceName string, idpName string) string {
	return fmt.Sprintf(`
	resource "sci_corporate_idp" "%s" {
		display_name = "%s"
		login_hint_config = {}
	}
	`, resourceName, idpName)
}

func ResourceCorporateIdPWithoutNameType(resourceName string, idpName string, config string) string {
	return fmt.Sprintf(`
	resource "sci_corporate_idp" "%s" {
		display_name = "%s"
		%s
	}	
	`, resourceName, idpName, config)
}

func ResourceCorporateIdPConfigTypeMismtach(resourceName string, idpName string, idpType string, idpConfig string) string {
	return fmt.Sprintf(`
	resource "sci_corporate_idp" "%s" {
		display_name = "%s"
		name = "%s"
		type = "%s"
		%s
	}
	`, resourceName, idpName, idpName, idpType, idpConfig)
}

func ResourceSaml2CorporateIdPWithAssertionAttributes(resourceName string, idpName string, saml2Attribute string) string {
	return fmt.Sprintf(`
	resource "sci_corporate_idp" "%s" {
		display_name = "%s"
		name = "%s"
		type = "saml2"
		saml2_config = {
			assertion_attributes = [
				{
					%s
				}
			]
		}
	}
	`, resourceName, idpName, idpName, saml2Attribute)
}

func ResourceSaml2CorporateIdPWithSsoEndpoints(resourceName string, idpName string, ssoEndpoints string) string {
	return fmt.Sprintf(`
	resource "sci_corporate_idp" "%s" {
		display_name = "%s"
		name = "%s"
		type = "saml2"
		saml2_config = {
			sso_endpoints = [
				{
					%s
				}
			]
		}
	}
	`, resourceName, idpName, idpName, ssoEndpoints)
}

func ResourceSaml2CorporateIdPWithSloEndpoints(resourceName string, idpName string, sloEndpoints string) string {
	return fmt.Sprintf(`
	resource "sci_corporate_idp" "%s" {
		display_name = "%s"
		name = "%s"
		type = "saml2"
		saml2_config = {
			slo_endpoints = [
				{
					%s
				}
			]
		}
	}
	`, resourceName, idpName, idpName, sloEndpoints)
}

func ResourceOidcCorporateIdP(resourceName string, idpDisplayName string, idpName string, idpType string, oidcConfig string) string {
	return fmt.Sprintf(`
		resource "sci_corporate_idp" "%s" {
			display_name = "%s"
			type = "%s"
			name = "%s"
			oidc_config = {
				%s
			}
		}
	`, resourceName, idpDisplayName, idpType, idpName, oidcConfig)
}
