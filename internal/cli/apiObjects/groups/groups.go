package groups

import "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/users"

type GroupExtension struct {
	// AdditionalId string `json:"additionalId,omitempty"`
	Name        string `json:"name,omitempty" tfsdk:"name"`
	Description string `json:"description,omitempty" tfsdk:"description"`
}

type GroupMember struct {
	Value string `json:"value,omitempty" tfsdk:"value"`
	Type  string `json:"type,omitempty" tfsdk:"type"`
	//ref
}

type Group struct {
	Id             string          `json:"id,omitempty"`
	Meta           users.Meta      `json:"meta,omitempty"`
	Schemas        []string        `json:"schemas,omitempty"`
	DisplayName    string          `json:"displayName,omitempty"`
	GroupMembers   []GroupMember   `json:"members,omitempty"`
	GroupExtension *GroupExtension `json:"urn:sap:cloud:scim:schemas:extension:custom:2.0:Group,omitempty"`
}

type GroupsResponse struct {
	Resources    []Group  `json:"Resources"`
	Schemas      []string `json:"schemas,omitempty"`
	TotalResults int      `json:"totalResults,omitempty"`
	ItemsPerPage int      `json:"itemsPerPage,omitempty"`
	//startIndex, startId, nextId, nextCursor
}
