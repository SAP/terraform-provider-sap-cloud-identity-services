package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"

	"os"
	"regexp"
	"strings"

	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

type User struct {
	Username string
	Password string
}

var redactedTestUser = User{
	Username: "test-user",
	Password: "test-password",
}

type RoundTripperFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func providerConfig(_ string, testUser User) string {
	return fmt.Sprintf(`
	provider "sci" {
	  tenant_url = "https://iasprovidertestblr.accounts400.ondemand.com/"
	  username   = "%s"
	  password   = "%s"
	}
	`, testUser.Username, testUser.Password)
}

func getTestProviders(httpClient *http.Client) map[string]func() (tfprotov6.ProviderServer, error) {
	sciProvider := NewWithClient(httpClient)
	return map[string]func() (tfprotov6.ProviderServer, error){
		"sci": providerserver.NewProtocol6WithError(sciProvider),
	}
}

func setupVCR(t *testing.T, cassetteName string) (*recorder.Recorder, User) {
	t.Helper()
	mode := recorder.ModeRecordOnce
	if testRecord, _ := strconv.ParseBool(os.Getenv("TEST_RECORD")); testRecord {
		mode = recorder.ModeRecordOnly
	}
	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       cassetteName,
		Mode:               mode,
		SkipRequestLatency: true,
		RealTransport:      http.DefaultTransport,
	})
	if err != nil {
		t.Fatal(err)
	}

	testUser := redactedTestUser
	if rec.IsRecording() {
		t.Logf("Recording new interactions in '%s'", cassetteName)
		testUser.Username = os.Getenv("SCI_USERNAME")
		testUser.Password = os.Getenv("SCI_PASSWORD")
		if testUser.Username == "" || testUser.Password == "" {
			t.Fatal("SCI_USERNAME and SCI_PASSWORD must be set for recording")
		}
	} else {
		t.Logf("Replaying cassette '%s'", cassetteName)
	}

	rec.SetMatcher(requestMatcher(t))
	rec.AddHook(redactCredentials(), recorder.BeforeSaveHook)

	return rec, testUser
}

func requestMatcher(t *testing.T) cassette.MatcherFunc {
	return func(r *http.Request, i cassette.Request) bool {
		if r.Method != i.Method || r.URL.String() != i.URL {
			return false
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal("Unable to read request body")
		}

		requestBody := string(body)
		return requestBody == i.Body
	}
}

func redactCredentials() recorder.HookFunc {
	return func(i *cassette.Interaction) error {
		redact := func(headers map[string][]string) {
			for key := range headers {
				if strings.Contains(strings.ToLower(key), "authorization") {
					headers[key] = []string{"REDACTED"}
				}
			}
		}
		redact(i.Request.Headers)
		redact(i.Response.Headers)

		if strings.Contains(i.Response.Body, "base64Certificate") {
			reBindingSecret := regexp.MustCompile(`"base64Certificate" : "(.*?)"`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `"base64Certificate" : "-----BEGIN CERTIFICATE-----\nredacted\n-----END CERTIFICATE-----"`)
		}

		if strings.Contains(i.Request.Body, "base64Certificate") {
			reBindingSecret := regexp.MustCompile(`"base64Certificate":"(.*?)"`)
			i.Request.Body = reBindingSecret.ReplaceAllString(i.Request.Body, `"base64Certificate":"-----BEGIN CERTIFICATE-----\nredacted\n-----END CERTIFICATE-----"`)
		}

		return nil
	}
}

func stopQuietly(rec *recorder.Recorder) {
	if err := rec.Stop(); err != nil {
		panic(err)
	}
}

func TestSciProvider_AllResources(t *testing.T) {
	expectedResources := []string{
		"sci_application",
		"sci_user",
		"sci_group",
		"sci_schema",
		"sci_corporate_idp",
	}
	ctx := context.Background()
	var registeredResources []string
	for _, resourceFunc := range New().Resources(ctx) {
		var resp tfresource.MetadataResponse
		resourceFunc().Metadata(ctx, tfresource.MetadataRequest{ProviderTypeName: "sci"}, &resp)
		registeredResources = append(registeredResources, resp.TypeName)
	}
	assert.ElementsMatch(t, expectedResources, registeredResources)
}

func TestSciProvider_AllDataSources(t *testing.T) {
	expectedDataSources := []string{
		"sci_application",
		"sci_applications",
		"sci_user",
		"sci_users",
		"sci_group",
		"sci_groups",
		"sci_schema",
		"sci_schemas",
		"sci_corporate_idp",
		"sci_corporate_idps",
	}
	ctx := context.Background()
	var registeredDataSources []string
	for _, dsFunc := range New().DataSources(ctx) {
		var resp datasource.MetadataResponse
		dsFunc().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "sci"}, &resp)
		registeredDataSources = append(registeredDataSources, resp.TypeName)
	}
	assert.ElementsMatch(t, expectedDataSources, registeredDataSources)
}

