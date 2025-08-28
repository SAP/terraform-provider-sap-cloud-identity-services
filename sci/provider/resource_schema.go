package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	defaultSchemaSchemas = []attr.Value{
		types.StringValue("urn:ietf:params:scim:schemas:core:2.0:Schema"),
	}
	attributeDataTypes        = []string{"string", "boolean", "decimal", "integer", "dateTime", "binary", "reference", "complex"}
	attributeMutabilityValues = []string{"readOnly", "readWrite", "writeOnly", "immutable"}
	attributeReturnValues     = []string{"always", "never", "default", "request"}
	attributeUniquenessValues = []string{"none", "server", "global"}
)

func newSchemaResource() resource.Resource {
	return &schemaResource{}
}

type schemaResource struct {
	cli *cli.SciClient
}

func (r *schemaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.cli = req.ProviderData.(*cli.SciClient)
}

func (r *schemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *schemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a schema in the SAP Cloud Identity Services.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "A unique id by which the schema can be referenced in other entities. The ID must follow the `urn:<namespace-identifier>:<resource-type>` pattern.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A unique name for the schema",
				Required:            true,
			},
			"attributes": schema.ListNestedAttribute{
				MarkdownDescription: "The list of attribites that comprise the schema",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(20),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The attribute name. Only alphanumeric characters and underscores are allowed.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(2, 30),
								utils.ValidAttributeName(),
							},
						},
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The attribute data type. " + utils.ValidValuesString(attributeDataTypes),
							Validators: []validator.String{
								stringvalidator.OneOf(attributeDataTypes...),
							},
						},
						"mutability": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Control the Read or Write access of the attribute. " + utils.ValidValuesString(attributeMutabilityValues),
							Validators: []validator.String{
								stringvalidator.OneOf(attributeMutabilityValues...),
							},
						},
						"returned": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Configure how the attribute's value must be returned. " + utils.ValidValuesString(attributeReturnValues),
							Validators: []validator.String{
								stringvalidator.OneOf(attributeReturnValues...),
							},
						},
						"uniqueness": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Define the context in which the attribute must be unique. " + utils.ValidValuesString(attributeUniquenessValues),
							Validators: []validator.String{
								stringvalidator.OneOf(attributeUniquenessValues...),
							},
						},
						"multivalued": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Confgire if the attribute can have more than one value.",
						},
						"required": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Configure if the attribute must be mandatory or not.",
						},
						"case_exact": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Configure if the attribute must be case-sensitive or not.",
						},
						"canonical_values": schema.ListAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							MarkdownDescription: "A collection of suggested canonical values that may be used",
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
						},
						"description": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "A brief description for the attribute",
						},
					},
				},
			},
			"schemas": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				MarkdownDescription: "List of SCIM schemas to configure schemas. The attribute is configured with default values :\n" +
					utils.PrintDefaultSchemas(defaultSchemaSchemas),
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						defaultSchemaSchemas,
					),
				),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					utils.DefaultValuesChecker(defaultSchemaSchemas),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description for the schema",
			},
		},
	}
}

func (r *schemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var config schemaData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Schema.GetBySchemaId(ctx, config.Id.ValueString())
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

func (r *schemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan schemaData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	args, diags := getSchemaRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Schema.Create(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error creating schema", fmt.Sprintf("%s", err))
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

func (r *schemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Invalid Operation", "The resource \"Schema\" does not support updates.")
}

func (r *schemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var config schemaData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Schema ID is missing", "Please provide a valid ID")
		return
	}

	err := r.cli.Schema.Delete(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting schema", fmt.Sprintf("%s", err))
		return
	}
}
