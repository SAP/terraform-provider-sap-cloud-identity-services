package provider

import (
	"fmt"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/schemas"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSchema(t *testing.T) {

	t.Parallel()

	schema := schemas.Schema{
		Id:   "urn:ietf:scim:schemas:Terraform",
		Name: "Terraform",
		Attributes: []schemas.Attribute{
			{
				Name:       "test_attribute",
				Type:       "string",
				Mutability: "readWrite",
				Returned:   "never",
				Uniqueness: "none",
				CanonicalValues: []string{
					"test",
					"attr",
				},
				Multivalued: true,
				Description: "Test Attribute",
				Required:    true,
				CaseExact:   false,
			},
		},
		Description: "Test Schema",
	}

	t.Run("happy path", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_schema")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceSchema("testSchema", schema),
					Check: resource.ComposeAggregateTestCheckFunc(
						// TODO add test check if regex is added
						resource.TestCheckResourceAttr("sci_schema.testSchema", "id", schema.Id),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "name", schema.Name),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.name", schema.Attributes[0].Name),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.type", schema.Attributes[0].Type),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.mutability", schema.Attributes[0].Mutability),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.returned", schema.Attributes[0].Returned),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.uniqueness", schema.Attributes[0].Uniqueness),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.canonical_values.0", schema.Attributes[0].CanonicalValues[0]),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.canonical_values.1", schema.Attributes[0].CanonicalValues[1]),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.multivalued", fmt.Sprintf("%t", schema.Attributes[0].Multivalued)),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.description", schema.Attributes[0].Description),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.required", fmt.Sprintf("%t", schema.Attributes[0].Required)),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.case_exact", fmt.Sprintf("%t", schema.Attributes[0].CaseExact)),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "description", schema.Description),
					),
				},
			},
		})
	})

	t.Run("error path - schemas cannot be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchemaWithoutSchemas("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform"),
					ExpectError: regexp.MustCompile("Attribute schemas set must contain at least 1 elements, got: 0"),
				},
			},
		})
	})

	t.Run("error path - schema ID is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchemaWithoutSchemaId("testSchema", "Terraform"),
					ExpectError: regexp.MustCompile("The argument \"id\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - schema name is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchemaWithoutSchemaName("testSchema", "urn:ietf:scim:schemas:Terraform"),
					ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - schema attributes has mandatory parameters : name, type, mutability, returned, uniqueness", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchemaWithoutAttributes("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform"),
					ExpectError: regexp.MustCompile("Inappropriate value for attribute \"attributes\": element 0: attributes\n\"mutability\", \"name\", \"returned\", \"type\", and \"uniqueness\" are required."),
				},
			},
		})
	})

	t.Run("error path - schema attributes.name must be a valid value", func(t *testing.T) {

		schema.Attributes = []schemas.Attribute{
			{
				Name:       "an-invalid-name",
				Type:       "string",
				Mutability: "readWrite",
				Returned:   "never",
				Uniqueness: "none",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchema("testSchema", schema),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].name value must be a valid name. Must start with an\nalphabet and should contain only alphanumeric characters and underscores,\ngot: %s", schema.Attributes[0].Name)),
				},
			},
		})

		schema.Attributes = []schemas.Attribute{
			{
				Name:       "@n_invalid_name",
				Type:       "string",
				Mutability: "readWrite",
				Returned:   "never",
				Uniqueness: "none",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchema("testSchema", schema),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].name value must be a valid name. Must start with an\nalphabet and should contain only alphanumeric characters and underscores,\ngot: %s", schema.Attributes[0].Name)),
				},
			},
		})

		schema.Attributes = []schemas.Attribute{
			{
				Name:       "1_an_invalid_name",
				Type:       "string",
				Mutability: "readWrite",
				Returned:   "never",
				Uniqueness: "none",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchema("testSchema", schema),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].name value must be a valid name. Must start with an\nalphabet and should contain only alphanumeric characters and underscores,\ngot: %s", schema.Attributes[0].Name)),
				},
			},
		})
	})

	t.Run("error path - schema attributes.type must be a valid value", func(t *testing.T) {

		schema.Attributes = []schemas.Attribute{
			{
				Name:       "test_attribute",
				Type:       "an_invalid_value",
				Mutability: "readWrite",
				Returned:   "never",
				Uniqueness: "none",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchema("testSchema", schema),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].type value must be one of: \\[\"string\" \"boolean\"\n\"decimal\" \"integer\" \"dateTime\" \"binary\" \"reference\" \"complex\"], got:\n\"%s\"", schema.Attributes[0].Type)),
				},
			},
		})
	})

	t.Run("error path - schema attributes.mutability must be a valid value", func(t *testing.T) {

		schema.Attributes = []schemas.Attribute{
			{
				Name:       "test_attribute",
				Type:       "integer",
				Mutability: "an_invalid_value",
				Returned:   "never",
				Uniqueness: "none",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchema("testSchema", schema),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].mutability value must be one of: \\[\"readOnly\"\n\"readWrite\" \"writeOnly\" \"immutable\"], got: \"%s\"", schema.Attributes[0].Mutability)),
				},
			},
		})
	})

	t.Run("error path - schema attributes.returned must be a valid value", func(t *testing.T) {

		schema.Attributes = []schemas.Attribute{
			{
				Name:       "test_attribute",
				Type:       "integer",
				Mutability: "readWrite",
				Returned:   "an_invalid_value",
				Uniqueness: "none",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchema("testSchema", schema),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].returned value must be one of: \\[\"always\" \"never\"\n\"default\" \"request\"], got: \"%s\"", schema.Attributes[0].Returned)),
				},
			},
		})
	})

	t.Run("error path - schema attributes.uniqueness must be a valid value", func(t *testing.T) {

		schema.Attributes = []schemas.Attribute{
			{
				Name:       "test_attribute",
				Type:       "integer",
				Mutability: "readWrite",
				Returned:   "never",
				Uniqueness: "an_invalid_value",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceSchema("testSchema", schema),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].uniqueness value must be one of: \\[\"none\" \"server\"\n\"global\"], got: \"%s\"", schema.Attributes[0].Uniqueness)),
				},
			},
		})
	})

}

