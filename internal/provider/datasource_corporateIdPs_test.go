package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceCorporateIdPs(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/datasource_corporate_idps")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceCorporateIdPs("testIdPs"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.sci_corporate_idps.testIdPs", "values.#", "5"),
					),
				},
			},
		})

	})

}

func DataSourceCorporateIdPs(datasourceName string) string {
	return fmt.Sprintf(`
		data "sci_corporate_idps" "%s" {}
	`, datasourceName)
}
