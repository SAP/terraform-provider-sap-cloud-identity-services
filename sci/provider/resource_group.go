package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	defaultGroupSchemas = []attr.Value{
		types.StringValue("urn:ietf:params:scim:schemas:core:2.0:Group"),
		types.StringValue("urn:sap:cloud:scim:schemas:extension:custom:2.0:Group"),
	}
	memberTypeValues = []string{"User", "Group"}
)

func newGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	cli *cli.SciClient
}

func (d *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
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
			"group_members": schema.SetNestedAttribute{
				MarkdownDescription: "Specify the members to be part of the group.",
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.AlsoRequires(
						path.MatchRoot("group_members").AtAnySetValue().AtName("value"),
					),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "SCIM ID of the user or the group",
							Validators: []validator.String{
								utils.ValidUUID(),
							},
						},
						"type": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Type of the member added to the group. " + utils.ValidValuesString(memberTypeValues),
							Validators: []validator.String{
								stringvalidator.OneOf(memberTypeValues...),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
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

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan groupData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	args, diags := r.GetGroupRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Group.Create(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", fmt.Sprintf("%s", err))
		return
	}

	state, diags := groupValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var config groupData
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := r.cli.Group.GetByGroupId(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", fmt.Sprintf("%s", err))
		return
	}

	state, diags := groupValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan groupData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state groupData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() {
		resp.Diagnostics.AddError("Group ID is missing", "Please provide a valid ID")
		return
	}

	args, diags := r.GetGroupRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	args.Id = state.Id.ValueString()

	res, _, err := r.cli.Group.Update(ctx, args)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := groupValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var config groupData
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
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	args := &groups.Group{
		Schemas:     schemas,
		DisplayName: plan.DisplayName.ValueString(),
	}

	if !plan.GroupMembers.IsNull() {

		var members []memberData
		diags = plan.GroupMembers.ElementsAs(ctx, &members, true)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		// the mapping is done manually in order to carry out the member validation
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

		var groupExtension groups.GroupExtension
		diags = plan.GroupExtension.As(ctx, &groupExtension, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		args.GroupExtension = &groupExtension
	}

	return args, diagnostics
}

func validateMembers(ctx context.Context, client *cli.SciClient, member string) error {

	// do a GET call for both the users and groups to check if the member exists
	_, _, userErr := client.User.GetByUserId(ctx, member)
	_, _, groupErr := client.Group.GetByGroupId(ctx, member)

	if userErr != nil && groupErr != nil {
		return fmt.Errorf("member %s is not found", member)
	}

	return nil
}
