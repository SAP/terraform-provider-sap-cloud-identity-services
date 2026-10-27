package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
)

func TestApplicationValueFrom_NoAssertionAttributes(t *testing.T) {
	ctx := context.Background()

	// the API omits assertionAttributes when the application has none
	app := applications.Application{
		Id:   "test-id",
		Name: "test-app",
		AuthenticationSchema: &applications.AuthenticationSchema{
			SsoType:               "saml2",
			SubjectNameIdentifier: "mail",
		},
	}

	application, diags := applicationValueFrom(ctx, app)
	require.False(t, diags.HasError(), "value from: %v", diags)

	var schema authenticationSchemaData
	diags = application.AuthenticationSchema.As(ctx, &schema, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError(), "authentication schema: %v", diags)

	// a null list would not match a configured empty list and apply would fail with "inconsistent result"
	assert.False(t, schema.AssertionAttributes.IsNull(), "assertion_attributes must not be null")
	assert.Empty(t, schema.AssertionAttributes.Elements())
}
