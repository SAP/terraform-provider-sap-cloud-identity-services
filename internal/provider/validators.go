package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var UuidRegexp   		= regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
var IPRegexp     	 	= regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])$`)
var EmailDomainRegexp 	= regexp.MustCompile(`^$|^(((\*|([a-zA-Z0-9_\-]{1,63}))\.)(?:[a-zA-Z0-9_\-]{1,63}\.)*(?:[a-zA-Z]{2,})|((?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))|(localhost))$`)

// Checks that the String held in the attribute is a valid UUID
func ValidUUID() validator.String {
	return stringvalidator.RegexMatches(UuidRegexp, "value must be a valid UUID")
}

// Checks that the String held in the attribute is a valid IP Address
func ValidIPAddress() validator.String {
	return stringvalidator.RegexMatches(IPRegexp, "value must be a valid IP Address")
}

// Checks that the String held in the attribute is a valid Email Domain
func ValidEmailDomain() validator.String {
	return stringvalidator.RegexMatches(EmailDomainRegexp, "value must be a valid Email Domain")
}


