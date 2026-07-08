package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGroupAssignment(t *testing.T) {
	t.Parallel()

	groupName := "Test Group"
	userName := "Terraform User Assignment"
	mockUuid := "af2f7963-358d-4336-bc51-57099394dee7"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_group_assignment")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceGroupAssignment("testAssignment", groupName, userName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.sci_group_assignment.testAssignment", "group_id", regexpUUID),
						resource.TestMatchResourceAttr("data.sci_group_assignment.testAssignment", "group_member.value", regexpUUID),
						resource.TestCheckResourceAttr("data.sci_group_assignment.testAssignment", "group_member.type", "User"),
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
					Config:      DataSourceGroupAssignmentById("testAssignment", "invalid-uuid", mockUuid),
					ExpectError: regexp.MustCompile(`Attribute group_id value must be a valid UUID, got: invalid-uuid`),
				},
			},
		})
	})

	t.Run("error path - invalid member value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceGroupAssignmentById("testAssignment", mockUuid, "invalid-uuid"),
					ExpectError: regexp.MustCompile(`value must be a valid UUID`),
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
					Config:      DataSourceGroupAssignmentNoGroupId("testAssignment", mockUuid),
					ExpectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - group_member is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      DataSourceGroupAssignmentNoMember("testAssignment", mockUuid),
					ExpectError: regexp.MustCompile(`The argument "group_member" is required, but no definition was found.`),
				},
			},
		})
	})
}

func DataSourceGroupAssignment(datasourceName string, groupName string, userName string) string {
	return fmt.Sprintf(`
	
	data "sci_group_bases" "allGroups" {}

	data "sci_group_base" "group" {
		id = [for group in data.sci_group_bases.allGroups.values : group.id if group.group_extension.name == "%s"][0]
	}

	data "sci_users" "allUsers" {}

	data "sci_user" "user" {
		id = [for user in data.sci_users.allUsers.values : user.id if user.user_name == "%s"][0]
	}
	
	data "sci_group_assignment" "%s" {
		group_id = data.sci_group_base.group.id
		group_member = {
			value = data.sci_user.user.id
		}
	}
	`, groupName, userName, datasourceName)
}

func DataSourceGroupAssignmentNoGroupId(datasourceName string, memberValue string) string {
	return fmt.Sprintf(`
	data "sci_group_assignment" "%s" {
		group_member = {
			value = "%s"
		}
	}
	`, datasourceName, memberValue)
}

func DataSourceGroupAssignmentNoMember(datasourceName string, groupId string) string {
	return fmt.Sprintf(`
	data "sci_group_assignment" "%s" {
		group_id = "%s"
	}
	`, datasourceName, groupId)
}

func DataSourceGroupAssignmentById(datasourceName string, groupId string, userId string) string {
	return fmt.Sprintf(`
	data "sci_group_assignment" "%s" {
		group_id = "%s"
		group_member = {
			value = "%s"
		}
	}
	`, datasourceName, groupId, userId)
}
