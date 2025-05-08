package provider

import (
	"fmt"
	"regexp"
	"terraform-provider-sci/internal/cli/apiObjects/groups"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGroup(t *testing.T) {
	t.Parallel()

	schemas := []string{
		"urn:ietf:params:scim:schemas:core:2.0:Group",
		"urn:sap:cloud:scim:schemas:extension:custom:2.0:Group",
	}

	members := []groups.GroupMember{
		{
			Value: "0b35d8cf-722c-4151-951e-176b623c0b78",
			Type:  "User",
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
					Config: providerConfig("", user) + ResourceGroup("testGroup", "Terraform Group", schemas, "For testing purposes", members),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group.testGroup", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_group.testGroup", "display_name", "Terraform Group"),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_extension.description", "For testing purposes"),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.value", members[0].Value),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.type", members[0].Type),
					),
				},
			},
		})

	})

	t.Run("happy path - group update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_group_updated")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceGroup("testGroup", "Terraform Group", schemas, "For testing purposes", members),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group.testGroup", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_group.testGroup", "display_name", "Terraform Group"),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_extension.description", "For testing purposes"),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.value", members[0].Value),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.type", members[0].Type),
					),
				},
				{
					Config: providerConfig("", user) + ResourceGroup("testGroup", "Updated Terraform Group", schemas, "For testing purposes", members),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group.testGroup", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_group.testGroup", "display_name", "Updated Terraform Group"),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_extension.description", "For testing purposes"),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.value", members[0].Value),
						resource.TestCheckResourceAttr("sci_group.testGroup", "group_members.0.type", members[0].Type),
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

		members := []groups.GroupMember{
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
					Config:      ResourceGroup("testGroup", "Terraform Group", schemas, "For testing purposes", members),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute group_members\\[0].value value must be a valid UUID, got:\n%s", "this-is-not-a-valid-UUID")),
				},
			},
		})
	})

	t.Run("error path - group_members.value must be a valid member", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_group_invalid_group_member")
		defer stopQuietly(rec)

		members := []groups.GroupMember{
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
					Config:      providerConfig("", user) + ResourceGroup("testGroup", "Terraform Group", schemas, "For testing purposes", members),
					ExpectError: regexp.MustCompile(fmt.Sprintf("member %s is not found", members[0].Value)),
				},
			},
		})
	})

	t.Run("error path - group_members.type must be a valid value", func(t *testing.T) {

		members := []groups.GroupMember{
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
					Config:      ResourceGroup("testGroup", "Terraform Group", schemas, "For testing purposes", members),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute group_members\\[0].type value must be one of: \\[\"User\" \"Group\"], got:\n\"%s\"", "this-is-not-a-valid-member-type")),
				},
			},
		})
	})
}

func ResourceGroup(resoureName string, displayName string, schemas []string, description string, members []groups.GroupMember) string {
	return fmt.Sprintf(`
	resource "sci_group" "%s"{
		schemas = [
			"%s",
			"%s"
		]
		display_name = "%s"
		group_members = [
			{
				value = "%s",
				type = "%s"
			}
		]
		group_extension = {
			description = "%s"
		}
	}
	`, resoureName, schemas[0], schemas[1], displayName, members[0].Value, members[0].Type, description)
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
