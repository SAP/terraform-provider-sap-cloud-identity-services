package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

type User struct {
	Username 	string
	Password 	string
}

func providerConfig(_ string, testUser User) string {
	tenantUrl := "https://iasprovidertestblr.accounts400.ondemand.com/"
	return fmt.Sprintf(`
	provider "ias" {
		tenant_url = "%s"
		username = "%s"
		password = "%s"
	}
	`, tenantUrl, testUser.Username, testUser.Password)
}

func getTestProviders(httpClient *http.Client) map[string]func() (tfprotov6.ProviderServer, error) {
	iasProvider := NewWithClient(httpClient)
	
	return map[string]func() (tfprotov6.ProviderServer, error){
		"ias": providerserver.NewProtocol6WithError(iasProvider),
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

	var  testUser User
	if rec.IsRecording() {
		t.Logf("ATTENTION: Recording '%s'", cassetteName)
		testUser.Username = os.Getenv("ias_username")
		testUser.Password = os.Getenv("ias_password")
		if len(testUser.Username) == 0 || len(testUser.Password) == 0 {
			t.Fatal("Env vars ias_username and ias_password are required when recording test fixtures")
		}
	} else {
		t.Logf("Replaying '%s'", cassetteName)
	}

	if err != nil {
		t.Fatal()
	}

	rec.SetMatcher(requestMatcher(t))
	rec.AddHook(redactAuthorizationToken(), recorder.BeforeSaveHook)

	return rec, testUser
}

func requestMatcher(t *testing.T) (cassette.MatcherFunc) {
	return func(r *http.Request, i cassette.Request) bool {
		if r.Method != i.Method || r.URL.String() != i.URL {
			return false
		}

		//headers verification?

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal("Unable to read body from request")
		}

		requestBody := string(bytes)
		return requestBody == i.Body
	}
}

func redactAuthorizationToken() (recorder.HookFunc) {
	return func(i *cassette.Interaction) error {

		redact := func(headers map[string][]string) {
			for key := range headers {
				if strings.Contains(strings.ToLower(key), "authorization") {
					headers[key] = []string{"redacted"}
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

func TestIasProvider_AllResources (t *testing.T){

	expectedResources := []string{
		"ias_application",
	}

	ctx := context.Background()
	registeredResources := []string{}

	for _, resourceFunc := range New().Resources(ctx) {
		var resp resource.MetadataResponse

		resourceFunc().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "ias"}, &resp)

		registeredResources = append(registeredResources, resp.TypeName)
	}

	assert.ElementsMatch(t, expectedResources, registeredResources)
}

func TestIasProvider_AllDataSources (t *testing.T){

	expectedDataSources := []string{
		"ias_application",
		"ias_applications",
	}

	ctx := context.Background()
	registeredDataSources := []string{}

	for _, datasourceFunc := range New().DataSources(ctx) {
		var resp datasource.MetadataResponse

		datasourceFunc().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "ias"}, &resp)

		registeredDataSources = append(registeredDataSources, resp.TypeName)
	}

	assert.ElementsMatch(t, expectedDataSources, registeredDataSources)
}