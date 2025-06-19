package provider

import (
	"context"
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

func TestAccSciProvider_withP12(t *testing.T) {
	rec, _ := setupVCR(t, "fixtures/provider_p12")
	defer stopQuietly(rec)

	base64Content := os.Getenv("SCI_CERTIFICATE_CONTENT")
	password := os.Getenv("SCI_P12_PASSWORD")

	if rec.IsRecording() {
		if base64Content == "" || password == "" {
			t.Skip("SCI_CERTIFICATE_CONTENT and SCI_P12_PASSWORD must be set for recording")
		}
	} else {
		base64Content = "ZHVtbXk=" // base64 for testing
		password = "test-password"
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
			provider "sci" {
			  tenant_url               = "https://iasprovidertestblr.accounts400.ondemand.com/"
			  p12_certificate_content  = "%s"
			  p12_certificate_password = "%s"
			}`, base64Content, password),
		}},
	})
}

func TestProvider_InvalidConfiguration(t *testing.T) {
	// Test case for missing tenant_url
	config := `
	provider "sci" {
	  username = "test-user"
	  password = "test-password"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("tenant_url is required"),
			},
		},
	})
}

func TestProvider_AuthenticationFailure(t *testing.T) {
	// Test case for invalid credentials
	config := `
	provider "sci" {
	  tenant_url = "https://iasprovidertestblr.accounts400.ondemand.com/"
	  username   = "invalid-user"
	  password   = "invalid-password"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("authentication failed"),
			},
		},
	})
}

func TestProvider_InvalidCertificateContent(t *testing.T) {
	// Test case for malformed certificate content
	config := `
	provider "sci" {
	  tenant_url               = "https://iasprovidertestblr.accounts400.ondemand.com/"
	  p12_certificate_content  = "invalid-base64"
	  p12_certificate_password = "test-password"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("invalid certificate content"),
			},
		},
	})
}

func TestProvider_NetworkError(t *testing.T) {
	// Test case for unreachable tenant_url
	config := `
	provider "sci" {
	  tenant_url = "https://invalid-url.com/"
	  username   = "test-user"
	  password   = "test-password"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProviders(http.DefaultClient),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("unable to reach tenant_url"),
			},
		},
	})
}
