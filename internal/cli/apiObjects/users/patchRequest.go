package users

import "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"

type PatchRequestBody struct {
	Schemas    []string               `json:"schemas"`
	Operations []generic.PatchRequest `json:"Operations"`
}
