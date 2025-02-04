package schemas

import (
	"terraform-provider-ias/internal/cli/apiObjects/users"
)

type Attribute struct {
	Name            string   `json:"name,omitempty"`
	Type            string   `json:"type,omitempty"`
	Multivalued     bool     `json:"multivalued,omitempty"`
	Description     string   `json:"description,omitempty"`
	Required        bool     `json:"required,omitempty"`
	CanonicalValues []string `json:"canonicalValues,omitempty"`
	CaseExact       bool     `json:"caseExact,omitempty"`
	Mutability      string   `json:"mutability,omitempty"`
	Returned        string   `json:"returned,omitempty"`
	Uniqueness      string   `json:"uniqueness,omitempty"`
	ReferenceTypes  []string `json:"referenceTypes,omitempty"`
}

type Schema struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Meta        users.Meta  `json:"meta,omitempty"`
	Schemas     []string    `json:"schemas"`
	ExternalId  string      `json:"externalId,omitempty"`
	Attributes  []Attribute `json:"attributes"`
}

type SchemasResponse struct {
	Resources    []Schema `json:"Resources,omitempty"`
	Schemas      []string `json:"schemas,omitempty"`
	TotalResults int      `json:"totalResults,omitempty"`
	//rest of th response object
}
