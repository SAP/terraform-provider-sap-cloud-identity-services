package provider

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"

	pkcs12 "software.sslmate.com/src/go-pkcs12"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New() provider.Provider {
	return NewWithClient(http.DefaultClient)
}

func NewWithClient(httpClient *http.Client) *SciProvider {
	return &SciProvider{
		httpClient: httpClient,
	}
}

type SciProvider struct {
	httpClient *http.Client
}

type SciProviderData struct {
	TenantUrl              types.String `tfsdk:"tenant_url"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	P12CertificateContent  types.String `tfsdk:"p12_certificate_content"`
	P12CertificatePassword types.String `tfsdk:"p12_certificate_password"`
}

func (p *SciProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sci"
}

func (p *SciProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The Terraform provider for SAP Cloud Identity Services enables you to automate the provisioning, management, and configuration of resources in the [SAP Cloud Identity Services](https://help.sap.com/docs/cloud-identity-services).`,
		Attributes: map[string]schema.Attribute{
			"tenant_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the SCI tenant.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"p12_certificate_content": schema.StringAttribute{
				MarkdownDescription: "Base64-encoded content of the `.p12` (PKCS#12) certificate bundle file used for x509 authentication. You can use `filebase64(\"certifiacte.p12\")` to load it.",
				Optional:            true,
				Sensitive:           true,
			},
			"p12_certificate_password": schema.StringAttribute{
				MarkdownDescription: "Password to decrypt the `.p12` certificate content.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SciProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config SciProviderData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if config.TenantUrl.IsNull() {
		resp.Diagnostics.AddError("Tenant URL missing", "Please provide a valid tenant URL.")
		return
	}

	parsedUrl, err := url.Parse(config.TenantUrl.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to parse URL", fmt.Sprintf("%s", err))
		return
	}

	var httpClient *http.Client
	var cert *tls.Certificate

	if !config.P12CertificateContent.IsNull() && !config.P12CertificatePassword.IsNull() {
		decoded, err := base64.StdEncoding.DecodeString(config.P12CertificateContent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to decode base64 content", err.Error())
			return
		}

		privateKey, leafCert, caCerts, err := pkcs12.DecodeChain(decoded, config.P12CertificatePassword.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid .p12 certificate", err.Error())
			return
		}

		chain := [][]byte{leafCert.Raw}
		for _, ca := range caCerts {
			chain = append(chain, ca.Raw)
		}

		tlsCert := tls.Certificate{
			Certificate: chain,
			PrivateKey:  privateKey,
			Leaf:        leafCert,
		}

		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{tlsCert},
			InsecureSkipVerify: false,
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
		cert = &tlsCert
	} else {
		// Fallback to default HTTP client
		httpClient = p.httpClient
	}

	client := cli.NewSciClient(cli.NewClient(httpClient, parsedUrl))

	// Use basic auth if certificate is not used
	if cert == nil {
		var username string
		if config.Username.IsNull() {
			username = os.Getenv("SCI_USERNAME")
		} else {
			username = config.Username.ValueString()
		}

		var password string
		if config.Password.IsNull() {
			password = os.Getenv("SCI_PASSWORD")
		} else {
			password = config.Password.ValueString()
		}

		client.AuthorizationToken = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SciProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newApplicationDataSource,
		newApplicationsDataSource,
		newUsersDataSource,
		newUserDataSource,
		newSchemasDataSource,
		newSchemaDataSource,
		newGroupsDataSource,
		newGroupDataSource,
	}
}

func (p *SciProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newApplicationResource,
		newUserResource,
		newSchemaResource,
		newGroupResource,
	}
}
