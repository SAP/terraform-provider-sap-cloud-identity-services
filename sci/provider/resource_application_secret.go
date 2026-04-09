package provider

import (
	"context"
	"strings"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var authorizationScopeValues = []string{"manageApp", "oAuth", "readUserProfile", "manageUsers", "manageAMSPolicies"}

func newApplicationSecretResource() resource.Resource {
	return &applicationSecretResource{}
}

type applicationSecretResource struct {
	cli *cli.SciClient
}

func (r *applicationSecretResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.cli = req.ProviderData.(*cli.SciClient)
}

func (r *applicationSecretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_secret"
}

func (r *applicationSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates an API secret for a SAP Cloud Identity Services application. The secret value is available after creation and stored in state as sensitive. It cannot be retrieved again after the initial creation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the application secret.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the application this secret belongs to. Changing this value forces a new secret to be created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID of the application.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "The generated secret value. Only available after creation — not returned by subsequent API reads. Stored as sensitive in state.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hint": schema.StringAttribute{
				MarkdownDescription: "A short hint (first characters) of the secret value for identification.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Human-readable description of the secret.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"valid_to": schema.StringAttribute{
				MarkdownDescription: "Expiry date-time of the secret in UTC format (YYYY-MM-DDTHH:MM:SSZ).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					utils.ValidDateTime(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"authorization_scopes": schema.SetAttribute{
				MarkdownDescription: "API authorization scopes granted to this secret. " + utils.ValidValuesString(authorizationScopeValues),
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(authorizationScopeValues...),
					),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_apis_access": schema.BoolAttribute{
				MarkdownDescription: "If set to true, the secret grants access to all APIs regardless of the authorization_scopes.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *applicationSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan applicationSecretData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args, diags := getApplicationSecretRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.cli.ApplicationSecret.Create(ctx, plan.ApplicationId.ValueString(), args)
	if err != nil {
		resp.Diagnostics.AddError("Error creating application secret", err.Error())
		return
	}

	state, diags := applicationSecretValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ApplicationId = plan.ApplicationId

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *applicationSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config applicationSecretData
	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.cli.ApplicationSecret.GetById(ctx, config.ApplicationId.ValueString(), config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading application secret", err.Error())
		return
	}

	state, diags := applicationSecretValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ApplicationId = config.ApplicationId
	// The API does not return the secret value after creation — preserve it from prior state
	state.Secret = config.Secret
	state.ClientId = config.ClientId
	state.Hint = config.Hint

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *applicationSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state applicationSecretData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ops, diags := getApplicationSecretUpdateRequest(ctx, plan, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(ops) == 0 {
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	res, err := r.cli.ApplicationSecret.Update(ctx, state.ApplicationId.ValueString(), state.Id.ValueString(), ops)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application secret", err.Error())
		return
	}

	newState, diags := applicationSecretValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState.ApplicationId = state.ApplicationId
	// Preserve secret value — not returned by the API after creation
	newState.Secret = state.Secret
	newState.Hint = state.Hint
	newState.ClientId = state.ClientId

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *applicationSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config applicationSecretData
	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.cli.ApplicationSecret.Delete(ctx, config.ApplicationId.ValueString(), config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting application secret", err.Error())
	}
}

func (r *applicationSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ",", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: <application_id>,<secret_id>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
