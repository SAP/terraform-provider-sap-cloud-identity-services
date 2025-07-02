package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceCorporateIdP(t *testing.T) {

	t.Parallel()

	t.Run("happy path - saml2 corporate idp", func(t *testing.T) {

		testEndpoint := "https://test.com"

		rec, user := setupVCR(t, "fixtures/datasource_saml2_corporate_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceCorporateIdP("testIdP", "Terraform - SAML2"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_corporate_idp.testIdP", "id", regexpUUID),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "display_name", "Terraform - SAML2"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "name", "Saml2 - Test"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "type", "saml2"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "logout_url", testEndpoint),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "forward_all_sso_requests", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "identity_federation.use_local_user_store", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "identity_federation.allow_local_users_only", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "identity_federation.apply_local_idp_auth_and_checks", fmt.Sprintf("%t", false)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "identity_federation.required_groups.0", "Test Group"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "login_hint_config.login_hint_type", "userInput"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "login_hint_config.send_method", "authRequest"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.allow_create", "none"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.assertion_attributes.0.name", "test_attr"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.assertion_attributes.0.value", "test_value"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.digest_algorithm", "sha256"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.include_scoping", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.name_id_format", "default"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.signing_certificates.#", "1"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.slo_endpoints.0.binding_name", "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.slo_endpoints.0.default", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.slo_endpoints.0.location", testEndpoint),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.sso_endpoints.0.binding_name", "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.sso_endpoints.0.default", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "saml2_config.sso_endpoints.0.location", testEndpoint),
					),
				},
			},
		})
	})

	t.Run("happy path - oidc corporate idp", func(t *testing.T) {

		issuer := "https://token.actions.githubusercontent.com"
		testEndpoint := "https://test.com"

		rec, user := setupVCR(t, "fixtures/datasource_oidc_corporate_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceCorporateIdP("testIdP", "Terraform - OIDC"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_corporate_idp.testIdP", "id", regexpUUID),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "display_name", "Terraform - OIDC"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "name", "Oidc - Test"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "type", "openIdConnect"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "logout_url", testEndpoint),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "forward_all_sso_requests", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "identity_federation.use_local_user_store", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "identity_federation.allow_local_users_only", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "identity_federation.apply_local_idp_auth_and_checks", fmt.Sprintf("%t", false)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "identity_federation.required_groups.0", "Test Group"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "login_hint_config.login_hint_type", "mail"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "login_hint_config.send_method", "urlParam"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "oidc_config.client_id", "Oidc - Test"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "oidc_config.discovery_url", issuer),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "oidc_config.enable_pkce", fmt.Sprintf("%t", false)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "oidc_config.is_client_secret_configured", fmt.Sprintf("%t", false)),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "oidc_config.issuer", issuer),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "oidc_config.jwks_uri", issuer+"/.well-known/jwks"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "oidc_config.scopes.0", "openid"),
						resource.TestCheckResourceAttr("data.sci_corporate_idp.testIdP", "oidc_config.subject_name_identifier", "email"),
					),
				},
			},
		})
	})

	t.Run("error path - id is a mandatory attribute", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceCorporateIdPWithoutId("testIdP"),
					ExpectError: regexp.MustCompile("The argument \"id\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - id must be a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceCorporateIdPWithId("testIdP", "invalid-uuid"),
					ExpectError: regexp.MustCompile("Attribute id value must be a valid UUID, got: invalid-uuid"),
				},
			},
		})
	})
}

func DataSourceCorporateIdP(datasourceName string, idpName string) string {
	return fmt.Sprintf(`
	data "sci_corporate_idps" "allIdPs" {}
	data "sci_corporate_idp" "%s" {
		id = [for idp in data.sci_corporate_idps.allIdPs.values : idp.id if idp.display_name == "%s"][0]
	}
	`, datasourceName, idpName)
}

func DataSourceCorporateIdPWithoutId(datasourceName string) string {
	return fmt.Sprintf(`
		data "sci_corporate_idp" "%s" {
		}
	`, datasourceName)
}

func DataSourceCorporateIdPWithId(datasourceName string, idpId string) string {
	return fmt.Sprintf(`
		data "sci_corporate_idp" "%s" {
			id = "%s"
		}
	`, datasourceName, idpId)
}
