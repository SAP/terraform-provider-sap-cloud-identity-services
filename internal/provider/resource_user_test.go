package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"terraform-provider-sci/internal/cli/apiObjects/users"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceUser(t *testing.T) {

	sciUser := users.User{
		UserName: "jdoe",
		Emails: []users.Email{
			{
				Type:    "work",
				Value:   "joe.doe@test.com",
				Primary: true,
			},
		},
		Name: users.Name{
			FamilyName:      "Doe",
			GivenName:       "Joe",
			HonorificPrefix: "Mr.",
		},
		DisplayName: "Joe Doe",
		UserType:    "employee",
		Active:      true,
		SAPExtension: users.SAPExtension{
			SendMail:     false,
			MailVerified: true,
			Status:       "active",
		},
	}

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUser("testUser", sciUser),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_user.testUser", "user_name", sciUser.UserName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.family_name", sciUser.Name.FamilyName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.given_name", sciUser.Name.GivenName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.honorific_prefix", sciUser.Name.HonorificPrefix),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.type", sciUser.Emails[0].Type),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.value", sciUser.Emails[0].Value),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.primary", fmt.Sprintf("%t", sciUser.Emails[0].Primary)),
						resource.TestCheckResourceAttr("sci_user.testUser", "display_name", sciUser.DisplayName),
						resource.TestCheckResourceAttr("sci_user.testUser", "user_type", sciUser.UserType),
						resource.TestCheckResourceAttr("sci_user.testUser", "active", fmt.Sprintf("%t", sciUser.Active)),
						resource.TestCheckResourceAttr("sci_user.testUser", "sap_extension_user.send_mail", fmt.Sprintf("%t", sciUser.SAPExtension.SendMail)),
						resource.TestCheckResourceAttr("sci_user.testUser", "sap_extension_user.mail_verified", fmt.Sprintf("%t", sciUser.SAPExtension.MailVerified)),
						resource.TestCheckResourceAttr("sci_user.testUser", "sap_extension_user.status", sciUser.SAPExtension.Status),
					),
				},
			},
		})
	})

	t.Run("happy path - user update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user_updated")
		defer stopQuietly(rec)

		updatedSciUser := users.User{
			UserName: "jdoe_s",
			Emails: []users.Email{
				{
					Type:    "work",
					Value:   "joe.doe@test.com",
					Primary: true,
				},
				{
					Type:    "home",
					Value:   "joe.doe.s@test.com",
					Primary: false,
				},
			},
			Name: users.Name{
				FamilyName:      "Doe S",
				GivenName:       "Joe",
				HonorificPrefix: "Mr.",
			},
			DisplayName: "Joe Doe S",
			UserType:    "customer",
			Active:      true,
			SAPExtension: users.SAPExtension{
				SendMail:     false,
				MailVerified: false,
				Status:       "active",
			},
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUser("testUser", sciUser),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_user.testUser", "user_name", sciUser.UserName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.family_name", sciUser.Name.FamilyName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.given_name", sciUser.Name.GivenName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.honorific_prefix", sciUser.Name.HonorificPrefix),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.type", sciUser.Emails[0].Type),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.value", sciUser.Emails[0].Value),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.primary", fmt.Sprintf("%t", sciUser.Emails[0].Primary)),
						resource.TestCheckResourceAttr("sci_user.testUser", "display_name", sciUser.DisplayName),
						resource.TestCheckResourceAttr("sci_user.testUser", "user_type", sciUser.UserType),
						resource.TestCheckResourceAttr("sci_user.testUser", "active", fmt.Sprintf("%t", sciUser.Active)),
						resource.TestCheckResourceAttr("sci_user.testUser", "sap_extension_user.send_mail", fmt.Sprintf("%t", sciUser.SAPExtension.SendMail)),
						resource.TestCheckResourceAttr("sci_user.testUser", "sap_extension_user.mail_verified", fmt.Sprintf("%t", sciUser.SAPExtension.MailVerified)),
						resource.TestCheckResourceAttr("sci_user.testUser", "sap_extension_user.status", sciUser.SAPExtension.Status),
					),
				},
				{
					Config: providerConfig("", user) + ResourceUser("testUser", updatedSciUser),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_user.testUser", "user_name", updatedSciUser.UserName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.family_name", updatedSciUser.Name.FamilyName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.given_name", updatedSciUser.Name.GivenName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.honorific_prefix", updatedSciUser.Name.HonorificPrefix),
						// resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.type", updatedSciUser.Emails[0].Type),
						// resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.value", updatedSciUser.Emails[0].Value),
						// resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.primary", fmt.Sprintf("%t", updatedSciUser.Emails[0].Primary)),
						// resource.TestCheckResourceAttr("sci_user.testUser", "emails.1.type", updatedSciUser.Emails[1].Type),
						// resource.TestCheckResourceAttr("sci_user.testUser", "emails.1.value", updatedSciUser.Emails[1].Value),
						// resource.TestCheckResourceAttr("sci_user.testUser", "emails.1.primary", fmt.Sprintf("%t", updatedSciUser.Emails[1].Primary)),
						resource.TestCheckResourceAttr("sci_user.testUser", "display_name", updatedSciUser.DisplayName),
						resource.TestCheckResourceAttr("sci_user.testUser", "user_type", updatedSciUser.UserType),
						resource.TestCheckResourceAttr("sci_user.testUser", "active", fmt.Sprintf("%t", updatedSciUser.Active)),
						resource.TestCheckResourceAttr("sci_user.testUser", "sap_extension_user.send_mail", fmt.Sprintf("%t", updatedSciUser.SAPExtension.SendMail)),
						resource.TestCheckResourceAttr("sci_user.testUser", "sap_extension_user.mail_verified", fmt.Sprintf("%t", updatedSciUser.SAPExtension.MailVerified)),
						resource.TestCheckResourceAttr("sci_user.testUser", "sap_extension_user.status", updatedSciUser.SAPExtension.Status),
					),
				},
			},
		})

	})

	t.Run("happy path - custom schemas", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user_with_custom_schemas")
		defer stopQuietly(rec)

		sciUser.Schemas = []string{
			"urn:ietf:params:scim:schemas:extension:sap:2.0:User",
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:test:terraform:1.0:User",
		}

		customSchemas, _ := json.Marshal(
			map[string]interface{}{
				"urn:test:terraform:1.0:User": map[string]interface{}{
					"test1": "testValue",
					"test2": false,
					"test3": map[string]interface{}{
						"test3a": 12.33,
						"test3b": 12,
					},
				},
			},
		)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUserWithCustomSchemas("testUser", sciUser, string(customSchemas)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_user.testUser", "user_name", sciUser.UserName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.family_name", sciUser.Name.FamilyName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.given_name", sciUser.Name.GivenName),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.type", sciUser.Emails[0].Type),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.value", sciUser.Emails[0].Value),
						resource.TestCheckResourceAttr("sci_user.testUser", "custom_schemas", string(customSchemas)),
					),
				},
			},
		})

	})

	t.Run("happy path - update custom schemas", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_user_with_updated_custom_schemas")
		defer stopQuietly(rec)

		sciUser.Schemas = []string{
			"urn:ietf:params:scim:schemas:extension:sap:2.0:User",
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:test:terraform:1.0:User",
		}

		customSchemas, _ := json.Marshal(
			map[string]interface{}{
				"urn:test:terraform:1.0:User": map[string]interface{}{
					"test1": "testValue",
					"test2": false,
					"test3": map[string]interface{}{
						"test3a": 12.33,
						"test3b": 12,
					},
				},
			},
		)

		newCustomSchemas, _ := json.Marshal(
			map[string]interface{}{
				"urn:test:terraform:1.0:User": map[string]interface{}{
					"test1": "newTestValue",
					"test2": true,
				},
			},
		)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceUserWithCustomSchemas("testUser", sciUser, string(customSchemas)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_user.testUser", "user_name", sciUser.UserName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.family_name", sciUser.Name.FamilyName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.given_name", sciUser.Name.GivenName),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.type", sciUser.Emails[0].Type),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.value", sciUser.Emails[0].Value),
						resource.TestCheckResourceAttr("sci_user.testUser", "custom_schemas", string(customSchemas)),
					),
				},
				{
					Config: providerConfig("", user) + ResourceUserWithCustomSchemas("testUser", sciUser, string(newCustomSchemas)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_user.testUser", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_user.testUser", "user_name", sciUser.UserName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.family_name", sciUser.Name.FamilyName),
						resource.TestCheckResourceAttr("sci_user.testUser", "name.given_name", sciUser.Name.GivenName),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.type", sciUser.Emails[0].Type),
						resource.TestCheckResourceAttr("sci_user.testUser", "emails.0.value", sciUser.Emails[0].Value),
						resource.TestCheckResourceAttr("sci_user.testUser", "custom_schemas", string(newCustomSchemas)),
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
					Config:      ResourceUserWithoutSchemas("testUser", sciUser),
					ExpectError: regexp.MustCompile("Attribute schemas set must contain at least 1 elements, got: 0"),
				},
			},
		})

	})

	t.Run("error path - username is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithoutUserName("testUser", sciUser),
					ExpectError: regexp.MustCompile("The argument \"user_name\" is required, but no definition was found."),
				},
			},
		})
	})

	t.Run("error path - emails cannot be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithoutEmails("testUser", sciUser),
					ExpectError: regexp.MustCompile("Inappropriate value for attribute \"emails\": element 0: attributes \"type\" and\n\"value\" are required."),
				},
			},
		})
	})

	t.Run("error path - user type must be a valid value", func(t *testing.T) {

		sciUser.UserType = "this-is-not-a-valid-user-type"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithUserType("testUser", sciUser),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute user_type value must be one of: \\[\"public\" \"partner\" \"customer\"\n\"external\" \"onboardee\" \"employee\"\\], got: \"%s\"", sciUser.UserType)),
				},
			},
		})
	})

	t.Run("error path - status must be a valid value", func(t *testing.T) {

		sciUser.SAPExtension.Status = "this-is-not-a-valid-status"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithStatus("testUser", sciUser),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute sap_extension_user.status value must be one of: \\[\"active\"\n\"inactive\" \"new\"\\], got: \"%s\"", sciUser.SAPExtension.Status)),
				},
			},
		})
	})

	t.Run("error path - custom schemas must be a valid json string", func(t *testing.T) {

		sciUser.Schemas = []string{
			"urn:ietf:params:scim:schemas:extension:sap:2.0:User",
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:test:terraform:1.0:User",
		}

		invalidCustomSchemas := "this-is-not-a-valid-json-string"

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      ResourceUserWithCustomSchemas("testUser", sciUser, invalidCustomSchemas),
					ExpectError: regexp.MustCompile(fmt.Sprintf("Attribute custom_schemas value must be valid json, got:\n\"%s\"", invalidCustomSchemas)),
				},
			},
		})
	})

}

