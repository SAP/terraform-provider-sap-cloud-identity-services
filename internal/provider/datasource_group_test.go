package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGroup(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func (t *testing.T) {

		rec, user := setupVCR(t, "fixtures/datasource_group")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()), 
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceGroup("testGroup", "testGroup"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.ias_group.testGroup", "id", regexpUUID),
						resource.TestCheckResourceAttr("data.ias_group.testGroup", "display_name", "testGroup"),
						resource.TestMatchResourceAttr("data.ias_group.testGroup", "group_members.0.value",regexpUUID),
						resource.TestCheckResourceAttr("data.ias_group.testGroup", "group_members.0.type", "User"),
					),
				},
			},
		})

	})

	t.Run("error path - invalid group id", func(t *testing.T){

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: DataSourceGroupById("testGroup","invalid-uuid"),
					ExpectError: regexp.MustCompile(`Attribute id value must be a valid UUID, got: invalid-uuid`),
				},
			},
		})

	})

	t.Run("error path - group id is mandatory", func(t *testing.T){

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: DataSourceGroupNoId("testGroup"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		})
		
	})
}

func DataSourceGroup(datasourceName string, groupName string) string {
	return fmt.Sprintf(`
	data "ias_groups" "allGroups" {}
	data "ias_group" "%s" {
		id = [for group in data.ias_groups.allGroups.values : group.id if group.display_name == "%s"][0]
	}
	`,datasourceName, groupName)
}

func DataSourceGroupById(datasourceName string, groupId string) string {
	return fmt.Sprintf(`
	data "ias_group" "%s" {
		id = "%s"
	}
	`, datasourceName, groupId)
}

func DataSourceGroupNoId(datasourceName string) string {
	return fmt.Sprintf(`
	data "ias_group" "%s" {}
	`, datasourceName)
}