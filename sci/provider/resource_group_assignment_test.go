package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceGroupAssignment(t *testing.T) {
	t.Parallel()

	groupName := "Test Group"
	userName := "Terraform Test"
	mockUuid := "af2f7963-358d-4336-bc51-57099394dee7"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_group_assignment")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceGroupAssignmentByGroupNameAndUserName("testAssignment", groupName, userName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group_assignment.testAssignment", "group_id", regexpUUID),
						resource.TestMatchResourceAttr("sci_group_assignment.testAssignment", "group_member.value", regexpUUID),
						resource.TestCheckResourceAttr("sci_group_assignment.testAssignment", "group_member.type", "User"),
					),
				},
				{
					ResourceName:      "sci_group_assignment.testAssignment",
					ImportStateIdFunc: getGroupAssignmentImportStateId("sci_group_assignment.testAssignment"),
					ImportState:       true,
				},
			},
		})
	})

	t.Run("happy path - update group member", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_group_assignment_update_group_member")
		defer stopQuietly(rec)

		updatedUserName := "Terraform Custom Schemas"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceGroupAssignmentByGroupNameAndUserName("testAssignment", groupName, userName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group_assignment.testAssignment", "group_id", regexpUUID),
						resource.TestMatchResourceAttr("sci_group_assignment.testAssignment", "group_member.value", regexpUUID),
						resource.TestCheckResourceAttr("sci_group_assignment.testAssignment", "group_member.type", "User"),
					),
				},
				{
					Config: providerConfig("", user) + ResourceGroupAssignmentByGroupNameAndUserName("testAssignment", groupName, updatedUserName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group_assignment.testAssignment", "group_id", regexpUUID),
						resource.TestMatchResourceAttr("sci_group_assignment.testAssignment", "group_member.value", regexpUUID),
						resource.TestCheckResourceAttr("sci_group_assignment.testAssignment", "group_member.type", "User"),
					),
				},
			},
		})
	})

	t.Run("happy path - update assigned group", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_group_assignment_update_assigned_group")
		defer stopQuietly(rec)

		groupName := "Base Test Group"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceGroupAssignmentByGroupNameAndUserName("testAssignment", groupName, userName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group_assignment.testAssignment", "group_id", regexpUUID),
						resource.TestMatchResourceAttr("sci_group_assignment.testAssignment", "group_member.value", regexpUUID),
						resource.TestCheckResourceAttr("sci_group_assignment.testAssignment", "group_member.type", "User"),
					),
				},
			},
		})
	})

	t.Run("error path - group_id must be a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceGroupAssignmentByGroupIdAndUserId("testAssignment", "not-a-valid-uuid", mockUuid, "User"),
					ExpectError: regexp.MustCompile(`value must be a valid UUID, got: not-a-valid-uuid`),
				},
			},
		})
	})

	t.Run("error path - group_member.value must be a valid UUID", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceGroupAssignmentByGroupIdAndUserId("testAssignment", mockUuid, "not-a-valid-uuid", "User"),
					ExpectError: regexp.MustCompile(`value must be a valid UUID`),
				},
			},
		})
	})

	t.Run("error path - group_member.type must be a valid value", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceGroupAssignmentByGroupIdAndUserId("testAssignment", mockUuid, mockUuid, "invalid-user-type"),
					ExpectError: regexp.MustCompile(`value must be one of: \["User" "Group"\]`),
				},
			},
		})
	})

	t.Run("error path - invalid user assigned to group", func(t *testing.T) {
		unknownMember := groups.GroupMember{
			Value: mockUuid,
			Type:  "User",
		}

		rec, user := setupVCR(t, "fixtures/resource_group_assignment_invalid_member")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      providerConfig("", user) + ResourceGroupAssignmentByUserId("testAssignment", groupName, unknownMember),
					ExpectError: regexp.MustCompile(fmt.Sprintf("member %s is not found", unknownMember.Value)),
				},
			},
		})
	})

	t.Run("error path - user assigned to invalid group", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_group_assignment_invalid_group")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      providerConfig("", user) + ResourceGroupAssignmentByGroupId("testAssignment", mockUuid, userName, "User"),
					ExpectError: regexp.MustCompile("Error adding group member"),
				},
			},
		})
	})
}

func ResourceGroupAssignmentByGroupNameAndUserName(resourceName string, groupName string, userName string) string {
	return fmt.Sprintf(`

	data "sci_group_bases" "allGroups" {}

	data "sci_group_base" "group" {
		id = [for group in data.sci_group_bases.allGroups.values : group.id if group.group_extension.name == "%s"][0]
	}

	data "sci_users" "allUsers" {}

	data "sci_user" "user" {
		id = [for user in data.sci_users.allUsers.values : user.id if user.user_name == "%s"][0]
	}

	resource "sci_group_assignment" "%s" {
		group_id = data.sci_group_base.group.id
		group_member = {
			value = data.sci_user.user.id
			type  = "%s"
		}
	}
	`, groupName, userName, resourceName, "User")
}

func ResourceGroupAssignmentByGroupIdAndUserId(resourceName string, groupId string, userId string, userType string) string {
	return fmt.Sprintf(`
	resource "sci_group_assignment" "%s" {
		group_id = "%s"
		group_member = {
			value = "%s"
			type = "%s"
		}
	}
	`, resourceName, groupId, userId, userType)
}

func ResourceGroupAssignmentByUserId(resourceName string, groupName string, user groups.GroupMember) string {
	return fmt.Sprintf(`
	data "sci_group_bases" "allGroups" {}

	data "sci_group_base" "group" {
		id = [for group in data.sci_group_bases.allGroups.values : group.id if group.group_extension.name == "%s"][0]
	}

	resource "sci_group_assignment" "%s" {
		group_id = data.sci_group_base.group.id
		group_member = {
			value = "%s"
			type = "%s"
		}
	}
	`, groupName, resourceName, user.Value, user.Type)
}

func ResourceGroupAssignmentByGroupId(resourceName string, groupId string, userName string, userType string) string {
	return fmt.Sprintf(`
	data "sci_users" "allUsers" {}

	data "sci_user" "user" {
		id = [for user in data.sci_users.allUsers.values : user.id if user.user_name == "%s"][0]
	}

	resource "sci_group_assignment" "%s" {
		group_id = "%s"
		group_member = {
			value = data.sci_user.user.id
			type  = "%s"
		}
	}
	`, userName, resourceName, groupId, userType)
}

func getGroupAssignmentImportStateId(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["group_member.value"]), nil
	}
}
