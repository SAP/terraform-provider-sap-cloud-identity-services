package provider

import (
	"fmt"
	"regexp"
	"terraform-provider-sci/internal/cli/apiObjects/applications"
	"terraform-provider-sci/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var regexpUUID = utils.UuidRegexp

func TestResourceApplication(t *testing.T) {

	application := applications.Application{
		Name:        "basic-test-app",
		Description: "application for testing purposes",
		AuthenticationSchema: applications.AuthenticationSchema{
			SsoType:                       "openIdConnect",
			SubjectNameIdentifier:         "mail",
			SubjectNameIdentifierFunction: "lowerCase",
			// AssertionAttributes: &[]applications.AssertionAttribute{
			// 	{
			// 		AssertionAttributeName: "param1",
			// 		UserAttributeName:      "firstName",
			// 	},
			// 	{
			// 		AssertionAttributeName: "param2",
			// 		UserAttributeName:      "mail",
			// 	},
			// },
			AdvancedAssertionAttributes: []applications.AdvancedAssertionAttribute{
				{
					AttributeName:  "adv_param1",
					AttributeValue: "value1",
				},
				{
					AttributeName:  "adv_param2",
					AttributeValue: "value2",
				},
			},
			DefaultAuthenticatingIdpId: "664c660e25cff252c5c202dc",
			ConditionalAuthentication: []applications.AuthenicationRule{
				{
					UserType:           "employee",
					UserEmailDomain:    "gmail.com",
					IpNetworkRange:     "10.0.0.1/8",
					IdentityProviderId: "664c660e25cff252c5c202dc",
				},
			},
		},
	}

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_application")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceApplication("testApp", application),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", application.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", application.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("sci_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", application.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier.value", application.AuthenticationSchema.SubjectNameIdentifier),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier_function", application.AuthenticationSchema.SubjectNameIdentifierFunction),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_name", application.AuthenticationSchema.AssertionAttributes[0].AssertionAttributeName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", application.AuthenticationSchema.AssertionAttributes[0].UserAttributeName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.1.attribute_name", application.AuthenticationSchema.AssertionAttributes[1].AssertionAttributeName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.1.attribute_value", application.AuthenticationSchema.AssertionAttributes[1].UserAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_name", application.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_value", application.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_name", application.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_value", application.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.default_authenticating_idp", application.AuthenticationSchema.DefaultAuthenticatingIdpId),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.authentication_rules.0.user_type", application.AuthenticationSchema.ConditionalAuthentication[0].UserType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.authentication_rules.0.user_email_domain", application.AuthenticationSchema.ConditionalAuthentication[0].UserEmailDomain),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.authentication_rules.0.ip_network_range", application.AuthenticationSchema.ConditionalAuthentication[0].IpNetworkRange),
					),
				},
				{
					ResourceName:      "sci_application.testApp",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - application update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_application_updated")
		defer stopQuietly(rec)

		updatedApplication := applications.Application{
			Name:        "test-app-updated",
			Description: "application for testing purposes",
			AuthenticationSchema: applications.AuthenticationSchema{
				SsoType:                       "saml2",
				SubjectNameIdentifier:         "userUuid",
				SubjectNameIdentifierFunction: "upperCase",
				// AssertionAttributes: &[]applications.AssertionAttribute{
				// 	{
				// 		AssertionAttributeName: "param1",
				// 		UserAttributeName:      "lastName",
				// 	},
				// },
				AdvancedAssertionAttributes: []applications.AdvancedAssertionAttribute{
					{
						AttributeName:  "adv_param1",
						AttributeValue: "updated_value1",
					},
					{
						AttributeName:  "adv_param2",
						AttributeValue: "value2",
					},
					{
						AttributeName:  "adv_param3",
						AttributeValue: "value3",
					},
				},
				DefaultAuthenticatingIdpId: "664c660e25cff252c5c202dc",
				ConditionalAuthentication: []applications.AuthenicationRule{
					{
						UserType:           "customer",
						UserEmailDomain:    "sap.com",
						IpNetworkRange:     "192.168.1.1/24",
						IdentityProviderId: "664c660e25cff252c5c202dc",
					},
				},
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/", user) + ResourceApplication("testApp", application),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", application.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", application.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("sci_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", application.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier.value", application.AuthenticationSchema.SubjectNameIdentifier),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier_function", application.AuthenticationSchema.SubjectNameIdentifierFunction),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_name", application.AuthenticationSchema.AssertionAttributes[0].AssertionAttributeName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", application.AuthenticationSchema.AssertionAttributes[0].UserAttributeName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.1.attribute_name", application.AuthenticationSchema.AssertionAttributes[1].AssertionAttributeName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.1.attribute_value", application.AuthenticationSchema.AssertionAttributes[1].UserAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_name", application.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_value", application.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_name", application.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_value", application.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.default_authenticating_idp", application.AuthenticationSchema.DefaultAuthenticatingIdpId),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.authentication_rules.0.user_type", application.AuthenticationSchema.ConditionalAuthentication[0].UserType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.authentication_rules.0.user_email_domain", application.AuthenticationSchema.ConditionalAuthentication[0].UserEmailDomain),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.authentication_rules.0.ip_network_range", application.AuthenticationSchema.ConditionalAuthentication[0].IpNetworkRange),
					),
				},
				{
					Config: providerConfig("https://iasprovidertestblr.accounts400.ondemand.com/", user) + ResourceApplication("testApp", updatedApplication),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_application.testApp", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_application.testApp", "name", updatedApplication.Name),
						resource.TestCheckResourceAttr("sci_application.testApp", "description", updatedApplication.Description),
						resource.TestCheckResourceAttr("sci_application.testApp", "multi_tenant_app", "false"),
						resource.TestCheckResourceAttr("sci_application.testApp", "global_account", "unknown"),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.sso_type", updatedApplication.AuthenticationSchema.SsoType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier.value", updatedApplication.AuthenticationSchema.SubjectNameIdentifier),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier_function", updatedApplication.AuthenticationSchema.SubjectNameIdentifierFunction),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_name", updatedApplication.AuthenticationSchema.AssertionAttributes[0].AssertionAttributeName),
						// resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.0.attribute_value", updatedApplication.AuthenticationSchema.AssertionAttributes[0].UserAttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_name", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.0.attribute_value", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[0].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_name", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.1.attribute_value", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[1].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.2.attribute_name", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[2].AttributeName),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.advanced_assertion_attributes.2.attribute_value", updatedApplication.AuthenticationSchema.AdvancedAssertionAttributes[2].AttributeValue),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.default_authenticating_idp", updatedApplication.AuthenticationSchema.DefaultAuthenticatingIdpId),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.authentication_rules.0.user_type", updatedApplication.AuthenticationSchema.ConditionalAuthentication[0].UserType),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.authentication_rules.0.user_email_domain", updatedApplication.AuthenticationSchema.ConditionalAuthentication[0].UserEmailDomain),
						resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.authentication_rules.0.ip_network_range", updatedApplication.AuthenticationSchema.ConditionalAuthentication[0].IpNetworkRange),
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

	t.Run("error path - subject_name_identifier_function needs to be a valid value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithSubjectNameIdentifierFunction("testApp", "test-app", "application for testing purposes", "invalid-function-name"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.subject_name_identifier_function value must\nbe one of: \\[\"none\" \"upperCase\" \"lowerCase\"], got: \"%s\"", "invalid-function-name")),
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

	t.Run("error path - conditional_authentication requires sub-attribute: identity_provider_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "user_type = \"user\""),
					ExpectError: regexp.MustCompile("Attribute\n\"authentication_schema.conditional_authentication\\[0].identity_provider_id\" must be\nspecified when \"authentication_schema.conditional_authentication\" is specified"),
				},
			},
		})
	})

	t.Run("error path - conditional_authentication requires atleast one of the following sub-attribute: user_type, user_group, user_email_domain, ip_network_range", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceApplicationWithAuthenticationRules("testApp", "test-app", "application for testing purposes", "identity_provider_id = \"664c660e25cff252c5c202dc\""),
					ExpectError: regexp.MustCompile("At least one attribute out of\n\\[authentication_schema.conditional_authentication\\[\\*].user_type,authentication_schema.conditional_authentication\\[\\*].user_group,authentication_schema.conditional_authentication\\[\\*].user_email_domain"),
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
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.conditional_authentication\\[0].user_email_domain\nvalue must be a valid Email Domain, got: %s", "this-is-not-email-domain")),
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
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute authentication_schema.conditional_authentication\\[0].ip_network_range\nvalue must be a valid IP Address, got: %s", "this-is-not-ip-address")),
				},
			},
		})
	})
}

