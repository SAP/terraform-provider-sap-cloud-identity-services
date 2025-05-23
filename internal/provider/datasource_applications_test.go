package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceApplications(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_applications_all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceApplications("allApps"),
					Check: resource.ComposeAggregateTestCheckFunc(
						//fixture needs to re-recorded after cleanup, no. of apps needs to be modified accordingly
						resource.TestCheckResourceAttr("data.sci_applications.allApps", "values.#", "21"),
					),
				},
			},
		})

	})

}

func DataSourceApplications(datasourceName string) string {
	return fmt.Sprintf(`
	data "sci_applications" "%s"{

	}
	`, datasourceName)
}
