package utils

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var UuidRegexp = regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
var AttributeNameRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
var IPRegexp = regexp.MustCompile(`^$|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/([0-9]|[1-2][0-9]|3[0-2]))$`)
var EmailDomainRegexp = regexp.MustCompile(`^$|^(((\*|([a-zA-Z0-9_\-]{1,63}))\.)(?:[a-zA-Z0-9_\-]{1,63}\.)*(?:[a-zA-Z]{2,})|((?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))|(localhost))$`)

// Checks that the String held in the attribute is a valid UUID
func ValidUUID() validator.String {
	return stringvalidator.RegexMatches(UuidRegexp, "value must be a valid UUID")
}

// Checks that the String held in the attribute is a valid s\attribute name
func ValidAttributeName() validator.String {
	return stringvalidator.RegexMatches(AttributeNameRegexp, "value must be a valid name. Must start with an alphabet and should contain only alphanumeric characters and underscores")
}

// Checks that the String held in the attribute is a valid IP Address
func ValidIPAddress() validator.String {
	return stringvalidator.RegexMatches(IPRegexp, "value must be a valid IP Address")
}

// Checks that the String held in the attribute is a valid Email Domain
func ValidEmailDomain() validator.String {
	return stringvalidator.RegexMatches(EmailDomainRegexp, "value must be a valid Email Domain")
}
