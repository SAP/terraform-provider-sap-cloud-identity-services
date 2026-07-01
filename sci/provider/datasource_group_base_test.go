package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGroupBase(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_group_base")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceGroupBase("testGroupBase", "Base Test Group"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_group_base.testGroupBase", "id", regexpUUID),
						resource.TestCheckResourceAttr("data.sci_group_base.testGroupBase", "display_name", "Terraform Test Base Group"),
						resource.TestCheckResourceAttr("data.sci_group_base.testGroupBase", "group_extension.name", "Base Test Group"),
					),
				},
			},
		})
	})

	t.Run("error path - invalid group id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceGroupBaseById("testGroupBase", "invalid-uuid"),
					ExpectError: regexp.MustCompile(`Attribute id value must be a valid UUID, got: invalid-uuid`),
				},
			},
		})
	})

	t.Run("error path - group id is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceGroupBaseNoId("testGroupBase"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func DataSourceGroupBase(datasourceName string, groupDisplayName string) string {
	return fmt.Sprintf(`
	data "sci_group_bases" "allGroups" {}
	data "sci_group_base" "%s" {
		id = [for group in data.sci_group_bases.allGroups.values : group.id if group.group_extension.name == "%s"][0]
	}
	`, datasourceName, groupDisplayName)
}

func DataSourceGroupBaseById(datasourceName string, groupId string) string {
	return fmt.Sprintf(`
	data "sci_group_base" "%s" {
		id = "%s"
	}
	`, datasourceName, groupId)
}

func DataSourceGroupBaseNoId(datasourceName string) string {
	return fmt.Sprintf(`
	data "sci_group_base" "%s" {}
	`, datasourceName)
}
