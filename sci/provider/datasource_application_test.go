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
