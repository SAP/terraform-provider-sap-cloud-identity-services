package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var UuidRegexp 				= regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
var AttributeNameRegexp   = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

// Checks that the String held in the attribute is a valid UUID
func ValidUUID() validator.String {
	return stringvalidator.RegexMatches(UuidRegexp, "value must be a valid UUID")
}

// Checks that the String held in the attribute is a valid s\attribute name
func ValidAttributeName() validator.String {
	return stringvalidator.RegexMatches(AttributeNameRegexp, "value must be a valid name. Must start with an alphabet and should contain only alphanumeric characters and underscores")
}
