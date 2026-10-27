package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestResourceApplication_RemoveFallbackAttribute(t *testing.T) {
	const appID = "f2a1c3d4-0000-4000-8000-000000000001"
	const authSchema = "urn:sap:identity:application:schemas:extension:sci:1.0:Authentication"

	// Mock of the Applications API that keeps the fallback subject name identifier of one application.
	var mu sync.Mutex
	fallback := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		body, _ := io.ReadAll(r.Body)

		switch {
		case r.URL.Path == "/oauth2/token" && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"mocked-token","token_type":"Bearer"}`))
		case r.URL.Path == "/Applications/v1/" && r.Method == http.MethodPost:
			var app map[string]any
			if err := json.Unmarshal(body, &app); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			auth, _ := app[authSchema].(map[string]any)
			fallback, _ = auth["fallbackSubjectNameIdentifier"].(string)
			w.Header().Set("Location", "/Applications/v1/"+appID)
			w.WriteHeader(http.StatusCreated)
		case r.URL.Path == "/Applications/v1/"+appID && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":          appID,
				"name":        "fallback-removal-app",
				"description": "application for testing the removal of the fallback attribute",
				authSchema: map[string]any{
					"subjectNameIdentifier":         "mail",
					"fallbackSubjectNameIdentifier": fallback,
				},
			})
		case r.URL.Path == "/Applications/v1/"+appID && r.Method == http.MethodPatch:
			var req struct {
				Operations []struct {
					Path  string `json:"path"`
					Value any    `json:"value"`
				} `json:"operations"`
			}
			_ = json.Unmarshal(body, &req)
			for _, op := range req.Operations {
				if strings.HasSuffix(op.Path, "/fallbackSubjectNameIdentifier") {
					fallback, _ = op.Value.(string)
				}
			}
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/Applications/v1/"+appID && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	config := func(fallbackBlock string) string {
		return fmt.Sprintf(`
			provider "sci" {
				tenant_url    = "%s"
				client_id     = "test-client-id"
				client_secret = "test-client-secret"
			}

			resource "sci_application" "testApp" {
				name        = "fallback-removal-app"
				description = "application for testing the removal of the fallback attribute"
				authentication_schema = {
					subject_name_identifier = {
						source = "Identity Directory"
						value  = "mail"
						%s
					}
				}
			}
		`, mockServer.URL, fallbackBlock)
	}

	fallbackIs := func(want string) resource.TestCheckFunc {
		return func(*terraform.State) error {
			mu.Lock()
			defer mu.Unlock()
			if fallback != want {
				return fmt.Errorf("fallbackSubjectNameIdentifier in the API is %q, want %q", fallback, want)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: getTestProviders(mockServer.Client()),
		Steps: []resource.TestStep{
			{
				Config:      config(`fallback_attribute = { value = "None" }`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				Config:      config(`fallback_attribute = {}`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
			},
			{
				Config: config(`fallback_attribute = { value = "uid" }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier.fallback_attribute.value", "uid"),
					fallbackIs("uid"),
				),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("sci_application.testApp", tfjsonpath.New("authentication_schema").AtMapKey("subject_name_identifier").AtMapKey("fallback_attribute"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("sci_application.testApp", "authentication_schema.subject_name_identifier.fallback_attribute.value"),
					fallbackIs("none"),
				),
			},
		},
	})
}
