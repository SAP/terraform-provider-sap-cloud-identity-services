package utils

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// checks that when OIDC is configured for the corporate IDP, the type of the IDP is "openIdConnect"
type typeValidator struct {
	typeExpr    path.Expression
	validValues []string
}

func (v typeValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v typeValidator) MarkdownDescription(_ context.Context) string {
	return ": value of attribute \"type\" must be modified to match the IDP configuration provided. " + ValidValuesString(v.validValues)
}

func (v typeValidator) ValidateObject(ctx context.Context, request validator.ObjectRequest, response *validator.ObjectResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	// get the path for attribute type from the expression
	typePath, _ := request.Config.PathMatches(ctx, v.typeExpr)

	// get the value of the attribute type from the path
	var typeVal attr.Value
	_ = request.Config.GetAttribute(ctx, typePath[0], &typeVal)

	if typeVal.IsNull() || typeVal.IsUnknown() {
		return
	}

	val, ok := typeVal.(types.String)
	if !ok {
		return
	}
	rawVal := val.ValueString() // safely extract the raw string value


	validValFound := false

	// check value of type is one of the valid values
	validValFound = slices.Contains(v.validValues, rawVal)

	if !validValFound {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			rawVal,
		))
	}

}

func ValidType(typeExpr path.Expression, validValues []string) validator.Object {
	return typeValidator{
		typeExpr:    typeExpr,
		validValues: validValues,
	}
}