func TestProviderConfig_MissingTenantURL(t *testing.T) {
	config := `
		provider "sci" {
			username = "test"
			password = "test"
		}

		data "sci_users" "test" {}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(nil),
		Steps: []resource.TestStep{{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "tenant_url" is required, but no definition was found.`),
		}},
	})
}

func TestProviderConfig_MissingAuthCredentials(t *testing.T) {
	config := `
		provider "sci" {
			tenant_url = "https://example.com/"
		}

		data "sci_users" "test" {}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Please provide either : \n- client_id and client_secret for OAuth2 Authentication \n- p12_certificate_content and p12_certificate_password for X.509\nAuthentication \n- username and password for Basic Authentication"),
			},
		},
	})
}

func TestProviderConfig_IncompleteAuthCredentials(t *testing.T) {

	config := `
		provider "sci" {
			tenant_url = "https://example.com/"
		}

		data "sci_users" "test" {}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(nil),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv("SCI_CLIENT_ID", "client-id")
					t.Cleanup(func() {
						t.Setenv("SCI_CLIENT_ID", "")
					})
				},
				Config:      config,
				ExpectError: regexp.MustCompile("Please provide the required OAuth Credentials : Client ID and Client Secret"),
			},
			{
				PreConfig: func() {
					t.Setenv("SCI_CLIENT_ID", "")
					t.Setenv("SCI_CLIENT_SECRET", "client-secret")
				},
				Config:      config,
				ExpectError: regexp.MustCompile("Please provide the required OAuth Credentials : Client ID and Client Secret"),
			},
			{
				PreConfig: func() {
					t.Setenv("SCI_CLIENT_ID", "")
					t.Setenv("SCI_CLIENT_SECRET", "")
					t.Setenv("SCI_USERNAME", "username")
				},
				Config:      config,
				ExpectError: regexp.MustCompile("Please provide the required Basic Authentication Credentials : Username and\nPassword"),
				
			},
			{
				PreConfig: func() {
					t.Setenv("SCI_CLIENT_ID", "")
					t.Setenv("SCI_CLIENT_SECRET", "")
					t.Setenv("SCI_USERNAME", "")
					t.Setenv("SCI_PASSWORD", "password")
				},
				Config:      config,
				ExpectError: regexp.MustCompile("Please provide the required Basic Authentication Credentials : Username and\nPassword"),
			},
			{
				PreConfig: func() {
					t.Setenv("SCI_CLIENT_ID", "")
					t.Setenv("SCI_CLIENT_SECRET", "")
					t.Setenv("SCI_USERNAME", "")
					t.Setenv("SCI_PASSWORD", "")
					t.Setenv("SCI_P12_CERTIFICATE_PASSWORD", "password")
				},
				Config: config,
				ExpectError: regexp.MustCompile("Please provide the required X.509 Authentication Credentials : P12\nCertificate and P12 Certificate Password"),
			},
			{
				PreConfig: func() {
					t.Setenv("SCI_CLIENT_ID", "")
					t.Setenv("SCI_CLIENT_SECRET", "")
					t.Setenv("SCI_USERNAME", "")
					t.Setenv("SCI_PASSWORD", "")
					t.Setenv("SCI_P12_CERTIFICATE_PASSWORD", "")
				},
				Config: `
					provider "sci" {
						tenant_url = "https://example.com/"
						p12_certificate_content = "certificate-content"
					}

					data "sci_users" "test" {}
				`,
				ExpectError: regexp.MustCompile("Please provide the required X.509 Authentication Credentials : P12\nCertificate and P12 Certificate Password"),
			},
		},
	})

}

