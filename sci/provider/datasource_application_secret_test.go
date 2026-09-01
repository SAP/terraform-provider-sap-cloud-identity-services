package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceClientSecret(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_application_secret")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceClientSecret("testSecret", "ca60ce62-efcf-4a04-8ce7-86405c7d1926", "88e88cc8-e6a8-4b19-be68-a6ea3f8c3e5e"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.sci_application_secret.testSecret", "id", "88e88cc8-e6a8-4b19-be68-a6ea3f8c3e5e"),
						resource.TestCheckResourceAttr("data.sci_application_secret.testSecret", "application_id", "ca60ce62-efcf-4a04-8ce7-86405c7d1926"),
						resource.TestCheckResourceAttr("data.sci_application_secret.testSecret", "hint", "7SA:"),
						resource.TestCheckResourceAttr("data.sci_application_secret.testSecret", "description", "terraform-test"),
					),
				},
			},
		})
	})

	t.Run("error path - id is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceClientSecretMissingId("testSecret", "64bdf294-d1f3-47ac-8a4e-b1bb88380ffc"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - application_id is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceClientSecretMissingAppId("testSecret", "996cee6e-ca22-4b9e-9e4d-8cad11ad3e10"),
					ExpectError: regexp.MustCompile(`The argument "application_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func DataSourceClientSecret(datasourceName, appId, secretId string) string {
	return fmt.Sprintf(`
data "sci_application_secret" "%s" {
  application_id = "%s"
  id             = "%s"
}`, datasourceName, appId, secretId)
}

func DataSourceClientSecretMissingId(datasourceName, appId string) string {
	return fmt.Sprintf(`
data "sci_application_secret" "%s" {
  application_id = "%s"
}`, datasourceName, appId)
}

func DataSourceClientSecretMissingAppId(datasourceName, secretId string) string {
	return fmt.Sprintf(`
data "sci_application_secret" "%s" {
  id = "%s"
}`, datasourceName, secretId)
}
