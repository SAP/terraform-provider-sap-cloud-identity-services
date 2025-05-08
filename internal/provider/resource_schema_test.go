package provider

import (
	"fmt"
	"regexp"
	"terraform-provider-sci/internal/cli/apiObjects/schemas"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSchema(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {

		schemaAttributes := schemas.Attribute{
			Name:       "test_attribute",
			Type:       "string",
			Mutability: "readWrite",
			Returned:   "never",
			Uniqueness: "none",
		}

		rec, user := setupVCR(t, "fixtures/resource_schema")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceSchema("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform", schemaAttributes),
					Check: resource.ComposeAggregateTestCheckFunc(
						//add if regex is added
						// resource.TestMatchResourceAttr("sci_user.testSchema", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "id", "urn:ietf:scim:schemas:Terraform"),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "name", "Terraform"),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.name", schemaAttributes.Name),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.type", schemaAttributes.Type),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.mutability", schemaAttributes.Mutability),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.returned", schemaAttributes.Returned),
						resource.TestCheckResourceAttr("sci_schema.testSchema", "attributes.0.uniqueness", schemaAttributes.Uniqueness),
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

		schemaAttributes := []schemas.Attribute{
			{
				Name:       "an-invalid-name",
				Type:       "string",
				Mutability: "readWrite",
				Returned:   "never",
				Uniqueness: "none",
			},
			{
				Name:       "@n_invalid_name",
				Type:       "string",
				Mutability: "readWrite",
				Returned:   "never",
				Uniqueness: "none",
			},
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
					Config:      ResourceSchema("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform", schemaAttributes[0]),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].name value must be a valid name. Must start with an\nalphabet and should contain only alphanumeric characters and underscores,\ngot: %s", schemaAttributes[0].Name)),
				},
				{
					Config:      ResourceSchema("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform", schemaAttributes[1]),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].name value must be a valid name. Must start with an\nalphabet and should contain only alphanumeric characters and underscores,\ngot: %s", schemaAttributes[1].Name)),
				},
				{
					Config:      ResourceSchema("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform", schemaAttributes[1]),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].name value must be a valid name. Must start with an\nalphabet and should contain only alphanumeric characters and underscores,\ngot: %s", schemaAttributes[1].Name)),
				},
			},
		})
	})

	t.Run("error path - schema attributes.type must be a valid value", func(t *testing.T) {

		schemaAttributes := []schemas.Attribute{
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
					Config:      ResourceSchema("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform", schemaAttributes[0]),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].type value must be one of: \\[\"string\" \"boolean\"\n\"decimal\" \"integer\" \"dateTime\" \"binary\" \"reference\" \"complex\"], got:\n\"%s\"", schemaAttributes[0].Type)),
				},
			},
		})
	})

	t.Run("error path - schema attributes.mutability must be a valid value", func(t *testing.T) {

		schemaAttributes := []schemas.Attribute{
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
					Config:      ResourceSchema("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform", schemaAttributes[0]),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].mutability value must be one of: \\[\"readOnly\"\n\"readWrite\" \"writeOnly\" \"immutable\"], got: \"%s\"", schemaAttributes[0].Mutability)),
				},
			},
		})
	})

	t.Run("error path - schema attributes.returned must be a valid value", func(t *testing.T) {

		schemaAttributes := []schemas.Attribute{
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
					Config:      ResourceSchema("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform", schemaAttributes[0]),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].returned value must be one of: \\[\"always\" \"never\"\n\"default\" \"request\"], got: \"%s\"", schemaAttributes[0].Returned)),
				},
			},
		})
	})

	t.Run("error path - schema attributes.uniqueness must be a valid value", func(t *testing.T) {

		schemaAttributes := []schemas.Attribute{
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
					Config:      ResourceSchema("testSchema", "urn:ietf:scim:schemas:Terraform", "Terraform", schemaAttributes[0]),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute attributes\\[0].uniqueness value must be one of: \\[\"none\" \"server\"\n\"global\"], got: \"%s\"", schemaAttributes[0].Uniqueness)),
				},
			},
		})
	})

}

func ResourceSchema(resourceName string, schemaId string, schemaName string, schemaAttributes schemas.Attribute) string {
	return fmt.Sprintf(`
	resource "sci_schema" "%s"{
		id = "%s"
		name = "%s"
		attributes = [
			{
				name = "%s"
				mutability = "%s"
				returned = "%s"
				type = "%s"
				uniqueness = "%s"
			}
		]
	}
	`, resourceName, schemaId, schemaName, schemaAttributes.Name, schemaAttributes.Mutability, schemaAttributes.Returned, schemaAttributes.Type, schemaAttributes.Uniqueness)
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
