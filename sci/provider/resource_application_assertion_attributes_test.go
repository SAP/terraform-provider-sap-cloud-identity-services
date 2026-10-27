package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceApplication_EmptyAssertionAttributes(t *testing.T) {
	const appID = "f2a1c3d4-0000-4000-8000-000000000002"
	const authSchema = "urn:sap:identity:application:schemas:extension:sci:1.0:Authentication"

	// Mock of the Applications API for one application. Like the API, it omits assertionAttributes
	// when the list is empty. The API does this only while the application has advanced assertion
	// attributes, otherwise it sets default attributes, so the configuration below has one.
	var mu sync.Mutex
	var assertionAttributes, advancedAssertionAttributes []any
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
			assertionAttributes, _ = auth["assertionAttributes"].([]any)
			advancedAssertionAttributes, _ = auth["advancedAssertionAttributes"].([]any)
			w.Header().Set("Location", "/Applications/v1/"+appID)
			w.WriteHeader(http.StatusCreated)
		case r.URL.Path == "/Applications/v1/"+appID && r.Method == http.MethodGet:
			auth := map[string]any{
				"subjectNameIdentifier":       "mail",
				"advancedAssertionAttributes": advancedAssertionAttributes,
			}
			if len(assertionAttributes) > 0 {
				auth["assertionAttributes"] = assertionAttributes
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":          appID,
				"name":        "empty-assertion-attributes-app",
				"description": "application for testing empty assertion attributes",
				authSchema:    auth,
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
				if strings.HasSuffix(op.Path, "/assertionAttributes") {
					assertionAttributes, _ = op.Value.([]any)
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

	config := func(assertionAttributesValue string) string {
		return fmt.Sprintf(`
			provider "sci" {
				tenant_url    = "%s"
				client_id     = "test-client-id"
				client_secret = "test-client-secret"
			}

			resource "sci_application" "testApp" {
				name        = "empty-assertion-attributes-app"
				description = "application for testing empty assertion attributes"
				authentication_schema = {
					advanced_assertion_attributes = [
						{
							source          = "Expression"
							attribute_name  = "Groups"
							attribute_value = "sac"
						},
					]
					assertion_attributes = %s
				}
			}
		`, mockServer.URL, assertionAttributesValue)
	}

	assertionAttributesInAPI := func(want int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			mu.Lock()
			defer mu.Unlock()
			if len(assertionAttributes) != want {
				return fmt.Errorf("the API has %d assertion attributes, want %d", len(assertionAttributes), want)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: getTestProviders(mockServer.Client()),
		Steps: []resource.TestStep{
			{
				Config: config(`[{ attribute_name = "firstName", attribute_value = "firstName" }]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.#", "1"),
					assertionAttributesInAPI(1),
				),
			},
			// The framework checks after apply that the state matches the planned empty list,
			// and the step fails if the next plan is not empty.
			{
				Config: config(`[]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sci_application.testApp", "authentication_schema.assertion_attributes.#", "0"),
					assertionAttributesInAPI(0),
				),
			},
		},
	})
}
