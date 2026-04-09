package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceClientSecrets(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_client_secrets")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceClientSecrets("testSecrets", "5fd22812-53c5-4803-8285-94cd1fb3b301"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.sci_client_secrets.testSecrets", "values.#"),
						resource.TestCheckResourceAttr("data.sci_client_secrets.testSecrets", "application_id", "5fd22812-53c5-4803-8285-94cd1fb3b301"),
					),
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
					Config:      DataSourceClientSecretsMissingAppId("testSecrets"),
					ExpectError: regexp.MustCompile(`The argument "application_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func DataSourceClientSecrets(datasourceName, appId string) string {
	return fmt.Sprintf(`
data "sci_client_secrets" "%s" {
  application_id = "%s"
}`, datasourceName, appId)
}

func DataSourceClientSecretsMissingAppId(datasourceName string) string {
	return fmt.Sprintf(`
data "sci_client_secrets" "%s" {
}`, datasourceName)
}
