package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SchemaSetValidator struct {
	DefaultSchemas types.Set
}

func (s SchemaSetValidator) Description(_ context.Context) string {
	return "Ensures the schema attribute contains all required default schemas."
}

func (s SchemaSetValidator) MarkdownDescription(_ context.Context) string {
	return "Ensures the schema attribute contains all required default schemas."
}

func (s SchemaSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {

	// check for the missing schemas
	missingSchemas, containsAllDefaultSchemas := s.CheckSchemas(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	// add an error if not all deafault schemas are present
	if !containsAllDefaultSchemas {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing schema values",
			fmt.Sprintf("Please add the schemas :\n%v", missingSchemas),
		)
	}
}

// checks if all the default schemas are present in the schema attribute
func (s SchemaSetValidator) CheckSchemas(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) ([]string, bool) {

	// if attribute is not configured, the default schemas will be added by the setdefault
	if req.ConfigValue.IsNull() {
		return nil, true
	}

	// extract config schemas
	var schemas []string
	diags := req.ConfigValue.ElementsAs(ctx, &schemas, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return nil, false
	}

	// extract default schemas
	var defaultSchemas []string
	diags = s.DefaultSchemas.ElementsAs(ctx, &defaultSchemas, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return nil, false
	}

	// create a map of default schemas
	schemasMap := make(map[string]bool)
	for _, schema := range defaultSchemas {
		schemasMap[schema] = false
	}

	// check if the config schemas are present in the default schemas
	for _, schema := range schemas {
		if _, ok := schemasMap[schema]; ok {
			schemasMap[schema] = true
		}
	}

	// check for missing schemas
	var missingSchemas []string
	for key, value := range schemasMap {
		if !value {
			missingSchemas = append(missingSchemas, key)
		}
	}

	// return missing schemas, if any
	return missingSchemas, len(missingSchemas) == 0
}

func SchemaValidator(defaultSchemas []attr.Value) validator.Set {
	// convert the default schemas to a set
	defaultSchemasSet := types.SetValueMust(types.StringType, defaultSchemas)

	return SchemaSetValidator{
		DefaultSchemas: defaultSchemasSet,
	}
}
