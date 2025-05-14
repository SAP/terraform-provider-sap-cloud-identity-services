package utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// string together the array of valid values for the attribute
func ValidValuesString(values []string) string {

	valString := "Acceptable values are : "

	for _, val := range values {
		valString += fmt.Sprintf("`%s`, ", val)
	}

	valString = valString[:len(valString)-2] + "." // remove the last comma and space and add a fullstop
	return valString
}

// string together the default schemas
func PrintDefaultSchemas(schemas []attr.Value) string {
	schemasString := ""
	for _, schema := range schemas {
		str := schema.String()
		str = str[1 : len(str)-1] // remove the quotes
		schemasString += fmt.Sprintf("\t- `%s` \n", str)
	}
	return schemasString
}