func ResourceUser(resourceName string, user users.User) string {

	return fmt.Sprintf(`
	resource "sci_user" "%s"{
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
			honorific_prefix = "%s"
		}
		emails = [%s]
		display_name = "%s"
		user_type = "%s"
		active = %t
		sap_extension_user = {
			send_mail = %t
			mail_verified = %t
			status = "%s"
		}
	}
	`, resourceName, user.UserName, user.Name.FamilyName, user.Name.GivenName, user.Name.HonorificPrefix, getEmails(user.Emails), user.DisplayName, user.UserType, user.Active, user.SAPExtension.SendMail, user.SAPExtension.MailVerified, user.SAPExtension.Status)
}

func ResourceUserWithoutSchemas(resourceName string, user users.User) string {
	return fmt.Sprintf(`
	resource "sci_user" "%s"{
		schemas = []
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [%s]
	}
	`, resourceName, user.UserName, user.Name.FamilyName, user.Name.GivenName, getEmails(user.Emails))
}

func ResourceUserWithoutUserName(resourceName string, user users.User) string {
	return fmt.Sprintf(`
	resource "sci_user" "%s"{
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [%s]
	}
	`, resourceName, user.Name.FamilyName, user.Name.GivenName, getEmails(user.Emails))
}

func ResourceUserWithoutName(resourceName string, userName string, emails []users.Email) string {
	return fmt.Sprintf(`
	resource "sci_user" "%s"{
		user_name = "%s"
		name = {}
		emails = [
			{
				type = "%s"
				value = "%s"
			}
		]
	}
	`, resourceName, userName, emails[0].Type, emails[0].Value)
}

