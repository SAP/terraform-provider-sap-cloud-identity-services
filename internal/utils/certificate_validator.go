package utils

import (
	"context"
	"encoding/pem"
	"log"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Certificate validator, checks that the attribute is a valid PEM string
type certificateValidator struct {
}

func (v certificateValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v certificateValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid PEM string in the format -----BEGIN CERTIFICATE-----\\n<certificate-content>\\n-----END CERTIFICATE-----"
}

func (v certificateValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	log.Default().Println("Value of the certificate " + value)

	if decodedString, _ := pem.Decode([]byte(value)); decodedString == nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

func ValidCertificate() validator.String {
	return certificateValidator{}
}
