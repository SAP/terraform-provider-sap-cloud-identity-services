package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var clientSecretObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":             types.StringType,
		"application_id": types.StringType,
		"client_id":      types.StringType,
		"secret":         types.StringType,
		"hint":           types.StringType,
		"description":    types.StringType,
		"valid_to":       types.StringType,
		"authorization_scopes": types.SetType{
			ElemType: types.StringType,
		},
		"all_apis_access": types.BoolType,
	},
}

func newClientSecretsDataSource() datasource.DataSource {
	return &clientSecretsDataSource{}
}

type clientSecretsDataSource struct {
	cli *cli.SciClient
}

type clientSecretsData struct {
	ApplicationId types.String `tfsdk:"application_id"`
	Values        types.List   `tfsdk:"values"`
}

func (d *clientSecretsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.cli = req.ProviderData.(*cli.SciClient)
}

func (d *clientSecretsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_secrets"
}

func (d *clientSecretsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Gets all API secrets for a SAP Cloud Identity Services application. Note: the secret value is not returned — it is only available at creation time.",
		Attributes: map[string]schema.Attribute{
			"application_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the application whose secrets to list.",
				Required:            true,
			},
			"values": schema.ListNestedAttribute{
				MarkdownDescription: "List of API secrets for the application.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier of the secret.",
							Computed:            true,
						},
						"application_id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier of the application this secret belongs to.",
							Computed:            true,
						},
						"client_id": schema.StringAttribute{
							MarkdownDescription: "Client ID of the application.",
							Computed:            true,
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "The generated secret value. Not returned by the API after creation — always null in this data source.",
							Computed:            true,
							Sensitive:           true,
						},
						"hint": schema.StringAttribute{
							MarkdownDescription: "A short hint (first characters) of the secret value for identification.",
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
				},
			},
		},
	}
}

func (d *clientSecretsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config clientSecretsData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.cli.ApplicationSecret.Get(ctx, config.ApplicationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving application secrets", fmt.Sprintf("%s", err))
		return
	}

	var secretItems []applicationSecretData
	for _, s := range res.Secrets {
		item, diags := applicationSecretValueFrom(ctx, s)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		item.ApplicationId = config.ApplicationId
		// Secret value is never returned by the read API
		item.Secret = types.StringNull()
		item.ClientId = types.StringNull()
		secretItems = append(secretItems, item)
	}

	values, diags := types.ListValueFrom(ctx, clientSecretObjType, secretItems)
	resp.Diagnostics.Append(diags...)

	config.Values = values
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
