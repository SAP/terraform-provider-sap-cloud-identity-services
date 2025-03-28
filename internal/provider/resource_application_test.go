package provider

import (
	"fmt"
	"regexp"
	"testing"
	"terraform-provider-ias/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var regexpUUID = utils.UuidRegexp

func TestResourceApplication(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_application")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceApplication("testApp", "basic-test-app", "application for testing purposes"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_application.testApp", "name", "basic-test-app"),
						resource.TestCheckResourceAttr("ias_application.testApp", "description", "application for testing purposes"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "authentication_schema.sso_type", "saml2"),
						resource.TestCheckResourceAttr("ias_application.testApp", "authentication_schema.subject_name_identifier.value", "uid"),
						resource.TestCheckResourceAttr("ias_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", "firstName"),
					),
				},
				{
					ResourceName:      "ias_application.testApp",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - application update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_application_updated")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/", user) + ResourceApplication("testApp", "test-app", "application for testing purposes"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_application.testApp", "name", "test-app"),
						resource.TestCheckResourceAttr("ias_application.testApp", "description", "application for testing purposes"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "authentication_schema.sso_type", "saml2"),
						resource.TestCheckResourceAttr("ias_application.testApp", "authentication_schema.subject_name_identifier.value", "uid"),
						resource.TestCheckResourceAttr("ias_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", "firstName"),
					),
				},
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/", user) + ResourceApplication("testApp", "test-app-updated", "application for testing purposes"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("ias_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("ias_application.testApp", "name", "test-app-updated"),
						resource.TestCheckResourceAttr("ias_application.testApp", "description", "application for testing purposes"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("ias_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("ias_application.testApp", "authentication_schema.sso_type", "saml2"),
						resource.TestCheckResourceAttr("ias_application.testApp", "authentication_schema.subject_name_identifier.value", "uid"),
						resource.TestCheckResourceAttr("ias_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", "firstName"),
					),
				},
			},
		})
	})

	t.Run("error path - app_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAppId("testApp", "this-is-not-uuid", "test-app", "application for testing purposes"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute id value must be a valid UUID, got: %s", "this-is-not-uuid")),
				},
			},
		})
	})

	t.Run("error path - name is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithoutAppName("testApp"),
					ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - parent_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithParent("testApp", "test-app", "application for testing purposes", "this-is-not-uuid"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute parent_application_id value must be a valid UUID, got:\n%s", "this-is-not-uuid")),
				},
			},
		})
	})

	t.Run("error path - sso_type needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSsoType("testApp", "test-app", "application for testing purposes", "this-is-not-a-valid-sso_type"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.sso_type value must be one of:\n\\[\"openIdConnect\" \"saml2\"], got: \"%s\"", "this-is-not-a-valid-sso_type")),
				},
			},
		})
	})

	t.Run("error path - subject_name_identifier requires sub-attributes: source, value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSubjectNameIdentifier("testApp", "test-app", "application for testing purposes", "source = \"source\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.subject_name_identifier.value\" must be\nspecified when \"authentication_schema.subject_name_identifier\" is specified"),
				},
				{
					Config:      ResourceApplicationWithSubjectNameIdentifier("testApp", "test-app", "application for testing purposes", "value = \"value\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.subject_name_identifier.source\" must be\nspecified when \"authentication_schema.subject_name_identifier\" is specified"),
				},
			},
		})
	})

	t.Run("error path - subject_name_identifier.source needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSubjectNameIdentifier("testApp", "test-app", "application for testing purposes", "source = \"this-is-not-a-valid-source\", value = \"value\""),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.subject_name_identifier.source value must be\none of: \\[\"Identity Directory\" \"Corporate Identity Provider\" \"Expression\"],\ngot: \"%s\"", "this-is-not-a-valid-source")),
				},
			},
		})
	})

	t.Run("error path - assertion_attributes requires sub-attributes: attribute_name, attribute_value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAssertionAttributes("testApp", "test-app", "application for testing purposes", "attribute_value = \"value\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.assertion_attributes\\[0].attribute_name\" must\nbe specified when \"authentication_schema.assertion_attributes\" is specified"),
				},
				{
					Config:      ResourceApplicationWithAssertionAttributes("testApp", "test-app", "application for testing purposes", "attribute_name = \"name\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.assertion_attributes\\[0].attribute_value\"\nmust be specified when \"authentication_schema.assertion_attributes\" is\nspecified"),
				},
			},
		})
	})

	t.Run("error path - advanced_assertion_attributes requires sub-attributes: source, attribute_name, attribute_value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAdvancedAssertionAttributes("testApp", "test-app", "application for testing purposes", "source = \"source\", attribute_value = \"value\""),
					ExpectError: regexp.MustCompile("Attribute\n\"authentication_schema.advanced_assertion_attributes\\[0].attribute_name\" must\nbe specified when \"authentication_schema.advanced_assertion_attributes\" is\nspecified"),
				},
				{
					Config:      ResourceApplicationWithAdvancedAssertionAttributes("testApp", "test-app", "application for testing purposes", "attribute_name = \"name\", attribute_value = \"value\""),
					ExpectError: regexp.MustCompile("Attribute \"authentication_schema.advanced_assertion_attributes\\[0].source\"\nmust be specified when \"authentication_schema.advanced_assertion_attributes\"\nis specified"),
				},
				{
					Config:      ResourceApplicationWithAdvancedAssertionAttributes("testApp", "test-app", "application for testing purposes", "source = \"source\", attribute_name = \"name\""),
					ExpectError: regexp.MustCompile("Attribute\n\"authentication_schema.advanced_assertion_attributes\\[0].attribute_value\" must\nbe specified when \"authentication_schema.advanced_assertion_attributes\" is\nspecified"),
				},
			},
		})
	})

	t.Run("error path - advanced_assertion_attributes.source needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAdvancedAssertionAttributes("testApp", "test-app", "application for testing purposes", "source = \"this-is-not-a-valid-source\", attribute_name = \"name\", attribute_value = \"value\""),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.advanced_assertion_attributes\\[0].source value\nmust be one of: \\[\"Corporate Identity Provider\" \"Expression\"], got:\n\"%s\"", "this-is-not-a-valid-source")),
				},
			},
		})
	})

	t.Run("error path - authentication_rules requires sub-attribute: identity_provider_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "user_type = \"user\""),
					ExpectError: regexp.MustCompile("Attribute\n\"authentication_schema.authentication_rules\\[0].identity_provider_id\" must be\nspecified when \"authentication_schema.authentication_rules\" is specified"),
				},
			},
		})
	})

	t.Run("error path - authentication_rules requires atleast one of the following sub-attribute: user_type, user_group, user_email_domain, ip_network_range", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "identity_provider_id = \"664c660e25cff252c5c202dc\""),
					ExpectError: regexp.MustCompile("At least one attribute out of\n\\[authentication_schema.authentication_rules\\[\\*].user_type,authentication_schema.authentication_rules\\[\\*].user_group,authentication_schema.authentication_rules\\[\\*].user_email_domain"),
				},
			},
		})
	})

	t.Run("error path - email not a valid domain", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "identity_provider_id = \"664c660e25cff252c5c202dc\", user_email_domain=\"this-is-not-email-domain\""),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.authentication_rules\\[0].user_email_domain\nvalue must be a valid Email Domain, got: %s", "this-is-not-email-domain")),
				},
			},
		})
	})

	t.Run("error path - ip address is invalid", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "identity_provider_id = \"664c660e25cff252c5c202dc\", ip_network_range = \"this-is-not-ip-address\""),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.authentication_rules\\[0].ip_network_range\nvalue must be a valid IP Address, got: %s", "this-is-not-ip-address")),
				},
			},
		})
	})
}

