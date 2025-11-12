package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceApplication(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/datasource_application")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceApplication("testApp", "oac.accounts.sap.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "name", "oac.accounts.sap.com"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", "firstName"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.default_authenticating_idp", "664c660e25cff252c5c202dc"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.sso_type", "saml2oidc"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.subject_name_identifier.value", "uid"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.oidc_config.token_policy.jwt_validity", "3600"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_validity", "43200"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_parallel", "1"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.oidc_config.token_policy.max_exchange_period", "unlimited"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.oidc_config.token_policy.refresh_token_rotation_scenario", "off"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.oidc_config.token_policy.access_token_format", "default"),
						resource.TestCheckTypeSetElemAttr("data.sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "jwtBearer"),
						resource.TestCheckTypeSetElemAttr("data.sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "clientCredentials"),
						resource.TestCheckTypeSetElemAttr("data.sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "authorizationCode"),
						resource.TestCheckTypeSetElemAttr("data.sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "password"),
						resource.TestCheckTypeSetElemAttr("data.sci_application.testApp", "authentication_schema.oidc_config.restricted_grant_types.*", "refreshToken"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.location", "https://iasprovidertestblr.accounts400.ondemand.com/saml2/sp/acs/oac.accounts.sap.com"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.index", fmt.Sprintf("%d", 0)),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.acs_endpoints.0.default", fmt.Sprintf("%t", false)),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.binding_name", "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.0.location", "https://iasprovidertestblr.accounts400.ondemand.com/saml2/sp/slo/oac.accounts.sap.com"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.1.binding_name", "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.slo_endpoints.1.location", "https://iasprovidertestblr.accounts400.ondemand.com/saml2/sp/slo/oac.accounts.sap.com"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.signing_certificates.0.dn", "CN=iasprovidertestblr.accounts400.ondemand.com,O=SAP-SE,C=DE"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.signing_certificates.0.default", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.encryption_certificate.dn", "CN=iasprovidertestblr.accounts400.ondemand.com,O=SAP-SE,C=DE"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.response_elements_to_encrypt", "none"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.default_name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.sign_slo_messages", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.require_signed_slo_messages", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.require_signed_auth_requests", fmt.Sprintf("%t", false)),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.sign_assertions", fmt.Sprintf("%t", true)),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.sign_auth_responses", fmt.Sprintf("%t", false)),
						resource.TestCheckResourceAttr("data.sci_application.testApp", "authentication_schema.saml2_config.digest_algorithm", "sha256"),
					),
				},
			},
		})
	})

	t.Run("happy path - bundled application1", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/datasource_bundled_application1")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceApplication("testBundledApp", "XSUAA_b75a605d-151c-4485-83f4-64604378e4ec"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_application.testBundledApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("data.sci_application.testBundledApp", "name", "XSUAA_b75a605d-151c-4485-83f4-64604378e4ec"),
						resource.TestCheckResourceAttr("data.sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.type", "xsuaa"),
						resource.TestCheckResourceAttr("data.sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.app_tenant_id", "b75a605d-151c-4485-83f4-64604378e4ec"),
					),
				},
			},
		})
	})

	t.Run("happy path - bundled application2", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/datasource_bundled_application2")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceApplication("testBundledApp", "identity-subscription-c6c390f4-c9a2-4a6c-9cc7-01675a31e4f6-in-b75a605d-151c-4485-83f4-64604378e4ec"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_application.testBundledApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("data.sci_application.testBundledApp", "name", "identity-subscription-c6c390f4-c9a2-4a6c-9cc7-01675a31e4f6-in-b75a605d-151c-4485-83f4-64604378e4ec"),
						resource.TestCheckResourceAttr("data.sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.type", "subscription"),
						resource.TestCheckResourceAttr("data.sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.app_tenant_id", "b75a605d-151c-4485-83f4-64604378e4ec"),
						resource.TestCheckResourceAttr("data.sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.source_app_id", "3cc4b385-fe8b-423a-a8c0-34e15c9970cd"),
						resource.TestCheckResourceAttr("data.sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.source_tenant_id", "sapdas"),
						resource.TestCheckResourceAttr("data.sci_application.testBundledApp", "authentication_schema.sap_managed_attributes.service_instance_id", "c6c390f4-c9a2-4a6c-9cc7-01675a31e4f6"),
					),
				},
			},
		})
	})

	t.Run("error path - invalid app id", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceApplicationById("testApp", "invalid-uuid"),
					ExpectError: regexp.MustCompile(`Attribute id value must be a valid UUID, got: invalid-uuid`),
				},
			},
		})
	})

	t.Run("error path - app id is mandatory", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceApplicationNoId("testApp"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		})

	})

}

func DataSourceApplication(datasourceName string, appName string) string {
	return fmt.Sprintf(`
	data "sci_applications" "allApps" {}
	data "sci_application" "%s" {
		id = [for app in data.sci_applications.allApps.values : app.id if app.name == "%s"][0]
	}
	`, datasourceName, appName)
}

func DataSourceApplicationById(datasourceName string, appId string) string {
	return fmt.Sprintf(`
	data "sci_application" "%s" {
		id = "%s"
	}
	`, datasourceName, appId)
}

func DataSourceApplicationNoId(datasourceName string) string {
	return fmt.Sprintf(`
	data "sci_application" "%s" {
	}
	`, datasourceName)
}
