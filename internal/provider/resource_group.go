package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"
	"terraform-provider-ias/internal/cli/apiObjects/groups"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newGroupResource() resource.Resource{
	return &groupResource{}
}

type groupResource struct {
	cli *cli.IasClient
}

func (d *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) { 
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) { 
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { 
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				// MarkdownDescription: ,
				Validators: []validator.String{
					ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				Required: true,
				ElementType: types.StringType,
				// MarkdownDescription: ,
			},
			"display_name": schema.StringAttribute{
				Required: true,
				// MarkdownDescription: ,
			},
			"group_members": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "SCIM ID of the user or the group",
							//see if any check can be placed to validate the ID
						},
						"type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							MarkdownDescription: "Type of the member added to the group",
							//check to only add specific type
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				//MarkdownDescription
			},
			"external_id": schema.StringAttribute{
				Computed: true,
				// MarkdownDescription: ,
			},
		},
	}
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { 

	var plan groupData

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	args, diags := getGroupRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)

	res, err := r.cli.Group.Create(ctx, args)

	if err != nil {
		resp.Diagnostics.AddError("Error creating user", fmt.Sprintf("%s",err))
		return
	}

	state, diags := groupValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { 

	var config groupData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, err := r.cli.Group.GetByGroupId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", fmt.Sprintf("%s",err))
		return
	}

	state, diags := groupValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { 

	var plan groupData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	var state groupData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if state.Id.IsNull() {
		resp.Diagnostics.AddError("Group ID is missing", "Please provide a valid ID")
		return
	}

	args, diags := getGroupRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)

	args.Id = state.Id.ValueString()

	res, err := r.cli.Group.Update(ctx, args)

	if err != nil {
		resp.Diagnostics.AddError("Error updating application", fmt.Sprintf("%s",err))
		return
	}

	updatedState, diags := groupValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { 

	var config groupData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Group ID is missing", "Please provide a valid ID")
		return 
	}

	err := r.cli.Group.Delete(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", fmt.Sprintf("%s",err))
		return
	}	
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getGroupRequest(ctx context.Context, plan groupData) (*groups.Group, diag.Diagnostics) {

	var diagnostics diag.Diagnostics

	var schemas []string
	diags := plan.Schemas.ElementsAs(ctx, &schemas, true)
	diagnostics.Append(diags...)

	var members []memberData
	diags = plan.GroupMembers.ElementsAs(ctx, &members, true)
	diagnostics.Append(diags...)

	args := &groups.Group{
		Schemas: schemas,
		DisplayName: plan.DisplayName.ValueString(),
	}

	if !plan.Description.IsNull() {
		args.GroupExtension.Description = plan.Description.ValueString()
	}

	for _, member := range members {
		groupMember := groups.GroupMember{
			Value: member.Value.ValueString(),
		}

		if !member.Type.IsNull() {
			groupMember.Type = member.Type.ValueString()
		}

		args.GroupMembers = append(args.GroupMembers, groupMember)
	}

	return args, diagnostics
}