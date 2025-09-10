package provider

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	pkcs12 "software.sslmate.com/src/go-pkcs12"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	basicAuthConflicts = []path.Expression{
		path.MatchRoot("client_id"),
		path.MatchRoot("client_secret"),
		path.MatchRoot("p12_certificate_content"),
		path.MatchRoot("p12_certificate_password"),
	}
	oauthConflicts = []path.Expression{
		path.MatchRoot("username"),
		path.MatchRoot("password"),
		path.MatchRoot("p12_certificate_content"),
		path.MatchRoot("p12_certificate_password"),
	}
	x509Conflicts = []path.Expression{
		path.MatchRoot("username"),
		path.MatchRoot("password"),
		path.MatchRoot("client_id"),
		path.MatchRoot("client_secret"),
	}
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
	ClientID               types.String `tfsdk:"client_id"`
	ClientSecret           types.String `tfsdk:"client_secret"`
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
			// Basic Authentication
			"username": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Your user name for Basic Authentication.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(basicAuthConflicts...),
					stringvalidator.AlsoRequires(path.MatchRoot("password")),
				},
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Your password for Basic Authentication.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(basicAuthConflicts...),
					stringvalidator.AlsoRequires(path.MatchRoot("username")),
				},
			},

			// OAuth2 Client Credentials
			"client_id": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The client ID for OAuth2 authentication.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(oauthConflicts...),
					stringvalidator.AlsoRequires(path.MatchRoot("client_secret")),
				},
			},
			"client_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The client secret for OAuth2 authentication.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(oauthConflicts...),
					stringvalidator.AlsoRequires(path.MatchRoot("client_id")),
				},
			},

			// X.509 Certificate Auth
			"p12_certificate_content": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Base64-encoded content of the `.p12` (PKCS#12) certificate bundle file used for x509 authentication. For example you can use `filebase64(\"certifiacte.p12\")` to load the file content, But any source that provides a valid .p12 certificate base64 string is accepted.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(x509Conflicts...),
				},
			},
			"p12_certificate_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password to decrypt the `.p12` certificate content.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(x509Conflicts...),
					stringvalidator.AlsoRequires(path.MatchRoot("p12_certificate_content")),
				},
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

	// OAuth2 Authentication using client_id and client_secret
	clientID := config.ClientID.ValueString()
	clientSecret := config.ClientSecret.ValueString()

	if clientID == "" {
		clientID = os.Getenv("SCI_CLIENT_ID")
	}

	if clientSecret == "" {
		clientSecret = os.Getenv("SCI_CLIENT_SECRET")
	}

	// X.509 Certificate Authentication
	p12CertificatePassword := config.P12CertificatePassword.ValueString()
	if p12CertificatePassword == "" {
		p12CertificatePassword = os.Getenv("SCI_P12_CERTIFICATE_PASSWORD")
	}

	p12CertificateContent := config.P12CertificateContent.ValueString()

	// Basic Auth (username + password)
	username := config.Username.ValueString()
	password := config.Password.ValueString()

	if username == "" {
		username = os.Getenv("SCI_USERNAME")
	}

	if password == "" {
		password = os.Getenv("SCI_PASSWORD")
	}

	client := cli.NewSciClient(cli.NewClient(p.httpClient, parsedUrl))

	switch {
	case len(clientID) != 0 && len(clientSecret) != 0:
		// OAuth2 authentication
		token, err := fetchOAuthToken(p.httpClient, parsedUrl.String(), clientID, clientSecret)
		if err != nil {
			resp.Diagnostics.AddError("OAuth2 Authentication Failed", err.Error())
			return
		}
		client.AuthorizationToken = "Bearer " + token

	case len(p12CertificateContent) != 0 && len(p12CertificatePassword) != 0:
		// X.509 authentication will be handled below
		decoded, err := base64.StdEncoding.DecodeString(p12CertificateContent)
		if err != nil {
			resp.Diagnostics.AddError("Failed to decode base64 content", err.Error())
			return
		}

		privateKey, leafCert, caCerts, err := pkcs12.DecodeChain(decoded, p12CertificatePassword)
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

		httpClient := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}

		client = cli.NewSciClient(cli.NewClient(httpClient, parsedUrl))

	case len(username) != 0 && len(password) != 0:
		// Basic authentication will be handled below
		client.AuthorizationToken = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	default:
		incompleteCreds, err := checkIncompleteCredentials(username, password, clientID, clientSecret, p12CertificateContent, p12CertificatePassword)

		if incompleteCreds {
			resp.Diagnostics.AddError("Incomplete Authentication Credentials", err)
			return
		}

		resp.Diagnostics.AddError("Authentication Details Missing", "Please provide either : \n- client_id and client_secret for OAuth2 Authentication \n- p12_certificate_content and p12_certificate_password for X.509 Authentication \n- username and password for Basic Authentication")
		return
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
		newCorporateIdPDataSource,
		newCorporateIdPsDataSource,
	}
}

func (p *SciProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newApplicationResource,
		newUserResource,
		newSchemaResource,
		newGroupResource,
		newCorporateIdPResource,
	}
}

func fetchOAuthToken(httpClient *http.Client, tenantURL, clientID, clientSecret string) (string, error) {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     strings.TrimSuffix(tenantURL, "/") + "/oauth2/token",
		AuthStyle:    oauth2.AuthStyleInParams,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	token, err := config.Token(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve token: %w", err)
	}

	if token.AccessToken == "" {
		return "", fmt.Errorf("empty access token received")
	}

	return token.AccessToken, nil
}

func checkIncompleteCredentials(username, password, clientID, clientSecret, p12CertificateContent, p12CertificatePassword string) (bool, string) {

	switch {
	case len(clientID) != 0 || len(clientSecret) != 0:
		return true, "Please provide the required OAuth Credentials : Client ID and Client Secret"
	case len(p12CertificateContent) != 0 || len(p12CertificatePassword) != 0:
		return true, "Please provide the required X.509 Authentication Credentials : P12 Certificate and P12 Certificate Password"
	case len(username) != 0 || len(password) != 0:
		return true, "Please provide the required Basic Authentication Credentials : Username and Password"
	default:
		return false, ""
	}
}
