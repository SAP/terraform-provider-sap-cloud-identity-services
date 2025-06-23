package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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

	var testUser User
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
	rec.AddHook(redactAuthorizationToken(), recorder.BeforeSaveHook)

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
		return string(body) == i.Body
	}
}

func redactAuthorizationToken() recorder.HookFunc {
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

func TestAccSciProvider_withOAuth2(t *testing.T) {
	rec, _ := setupVCR(t, "fixtures/provider_oauth")
	defer stopQuietly(rec)

	clientID := os.Getenv("SCI_CLIENT_ID")
	clientSecret := os.Getenv("SCI_CLIENT_SECRET")

	if rec.IsRecording() {
		if clientID == "" || clientSecret == "" {
			t.Skip("SCI_CLIENT_ID and SCI_CLIENT_SECRET must be set for recording")
		}
	} else {
		clientID = "test-client-id"
		clientSecret = "test-client-secret"
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
			provider "sci" {
			  tenant_url    = "https://iasprovidertestblr.accounts400.ondemand.com/"
			  client_id     = "%s"
			  client_secret = "%s"
			}`, clientID, clientSecret),
		}},
	})
}

func TestFetchOAuthToken_Failure(t *testing.T) {

	httpClient := &http.Client{
		Transport: RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			// Simulate a 401 Unauthorized response
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"error": "invalid_client"}`)),
			}, nil
		}),
	}

	tenantURL := "https://iasprovidertestblr.accounts400.ondemand.com/"
	clientID := "invalid-client-id"
	clientSecret := "invalid-client-secret"

	token, err := fetchOAuthToken(httpClient, tenantURL, clientID, clientSecret)

	// Assertion that the token is empty and an error is returned
	assert.Empty(t, token, "Expected token to be empty for invalid credentials")
	assert.Error(t, err, "Expected error for invalid credentials")
	assert.Contains(t, err.Error(), "token request failed with status 401", "Expected 401 Unauthorized error")
}

func TestAccSciProvider_withP12(t *testing.T) {
	mode := recorder.ModeRecordOnce
	if testRecord, _ := strconv.ParseBool(os.Getenv("TEST_RECORD")); testRecord {
		mode = recorder.ModeRecordOnly
	}

	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       "fixtures/provider_p12_success",
		Mode:               mode,
		SkipRequestLatency: true,
		RealTransport:      http.DefaultTransport,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stopQuietly(rec)

	// Use real env vars if recording
	var base64Content, password string
	if rec.IsRecording() {
		base64Content = os.Getenv("SCI_CERTIFICATE_CONTENT")
		password = os.Getenv("SCI_P12_PASSWORD")

		if base64Content == "" || password == "" {
			t.Skip("SCI_CERTIFICATE_CONTENT and SCI_P12_PASSWORD must be set for recording")
		}
	} else {
		// Values ignored during replay
		base64Content = base64.StdEncoding.EncodeToString([]byte("dummy-p12-content"))
		password = "12345678"
	}

	rec.SetMatcher(requestMatcher(t))
	rec.AddHook(redactAuthorizationToken(), recorder.BeforeSaveHook)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
			provider "sci" {
			  tenant_url               = "https://iasprovidertestblr.accounts400.ondemand.com/"
			  p12_certificate_content  = "%s"
			  p12_certificate_password = "%s"
			}

			data "sci_users" "dummy" {}
			`, base64Content, password),
		}},
	})
}

func TestProvider_AuthenticationFailure(t *testing.T) {
	config := `
	provider "sci" {
	  tenant_url = "https://iasprovidertestblr.accounts400.ondemand.com/"
	  username   = "invalid-user"
	  password   = "invalid-password"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{{
			Config:             config,
			ExpectNonEmptyPlan: false,
		}},
	})
}

func TestProvider_InvalidCertificateContent(t *testing.T) {
	config := `
	provider "sci" {
	  tenant_url               = "https://iasprovidertestblr.accounts400.ondemand.com/"
	  p12_certificate_content  = "not-a-valid-base64-string"
	  p12_certificate_password = "test-password"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{{
			Config:             config,
			ExpectNonEmptyPlan: false,
		}},
	})
}

func TestProvider_MissingTenantURL(t *testing.T) {
	config := `
	provider "sci" {
	  username = "test"
	  password = "test"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{{
			Config:             config,
			ExpectNonEmptyPlan: false,
		}},
	})
}

func TestProvider_InvalidTenantURL(t *testing.T) {
	config := `
	provider "sci" {
	  tenant_url = "ht@tp://bad_url"
	  username   = "test"
	  password   = "test"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{{
			Config:             config,
			ExpectNonEmptyPlan: false,
		}},
	})
}

func TestProvider_InvalidCertificateBase64(t *testing.T) {
	config := `
	provider "sci" {
	  tenant_url               = "https://iasprovidertestblr.accounts400.ondemand.com/"
	  p12_certificate_content  = "invalid-base64-@@@"
	  p12_certificate_password = "any"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{{
			Config:             config,
			ExpectNonEmptyPlan: false,
		}},
	})
}

func TestProvider_InvalidCertificateFormat(t *testing.T) {
	invalidP12Content := base64.StdEncoding.EncodeToString([]byte("not-a-valid-p12-content"))
	config := fmt.Sprintf(`
	provider "sci" {
	  tenant_url               = "https://iasprovidertestblr.accounts400.ondemand.com/"
	  p12_certificate_content  = "%s"
	  p12_certificate_password = "wrong-password"
	}
	`, invalidP12Content)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{{
			Config:             config,
			ExpectNonEmptyPlan: false,
		}},
	})
}

func TestFetchOAuthToken_HTTPErrorStatus(t *testing.T) {
	httpClient := &http.Client{
		Transport: RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"error": "unauthorized"}`)),
			}, nil
		}),
	}

	token, err := fetchOAuthToken(httpClient, "https://example.com", "id", "secret")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "token request failed with status 401")
}

func TestFetchOAuthToken_InvalidJSONResponse(t *testing.T) {
	httpClient := &http.Client{
		Transport: RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("not a json")),
			}, nil
		}),
	}

	token, err := fetchOAuthToken(httpClient, "https://example.com", "id", "secret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse token response")
	assert.Empty(t, token)
}

func TestFetchOAuthToken_EmptyAccessToken(t *testing.T) {
	httpClient := &http.Client{
		Transport: RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"token_type":"Bearer"}`)),
			}, nil
		}),
	}

	token, err := fetchOAuthToken(httpClient, "https://example.com", "id", "secret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty access token")
	assert.Empty(t, token)
}

func TestAccConfigure_Error_InvalidBase64Content(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{
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
		},
	})
}

func TestAccConfigure_Error_InvalidP12Certificate(t *testing.T) {
	badP12 := base64.StdEncoding.EncodeToString([]byte("not-a-valid-p12"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				provider "sci" {
				  tenant_url               = "https://example.com"
				  p12_certificate_content  = "%s"
				  p12_certificate_password = "wrong-password"
				}

				data "sci_users" "test" {}
				`, badP12),
				ExpectError: regexp.MustCompile("Invalid .p12 certificate"),
			},
		},
	})
}
