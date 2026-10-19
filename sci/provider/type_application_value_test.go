package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"

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
	assert.False(t, diags.HasError(), "value from: %v", diags)

	var schema authenticationSchemaData
	diags = application.AuthenticationSchema.As(ctx, &schema, basetypes.ObjectAsOptions{})
	assert.False(t, diags.HasError(), "authentication schema: %v", diags)

	// must match a configured empty list, otherwise apply fails with "inconsistent result"
	assert.False(t, schema.AssertionAttributes.IsNull(), "assertion_attributes must not be null")
	assert.Empty(t, schema.AssertionAttributes.Elements())
}