func ResourceApplication(resourceName string, app applications.Application) string {

	var assertionAttributes string
	// for _, attribute := range *app.AuthenticationSchema.AssertionAttributes {
	// 	assertionAttributes += fmt.Sprintf(`
	// 			{
	// 				attribute_name = "%s"
	// 				attribute_value = "%s"
	// 			},`, attribute.AssertionAttributeName, attribute.UserAttributeName)
	// }

	var advancedAssertionAttributes string
	for _, attribute := range app.AuthenticationSchema.AdvancedAssertionAttributes {
		advancedAssertionAttributes += fmt.Sprintf(`
				{
					source = "Corporate Identity Provider",
					attribute_name = "%s",
					attribute_value = "%s",

				},`, attribute.AttributeName, attribute.AttributeValue)
	}

	var authenticationRules string
	for _, rule := range app.AuthenticationSchema.ConditionalAuthentication {
		authenticationRules += fmt.Sprintf(`{
					user_type = "%s"
					user_email_domain = "%s"
					ip_network_range = "%s"
					identity_provider_id = "%s"
					},`, rule.UserType, rule.UserEmailDomain, rule.IpNetworkRange, rule.IdentityProviderId)
	}

	return fmt.Sprintf(
		`resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			sso_type = "%s"
			subject_name_identifier = {
				source = "Identity Directory"
				value = "%s"
			}
			subject_name_identifier_function = "%s"
			assertion_attributes = [%s]
			advanced_assertion_attributes = [%s]
			default_authenticating_idp = "%s"
			conditional_authentication = [%s]
		}
	}`, resourceName, app.Name, app.Description, app.AuthenticationSchema.SsoType, app.AuthenticationSchema.SubjectNameIdentifier, app.AuthenticationSchema.SubjectNameIdentifierFunction, assertionAttributes, advancedAssertionAttributes, app.AuthenticationSchema.DefaultAuthenticatingIdpId, authenticationRules)
}

func ResourceApplicationWithParent(resourceName string, appName string, description string, parentAppId string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		parent_application_id = "%s"
	}
	`, resourceName, appName, description, parentAppId)
}

func ResourceApplicationWithAppId(resourceName string, appID string, appName string, description string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		id = "%s"
		name = "%s"
		description = "%s"
	}
	`, resourceName, appID, appName, description)
}

func ResourceApplicationWithoutAppName(resourceName string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
	}
	`, resourceName)
}

func ResourceApplicationWithSubjectNameIdentifier(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
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

func ResourceApplicationWithSubjectNameIdentifierFunction(resourceName string, appName string, description string, functionName string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			subject_name_identifier_function = "%s"
		}
	}
	`, resourceName, appName, description, functionName)
}

func ResourceApplicationWithAssertionAttributes(resourceName string, appName string, description string, subAttribute string) string {
	return fmt.Sprintf(`
	resource "sci_application" "%s" {
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
	resource "sci_application" "%s" {
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
	resource "sci_application" "%s" {
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
	resource "sci_application" "%s" {
		name = "%s"
		description = "%s"
		authentication_schema = {
			conditional_authentication = [
				{
					%s
				}
			]
		}
	}
	`, resourceName, appName, description, subAttribute)
}
