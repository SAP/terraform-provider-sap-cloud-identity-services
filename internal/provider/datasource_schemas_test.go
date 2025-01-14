package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)


func TestDataSourceSchemas(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func (t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_schemas_all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceSchemas("allSchemas"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.ias_schemas.allSchemas", "values.#", "13"),
					),
				},
			},
		})

	})

}

func DataSourceSchemas(datasourceName string) string {
	return fmt.Sprintf(`
	data "ias_schemas" "%s"{

	}
	`, datasourceName)
}