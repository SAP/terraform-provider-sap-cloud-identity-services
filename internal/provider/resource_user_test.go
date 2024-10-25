package provider

import (
	"fmt"
	"regexp"
	"terraform-provider-ias/internal/cli/apiObjects/users"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceUser (t *testing.T) {

	schemas := []string{
		"urn:ietf:params:scim:schemas:extension:sap:2.0:User",
		"urn:ietf:params:scim:schemas:core:2.0:User",
	}

	name := users.Name{
		FamilyName: "Doe",
		GivenName: "Joe",
	}

	emails := []users.Email{
		{
			Type: "work",
			Value: "joe.doe@test.com",
		},
	}

	updatedName := users.Name{
		FamilyName: "Doe S",
		GivenName: "Joe",
	}

	t. Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUser("testUser", "Joe Doe", schemas, name, emails),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_name", "Joe Doe"),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.family_name", name.FamilyName),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.given_name", name.GivenName),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.type", emails[0].Type),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.value", emails[0].Value),
					),
				},
			},
		})

	})

	t.Run("happy path - user update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user_updated")
		defer stopQuietly(rec)


		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUser("testUser", "Joe Doe", schemas, name, emails),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_name", "Joe Doe"),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.family_name", name.FamilyName),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.given_name", name.GivenName),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.type", emails[0].Type),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.value", emails[0].Value),
					),
				},
				{
					Config: providerConfig("", user) + ResourceUser("testUser", "Joe Doe S", schemas, updatedName, emails),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_user.testUser", "user_name", "Joe Doe S"),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.family_name", updatedName.FamilyName),
						resource.TestCheckResourceAttr("ias_user.testUser", "name.given_name", updatedName.GivenName),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.type", emails[0].Type),
						resource.TestCheckResourceAttr("ias_user.testUser", "emails.0.value", emails[0].Value),
					),
				},
			},
		})

	})

	t.Run("error path - schemas cannot be empty", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user_empty_schemas")
		defer stopQuietly(rec)
		
		resource.Test(t, resource.TestCase{
			IsUnitTest: 			  true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: 	 providerConfig("", user) + ResourceUserWithoutSchemas("testUser", "Joe Doe", name, emails),
					ExpectError: regexp.MustCompile("Provide a valid value for \"schemas\""),
				},
			},
		})

	})	

	t.Run("error path - username is mandatory", func(t *testing.T){
		resource.Test(t, resource.TestCase{
			IsUnitTest: 			  true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: 	 ResourceUserWithoutUserName("testUser", schemas, name, emails),
					ExpectError: regexp.MustCompile("The argument \"user_name\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - name cannot be empty", func(t *testing.T){
		resource.Test(t, resource.TestCase{
			IsUnitTest: 			  true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: 	 ResourceUserWithoutName("testUser", "Joe Doe", schemas, emails),
					ExpectError: regexp.MustCompile("Inappropriate value for attribute \"name\": attributes \"family_name\" and\n\"given_name\" are required."),
				},
			},
		})
	})

	t.Run("error path - emails cannot be empty", func(t *testing.T){
		resource.Test(t, resource.TestCase{
			IsUnitTest: 			  true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: 	 ResourceUserWithoutEmails("testUser", "Joe Doe", schemas, name),
					ExpectError: regexp.MustCompile("Inappropriate value for attribute \"email\": element 0: attributes \"type\" and\n \"value\" are required."),
				},
			},
		})
	})

}

func ResourceUser (resourceName string, userName string, schemas []string, name users.Name, emails []users.Email) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		schemas = [
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
	}
	`, resourceName, schemas[0], schemas[1], userName, name.FamilyName, name.GivenName, emails[0].Type, emails[0].Value)
}

func ResourceUserWithoutSchemas (resourceName string, userName string, name users.Name, emails []users.Email) string {
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

func ResourceUserWithoutUserName (resourceName string, schemas []string, name users.Name, emails []users.Email) string { 
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		schemas = [
			"%s",
			"%s"
		]
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
	`,resourceName, schemas[0], schemas[1], name.FamilyName, name.GivenName, emails[0].Type, emails[0].Value)
}

func ResourceUserWithoutName (resourceName string, userName string, schemas []string, emails []users.Email) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		schemas = [
			"%s",
			"%s"
		]
		user_name = "%s"
		name = {}
		emails = [
			{
				type = "%s"
				value = "%s"
			}
		]
	}
	`, resourceName, schemas[0], schemas[1], userName, emails[0].Type, emails[0].Value)
}

func ResourceUserWithoutEmails (resourceName string, userName string, schemas []string, name users.Name) string {
	return fmt.Sprintf(`
	resource "ias_user" "%s"{
		schemas = [
			"%s",
			"%s"
		]
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [
			{}
		]
	}
	`, resourceName, schemas[0], schemas[1], userName, name.FamilyName, name.GivenName)
}