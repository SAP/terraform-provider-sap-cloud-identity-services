package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newGroupBaseDataSource() datasource.DataSource {
	return &groupBaseDataSource{}
}

type groupBaseDataSource struct {
	cli *cli.SciClient
}

func (d *groupBaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *groupBaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_base"
}

func (d *groupBaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a group without its member assignments from the SAP Cloud Identity Services tenant.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique ID of the group.",
				Validators: []validator.String{
					utils.ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				MarkdownDescription: "List of SCIM schemas to configure groups. The attribute is configured with default values :\n" +
					utils.PrintDefaultSchemas(defaultGroupSchemas),
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Display Name of the group.",
			},
			"group_extension": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure attributes particular to the schema `" + defaultGroupSchemas[1].String() + "`.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "Provide a unique name for the group.",
						Computed:            true,
					},
					"description": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Briefly describe the nature of the group.",
					},
				},
			},
		},
	}
}

func (d *groupBaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config groupBaseData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Group ID is missing", "Please provide a valid ID")
		return
	}

	res, _, err := d.cli.Group.GetByGroupId(ctx, config.Id.ValueString())
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
