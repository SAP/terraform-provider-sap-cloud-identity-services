package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var regexpUUID = UuidRegexp

func TestResourceApplication (t *testing.T) { 
	
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_application")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/", user) + ResourceApplication("testApp", "basic-test-app", "application for testing purposes"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_application.testApp", "name", "basic-test-app"),
						resource.TestCheckResourceAttr("ias_application.testApp", "description", "application for testing purposes"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
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
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute parent_application_id value must be a valid UUID, got:\n%s","this-is-not-uuid")),
				},
			},
		})
	})
	
}

func ResourceApplication (resourceName string, appName string, description string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
	}
	`, resourceName, appName, description )
}

func ResourceApplicationWithParent (resourceName string, appName string, description string, parentAppId string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		parent_application_id = "%s"
	}
	`, resourceName, appName, description, parentAppId)
}

func ResourceApplicationWithAppId (resourceName string, appID string, appName string, description string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		id = "%s"
		name = "%s"
		description = "%s"
	}
	`, resourceName, appID, appName, description )
}

func ResourceApplicationWithoutAppName(resourceName string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
	}
	`, resourceName)
}
