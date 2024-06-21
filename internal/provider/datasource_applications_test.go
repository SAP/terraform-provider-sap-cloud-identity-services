package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)


func TestDataSourceApplications(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func (t *testing.T) {
		rec, _ := setupVCR(t, "fixtures/datasource_applications_all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/") + DataSourceApplication("allApps"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.ias_applications.allApps", "values.#", "21"),
					),
				},
			},
		})

	})

}

func DataSourceApplication(datasourceName string) string {
	return fmt.Sprintf(`
	data "ias_applications" "%s"{

	}
	`, datasourceName)
}