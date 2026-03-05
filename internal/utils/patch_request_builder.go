package utils

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
)

var replaceOperation = "replace"

func GetPatchRequest(attrName string, path string, value any, argsType reflect.Type) (generic.PatchRequest, diag.Diagnostics) {

	tag, diags := GetAttributeTag(attrName, argsType)
	if diags.HasError() {
		return generic.PatchRequest{}, diags
	}

	tag = fmt.Sprintf("/%s", tag)

	if path != "" {
		tag = fmt.Sprintf("/%s%s", path, tag)
	}

	return generic.PatchRequest{
		Op:    replaceOperation,
		Path:  tag,
		Value: value,
	}, nil
}

func GetAttributeTag(attrName string, argsType reflect.Type) (string, diag.Diagnostics) {

	var diags diag.Diagnostics
	arg, ok := argsType.FieldByName(attrName)
	if !ok {
		diags.AddError("error fetching attribute", fmt.Sprintf("field '%s' not found in type", attrName))
		return "", diags
	}
	tag := arg.Tag.Get("json")
	if tag == "" {
		diags.AddError("error fetching json tag", fmt.Sprintf("field '%s' has no json tag", attrName))
		return "", diags
	}

	return tag, nil

}