package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
)

func saml2TestApplication(grantTypes []applications.GrantType, nameIdFormat string) applications.Application {
	return applications.Application{
		Id:   "test-id",
		Name: "test-app",
		AuthenticationSchema: &applications.AuthenticationSchema{
			SsoType:               "saml2",
			SubjectNameIdentifier: "mail",
			OidcConfig: &applications.OidcConfig{
				RestrictedGrantTypes: grantTypes,
			},
			Saml2Configuration: &applications.SamlConfiguration{
				DefaultNameIdFormat: nameIdFormat,
			},
		},
	}
}

func TestGetApplicationUpdateRequest_OidcAndSaml2Changed(t *testing.T) {
	ctx := context.Background()

	state, diags := applicationValueFrom(ctx, saml2TestApplication(
		[]applications.GrantType{applications.AUTHORIZATION_CODE_GrantType},
		"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"))
	assert.False(t, diags.HasError(), "state: %v", diags)

	plan, diags := applicationValueFrom(ctx, saml2TestApplication(
		[]applications.GrantType{applications.AUTHORIZATION_CODE_GrantType, applications.REFRESH_TOKEN_GrantType},
		"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"))
	assert.False(t, diags.HasError(), "plan: %v", diags)

	reqs, diags := getApplicationUpdateRequest(ctx, plan, state)
	assert.False(t, diags.HasError(), "update request: %v", diags)

	var oidcPatched, saml2Patched bool
	for _, req := range reqs {
		oidcPatched = oidcPatched || strings.Contains(req.Path, "openIdConnectConfiguration")
		saml2Patched = saml2Patched || strings.Contains(req.Path, "saml2Configuration")
	}
	assert.True(t, oidcPatched, "expected a patch for the OIDC configuration")
	assert.True(t, saml2Patched, "expected a patch for the SAML 2.0 configuration")
}
