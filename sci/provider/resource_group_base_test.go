package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGroupBase(t *testing.T) {
	t.Parallel()

	group := groups.Group{
		DisplayName: "Test Base Group",
		GroupExtension: &groups.GroupExtension{
			Name:        "Test-Base-Group",
			Description: "Base group for testing",
		},
	}

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_group_base")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceGroupBase("testGroupBase", group),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group_base.testGroupBase", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_group_base.testGroupBase", "display_name", group.DisplayName),
						resource.TestCheckResourceAttr("sci_group_base.testGroupBase", "group_extension.name", group.GroupExtension.Name),
						resource.TestCheckResourceAttr("sci_group_base.testGroupBase", "group_extension.description", group.GroupExtension.Description),
					),
				},
				{
					ResourceName:      "sci_group_base.testGroupBase",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - group_base update", func(t *testing.T) {
		updatedGroup := groups.Group{
			DisplayName: "Test Base Group - Updated",
			GroupExtension: &groups.GroupExtension{
				Name:        "Test-Base-Group",
				Description: "Updated description",
			},
		}

		rec, user := setupVCR(t, "fixtures/resource_group_base_updated")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceGroupBase("testGroupBase", group),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group_base.testGroupBase", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_group_base.testGroupBase", "display_name", group.DisplayName),
						resource.TestCheckResourceAttr("sci_group_base.testGroupBase", "group_extension.name", group.GroupExtension.Name),
						resource.TestCheckResourceAttr("sci_group_base.testGroupBase", "group_extension.description", group.GroupExtension.Description),
					),
				},
				{
					Config: providerConfig("", user) + ResourceGroupBase("testGroupBase", updatedGroup),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_group_base.testGroupBase", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_group_base.testGroupBase", "display_name", updatedGroup.DisplayName),
						resource.TestCheckResourceAttr("sci_group_base.testGroupBase", "group_extension.name", updatedGroup.GroupExtension.Name),
						resource.TestCheckResourceAttr("sci_group_base.testGroupBase", "group_extension.description", updatedGroup.GroupExtension.Description),
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
					Config:      ResourceGroupBaseWithoutSchemas("testGroupBase", "Some Group"),
					ExpectError: regexp.MustCompile("Attribute schemas set must contain at least 1 elements, got: 0"),
				},
			},
		})
	})

	t.Run("error path - display_name is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceGroupBaseWithoutDisplayName("testGroupBase"),
					ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
				},
			},
		})
	})
}

func ResourceGroupBase(resourceName string, group groups.Group) string {
	return fmt.Sprintf(`
	resource "sci_group_base" "%s" {
		display_name = "%s"
		group_extension = {
			name        = "%s"
			description = "%s"
		}
	}
	`, resourceName, group.DisplayName, group.GroupExtension.Name, group.GroupExtension.Description)
}

func ResourceGroupBaseWithoutSchemas(resourceName string, displayName string) string {
	return fmt.Sprintf(`
	resource "sci_group_base" "%s" {
		schemas      = []
		display_name = "%s"
	}
	`, resourceName, displayName)
}

func ResourceGroupBaseWithoutDisplayName(resourceName string) string {
	return fmt.Sprintf(`
	resource "sci_group_base" "%s" {
	}
	`, resourceName)
}
