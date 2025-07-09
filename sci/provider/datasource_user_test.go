package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var checkCustomSchemas resource.CheckResourceAttrWithFunc = func(value string) error {
	var err error
	if len(value) == 0 {
		err = fmt.Errorf("%s has length 0", value)
	}
	return err
}

func TestDataSourceUser(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/datasource_user")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceUser("testUser", "Terraform Test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_user.testUser", "id", regexpUUID),

						resource.TestCheckResourceAttr("data.sci_user.testUser", "name.given_name", "Terraform"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "name.family_name", "Test User"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "user_name", "Terraform Test"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "emails.0.value", "test.user1@sap.com"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "emails.0.primary", "false"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "emails.1.value", "test.user2@gmail.com"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "emails.1.type", "home"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "sap_extension_user.status", "inactive"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "user_type", "public"),
					),
				},
			},
		})
	})

	t.Run("happy path - user with custom schemas", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/datasource_user_with_custom_schemas")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceUser("testUser", "Terraform Custom Schemas"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_user.testUser", "id", regexpUUID),

						resource.TestCheckResourceAttr("data.sci_user.testUser", "name.given_name", "Terraform Custom Schemas"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "name.family_name", "User"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "emails.0.value", "custom.user@test.com"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "emails.0.primary", "true"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "sap_extension_user.status", "active"),
						resource.TestCheckResourceAttr("data.sci_user.testUser", "user_type", "employee"),
						resource.TestCheckResourceAttrWith("data.sci_user.testUser", "custom_schemas", checkCustomSchemas),
					),
				},
			},
		})
	})

	t.Run("error path - invalid user id", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceUserById("testUser", "invalid-user-id"),
					ExpectError: regexp.MustCompile(`Attribute id value must be a valid UUID, got: invalid-user-id`),
				},
			},
		})

	})

	t.Run("error path - user id is mandatory", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceUserNoId("testUser"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		})

	})
}

func DataSourceUser(resourceName string, userName string) string {
	return fmt.Sprintf(`
	data "sci_users" "allUsers" {}
	data "sci_user" "%s" {
		id = [for user in data.sci_users.allUsers.values : user.id if user.user_name == "%s"][0]
	}
	`, resourceName, userName)
}

func DataSourceUserById(resourceName string, userId string) string {
	return fmt.Sprintf(`
	data "sci_user" "%s"{
		id = "%s"
	}
	`, resourceName, userId)
}

func DataSourceUserNoId(resourceName string) string {
	return fmt.Sprintf(`
	data "sci_user" "%s"{
	}
	`, resourceName)
}
