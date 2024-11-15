package provider

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newApplicationsDataSource() datasource.DataSource {
	return &applicationsDataSource{}
}

type applicationsDataSource struct {
	cli *cli.IasClient
}

type applicationsData struct{
	Id 		types.String	`tfsdk:"id"`
	Values 	types.List 		`tfsdk:"values"`
}

var appObjType = types.ObjectType {
	AttrTypes: map[string]attr.Type{
		"id": types.StringType,
		"name": types.StringType,
		"description": types.StringType,
		"parent_application_id": types.StringType,
		"multi_tenant_app": types.BoolType,
		"global_account": types.StringType,
	},
}


func (d *applicationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.cli = req.ProviderData.(*cli.IasClient)
}

func (d *applicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

func (d *applicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed : true,
			},

			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								ValidUUID(),
							},
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description for the application",
							Computed: true,
						},
						"parent_application_id": schema.StringAttribute{
							MarkdownDescription: "ID of the parent, from which the application will inherit its configurations",
							Computed: true,
						},
						"multi_tenant_app": schema.BoolAttribute{
							// MarkdownDescription: "Show whether the application ",
							Computed: true,
						},
						"global_account": schema.StringAttribute{
							// MarkdownDescription: "The ",
							Computed: true,
						},
					},
				},
				Computed: true,
			},

		},
	}
}

func (d *applicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config applicationsData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	res, err := d.cli.Application.Get(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application", fmt.Sprintf("%s", err))
		return
	}

	config.Id = types.StringValue("dummy")

	resApps := applicationsValueFrom(ctx, res)

	config.Values, diags = types.ListValueFrom(ctx, appObjType, resApps)
	resp.Diagnostics.Append(diags...)
	
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
