package utils

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Certificate validator, checks that the attribute is a valid PEM string
type dateTimeValidator struct {
}

func (v dateTimeValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v dateTimeValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid date-time string in the UTC format YYYY-MM-DDTHH:MM:SSZ"
}

func (v dateTimeValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if _, err := time.Parse(time.RFC3339, value); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

func ValidDateTime() validator.String {
	return dateTimeValidator{}
}
