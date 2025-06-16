package provider

import (
	"fmt"
	"testing"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceCorporateIdP(t *testing.T) {

	corporateIdP := corporateidps.IdentityProvider{
		ForwardAllSsoRequests: true,
		IdentityFederation: corporateidps.IdentityFederation{
			UseLocalUserStore:        true,
			AllowLocalUsersOnly:      true,
			ApplyLocalIdPAuthnChecks: true,
			RequiredGroups: []string{
				"Test Group",
			},
		},
		LoginHintConfiguration: corporateidps.LoginHintConfiguration{
			LoginHintType: "mail",
			SendMethod: "urlParam",
		},
		LogoutUrl: "https://example.com/logout",
	}

	saml2Config := corporateidps.SAML2Configuration {
		SamlMetadataUrl: "https://example.com/saml2-metadata",
		AssertionAttributes: []corporateidps.AssertionAttribute{
			{
				Name: "attr_name",
				Value: "attr_value",
			},
		},
		DigestAlgorithm: "sha1",
		IncludeScoping: true,
		DefaultNameIdFormat: "email",
		AllowCreate: "true",
	}

	// oidcCoporateIdP := corporateidps.IdentityProvider{
	// 	DisplayName: 		 "OIDC - Test Corporate IdP",
	// 	Name: 			  	 "OIDC - Test IdP",
	// 	Type: 			  	 "openIdConnect",
	// }

	t.Parallel()

	t.Run("happy path - saml2 corporate idp", func(t *testing.T) {
		
		saml2IdP := corporateIdP
		saml2IdP.DisplayName = "SAML2 - Test Corporate IdP"
		saml2IdP.Name = "SAML2 - Test IdP"
		saml2IdP.Type = "saml2"
		saml2IdP.Saml2Configuration = saml2Config

		rec, user := setupVCR(t, "fixtures/resource_corporateIdP_saml2")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + ResourceSaml2CorporateIdP("testIdP", saml2IdP),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("sci_corporate_idp.testIdP", "id", regexpUUID),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "display_name", saml2IdP.DisplayName),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "name", saml2IdP.Name),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "type", saml2IdP.Type),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "logout_url", saml2IdP.LogoutUrl),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "forward_all_sso_requests", fmt.Sprintf("%t",saml2IdP.ForwardAllSsoRequests)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.use_local_user_store", fmt.Sprintf("%t",saml2IdP.IdentityFederation.UseLocalUserStore)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.allow_local_users_only", fmt.Sprintf("%t",saml2IdP.IdentityFederation.AllowLocalUsersOnly)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.apply_local_idp_auth_and_checks", fmt.Sprintf("%t",saml2IdP.IdentityFederation.ApplyLocalIdPAuthnChecks)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "identity_federation.required_groups.0", saml2IdP.IdentityFederation.RequiredGroups[0]),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "login_hint_config.login_hint_type", saml2IdP.LoginHintConfiguration.LoginHintType),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "login_hint_config.send_method", saml2IdP.LoginHintConfiguration.SendMethod),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.saml_metadata_url", saml2IdP.Saml2Configuration.SamlMetadataUrl),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.assertion_attributes.0.name", saml2IdP.Saml2Configuration.AssertionAttributes[0].Name),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.assertion_attributes.0.value", saml2IdP.Saml2Configuration.AssertionAttributes[0].Value),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.digest_algorithm", saml2IdP.Saml2Configuration.DigestAlgorithm),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.include_scoping", fmt.Sprintf("%t", saml2IdP.Saml2Configuration.IncludeScoping)),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.name_id_format", saml2IdP.Saml2Configuration.DefaultNameIdFormat),
						resource.TestCheckResourceAttr("sci_corporate_idp.testIdP", "saml2_config.allow_create", saml2IdP.Saml2Configuration.AllowCreate),
					),
				},
			},
		})
	})
}

func ResourceSaml2CorporateIdP(resourceName string, idp corporateidps.IdentityProvider) string {

	var groups string
	for _, group := range idp.IdentityFederation.RequiredGroups {
		groups += fmt.Sprintf(`"%s",`, group)
	}

	resourceIdP := fmt.Sprintf(`
	resource "sci_corporate_idp" "%s" {
		display_name = "%s"
		name = "%s"
		type = "%s"
		logout_url = "%s"
		forward_all_sso_requests = %t
		identity_federation = {
			use_local_user_store = %t
			allow_local_users_only = %t
			apply_local_idp_auth_and_checks = %t
			required_groups = [%s]
		}
		login_hint_config = {
			login_hint_type = "%s"
			send_method = "%s"
		}
	`, resourceName, idp.DisplayName, idp.Name, idp.Type, idp.LogoutUrl, idp.ForwardAllSsoRequests, idp.IdentityFederation.UseLocalUserStore, idp.IdentityFederation.AllowLocalUsersOnly, idp.IdentityFederation.ApplyLocalIdPAuthnChecks, groups, idp.LoginHintConfiguration.LoginHintType, idp.LoginHintConfiguration.SendMethod) 

	if idp.Type == "saml2" {

		saml2Config := idp.Saml2Configuration
		
		var assertionAttributes string
		for _, attribute := range saml2Config.AssertionAttributes {
			assertionAttributes += `
				{
					name = "` + attribute.Name + `"
					value = "` + attribute.Value + `"
				},
			`
		}

		resourceIdP += fmt.Sprintf(`
			saml2_config = {
				saml_metadata_url = "%s"
				assertion_attributes = [%s]
				digest_algorithm = "%s"
				include_scoping = %t
				name_id_format = "%s"
				allow_create = "%s"
			}
		`, saml2Config.SamlMetadataUrl, assertionAttributes, saml2Config.DigestAlgorithm, saml2Config.IncludeScoping, saml2Config.DefaultNameIdFormat, saml2Config.AllowCreate)

	} else {}
	
	resourceIdP += `}`

	return resourceIdP
}