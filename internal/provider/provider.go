package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"terraform-provider-ias/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New() provider.Provider {
	return &IasProvider{}
}

type IasProvider struct {
}

type IasProviderData struct {
	TenantUrl types.String `tfsdk:"tenant_url"`
}

func (p *IasProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ias"
}

func (p *IasProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the IAS tenant",
				Required:            true,
			},
		},
	}
}

func (p *IasProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var config IasProviderData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.TenantUrl.IsNull() {
		resp.Diagnostics.AddError("Tenant URL missing", "Please provide a valid tenant URL")
		return
	}

	pasrsedUrl, err := url.Parse(config.TenantUrl.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Unable to parse URL", fmt.Sprintf("%s", err))
		return
	}

	client := cli.NewIasClient(cli.NewClient(http.DefaultClient, pasrsedUrl))

	username := os.Getenv("ias_username")
	password := os.Getenv("ias_password")

	client.AuthorizationToken = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *IasProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newApplicationDataSource,
	}
}

func (p *IasProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newApplicationResource,
	}
}
