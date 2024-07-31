package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceUsers(t *testing.T){

	t.Parallel()

	t.Run("happy path", func (t *testing.T){
		rec, user := setupVCR(t, "fixtures/datasource_users_all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceUsers("allUsers"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.ias_users.allUsers","values.#","4"),
					),
				},
			},
		})

	})

}


func DataSourceUsers(datasourceName string) string {
	return fmt.Sprintf(`
	data "ias_users" "%s"{

	}
	`, datasourceName)
}