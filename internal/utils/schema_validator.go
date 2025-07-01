package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DefaultValuesValidator struct {
	DefaultParamValues types.Set
}

func (s DefaultValuesValidator) Description(_ context.Context) string {
	return "Ensures the attribute contains all required default values."
}

func (s DefaultValuesValidator) MarkdownDescription(_ context.Context) string {
	return "Ensures the attribute contains all required default values."
}

func (s DefaultValuesValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {

	// check for the missing values
	missingValues, containsAllDefaultValues := s.CheckDefaultValues(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	// add an error if not all deafault values are present
	if !containsAllDefaultValues {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing values",
			fmt.Sprintf("Please add the values :\n%v", missingValues),
		)
	}
}

// checks if all the default values are present in the attribute
func (s DefaultValuesValidator) CheckDefaultValues(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) ([]string, bool) {

	if req.ConfigValue.IsNull() {
		return nil, true
	}

	// extract config values
	var values []string
	diags := req.ConfigValue.ElementsAs(ctx, &values, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return nil, false
	}

	// extract default values
	var defaultValues []string
	diags = s.DefaultParamValues.ElementsAs(ctx, &defaultValues, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return nil, false
	}

	// create a map of default values
	valuesMap := make(map[string]bool)
	for _, value := range defaultValues {
		valuesMap[value] = false
	}

	// check if the config values are present in the default values
	for _, value := range values {
		if _, ok := valuesMap[value]; ok {
			valuesMap[value] = true
		}
	}

	// check for missing values
	var missingValues []string
	for key, value := range valuesMap {
		if !value {
			missingValues = append(missingValues, key)
		}
	}

	// return missing values, if any
	return missingValues, len(missingValues) == 0
}

func DefaultValuesChecker(defaultParamValues []attr.Value) validator.Set {
	// convert the default values to a set
	defaultParamsSet := types.SetValueMust(types.StringType, defaultParamValues)

	return DefaultValuesValidator{
		DefaultParamValues: defaultParamsSet,
	}
}
