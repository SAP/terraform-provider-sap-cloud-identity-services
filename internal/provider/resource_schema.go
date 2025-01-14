package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ias/internal/cli"
	"terraform-provider-ias/internal/cli/apiObjects/schemas"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var attributeDataTypes 			 = []string{"string", "boolean", "decimal", "integer", "dateTime", "binary", "reference", "complex"}
var attributeMutabilityValues    = []string{"readOnly", "readWrite", "writeOnly", "immutable"}
var attributeReturnValues 		 = []string{"always", "never", "default", "request"}
var attributeUniquenessValues    = []string{"none", "server", "global"}

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
				MarkdownDescription: "A unique id by which the schema can be referenced in other entities",
				Required: true,
				//add a regex check
			},
			"name" : schema.StringAttribute{
				MarkdownDescription: "A unique name for the schema",
				Required: true,
			},
			"attributes" : schema.ListNestedAttribute{
				MarkdownDescription: "The list of attribites that comprise the schema",
				Required: true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(20),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name" : schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The attribute name. Only alphanumeric characters and underscores are allowed.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(2,20),
								ValidAttributeName(),
							},
						},
						"type" : schema.StringAttribute{
							Required: true,
							MarkdownDescription: fmt.Sprintf("The attribute data type. Valid values include : %s", strings.Join(attributeDataTypes, ",")),
							Validators: []validator.String{
								stringvalidator.OneOf(attributeDataTypes...),
							},
						},
						"mutability": schema.StringAttribute{
							Required: true,
							MarkdownDescription: fmt.Sprintf("Control the Read or Write access of the attribute. Valid values include : %s", strings.Join(attributeMutabilityValues,",")),
							Validators: []validator.String{
								stringvalidator.OneOf(attributeMutabilityValues...),
							},
						},
						"returned": schema.StringAttribute{
							Required: true,
							//description must be enhanced
							MarkdownDescription: fmt.Sprintf("Valid values include : %s", strings.Join(attributeReturnValues,",")),
							Validators: []validator.String{
								stringvalidator.OneOf(attributeReturnValues...),
							},
						},
						"uniqueness": schema.StringAttribute{
							Required: true,
							// description must be enhanced
							MarkdownDescription: fmt.Sprintf("Define the context in which the attribute must be unique. Valid values include : %s", strings.Join(attributeUniquenessValues,",")),
							Validators: []validator.String{
								stringvalidator.OneOf(attributeUniquenessValues...),
							},
						},
						"canonical_values": schema.ListAttribute{
							ElementType: types.StringType,
							Optional: true,
							MarkdownDescription: "A collection of suggested canonical values that may be used",
						},
						"multivalued" : schema.BoolAttribute{
							Optional: true,
							Computed: true,
							// MarkDownDescription
						},
						"description" : schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "A brief description for the attribute",
						},
						"required": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							//enhance description
							MarkdownDescription: "Set a restriction on whether the attribute may be mandatory or not",
						},
						"case_exact": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							//enhance description
							MarkdownDescription: "Set a restriction on whether the attribute may be case-sensitive or not",
						},
					},
				},
			},
			"schemas" : schema.SetAttribute{
				ElementType: types.StringType,
				Optional: true,
				Computed: true,
				//MarkdownDescription
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("urn:ietf:params:scim:schemas:core:2.0:Schema"),
						},
					),
				),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"description" : schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "A description for the schema",
			},
			"external_id" : schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Unique and global identifier for the given schema",
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
	resp.Diagnostics.AddError("Invalid Operation","The resource \"Schema\" does not support updates.")
	return
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

	args.Attributes = []schemas.Attribute{}
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