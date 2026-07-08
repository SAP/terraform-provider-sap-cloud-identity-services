package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newGroupBaseResource() resource.Resource {
	return &groupBaseResource{}
}

type groupBaseResource struct {
	cli *cli.SciClient
}

func (d *groupBaseResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (r *groupBaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_base"
}

func (r *groupBaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
		Create and manage a group without member assignments in the SAP Cloud Identity Services tenant.
		
### Conflict Warning
There are 2 ways to manage members assigned to a group:
- the **sci_group** resource which manages the group and all its assignments together
- the **sci_group_base** resource in combination with **sci_group_assignment** which manages the group and individual assignments

If both the monolithic resource and the individual base/assignment resources are used against the same Group, spurious changes and conflicting state updates will occur.

		`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique ID of the group.",
				Validators: []validator.String{
					utils.ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				MarkdownDescription: "List of SCIM schemas to configure groups. The attribute is configured with default values :\n" +
					utils.PrintDefaultSchemas(defaultGroupSchemas),
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						defaultGroupSchemas,
					),
				),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					utils.DefaultValuesChecker(defaultGroupSchemas),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display Name of the group.",
				Required:            true,
			},
			"group_extension": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure attributes particular to the schema `" + defaultGroupSchemas[1].String() + "`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "Provide a unique name for the group.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
							stringplanmodifier.RequiresReplace(),
						},
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "Briefly describe the nature of the group.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *groupBaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan groupBaseData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	args, diags := getGroupBaseRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Group.Create(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error creating group", fmt.Sprintf("%s", err))
		return
	}

	state, diags := groupBaseValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupBaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var config groupBaseData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Group.GetByGroupId(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving group", fmt.Sprintf("%s", err))
		return
	}

	state, diags := groupBaseValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupBaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan groupBaseData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state groupBaseData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() {
		resp.Diagnostics.AddError("Group ID is missing", "Please provide a valid ID")
		return
	}

	args, diag := getGroupBaseUpdateRequest(ctx, plan, state)
	resp.Diagnostics.Append(diag.Errors()...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Group.Update(ctx, args, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating group", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := groupBaseValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (r *groupBaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var config groupBaseData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Group ID is missing", "Please provide a valid ID")
		return
	}

	err := r.cli.Group.Delete(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting group", fmt.Sprintf("%s", err))
		return
	}
}

func (r *groupBaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
