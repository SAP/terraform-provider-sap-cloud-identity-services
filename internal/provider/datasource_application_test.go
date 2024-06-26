package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceApplication(t *testing.T){

	t.Parallel()

	t.Run("happy path", func (t *testing.T){

		rec, _ := setupVCR(t, "fixtures/datasource_application")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/") + DataSourceApplication("testApp", "oac.accounts.sap.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.ias_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("data.ias_application.testApp", "description", ""),
						resource.TestCheckResourceAttr("data.ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("data.ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("data.ias_application.testApp", "name", "oac.accounts.sap.com"),
						resource.TestCheckResourceAttr("data.ias_application.testApp", "parent_application_id", ""),
					),
				},
			},
		})
	})

	t.Run("error path - invalid app id", func(t *testing.T){
		
		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/") + DataSourceApplicationById("testApp", "invalid-uuid"),
					ExpectError: regexp.MustCompile(`Attribute id value must be a valid UUID, got: invalid-uuid`),
				},
			},
		})
	})

	t.Run("error path - app id is mandatory", func(t *testing.T){
		
		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/") + DataSourceApplicationNoId("testApp"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		})
	
	})

}

func DataSourceApplication (datasourceName string, appName string) string {
	return fmt.Sprintf(`
	data "ias_applications" "allApps" {}
	data "ias_application" "%s" {
		id = [for app in data.ias_applications.allApps.values : app.id if app.name == "%s"][0]
	}
	`, datasourceName, appName)
}

func DataSourceApplicationById (datasourceName string, appId string) string {
	return fmt.Sprintf(`
	data "ias_application" "%s" {
		id = "%s"
	}
	`, datasourceName, appId)
}

func DataSourceApplicationNoId (datasourceName string) string {
	return fmt.Sprintf(`
	data "ias_application" "%s" {
	}
	`, datasourceName)
}