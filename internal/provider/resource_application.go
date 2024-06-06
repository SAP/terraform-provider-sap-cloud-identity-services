package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the application",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for the application",
				Optional:            true,
			},
			"parentApplicationId": schema.StringAttribute{
				MarkdownDescription: "ID of the parent, from which the application will inherit its configurations",
				Optional:			 true,
			},
			"multiTenantApp": schema.BoolAttribute{
				// MarkdownDescription: "Show whether the application ",
				Optional: true,
				Computed: true,
			},
			"globalAccount": schema.StringAttribute{
				// MarkdownDescription: "The ",
				Optional: true,
				Computed: true,
			},
			
		},
	}
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var config applicationData

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	args := &cli.ApplicationCreateInput{}

	if !config.Id.IsUnknown() && !config.Id.IsNull() {
		args.Id = config.Id.ValueString()
	}

	if !config.Name.IsUnknown() {
		args.Name = config.Name.ValueString()
	}

	if !config.Description.IsUnknown() {
		args.Description = config.Description.ValueString()
	}
	id, err := r.cli.ApplicationConfiguration.Application.Create(ctx, args)

	if err != nil {
		resp.Diagnostics.AddError("Error creating application", fmt.Sprintf("%s", err))
		return
	}

	id = strings.Split(id, "/")[3]
	config.Id = types.StringValue(id)
	diags = resp.State.Set(ctx, &config)

	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
