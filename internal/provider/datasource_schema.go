package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	// "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newSchemaDataSource() datasource.DataSource {
	return &schemaDataSource{}
}

type schemaDataSource struct {
	cli *cli.IasClient
}

func (d *schemaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (d *schemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (d *schemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{
			"id" : schema.StringAttribute{
				Required: true,
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
	}
}

func (d *schemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config schemaData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, err := d.cli.Directory.Schema.GetBySchemaId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving schema", fmt.Sprintf("%s", err))
		return
	}

	state, diags := schemaValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}