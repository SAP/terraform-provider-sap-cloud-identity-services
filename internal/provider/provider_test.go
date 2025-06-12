package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"strconv"

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
	rec, _ := setupVCR(t, "provider_oauth")
	defer stopQuietly(rec)

	clientID := os.Getenv("SCI_CLIENT_ID")
	clientSecret := os.Getenv("SCI_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skip("SCI_CLIENT_ID and SCI_CLIENT_SECRET must be set")
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
	rec, _ := setupVCR(t, "provider_p12")
	defer stopQuietly(rec)

	content, err := os.ReadFile("test-fixtures/cert.p12")
	if err != nil {
		t.Skipf("skipping test: failed to read cert.p12: %v", err)
	}
	base64Content := base64.StdEncoding.EncodeToString(content)
	password := os.Getenv("SCI_P12_PASSWORD")
	if password == "" {
		t.Skip("SCI_P12_PASSWORD must be set")
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
