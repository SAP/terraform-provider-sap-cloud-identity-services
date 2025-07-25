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

	valString = valString[:len(valString)-2] // remove the last comma and space
	return valString
}

// string together the default schemas
func PrintDefaultSchemas(schemas []attr.Value) string {
	schemasString := ""
	for _, schema := range schemas {
		str := schema.String()

		// check and remove beginning and ending quotes
		if str[0] == '"' {
			str = str[1:]
		}
		if str[len(str)-1] == '"' {
			str = str[:len(str)-1]
		}

		schemasString += fmt.Sprintf("\t- `%s` \n", str)
	}
	return schemasString
}
