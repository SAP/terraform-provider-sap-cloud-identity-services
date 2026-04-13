package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceApplicationSecret(t *testing.T) {

	t.Parallel()

	t.Run("happy path - create", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_client_secret")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceApplicationSecretByAppName("testSecret", "5fd22812-53c5-4803-8285-94cd1fb3b301", []string{"manageApp"}, "test secret", "2029-10-12T10:00:00Z"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_client_secret.testSecret", "id", regexpUUID),
						resource.TestMatchResourceAttr("sci_client_secret.testSecret", "application_id", regexpUUID),
						resource.TestCheckResourceAttr("sci_client_secret.testSecret", "description", "test secret"),
						resource.TestCheckResourceAttr("sci_client_secret.testSecret", "valid_to", "2029-10-12T10:00:00Z"),
						resource.TestCheckTypeSetElemAttr("sci_client_secret.testSecret", "authorization_scopes.*", "manageApp"),
						resource.TestCheckResourceAttrSet("sci_client_secret.testSecret", "secret"),
						resource.TestCheckResourceAttrSet("sci_client_secret.testSecret", "client_id"),
					),
				},
				{
					ResourceName: "sci_client_secret.testSecret",
					ImportState:  true,
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						rs := s.RootModule().Resources["sci_client_secret.testSecret"]
						return rs.Primary.Attributes["application_id"] + "," + rs.Primary.ID, nil
					},
					ImportStateVerify: true,
					// secret, hint, and client_id are only returned at creation time and cannot be retrieved again
					ImportStateVerifyIgnore: []string{"secret", "hint", "client_id"},
				},
			},
		})
	})

	t.Run("happy path - update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_client_secret_updated")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceApplicationSecretByAppName("testSecret", "5fd22812-53c5-4803-8285-94cd1fb3b301", []string{"manageApp"}, "test secret", "2029-10-12T10:00:00Z"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_client_secret.testSecret", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_client_secret.testSecret", "description", "test secret"),
					),
				},
				{
					Config: providerConfig("", user) + ResourceApplicationSecretByAppName("testSecret", "5fd22812-53c5-4803-8285-94cd1fb3b301", []string{"manageApp", "oAuth"}, "updated secret", "2030-01-01T00:00:00Z"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_client_secret.testSecret", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_client_secret.testSecret", "description", "updated secret"),
						resource.TestCheckResourceAttr("sci_client_secret.testSecret", "valid_to", "2030-01-01T00:00:00Z"),
						resource.TestCheckTypeSetElemAttr("sci_client_secret.testSecret", "authorization_scopes.*", "manageApp"),
						resource.TestCheckTypeSetElemAttr("sci_client_secret.testSecret", "authorization_scopes.*", "oAuth"),
					),
				},
			},
		})
	})

	t.Run("error path - application not found", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_client_secret_app_not_found")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      providerConfig("", user) + ResourceApplicationSecret("testSecret", "00000000-0000-0000-0000-000000000000", []string{"manageApp"}, "test secret", "2029-10-12T10:00:00Z"),
					ExpectError: regexp.MustCompile(`Error creating application secret`),
				},
			},
		})
	})

	t.Run("error path - authorization_scopes must be valid values", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationSecret("testSecret", "64bdf294-d1f3-47ac-8a4e-b1bb88380ffc", []string{"invalidScope"}, "test secret", "2029-10-12T10:00:00Z"),
					ExpectError: regexp.MustCompile(`(?s)value must be one of:.*"manageApp".*"oAuth".*"readUserProfile".*"manageUsers"`),
				},
			},
		})
	})

}

// ResourceApplicationSecretByAppName looks up an application by name and creates a secret for it.
func ResourceApplicationSecretByAppName(resourceName, appID string, scopes []string, description, validTo string) string {
	scopesList := ""
	for _, s := range scopes {
		scopesList += fmt.Sprintf(`"%s", `, s)
	}

	return fmt.Sprintf(`
resource "sci_client_secret" "%s" {
  application_id       = "%s"
  description          = "%s"
  valid_to             = "%s"
  authorization_scopes = [%s]
}`, resourceName, appID, description, validTo, scopesList)
}

// ResourceApplicationSecret creates a secret with a hard-coded application_id. Used for error-path tests.
func ResourceApplicationSecret(resourceName, applicationId string, scopes []string, description, validTo string) string {
	var scopesList strings.Builder
	for _, s := range scopes {
		fmt.Fprintf(&scopesList, `"%s", `, s)
	}

	return fmt.Sprintf(`
resource "sci_client_secret" "%s" {
  application_id       = "%s"
  description          = "%s"
  valid_to             = "%s"
  authorization_scopes = [%s]
}`, resourceName, applicationId, description, validTo, scopesList.String())
}
