package applications

import "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"

type ApplicationSecret struct {
	Id                  string   `json:"id,omitempty"`
	ClientId            string   `json:"clientId,omitempty"`
	Secret              string   `json:"secret,omitempty"`
	Hint                string   `json:"hint,omitempty"`
	Description         string   `json:"description,omitempty"`
	ValidTo             string   `json:"validTo,omitempty"`
	AuthorizationScopes []string `json:"authorizationScopes,omitempty"`
	AllApisAccess       bool     `json:"allApisAccess,omitempty"`
}

type ApplicationSecretRequest struct {
	Description         string   `json:"description,omitempty"`
	ValidTo             string   `json:"validTo,omitempty"`
	AuthorizationScopes []string `json:"authorizationScopes,omitempty"`
	AllApisAccess       bool     `json:"allApisAccess,omitempty"`
}

type ApplicationSecretPatchRequestBody struct {
	Operations []generic.PatchRequest `json:"operations"`
}

type ApplicationSecretsListResponse struct {
	Secrets []ApplicationSecret `json:"secrets"`
}
