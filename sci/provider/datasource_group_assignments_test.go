package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGroupAssignments(t *testing.T) {
	t.Parallel()

	groupName := "Test Group"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_group_assignments")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceGroupAssignments("assignments", groupName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_group_assignments.assignments", "group_id", regexpUUID),
						resource.TestCheckResourceAttr("data.sci_group_assignments.assignments", "values.#", "2"),
					),
				},
			},
		})
	})

	t.Run("error path - invalid group_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceGroupAssignmentsByGroupId("assignments", "invalid-uuid"),
					ExpectError: regexp.MustCompile(`Attribute group_id value must be a valid UUID, got: invalid-uuid`),
				},
			},
		})
	})

	t.Run("error path - group_id is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceGroupAssignmentsNoGroupId("assignments"),
					ExpectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func DataSourceGroupAssignments(datasourceName string, groupName string) string {
	return fmt.Sprintf(`
	data "sci_group_bases" "allGroups" {}
	data "sci_group_base" "group" {
		id = [for group in data.sci_group_bases.allGroups.values : group.id if group.group_extension.name == "%s"][0]
	}
	data "sci_group_assignments" "%s" {
		group_id = data.sci_group_base.group.id
	}
	`, groupName, datasourceName)
}

func DataSourceGroupAssignmentsByGroupId(datasourceName string, groupId string) string {
	return fmt.Sprintf(`
	data "sci_group_assignments" "%s" {
		group_id = "%s"
	}
	`, datasourceName, groupId)
}

func DataSourceGroupAssignmentsNoGroupId(datasourceName string) string {
	return fmt.Sprintf(`
	data "sci_group_assignments" "%s" {}
	`, datasourceName)
}
