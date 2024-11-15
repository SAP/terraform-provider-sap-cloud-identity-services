package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	cli *cli.IasClient
}

func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) { 
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) { 
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { 
	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				//MarkdownDescription:
				Validators: []validator.String{
					ValidUUID(),
				},
			},
			"schemas": schema.SetAttribute{
				ElementType: types.StringType,
				Computed: true,
				//MarkdownDescription:
			},
			"display_name": schema.StringAttribute{
				Computed: true,
				// MarkdownDescription: ,
			},
			"group_members": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Computed: true,
							// MarkdownDescription: ,
						},
						"type": schema.StringAttribute{
							Computed: true,
							// MarkdownDescription: ,
						},
					},
				},
				Computed: true,
			},
			"external_id": schema.StringAttribute{
				Computed: true,
				// MarkdownDescription: ,
			},
			"description": schema.StringAttribute{
				Computed: true,
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

	res, err := d.cli.Group.GetByGroupId(ctx, config.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving group", fmt.Sprintf("%s", err))
		return
	}

	state, diags := groupValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}