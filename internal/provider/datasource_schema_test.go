package provider

import (
	"fmt"
	// "regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSchema(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/datasource_schema")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceSchema("testSchema", ""),
				},
			},
		})

	})

}

func DataSourceSchema (datasourceName string, schemaName string) string {
	return fmt.Sprintf(`
	data "ias_schemas" "allSchemas" {}
	data "ias_schema" "%s" {
		id = [for schema in data.ias_schemas.allSchemas.values : schema.id if schema.name == "%s"][0]
	}
	`, datasourceName,schemaName)
}