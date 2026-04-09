package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newClientSecretDataSource() datasource.DataSource {
	return &clientSecretDataSource{}
}

type clientSecretDataSource struct {
	cli *cli.SciClient
}

func (d *clientSecretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *clientSecretDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_secret"
}

func (d *clientSecretDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Gets a single API secret for a SAP Cloud Identity Services application. Note: the secret value is not returned — it is only available at creation time.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the secret.",
				Required:            true,
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the application this secret belongs to.",
				Required:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID of the application.",
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Client ID of the application.",
				Computed:            true,
			},
			"hint": schema.StringAttribute{
				MarkdownDescription: "A short hint (last characters) of the secret value for identification.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Human-readable description of the secret.",
				Computed:            true,
			},
			"valid_to": schema.StringAttribute{
				MarkdownDescription: "Expiry date-time of the secret in UTC format (YYYY-MM-DDTHH:MM:SSZ).",
				Computed:            true,
			},
			"authorization_scopes": schema.SetAttribute{
				MarkdownDescription: "API authorization scopes granted to this secret.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"all_apis_access": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether this secret has access to all APIs.",
				Computed:            true,
			},
		},
	}
}

func (d *clientSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config applicationSecretData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.cli.ApplicationSecret.GetById(ctx, config.ApplicationId.ValueString(), config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application secret", fmt.Sprintf("%s", err))
		return
	}

	state, diags := applicationSecretValueFrom(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ApplicationId = config.ApplicationId
	// Secret value is never returned by the read API — omit it from data source state
	state.Secret = types.StringNull()
	state.ClientId = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
