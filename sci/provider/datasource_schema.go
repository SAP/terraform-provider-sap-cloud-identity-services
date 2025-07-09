package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newSchemaDataSource() datasource.DataSource {
	return &schemaDataSource{}
}

type schemaDataSource struct {
	cli *cli.SciClient
}

func (d *schemaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *schemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (d *schemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a schema from the SAP Cloud Identity services.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id by which the schema can be referenced in other entities",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name for the schema",
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
							MarkdownDescription: "The attribute data type",
						},
						"mutability": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Control the Read or Write access of the attribute",
						},
						"returned": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Configure how the attribute's value must be returned",
						},
						"uniqueness": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Define the context in which the attribute must be unique.",
						},
						"canonical_values": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "A collection of suggested canonical values that may be used",
						},
						"multivalued": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Confgire if the attribute can have more than one value.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A brief description for the attribute",
						},
						"required": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Configure if the attribute must be mandatory or not.",
						},
						"case_exact": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Configure if the attribute must be case-sensitive or not.",
						},
					},
				},
			},
			"schemas": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				MarkdownDescription: "List of SCIM schemas to configure schemas. The attribute is configured with default values :\n" +
					utils.PrintDefaultSchemas(defaultSchemaSchemas),
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description for the schema",
			},
		},
	}
}

func (d *schemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config schemaData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := d.cli.Schema.GetBySchemaId(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving schema", fmt.Sprintf("%s", err))
		return
	}

	state, diags := schemaValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
