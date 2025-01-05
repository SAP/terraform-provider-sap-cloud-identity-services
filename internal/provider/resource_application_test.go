package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var regexpUUID = UuidRegexp

func TestResourceApplication(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_application")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceApplication("testApp", "basic-test-app", "application for testing purposes"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_application.testApp", "name", "basic-test-app"),
						resource.TestCheckResourceAttr("ias_application.testApp", "description", "application for testing purposes"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "sso_type", "saml2"),
						resource.TestCheckResourceAttr("ias_application.testApp", "subject_name_identifier.value", "uid"),
					),
				},
				{
					ResourceName:      "ias_application.testApp",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - application update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_application_updated")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/", user) + ResourceApplication("testApp", "test-app", "application for testing purposes"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_application.testApp", "name", "test-app"),
						resource.TestCheckResourceAttr("ias_application.testApp", "description", "application for testing purposes"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "sso_type", "saml2"),
						resource.TestCheckResourceAttr("ias_application.testApp", "subject_name_identifier.value", "uid"),
					),
				},
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/", user) + ResourceApplication("testApp", "test-app-updated", "application for testing purposes"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_application.testApp", "name", "test-app-updated"),
						resource.TestCheckResourceAttr("ias_application.testApp", "description", "application for testing purposes"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "sso_type", "saml2"),
						resource.TestCheckResourceAttr("ias_application.testApp", "subject_name_identifier.value", "uid"),
					),
				},
			},
		})
	})

	t.Run("error path - app_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAppId("testApp", "this-is-not-uuid", "test-app", "application for testing purposes"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute id value must be a valid UUID, got: %s","this-is-not-uuid")),
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
					Config:      ResourceApplicationWithSsoType("testApp", "test-app", "application for testing purposes", "this-is-not-a-valid-sso_type"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute sso_type value must be one of: [\"openIdConnect\" \"saml2\"], got:\n \"%s\"","this-is-not-a-valid-sso_type")),
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
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "this-is-not-email-domain", "0.0.0.0"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_rules[0].user_email_domain value must be a valid\nEmail Domain, got: %s", "this-is-not-email-domain")),
				},
			},
		})
	})

	t.Run("error path - ip is not a valid address", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "test.com", "this-is-not-ip-address"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_rules[0].ip_network_range value must be a valid IP\nAddress, got: %s", "this-is-not-ip-address")),
				},
			},
		})
	})
}

func ResourceApplication(resourceName string, appName string, description string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
	}
	`, resourceName, appName, description)
}

func ResourceApplicationWithParent(resourceName string, appName string, description string, parentAppId string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		parent_application_id = "%s"
	}
	`, resourceName, appName, description, parentAppId)
}

func ResourceApplicationWithAppId(resourceName string, appID string, appName string, description string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		id = "%s"
		name = "%s"
		description = "%s"
	}
	`, resourceName, appID, appName, description)
}

func ResourceApplicationWithoutAppName(resourceName string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
	}
	`, resourceName)
}

func ResourceApplicationWithSsoType(resourceName string, appName string, description string, ssoType string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		sso_type = "%s"
	}
	`, resourceName, appName, description, ssoType)
}

func ResourceApplicationWithAuthenticationRules(resourceName string, appName string, description string, emailDomain string, ipAddress string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_rules = [
			{
				identity_provider_id = "664c660e25cff252c5c202dc",
				user_email_domain = "%s"
				ip_network_range = "%s"
			}
		]
	}
	`, resourceName, appName, description, emailDomain, ipAddress)
}
