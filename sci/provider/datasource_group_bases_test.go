package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGroupBases(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_group_bases")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceGroupBases("allGroupBases"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.sci_group_bases.allGroupBases", "values.#", "10"),
					),
				},
			},
		})
	})
}

func DataSourceGroupBases(datasourceName string) string {
	return fmt.Sprintf(`
	data "sci_group_bases" "%s" {}
	`, datasourceName)
}
