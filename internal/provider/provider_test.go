package provider

import (
	"fmt"
	"io"
	"net/http"
	"os"
	// "strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

type User struct {
	Username 	string
	Password 	string
}

func providerConfig(tenantUrl string) string {
	return fmt.Sprintf(`
	provider "ias" {
		tenant_url = "%s"
	}
	`, tenantUrl)
}

func getTestProviders(httpClient *http.Client) map[string]func() (tfprotov6.ProviderServer, error) {
	iasProvider := NewWithClient(httpClient)
	
	return map[string]func() (tfprotov6.ProviderServer, error){
		"ias": providerserver.NewProtocol6WithError(iasProvider),
	}
}

func setupVCR(t *testing.T, cassetteName string) (*recorder.Recorder, User) {
	t.Helper()

	mode := recorder.ModeRecordOnly
	// if testRecord, _ := strconv.ParseBool(os.Getenv("TEST_RECORD")); testRecord {
	// 	mode = recorder.ModeRecordOnly
	// }

	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       cassetteName,
		Mode:               mode,
		SkipRequestLatency: true,
		RealTransport:      http.DefaultTransport,
	})

	var testUser User
	if rec.IsRecording() {
		t.Logf("ATTENTION: Recording '%s'", cassetteName)
		testUser.Username = os.Getenv("ias_username")
		testUser.Password = os.Getenv("ias_password")
		// if len(testUser.Username) == 0 || len(testUser.Password) == 0 {
		// 	t.Fatal("Env vars ias_username and ias_password are required when recording test fixtures")
		// }
	} else {
		t.Logf("Replaying '%s'", cassetteName)
	}

	if err != nil {
		t.Fatal()
	}

	rec.SetMatcher(RequestMatcher(t))

	//any hooks to be added?

	return rec, testUser
}

func RequestMatcher(t *testing.T) (cassette.MatcherFunc) {
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

func stopQuietly(rec *recorder.Recorder) {
	if err := rec.Stop(); err != nil {
		panic(err)
	}
}