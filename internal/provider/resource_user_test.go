package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"terraform-provider-ias/internal/cli/apiObjects/users"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceUser(t *testing.T) {

	name := users.Name{
		FamilyName: "Doe",
		GivenName:  "Joe",
	}

	emails := []users.Email{
		{
			Type:  "work",
			Value: "joe.doe@test.com",
		},
	}

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUser("testUser", "Joe Doe", name, emails),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_name", "Joe Doe"),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.family_name", name.FamilyName),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.given_name", name.GivenName),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.type", emails[0].Type),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.value", emails[0].Value),
						resource.TestCheckResourceAttr("ias_user.testUser", "active", "false"),
						resource.TestCheckResourceAttr("ias_user.testUser", "sap_extension_user.status", "inactive"),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_type", "public"),
					),
				},
			},
		})

	})

	t.Run("happy path - user update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user_updated")
		defer stopQuietly(rec)

		updatedName := users.Name{
			FamilyName: "Doe S",
			GivenName:  "Joe",
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUser("testUser", "Joe Doe", name, emails),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_name", "Joe Doe"),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.family_name", name.FamilyName),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.given_name", name.GivenName),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.type", emails[0].Type),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.value", emails[0].Value),
						resource.TestCheckResourceAttr("ias_user.testUser", "active", "false"),
						resource.TestCheckResourceAttr("ias_user.testUser", "sap_extension_user.status", "inactive"),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_type", "public"),
					),
				},
				{
					Config: providerConfig("", user) + ResourceUser("testUser", "Joe Doe S", updatedName, emails),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_name", "Joe Doe S"),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.family_name", updatedName.FamilyName),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.given_name", updatedName.GivenName),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.type", emails[0].Type),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.value", emails[0].Value),
						resource.TestCheckResourceAttr("ias_user.testUser", "active", "false"),
						resource.TestCheckResourceAttr("ias_user.testUser", "sap_extension_user.status", "inactive"),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_type", "public"),
					),
				},
			},
		})

	})

	t.Run("happy path - custom schemas", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user_with_custom_schemas")
		defer stopQuietly(rec)

		schemas := []string{
			"urn:ietf:params:scim:schemas:extension:sap:2.0:User",
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:test:terraform:1.0:User",
		}

		customSchemas, _ := json.Marshal(
			map[string]interface{}{
				"urn:test:terraform:1.0:User": map[string]interface{}{
					"test1": "testValue",
					"test2": false,
					"test3": map[string]interface{}{
						"test3a": 12.33,
						"test3b": 12,
					},
				},
			},
		)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUserWithCustomSchemas("testUser", "Joe Doe", schemas, name, emails, string(customSchemas)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_name", "Joe Doe"),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.family_name", name.FamilyName),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.given_name", name.GivenName),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.type", emails[0].Type),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.value", emails[0].Value),
						resource.TestCheckResourceAttr("ias_user.testUser", "active", "false"),
						resource.TestCheckResourceAttr("ias_user.testUser", "sap_extension_user.status", "inactive"),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_type", "public"),
						resource.TestCheckResourceAttr("ias_user.testUser", "custom_schemas", string(customSchemas)),
					),
				},
			},
		})

	})

	t.Run("happy path - update custom schemas", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user_with_updated_custom_schemas")
		defer stopQuietly(rec)

		schemas := []string{
			"urn:ietf:params:scim:schemas:extension:sap:2.0:User",
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:test:terraform:1.0:User",
		}

		customSchemas, _ := json.Marshal(
			map[string]interface{}{
				"urn:test:terraform:1.0:User": map[string]interface{}{
					"test1": "testValue",
					"test2": false,
					"test3": map[string]interface{}{
						"test3a": 12.33,
						"test3b": 12,
					},
				},
			},
		)

		newCustomSchemas, _ := json.Marshal(
			map[string]interface{}{
				"urn:test:terraform:1.0:User": map[string]interface{}{
					"test1": "newTestValue",
					"test2": true,
				},
			},
		)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUserWithCustomSchemas("testUser", "Joe Doe", schemas, name, emails, string(customSchemas)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_name", "Joe Doe"),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.family_name", name.FamilyName),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.given_name", name.GivenName),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.type", emails[0].Type),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.value", emails[0].Value),
						resource.TestCheckResourceAttr("ias_user.testUser", "active", "false"),
						resource.TestCheckResourceAttr("ias_user.testUser", "sap_extension_user.status", "inactive"),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_type", "public"),
						resource.TestCheckResourceAttr("ias_user.testUser", "custom_schemas", string(customSchemas)),
					),
				},
				{
					Config: providerConfig("", user) + ResourceUserWithCustomSchemas("testUser", "Joe Doe", schemas, name, emails, string(newCustomSchemas)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_name", "Joe Doe"),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.family_name", name.FamilyName),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.given_name", name.GivenName),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.type", emails[0].Type),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.value", emails[0].Value),
						resource.TestCheckResourceAttr("ias_user.testUser", "active", "false"),
						resource.TestCheckResourceAttr("ias_user.testUser", "sap_extension_user.status", "inactive"),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_type", "public"),
						resource.TestCheckResourceAttr("ias_user.testUser", "custom_schemas", string(newCustomSchemas)),
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
					Config:       ResourceUserWithoutSchemas("testUser", "Joe Doe", name, emails),
					ExpectError: regexp.MustCompile("Attribute schemas set must contain at least 1 elements, got: 0"),
				},
			},
		})

	})

	t.Run("error path - username is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithoutUserName("testUser", name, emails),
					ExpectError: regexp.MustCompile("The argument \"user_name\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - name cannot be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithoutName("testUser", "Joe Doe", emails),
					ExpectError: regexp.MustCompile("Inappropriate value for attribute \"name\": attributes \"family_name\" and\n\"given_name\" are required."),
				},
			},
		})
	})

	t.Run("error path - emails cannot be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithoutEmails("testUser", "Joe Doe", name),
					ExpectError: regexp.MustCompile("Inappropriate value for attribute \"emails\": element 0: attributes \"type\" and\n\"value\" are required."),
				},
			},
		})
	})

	t.Run("error path - user type must be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithUserType("testUser", "Joe Doe", name, emails, "this-is-not-a-valid-user-type"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute user_type value must be one of: \\[\"public\" \"partner\" \"customer\"\n\"external\" \"onboardee\" \"employee\"\\], got: \"%s\"", "this-is-not-a-valid-user-type")),
				},
			},
		})
	})

	t.Run("error path - status must be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithStatus("testUser", "Joe Doe", name, emails, "this-is-not-a-valid-status"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute sap_extension_user.status value must be one of: \\[\"active\"\n\"inactive\" \"new\"\\], got: \"%s\"", "this-is-not-a-valid-status")),
				},
			},
		})
	})

	t.Run("error path - custom schemas must be a valid json string", func(t *testing.T) {

		schemas := []string{
			"urn:ietf:params:scim:schemas:extension:sap:2.0:User",
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:test:terraform:1.0:User",
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithCustomSchemas("testUser", "Joe Doe", schemas, name, emails, "this-is-not-a-valid-json-string"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute custom_schemas value must be valid json, got:\n\"%s\"", "this-is-not-a-valid-json-string")),
				},
			},
		})
	})

}

