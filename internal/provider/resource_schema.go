package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"
	"terraform-provider-ias/internal/cli/apiObjects/schemas"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newSchemaResource() resource.Resource {
	return &schemaResource{}
}

type schemaResource struct {
	cli *cli.IasClient
}

func (r *schemaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.cli = req.ProviderData.(*cli.IasClient)
}

func (r *schemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *schemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{
			"id" : schema.StringAttribute{
				Required: true,
				//maybe add a regex
				MarkdownDescription: "A unique id by which the schema can be referenced in other entities",
			},
			"name" : schema.StringAttribute{
				Required: true,
				MarkdownDescription: "A unique name for the schema",
			},
			"attributes" : schema.ListNestedAttribute{
				//cap at 20
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name" : schema.StringAttribute{
							Required: true,
							MarkdownDescription: "Name of the attribute",
						},
						"type" : schema.StringAttribute{
							Required: true,
							MarkdownDescription: "Type of the attribute",
						},
						"mutability": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "Read or Write access",
						},
						"returned": schema.StringAttribute{
							Required: true,
							// 
						},
						"uniqueness": schema.StringAttribute{
							Required: true,
							// MarkdownDescription: ,
						},
						"canonical_values": schema.ListAttribute{
							ElementType: types.StringType,
							Optional: true,
							// MarkdownDescription: ,
						},
						"multivalued" : schema.BoolAttribute{
							Optional: true,
							Computed: true,
							// MarkDownDescription
						},
						"description" : schema.StringAttribute{
							Optional: true,
							Computed: true,
							MarkdownDescription: "Description for the attribute",
						},
						"required": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							MarkdownDescription: "Set a restriction on attribute, it it can be optional or not",
						},
						"case_exact": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							// MarkdownDescription: "Set a restriction on attribute",
						},
					},
				},
				Required: true,
			},
			"schemas" : schema.SetAttribute{
				ElementType: types.StringType,
				Required: true,
				//MarkDown
			},
			//meta
			"description" : schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "A description for the schema",
			},
			"external_id" : schema.StringAttribute{
				Optional: true,
				Computed: true,
				// MarkdownDescription: ,
			},
		},
	}
}

func (r *schemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var config schemaData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, err := r.cli.Schema.GetBySchemaId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving schema", fmt.Sprintf("%s", err))
		return
	}

	state, diags := schemaValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *schemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan schemaData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	args, diags := getSchemaRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)

	res, err := r.cli.Schema.Create(ctx, args)

	if err!=nil {
		resp.Diagnostics.AddError("Error creating schema", fmt.Sprintf("%s",err))
		return
	}

	state, diags := schemaValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *schemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//This resource does not suppport updates
}

func (r *schemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { 
	
	var config schemaData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Schema ID is missing", "Please provide a valid ID")
		return
	}

	err := r.cli.Schema.Delete(ctx, config.Id.ValueString())

	if err!=nil {
		resp.Diagnostics.AddError("Error deleting schema", fmt.Sprintf("%s", err))
		return
	}
}


func getSchemaRequest(ctx context.Context, plan schemaData) (*schemas.Schema, diag.Diagnostics) {
	
	var diagnostics  diag.Diagnostics

	var schemaList []string
	diags := plan.Schemas.ElementsAs(ctx, &schemaList, true)
	diagnostics.Append(diags...)

	var attributes []attributesData
	diags = plan.Attributes.ElementsAs(ctx, &attributes, true)
	diagnostics.Append(diags...)

	args := &schemas.Schema {
		Id: plan.Id.ValueString(),
		Name: plan.Name.ValueString(),
		Schemas: schemaList,
	}

	if !plan.Description.IsNull() {
		args.Description = plan.Description.ValueString()
	}

	if !plan.ExternalId.IsNull() {
		args.ExternalId = plan.ExternalId.ValueString()
	}

	for _, attribute := range attributes{
		schemaAttribute := schemas.Attribute{
			Name: attribute.Name.ValueString(),
			Type: attribute.Type.ValueString(),
			Mutability: attribute.Mutability.ValueString(),
			Returned: attribute.Returned.ValueString(),
			Uniqueness: attribute.Uniqueness.ValueString(),
			Multivalued: attribute.Multivalued.ValueBool(),
			Description: attribute.Description.ValueString(),
			Required: attribute.Required.ValueBool(),
			CaseExact: attribute.CaseExact.ValueBool(),
		}

		if !attribute.CanonicalValues.IsNull() {
			var canonicalValues []string
			diags := attribute.CanonicalValues.ElementsAs(ctx, &canonicalValues, true)
			diagnostics.Append(diags...)
			schemaAttribute.CanonicalValues = canonicalValues
		}
	
		args.Attributes = append(args.Attributes, schemaAttribute)
	}

	return args, diagnostics
}