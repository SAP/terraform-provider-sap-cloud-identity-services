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
					ExpectError: regexp.MustCompile("Attribute oidc_config : value of attribute \"type\" must be modified to match\nthe IDP configuration provided. Acceptable values are : `openIdConnect`, got:\n\"saml2\""),
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

	if idp.Type == "openIdConnect" {

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

	} else if idp.Type == "sapSSO" || idp.Type == "microsoftADFS" || idp.Type == "saml2" {
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
		type = "%s"
		%s
	}
	`, resourceName, idpName, idpType, idpConfig)
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
