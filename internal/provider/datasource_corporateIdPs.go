package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var saml2AssertionAttributeObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"value":       types.StringType,
	},
}