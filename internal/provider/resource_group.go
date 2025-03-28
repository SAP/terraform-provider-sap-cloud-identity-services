package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ias/internal/cli"
	"terraform-provider-ias/internal/cli/apiObjects/groups"
	"terraform-provider-ias/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var defaultGroupSchemas = []attr.Value{
	types.StringValue("urn:ietf:params:scim:schemas:core:2.0:Group"),
	types.StringValue("urn:sap:cloud:scim:schemas:extension:custom:2.0:Group"),
}

var memberTypeValues = []string{"User", "Group"}

func newGroupResource() resource.Resource {
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
		MarkdownDescription: `Creates a group in the SAP Cloud Identity Services.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique ID of the group.",
				Validators: []validator.String{
					utils.ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				// MarkdownDescription: ,
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
					utils.SchemaValidator(defaultGroupSchemas),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display Name of the group.",
				Required:            true,
			},
			"group_members": schema.ListNestedAttribute{
				MarkdownDescription: "Specify the members to be part of the group.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "SCIM ID of the user or the group",
							Validators: []validator.String{
								utils.ValidUUID(),
							},
						},
						"type": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: fmt.Sprintf("Type of the member added to the group. Valid Values can be one of the following : %s", strings.Join(memberTypeValues, ",")),
							Validators: []validator.String{
								stringvalidator.OneOf(memberTypeValues...),
							},
						},
					},
				},
			},
			"external_id": schema.StringAttribute{
				MarkdownDescription: "Unique and global identifier for the given group",
				Computed:            true,
			},
			"group_extension": schema.SingleNestedAttribute{
				// MarkdownDescription: ,
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "Provide a unique name for the group.",
						Optional:            true,
						Computed:            true,
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "Briefly describe the nature of the group.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan groupData

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	args, diags := r.GetGroupRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)

	res, _, err := r.cli.Group.Create(ctx, args)
	if resp.Diagnostics.HasError() {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error creating user", fmt.Sprintf("%s", err))
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

	res, _, err := r.cli.Group.GetByGroupId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", fmt.Sprintf("%s", err))
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

	args, diags := r.GetGroupRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)

	args.Id = state.Id.ValueString()

	res, _, err := r.cli.Group.Update(ctx, args)

	if err != nil {
		resp.Diagnostics.AddError("Error updating application", fmt.Sprintf("%s", err))
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
		resp.Diagnostics.AddError("Error deleting user", fmt.Sprintf("%s", err))
		return
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *groupResource) GetGroupRequest(ctx context.Context, plan groupData) (*groups.Group, diag.Diagnostics) {

	var diagnostics diag.Diagnostics

	var schemas []string
	diags := plan.Schemas.ElementsAs(ctx, &schemas, true)
	diagnostics.Append(diags...)

	args := &groups.Group{
		Schemas:     schemas,
		DisplayName: plan.DisplayName.ValueString(),
	}

	if !plan.GroupMembers.IsNull() {

		var members []memberData
		diags = plan.GroupMembers.ElementsAs(ctx, &members, true)
		diagnostics.Append(diags...)

		for _, member := range members {

			// validate the member as a valid user or group as the API does not handle this
			err := validateMembers(ctx, r.cli, member.Value.ValueString())
			if err != nil {
				diagnostics.AddError(
					fmt.Sprintf("%s", err),
					"please provide a valid member UUID",
				)
				return nil, diagnostics
			}

			groupMember := groups.GroupMember{
				Value: member.Value.ValueString(),
			}

			if !member.Type.IsNull() {
				groupMember.Type = member.Type.ValueString()
			}

			args.GroupMembers = append(args.GroupMembers, groupMember)
		}
	}

	if !plan.GroupExtension.IsNull() && !plan.GroupExtension.IsUnknown() {

		var groupExtension groupExtensionData
		diags = plan.GroupExtension.As(ctx, &groupExtension, basetypes.ObjectAsOptions{})
		diagnostics.Append(diags...)

		args.GroupExtension.Name = groupExtension.Name.ValueString()
		args.GroupExtension.Description = groupExtension.Description.ValueString()
	}

	return args, diagnostics
}

func validateMembers(ctx context.Context, client *cli.IasClient, member string) error {

	// do a GET call for both the users and groups to check if the member exists
	_, _, userErr := client.User.GetByUserId(ctx, member)
	_, _, groupErr := client.Group.GetByGroupId(ctx, member)

	if userErr != nil && groupErr != nil {
		return fmt.Errorf("member %s is not found", member)
	}

	return nil
}
