package provider

import (
	"context"
	"fmt"
	"strings"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	cli *cli.SciClient
}

func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a group from the SAP Cloud Identity services.`,
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
				//MarkdownDescription:
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Display Name of the group.",
			},
			"group_members": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Specify the members to be part of the group.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "SCIM ID of the user or the group",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: fmt.Sprintf("Type of the member added to the group. Valid Values can be one of the following : %s", strings.Join(memberTypeValues, ",")),
						},
					},
				},
			},
			"external_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique and global identifier for the given group",
			},
			"group_extension": schema.SingleNestedAttribute{
				// MarkdownDescription: ,
				Computed: true,
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

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config groupData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.Id.IsNull() {
		resp.Diagnostics.AddError("Group ID is missing", "Please provide a valid ID")
		return
	}

	res, _, err := d.cli.Group.GetByGroupId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving group", fmt.Sprintf("%s", err))
		return
	}

	state, diags := groupValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}