func ResourceUserWithoutEmails(resourceName string, user users.User) string {
	return fmt.Sprintf(`
	resource "sci_user" "%s"{
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [
			{}
		]
	}
	`, resourceName, user.UserName, user.Name.FamilyName, user.Name.GivenName)
}

func ResourceUserWithUserType(resourceName string, user users.User) string {
	return fmt.Sprintf(`
	resource "sci_user" "%s"{
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [%s]
		user_type = "%s"
	}
	`, resourceName, user.UserName, user.Name.FamilyName, user.Name.GivenName, getEmails(user.Emails), user.UserType)
}

func ResourceUserWithStatus(resourceName string, user users.User) string {
	return fmt.Sprintf(`
	resource "sci_user" "%s"{
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [%s]
		sap_extension_user = {
			status = "%s"
		}
	}
	`, resourceName, user.UserName, user.Name.FamilyName, user.Name.GivenName, getEmails(user.Emails), user.SAPExtension.Status)
}

func ResourceUserWithCustomSchemas(resourceName string, user users.User, customSchemas string) string {

	var schemas string
	for _, schema := range user.Schemas {
		schemas += fmt.Sprintf(`
			"%s" ,
		`, schema)
	}

	return fmt.Sprintf(`
	resource "sci_user" "%s"{
		schemas = [%s]
		user_name = "%s"
		name = {
			family_name = "%s"
			given_name = "%s"
		}
		emails = [%s]
		custom_schemas = %q
	}
	`, resourceName, schemas, user.UserName, user.Name.FamilyName, user.Name.GivenName, getEmails(user.Emails), customSchemas)
}

func getEmails(userEmails []users.Email) string {

	var emails string
	for _, email := range userEmails {
		emails += fmt.Sprintf(`
			{
				value = "%s",
				type = "%s",
				primary = "%t",
			},
		`, email.Value, email.Type, email.Primary)

	}
	return emails
}
