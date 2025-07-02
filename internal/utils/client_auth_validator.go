package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// checks that according to what oidc_config.token_endpoint_auth_method is configured for the corporate IDP, the oidc_config.client_secret parameter may be required or not
type clientAuthMethodValidator struct {
	typeExpr    path.Expression
	authMethods []string
}

func (v clientAuthMethodValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v clientAuthMethodValidator) MarkdownDescription(_ context.Context) string {
	return "validate whether the parameter oidc_config.client_secret is required depending on the value of oidc_config.token_endpoint_auth_method."
}

func (v clientAuthMethodValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if !request.ConfigValue.IsNull() && !request.ConfigValue.IsUnknown() {
		return
	}

	// get the path for attribute type from the expression
	typePath, _ := request.Config.PathMatches(ctx, v.typeExpr)

	// get the value of the attribute type from the path
	var typeVal attr.Value
	_ = request.Config.GetAttribute(ctx, typePath[0], &typeVal)

	if typeVal.IsNull() || typeVal.IsUnknown() {
		// We assume token_endpoint_auth_method to take default value "clientSecretPost" for which the client_secret is required
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"attribute \"oidc_config.client_secret\" must be specified",
			request.ConfigValue.String(),
		))
		return
	}

	val := typeVal.String()

	// if the value of token_endpoint_auth_method is either "privateKeyJwt" or "privateKeyJwtRfc7523", the client_secret is not required
	if val == fmt.Sprintf("\"%s\"", v.authMethods[0]) || val == fmt.Sprintf("\"%s\"", v.authMethods[1]) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"attribute \"oidc_config.client_secret\" must be specified when oidc_config.token_endpoint_auth_method is set to one of: "+v.authMethods[0]+", "+v.authMethods[1],
			request.ConfigValue.String(),
		))
		return
	}

}

func CheckClientAuthMethod(typeExpr path.Expression, authMethods []string) validator.String {
	return clientAuthMethodValidator{
		typeExpr:    typeExpr,
		authMethods: authMethods,
	}
}
