package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ias/internal/cli"
	"terraform-provider-ias/internal/cli/apiObjects/applications"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var sourceValues = []string {"Identity Directory", "Corporate Identity Provider", "Expression"}

func newApplicationResource() resource.Resource {
	return &applicationResource{}
}

type applicationResource struct {
	cli *cli.IasClient
}

func (d *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the application",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the application",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for the application",
				Optional:            true,
			},
			"parent_application_id": schema.StringAttribute{
				MarkdownDescription: "ID of the parent, from which the application will inherit its configurations",
				Optional:			 true,
				Validators: []validator.String{
					ValidUUID(),
				},
			},
			"multi_tenant_app": schema.BoolAttribute{
				// MarkdownDescription: "Show whether the application ",
				Optional: true,
				Computed: true,
			},
			"global_account": schema.StringAttribute{
				// MarkdownDescription: "The ",
				Optional: true,
				Computed: true,
			},
			"sso_type": schema.StringAttribute{
				//MarkdownDescription:
				Optional: true,
				Computed: true,
			},
			"subjectNameIdentifier" : schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"source": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "Acceptable values: \"Identity Directory\", \"Corporate Idenity Provider\", \"Expression\"",
						Validators: []validator.String{
							stringvalidator.OneOf(sourceValues...),
						},
					},
					"value": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "If the source is Identity Directory, the only acceptable values are \"Global User ID\", \"User ID\", \"Email\", \"Display Name\", \"Login Name\", \"Employee Number\"",

					},

				},
			},
			"assertion_attributes": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"assertion_attribute_name": schema.StringAttribute{
							Required: true,
						},
						"user_attribute_name": schema.StringAttribute{
							Required: true,
						},
						"inherited": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var config applicationData

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	args, diags := getApplicationRequest(ctx, config)
	resp.Diagnostics.Append(diags...)

	id, err := r.cli.ApplicationConfiguration.Application.Create(ctx, args)

	if err != nil {
		resp.Diagnostics.AddError("Error creating application", fmt.Sprintf("%s", err))
		return
	}

	id = strings.Split(id, "/")[3]

	res, err := r.cli.ApplicationConfiguration.Application.GetByAppId(ctx, id)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application", fmt.Sprintf("%s", err))
		return
	}

	state, diags := applicationValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config applicationData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	
	res, err := r.cli.ApplicationConfiguration.Application.GetByAppId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application", fmt.Sprintf("%s", err))
		return
	}

	state, diags := applicationValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config applicationData
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	// Retrieve the current state to get the existing application ID
	var state applicationData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if state.Id.IsNull() {
		resp.Diagnostics.AddError("Application ID is missing", "Please provide a valid ID")
		return
	}

	// Update the application details
	args, diags := getApplicationRequest(ctx, config)
	resp.Diagnostics.Append(diags...)

	args.Id = state.Id.ValueString()

	err := r.cli.ApplicationConfiguration.Application.Update(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", fmt.Sprintf("%s", err))
		return
	}

	// Refresh the state with the latest data
	res, err := r.cli.ApplicationConfiguration.Application.GetByAppId(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving updated application", fmt.Sprintf("%s", err))
		return
	}

	// **Fix: Update the state with the latest data from the resource**
	updatedState, diags := applicationValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config applicationData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Application ID is missing", "Please provide a valid ID")
		return
	}

	err := r.cli.ApplicationConfiguration.Application.Delete(ctx, config.Id.ValueString())

	if err!=nil{
		resp.Diagnostics.AddError("Error deleting application", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getApplicationRequest (ctx context.Context, plan applicationData) (*applications.Application, diag.Diagnostics){

	var diagnostics  diag.Diagnostics

	args := &applications.Application{
		Name: plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		ParentApplicationId: plan.ParentApplicationId.ValueString(),
		MultiTenantApp: plan.MultiTenantApp.ValueBool(),
		GlobalAccount: plan.GlobalAccount.String(),
		AuthenticationSchema: applications.AuthenticationSchema{
			SsoType: plan.SsoType.ValueString(),
		},
	}

	if !plan.AssertionAttributes.IsNull() {

		var attributes []assertionAttributesData
		diags := plan.AssertionAttributes.ElementsAs(ctx, &attributes, true)
		diagnostics.Append(diags ...)

		for _, attribute := range attributes {
			assertionAttribute := applications.AssertionAttribute{
				AssertionAttributeName: attribute.AssertionAttributeName.ValueString(),
				UserAttributeName: attribute.UserAttributeName.ValueString(),
			}
			args.AuthenticationSchema.AssertionAttributes = append(args.AuthenticationSchema.AssertionAttributes, assertionAttribute)
		}
	}

	return args, diagnostics
}