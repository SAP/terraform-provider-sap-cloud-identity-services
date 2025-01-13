package provider

import (
	"fmt"
	"terraform-provider-ias/internal/cli/apiObjects/users"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceUser (t *testing.T) {
	t. Parallel()

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

		updatedName := users.Name{
			FamilyName: "Doe S",
			GivenName: "Joe",
		}

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


}

func ResourceUser (resoureName string, userName string, schemas []string, name users.Name, emails []users.Email) string {
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
	`, resoureName, schemas[0], schemas[1], userName, name.FamilyName, name.GivenName, emails[0].Type, emails[0].Value)
}