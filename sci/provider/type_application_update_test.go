package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	require.False(t, diags.HasError(), "state: %v", diags)

	plan, diags := applicationValueFrom(ctx, saml2TestApplication(
		[]applications.GrantType{applications.AUTHORIZATION_CODE_GrantType, applications.REFRESH_TOKEN_GrantType},
		"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"))
	require.False(t, diags.HasError(), "plan: %v", diags)

	reqs, diags := getApplicationUpdateRequest(ctx, plan, state)
	require.False(t, diags.HasError(), "update request: %v", diags)

	// one patch per changed field, each with the path of its own nesting level
	const authSchemaPath = "/urn:sap:identity:application:schemas:extension:sci:1.0:Authentication"
	require.Len(t, reqs, 2)

	assert.Equal(t, "replace", reqs[0].Op)
	assert.Equal(t, authSchemaPath+"/openIdConnectConfiguration/restrictedGrantTypes", reqs[0].Path)
	assert.ElementsMatch(t, []string{
		string(applications.AUTHORIZATION_CODE_GrantType),
		string(applications.REFRESH_TOKEN_GrantType),
	}, reqs[0].Value)

	assert.Equal(t, "replace", reqs[1].Op)
	assert.Equal(t, authSchemaPath+"/saml2Configuration/defaultNameIdFormat", reqs[1].Path)
	assert.Equal(t, "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress", reqs[1].Value)
}