func ResourceSchema(resourceName string, schema schemas.Schema) string {
	return fmt.Sprintf(`
	resource "sci_schema" "%s"{
		id = "%s"
		name = "%s"
		attributes = [%s]
		description = "%s"
	}
	`, resourceName, schema.Id, schema.Name, getSchemaAttributes(schema.Attributes), schema.Description)
}

func ResourceSchemaWithoutSchemas(resourceName string, schemaId string, schemaName string) string {
	return fmt.Sprintf(`
	resource "sci_schema" "%s"{
		id = "%s"
		name = "%s"
		attributes = []
		schemas = []
	}
	`, resourceName, schemaId, schemaName)
}

func ResourceSchemaWithoutSchemaId(resourceName string, schemaName string) string {
	return fmt.Sprintf(`
	resource "sci_schema" "%s"{
		name = "%s"
		attributes = []
	}
	`, resourceName, schemaName)
}

func ResourceSchemaWithoutSchemaName(resourceName string, schemaId string) string {
	return fmt.Sprintf(`
	resource "sci_schema" "%s"{
		id = "%s"
		attributes = []
	}
	`, resourceName, schemaId)
}

func ResourceSchemaWithoutAttributes(resourceName string, schemaId string, schemaName string) string {
	return fmt.Sprintf(`
	resource "sci_schema" "%s"{
		id = "%s"
		name = "%s"
		attributes = [
			{}
		]
	}
	`, resourceName, schemaId, schemaName)
}

func getSchemaAttributes(schemaAttributes []schemas.Attribute) string {
	attributes := ""
	for _, attr := range schemaAttributes {

		canonicalValues := ""
		for _, val := range attr.CanonicalValues {
			canonicalValues += fmt.Sprintf(`
				"%s",
			`, val)
		}

		attributes += fmt.Sprintf(`{
			name = "%s"
			mutability = "%s"
			returned = "%s"
			type = "%s"
			uniqueness = "%s"
			canonical_values = [%s]
			multivalued = %t
			description = "%s"
			required = %t
			case_exact = %t
		},`, attr.Name, attr.Mutability, attr.Returned, attr.Type, attr.Uniqueness,
			canonicalValues, attr.Multivalued, attr.Description, attr.Required, attr.CaseExact)
	}
	return attributes
}
