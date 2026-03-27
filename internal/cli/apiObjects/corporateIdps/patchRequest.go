package corporateidps

import "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"

type PatchRequestBody struct {
	Operations []generic.PatchRequest `json:"operations"`
}