func ResourceUser(resourceName string, userName string, name users.Name, emails []users.Email) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [
			{
				type = "%s"
				value = "%s"
			}
		]
	}
	`, resourceName, userName, name.FamilyName, name.GivenName, emails[0].Type, emails[0].Value)
}

func ResourceUserWithoutSchemas(resourceName string, userName string, name users.Name, emails []users.Email) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		schemas = []
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [
			{
				type = "%s"
				value = "%s"
			}
		]
	}
	`, resourceName, userName, name.FamilyName, name.GivenName, emails[0].Type, emails[0].Value)
}

func ResourceUserWithoutUserName(resourceName string, name users.Name, emails []users.Email) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [
			{
				type = "%s"
				value = "%s"
			}
		]
	}
	`, resourceName, name.FamilyName, name.GivenName, emails[0].Type, emails[0].Value)
}

func ResourceUserWithoutName(resourceName string, userName string, emails []users.Email) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		user_name = "%s"
		name = {}
		emails = [
			{
				type = "%s"
				value = "%s"
			}
		]
	}
	`, resourceName, userName, emails[0].Type, emails[0].Value)
}

func ResourceUserWithoutEmails(resourceName string, userName string, name users.Name) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [
			{}
		]
	}
	`, resourceName, userName, name.FamilyName, name.GivenName)
}

func ResourceUserWithUserType(resourceName string, userName string, name users.Name, emails []users.Email, userType string) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [
			{
				type = "%s"
				value = "%s"
			}
		]
		user_type = "%s"
	}
	`, resourceName, userName, name.FamilyName, name.GivenName, emails[0].Type, emails[0].Value, userType)
}

func ResourceUserWithStatus(resourceName string, userName string, name users.Name, emails []users.Email, status string) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [
			{
				type = "%s"
				value = "%s"
			}
		]
		sap_extension_user = {
			status = "%s"
		}
	}
	`, resourceName, userName, name.FamilyName, name.GivenName, emails[0].Type, emails[0].Value, status)
}

func ResourceUserWithCustomSchemas(resourceName string, userName string, schemas []string, name users.Name, emails []users.Email, customSchemas string) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		schemas = [
			"%s",
			"%s",
			"%s"
		]
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [
			{
				type = "%s"
				value = "%s"
			}
		]
		custom_schemas = %q
	}
	`, resourceName, schemas[0], schemas[1], schemas[2], userName, name.FamilyName, name.GivenName, emails[0].Type, emails[0].Value, customSchemas)
}