func TestAuthentication_withOAuth2(t *testing.T) {

	// Setup mock OAuth2 server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/oauth2/token" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"mocked-token","token_type":"Bearer"}`))
			return
		}

		if r.URL.Path == "/scim/Users/" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/scim+json")
			_, _ = w.Write([]byte(`{"Resources": [], "totalResults": 0}`))
			return
		}

		http.NotFound(w, r)
	}))
	defer mockServer.Close()

	// Use the mock server's URL as tenant_url
	tenantURL := mockServer.URL

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: getTestProviders(mockServer.Client()),
		Steps: []resource.TestStep{

			// Test the credentials as env variables
			{
				PreConfig: func() {
					t.Setenv("SCI_CLIENT_ID", "env-client-id")
					t.Setenv("SCI_CLIENT_SECRET", "env-client-secret")
				},
				Config: fmt.Sprintf(`
					provider "sci" {
						tenant_url = "%s"
					}

					data "sci_users" "test" {}
				`, tenantURL),
				ExpectNonEmptyPlan: false,
			},

			// Test the credentials as schema parameters
			{
				Config: fmt.Sprintf(`
					provider "sci" {
						tenant_url = "%s"
						client_id  = "test-client-id"
						client_secret = "test-client-secret"
					}

					data "sci_users" "test" {}
				`, tenantURL),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func TestAuthenticationFailure_withOAuth2(t *testing.T) {

	tenantURL := "https://example.accounts.ondemand.com/"
	clientID := "invalid-client-id"
	clientSecret := "invalid-client-secret"

	errorCases := []struct {
		StatusCode   int
		ResponseBody string
		Errors       []string
	}{
		// Invalid Client Credentials
		{
			StatusCode:   http.StatusUnauthorized,
			ResponseBody: "invalid_client",
			Errors:       []string{"failed to retrieve token", "invalid_client"},
		},
		// Unauthorized Access
		{
			StatusCode:   http.StatusUnauthorized,
			ResponseBody: "unauthorized",
			Errors:       []string{"failed to retrieve token", "unauthorized"},
		},
		// Improper JSON Response
		{
			StatusCode:   http.StatusOK,
			ResponseBody: "not a json",
			Errors:       []string{"invalid character"},
		},
		// Missing Access Token in Response
		{
			StatusCode:   http.StatusOK,
			ResponseBody: `{"token_type":"Bearer"}`,
			Errors:       []string{"server response missing access_token"},
		},
	}

	for _, test := range errorCases {

		httpClient := &http.Client{
			Transport: RoundTripperFunc(func(req *http.Request) (*http.Response, error) {

				return &http.Response{
					StatusCode: test.StatusCode,
					Body:       io.NopCloser(strings.NewReader(test.ResponseBody)),
				}, nil
			}),
		}

		token, err := fetchOAuthToken(httpClient, tenantURL, clientID, clientSecret)

		assert.Empty(t, token, "Expected token to be empty for invalid credentials")
		assert.Error(t, err, "Expected error for invalid credentials")
		for _, errMsg := range test.Errors {
			assert.Contains(t, err.Error(), errMsg, fmt.Sprintf("Error message should contain '%s'", errMsg))
		}
	}

}

func TestAuthentication_withX509Auth(t *testing.T) {

	// Setup mock SCIM server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/scim/Users/" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/scim+json")
			_, _ = w.Write([]byte(`{"Resources": [], "totalResults": 0}`))
			return
		}
		http.NotFound(w, r)

	}))
	defer mockServer.Close()

	// Generate a mock p12 file and password
	tempDir := t.TempDir()
	p12File := tempDir + "/mock-cert.p12"
	p12Password := "mockpassword"

	// Set the file paths
	configFile := tempDir + "/openssl.cnf"
	keyFile := tempDir + "/key.pem"
	certFile := tempDir + "/cert.pem"

	_ = os.WriteFile(
		configFile,
		[]byte(`
			[mock p12 certifictae]
		`), 0644)

	// Generate a self-signed certificate and convert it to p12 format
	_ = exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048", "-keyout", keyFile, "-out", certFile, "-days", "1", "-nodes", "-subj", "/CN=Mock Cert", "-config", configFile).Run()
	_ = exec.Command("openssl", "pkcs12", "-export", "-out", p12File, "-inkey", keyFile, "-in", certFile, "-password", "pass:"+p12Password).Run()

	p12Bytes, err := os.ReadFile(p12File)
	assert.NoError(t, err)
	p12Base64 := base64.StdEncoding.EncodeToString(p12Bytes)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: getTestProviders(mockServer.Client()),
		Steps: []resource.TestStep{

			// Test the certificate content as schema parameter and password as env variable
			{
				PreConfig: func() {
					t.Setenv("SCI_P12_CERTIFICATE_PASSWORD", p12Password)
				},
				Config: fmt.Sprintf(`
                    provider "sci" {
                        tenant_url               = "%s"
                        p12_certificate_content  = "%s"
                    }

                    data "sci_users" "test" {}
                `, mockServer.URL, p12Base64),
				ExpectNonEmptyPlan: false,
			},

			// Test the credentials as schema parameters
			{
				Config: fmt.Sprintf(`
                    provider "sci" {
                        tenant_url               = "%s"
                        p12_certificate_content  = "%s"
                        p12_certificate_password = "%s"
                    }

                    data "sci_users" "test" {}
                `, mockServer.URL, p12Base64, p12Password),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAuthenticationFailure_withX509Auth(t *testing.T) {

	invalidP12 := base64.StdEncoding.EncodeToString([]byte("not-a-valid-p12"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(nil),
		Steps: []resource.TestStep{

			// Test invalid base64 content
			{
				Config: `
					provider "sci" {
						tenant_url               = "https://example.com"
						p12_certificate_content  = "%%%invalid-base64%%%"
						p12_certificate_password = "dummy"
					}

					data "sci_users" "test" {}
				`,
				ExpectError: regexp.MustCompile("Failed to decode base64 content"),
			},

			// Test invalid p12 content
			{
				Config: fmt.Sprintf(`
					provider "sci" {
						tenant_url               = "https://example.com"
						p12_certificate_content  = "%s"
						p12_certificate_password = "wrong-password"
					}

					data "sci_users" "test" {}
				`, invalidP12),
				ExpectError: regexp.MustCompile("Invalid .p12 certificate"),
			},
		},
	})
}