func ResourceApplication(resourceName string, appName string, description string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
	}
	`, resourceName, appName, description)
}

func ResourceApplicationWithParent(resourceName string, appName string, description string, parentAppId string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		parent_application_id = "%s"
	}
	`, resourceName, appName, description, parentAppId)
}

func ResourceApplicationWithAppId(resourceName string, appID string, appName string, description string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		id = "%s"
		name = "%s"
		description = "%s"
	}
	`, resourceName, appID, appName, description)
}

func ResourceApplicationWithoutAppName(resourceName string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
	}
	`, resourceName)
}

func ResourceApplicationWithSubjectNameIdentifier(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			subject_name_identifier = {
				%s
			}
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithAssertionAttributes(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			assertion_attributes = [
				{	%s	}
			]
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithAdvancedAssertionAttributes(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			advanced_assertion_attributes = [
				{	%s	}
			]
		}
	}
	`, resourceName, appName, description, subAttribute)
}

func ResourceApplicationWithSsoType(resourceName string, appName string, description string, ssoType string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			sso_type = "%s"
		}
	}
	`, resourceName, appName, description, ssoType)
}

func ResourceApplicationWithAuthenticationRules(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "ias_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			authentication_rules = [
				{
					%s
				}
			]
		}
	}
	`, resourceName, appName, description, subAttribute)
}
