package provider

import (
	"context"
	"terraform-provider-ias/internal/cli/apiObjects/schemas"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type attributesData struct {
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Multivalued     types.Bool   `tfsdk:"multivalued"`
	Description     types.String `tfsdk:"description"`
	Required        types.Bool   `tfsdk:"required"`
	CanonicalValues types.List   `tfsdk:"canonical_values"`
	CaseExact       types.Bool   `tfsdk:"case_exact"`
	Mutability      types.String `tfsdk:"mutability"`
	Returned        types.String `tfsdk:"returned"`
	Uniqueness      types.String `tfsdk:"uniqueness"`
}

type schemaData struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Schemas     types.Set    `tfsdk:"schemas"`
	ExternalId  types.String `tfsdk:"external_id"`
	Attributes  types.List   `tfsdk:"attributes"`
}

func schemaValueFrom(ctx context.Context, s schemas.Schema) (schemaData, diag.Diagnostics) {
	var diagnostics, diags diag.Diagnostics

	schema := schemaData{
		Id:   types.StringValue(s.Id),
		Name: types.StringValue(s.Name),
	}

	if len(s.Description) > 0 {
		schema.Description = types.StringValue(s.Description)
	}
	if len(s.ExternalId) > 0 {
		schema.ExternalId = types.StringValue(s.ExternalId)
	}

	schema.Schemas, diags = types.SetValueFrom(ctx, types.StringType, s.Schemas)
	diagnostics.Append(diags...)

	attributes := []attributesData{}

	for _, attributeRes := range s.Attributes {
		attribute := attributesData{
			Name:        types.StringValue(attributeRes.Name),
			Type:        types.StringValue(attributeRes.Type),
			Multivalued: types.BoolValue(attributeRes.Multivalued),
			Description: types.StringValue(attributeRes.Description),
			Required:    types.BoolValue(attributeRes.Required),
			CaseExact:   types.BoolValue(attributeRes.CaseExact),
			Mutability:  types.StringValue(attributeRes.Mutability),
			Returned:    types.StringValue(attributeRes.Returned),
			Uniqueness:  types.StringValue(attributeRes.Uniqueness),
		}

		if len(attributeRes.Description) > 0 {
			attribute.Description = types.StringValue(attributeRes.Description)
		} else {
			attribute.Description = types.StringNull()
		}

		attribute.CanonicalValues, diags = types.ListValueFrom(ctx, types.StringType, attributeRes.CanonicalValues)
		diagnostics.Append(diags...)

		attributes = append(attributes, attribute)
	}

	schema.Attributes, diags = types.ListValueFrom(ctx, attributeObjType, attributes)
	diagnostics.Append(diags...)

	return schema, diagnostics
}

func schemasValueFrom(ctx context.Context, s schemas.SchemasResponse) ([]schemaData, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	schemas := []schemaData{}

	for i := 0; i < len(s.Resources); i++ {
		schema, diags := schemaValueFrom(ctx, s.Resources[i])

		schemas = append(schemas, schema)
		diagnostics.Append(diags...)
	}

	return schemas, diagnostics
}
