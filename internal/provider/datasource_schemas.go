package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	// "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newSchemasDataSource() datasource.DataSource {
	return &schemasDataSource{}
}

type schemasDataSource struct {
	cli *cli.IasClient
}

type schemasData struct {
	Values types.List `tfsdk:"values"`
}

var attributeObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"type":        types.StringType,
		"multivalued": types.BoolType,
		"description": types.StringType,
		"required":    types.BoolType,
		"canonical_values": types.ListType{
			ElemType: types.StringType,
		},
		"case_exact": types.BoolType,
		"mutability": types.StringType,
		"returned":   types.StringType,
		"uniqueness": types.StringType,
	},
}

var schemaObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"schemas": types.SetType{
			ElemType: types.StringType,
		},
		"external_id": types.StringType,
		"attributes": types.ListType{
			ElemType: attributeObjType,
		},
	},
}

func (d *schemasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (d *schemasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schemas"
}

func (d *schemasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a list of schemas from the SAP Cloud Identity services.`,
		Attributes: map[string]schema.Attribute{

			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "A unique id by which the schema can be referenced in other entities",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "A unique name for the schema",
							Computed:            true,
						},
						"attributes": schema.ListNestedAttribute{
							MarkdownDescription: "The list of attribites that comprise the schema",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The attribute name. Only alphanumeric characters and underscores are allowed.",
									},
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: fmt.Sprintf("The attribute data type. Valid values include : %s", strings.Join(attributeDataTypes, ",")),
									},
									"mutability": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: fmt.Sprintf("Control the Read or Write access of the attribute. Valid values include : %s", strings.Join(attributeMutabilityValues, ",")),
									},
									"returned": schema.StringAttribute{
										Computed: true,
										//description must be enhanced
										MarkdownDescription: fmt.Sprintf("Valid values include : %s", strings.Join(attributeReturnValues, ",")),
									},
									"uniqueness": schema.StringAttribute{
										Computed: true,
										// description must be enhanced
										MarkdownDescription: fmt.Sprintf("Define the context in which the attribute must be unique. Valid values include : %s", strings.Join(attributeReturnValues, ",")),
									},
									"canonical_values": schema.ListAttribute{
										ElementType:         types.StringType,
										Computed:            true,
										MarkdownDescription: "A collection of suggested canonical values that may be used",
									},
									"multivalued": schema.BoolAttribute{
										Computed: true,
										// MarkDownDescription
									},
									"description": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "A brief description for the attribute",
									},
									"required": schema.BoolAttribute{
										Computed: true,
										//enhance description
										MarkdownDescription: "Set a restriction on whether the attribute may be mandatory or not",
									},
									"case_exact": schema.BoolAttribute{
										Computed: true,
										//enhance description
										MarkdownDescription: "Set a restriction on whether the attribute may be case-sensitive or not",
									},
								},
							},
						},
						"schemas": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							//MarkdownDescription
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A description for the schema",
						},
						"external_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique and global identifier for the given schema",
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *schemasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config schemasData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, _, err := d.cli.Schema.Get(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving schemas", fmt.Sprintf("%s", err))
		return
	}

	resSchemas, diags := schemasValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	config.Values, diags = types.ListValueFrom(ctx, schemaObjType, resSchemas)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
