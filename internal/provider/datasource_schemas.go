package provider

import (
	"context"
	"fmt"
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

type schemasData struct{
	Values 	types.List 		`tfsdk:"values"`
}

var attributeObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name" : types.StringType,
		"type" : types.StringType,
		"multivalued" : types.BoolType,
		"description" : types.StringType,
		"required" : types.BoolType,
		"canonical_values" : types.ListType{
			ElemType: types.StringType,
		},
		"case_exact" : types.BoolType,
		"mutability" : types.StringType,
		"returned" : types.StringType,
		"uniqueness" : types.StringType,
	},
}

var schemaObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id" : types.StringType,
		"name" : types.StringType,
		"description" : types.StringType,
		"schemas" : types.SetType{
			ElemType: types.StringType,
		},
		"external_id" : types.StringType,
		"attributes" : types.ListType{
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

		Attributes: map[string]schema.Attribute{
			
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id" : schema.StringAttribute{
							Computed: true,
							//maybe add a regex
							MarkdownDescription: "A unique id by which the schema can be referenced in other entities",
						},
						"name" : schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "A unique name for the schema",
						},
						//meta
						"description" : schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "A description for the schema",
						},
						"schemas" : schema.SetAttribute{
							ElementType: types.StringType,
							Computed: true,
							//MarkDown
						},
						"external_id" : schema.StringAttribute{
							Computed: true,
							// MarkdownDescription: ,
						},
						"attributes" : schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name" : schema.StringAttribute{
										Computed: true,
										MarkdownDescription: "Name of the attribute",
									},
									"type" : schema.StringAttribute{
										Computed: true,
										MarkdownDescription: "Type of the attribute",
									},
									"multivalued" : schema.BoolAttribute{
										Computed: true,
										// MarkDownDescription
									},
									"description" : schema.StringAttribute{
										Computed: true,
										MarkdownDescription: "Description for the attribute",
									},
									"required": schema.BoolAttribute{
										Computed: true,
										MarkdownDescription: "Set a restriction on attribute, it it can be optional or not",
									},
									"canonical_values": schema.ListAttribute{
										ElementType: types.StringType,
										Computed: true,
										// MarkdownDescription: ,
									},
									"case_exact": schema.BoolAttribute{
										Computed: true,
										// MarkdownDescription: "Set a restriction on attribute",
									},
									"mutability": schema.StringAttribute{
										Computed: true,
										MarkdownDescription: "Read or Write access",
									},
									"returned": schema.StringAttribute{
										Computed: true,
										// 
									},
									"uniqueness": schema.StringAttribute{
										Computed: true,
										// MarkdownDescription: ,
									},
								},
							},
							Computed: true,
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

	res, err := d.cli.Schema.Get(ctx)

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