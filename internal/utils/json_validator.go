package utils

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// JSON validator, checks that the attribute is a valid JSON string
type jsonValidator struct {
}

func (v jsonValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v jsonValidator) MarkdownDescription(_ context.Context) string {
	return "value must be valid json"
}

func (v jsonValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	if json.Valid([]byte(value.ValueString())) {
		return
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}

func ValidJSON() validator.String {
	return jsonValidator{}
}
