package provider

import (
	"fmt"
	"regexp"

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
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceSchema("testSchema", "Schema"),
					Check: resource.ComposeAggregateTestCheckFunc(
						//add a regex check
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "description", "Specifies the schema that describes a SCIM schema"),
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "name", "Schema"),
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "attributes.0.name", "name"),
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "attributes.0.description", "The schemas human-readable name. When applicable, service providers MUST specify the name, e.g., User."),
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "attributes.0.multivalued", "false"),
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "attributes.0.mutability", "readOnly"),
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "attributes.0.required", "true"),
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "attributes.0.returned", "default"),
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "attributes.0.type", "string"),
						resource.TestCheckResourceAttr("data.sci_schema.testSchema", "attributes.0.uniqueness", "none"),
					),
				},
			},
		})

	})

	// if regex is added for schema id, add a test

	t.Run("error path - schema id is mandatory", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceSchemaNoId("testSchema"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		})

	})

}

func DataSourceSchema(datasourceName string, schemaName string) string {
	return fmt.Sprintf(`
	data "sci_schemas" "allSchemas" {}
	data "sci_schema" "%s" {
		id = [for schema in data.sci_schemas.allSchemas.values : schema.id if schema.name == "%s"][0]
	}
	`, datasourceName, schemaName)
}

// func DataSourceSchemaById (datasourceName string, schemaId string) string {
// 	return fmt.Sprintf(`
// 	data "sci_schema" "%s" {
// 		id = "%s"
// 	}
// 	`, datasourceName, schemaId)
// }

func DataSourceSchemaNoId(datasourceName string) string {
	return fmt.Sprintf(`
	data "sci_schema" "%s" {
	}
	`, datasourceName)
}
