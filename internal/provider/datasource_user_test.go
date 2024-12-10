package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceUser(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func (t *testing.T){

		rec, user := setupVCR(t, "fixtures/datasource_user")
		defer stopQuietly(rec)
		
		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceUser("testUser", "Test User"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.ias_user.testUser","id",regexpUUID),
						resource.TestCheckResourceAttr("data.ias_user.testUser","name.given_name","Test"),
						resource.TestCheckResourceAttr("data.ias_user.testUser","name.family_name","User"),
						resource.TestCheckResourceAttr("data.ias_user.testUser","emails.0.value","test.user@sap.com"),
						resource.TestCheckResourceAttr("data.ias_user.testUser","emails.0.primary","true"),
					),
				},
			},
		})
	})

	t.Run("error path - invalid user id", func (t *testing.T){

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: DataSourceUserById("testUser","invalid-user-id"),
					ExpectError: regexp.MustCompile(`Attribute id value must be a valid UUID, got: invalid-user-id`),
				},
			},
		})

	})

	t.Run("error path - user id is mandatory", func (t *testing.T){

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: DataSourceUserNoId("testUser"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		})

	})
}

func DataSourceUser(resourceName string, userName string) string {
	return fmt.Sprintf(`
	data "ias_users" "allUsers" {}
	data "ias_user" "%s" {
		id = [for user in data.ias_users.allUsers.values : user.id if join(" ",[user.name.given_name,user.name.family_name]) == "%s"][0]
	}
	`, resourceName, userName)
}

func DataSourceUserById(resourceName string, userId string) string {
	return fmt.Sprintf(`
	data "ias_user" "%s"{
		id = "%s"
	}
	`, resourceName, userId)
}

func DataSourceUserNoId(resourceName string) string {
	return fmt.Sprintf(`
	data "ias_user" "%s"{
	}
	`, resourceName)
}