package provider

import (
	"fmt"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGroup(t *testing.T) {
	t.Parallel()

	group := groups.Group{
		DisplayName: "Test Group",
		GroupMembers: []groups.GroupMember{
			{
				Value: "0b35d8cf-722c-4151-951e-176b623c0b78",
				Type:  "User",
			},
		},
		GroupExtension: &groups.GroupExtension{
			Name:        "Test-Group",
			Description: "For testing purposes",
		},
	}

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_group")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceGroup("testGroup", group),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group.testGroup", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_group.testGroup", "display_name", group.DisplayName),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.value", group.GroupMembers[0].Value),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.type", group.GroupMembers[0].Type),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_extension.name", group.GroupExtension.Name),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_extension.description", group.GroupExtension.Description),
					),
				},
			},
		})

	})

	t.Run("happy path - group update", func(t *testing.T) {

		updatedGroup := groups.Group{
			DisplayName: "Test Group - New",
			GroupMembers: []groups.GroupMember{
				{
					Value: "0b35d8cf-722c-4151-951e-176b623c0b78",
					Type:  "User",
				},
				{
					Value: "59aeb87b-777a-4034-8f3e-709d39fb1a18",
					Type:  "Group",
				},
			},
			GroupExtension: &groups.GroupExtension{
				Name:        "Test-Group",
				Description: "For testing purposes",
			},
		}

		rec, user := setupVCR(t, "fixtures/resource_group_updated")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceGroup("testGroup", group),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group.testGroup", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_group.testGroup", "display_name", group.DisplayName),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.value", group.GroupMembers[0].Value),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.type", group.GroupMembers[0].Type),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_extension.name", group.GroupExtension.Name),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_extension.description", group.GroupExtension.Description),
					),
				},
				{
					Config: providerConfig("", user) + ResourceGroup("testGroup", updatedGroup),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group.testGroup", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_group.testGroup", "display_name", updatedGroup.DisplayName),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.value", updatedGroup.GroupMembers[1].Value),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.type", updatedGroup.GroupMembers[1].Type),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.1.value", updatedGroup.GroupMembers[0].Value),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.1.type", updatedGroup.GroupMembers[0].Type),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_extension.name", updatedGroup.GroupExtension.Name),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_extension.description", updatedGroup.GroupExtension.Description),
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
					Config:      ResourceGroupWithoutSchemas("testGroup", "Terraform Group"),
					ExpectError: regexp.MustCompile("Attribute schemas set must contain at least 1 elements, got: 0"),
				},
			},
		})
	})

	t.Run("error path - display name is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceGroupWithoutDisplayName("testGroup"),
					ExpectError: regexp.MustCompile("The argument \"display_name\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - group_members.value must be a valid UUID", func(t *testing.T) {

		group.GroupMembers = []groups.GroupMember{
			{
				Value: "this-is-not-a-valid-UUID",
				Type:  "User",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceGroup("testGroup", group),
					ExpectError: regexp.MustCompile(fmt.Sprintf("\nvalue must be a valid UUID, got: %s", group.GroupMembers[0].Value)),
				},
			},
		})
	})

	t.Run("error path - group_members.value must be a valid member", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_group_invalid_group_member")
		defer stopQuietly(rec)

		group.GroupMembers = []groups.GroupMember{
			{
				Value: "5b4e7391-67d2-419f-8f0e-46f46f1f67ec",
				Type:  "User",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      providerConfig("", user) + ResourceGroup("testGroup", group),
					ExpectError: regexp.MustCompile(fmt.Sprintf("member %s is not found", group.GroupMembers[0].Value)),
				},
			},
		})
	})

	t.Run("error path - group_members.type must be a valid value", func(t *testing.T) {

		group.GroupMembers = []groups.GroupMember{
			{
				Value: "user-id",
				Type:  "this-is-not-a-valid-member-type",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceGroup("testGroup", group),
					ExpectError: regexp.MustCompile(fmt.Sprintf("\nvalue must be one of: \\[\"User\" \"Group\"], got:\n\"%s\"", group.GroupMembers[0].Type)),
				},
			},
		})
	})
}

func ResourceGroup(resoureName string, group groups.Group) string {
	return fmt.Sprintf(`
	resource "sci_group" "%s"{
		display_name = "%s"
		group_members = [%s]
		group_extension = {
			name = "%s"
			description = "%s"
		}
	}
	`, resoureName, group.DisplayName, getGroupMembers(group.GroupMembers), group.GroupExtension.Name, group.GroupExtension.Description)
}

func ResourceGroupWithoutSchemas(resoureName string, displayName string) string {
	return fmt.Sprintf(`
	resource "sci_group" "%s"{
		schemas = []
		display_name = "%s"
	}
	`, resoureName, displayName)
}

func ResourceGroupWithoutDisplayName(resoureName string) string {
	return fmt.Sprintf(`
	resource "sci_group" "%s"{
	}
	`, resoureName)
}

func getGroupMembers(groupMembers []groups.GroupMember) string {

	members := ""
	for _, member := range groupMembers {
		members += fmt.Sprintf(`{
			value = "%s"
			type = "%s"
		},`, member.Value, member.Type)
	}
	return members
}
