package cli

import (
	"context"
	"encoding/json"

	"fmt"
	"net/http"
	"terraform-provider-ias/internal/cli/apiObjects/applications"
	"testing"

	"github.com/stretchr/testify/assert"
)
 
var applicationsPath = "/Applications/v1/"
var applicationsBody = applications.Application{
	Name: "Test Application",
	Description: "test app",
	AuthenticationSchema: applications.AuthenticationSchema{
		SsoType: "saml",
		SubjectNameIdentifier: "mail",
		AssertionAttributes: []applications.AssertionAttribute{
			{
				AssertionAttributeName: "attr_name_1",
				UserAttributeName: "user_name_1",
			},
			{
				AssertionAttributeName: "attr_name_2",
				UserAttributeName: "user_name_2",
			},
		},
		AdvancedAssertionAttributes: []applications.AdvancedAssertionAttribute{
			{
				AttributeName: "attr_name_1",
				AttributeValue: "attr_value_1",
			},
			{
				AttributeName: "attr_name_2",
				AttributeValue: "attr_value_2",
			},
		},
		DefaultAuthenticatingIdpId: "valid-idp-uuid",
		ConditionalAuthentication: []applications.AuthenicationRule{
			{
				UserType: "employee",
				IdentityProviderId: "valid-idp-uuid",
			},
			{
				UserGroup: "valid-group-id",
				UserEmailDomain: "test.com",
			},
		},
	},
}

var applicationsResponse []byte
 
func TestApplications_Create(t *testing.T) {

	applicationsResponse, _ = json.Marshal(applicationsBody)

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("location", "/Applications/v1/valid-app-id")
			w.Write(applicationsResponse)

			if r.Method != "GET" {
				assertCall[applications.Application](t, r, applicationsPath, "POST", applicationsBody)
			}
            
        }))
 
        defer srv.Close()
 
        _, _, err := client.Application.Create(context.TODO(), &applicationsBody)
 
		assert.NoError(t, err)
    }) 
}

func TestApplications_Get(t *testing.T){

    allApplications := []applications.Application{
		applicationsBody,
		applicationsBody,
	}

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        res, _ := json.Marshal(applications.ApplicationsResponse{
            Applications: allApplications,
        })

        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(res)
            assertCall[applications.Application](t, r, applicationsPath, "GET", nil)
        }))
 
        defer srv.Close()
 
        _, _, err := client.Application.Get(context.TODO())
 
		assert.NoError(t, err)
    }) 
}

func TestApplications_GetByAppId(t *testing.T){

    applicationsResponse, _ = json.Marshal(applicationsBody)

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(applicationsResponse)
            assertCall[applications.Application](t, r, fmt.Sprintf("%s%s", applicationsPath, "valid-app-id"), "GET", nil)
        }))
 
        defer srv.Close()
 
        _, _, err := client.Application.GetByAppId(context.TODO(), "valid-app-id")
		assert.NoError(t, err)
    }) 
}

func TestApplications_Update(t *testing.T) {

	applicationsBody.Id = "valid-app-id"
	applicationsResponse, _ := json.Marshal(applicationsBody)

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(applicationsResponse)

			if r.Method != "GET" {
				assertCall[applications.Application](t, r, fmt.Sprintf("%s%s", applicationsPath, "valid-app-id"), "PUT", applicationsBody)
			}
        }))
 
        defer srv.Close()
 
        _, _, err := client.Application.Update(context.TODO(), &applicationsBody)
 
		assert.NoError(t, err)
    }) 

}

func TestApplications_Delete(t *testing.T) {

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            assertCall[applications.Application](t, r, fmt.Sprintf("%s%s", applicationsPath, "valid-app-id"), "DELETE", nil)
        }))
 
        defer srv.Close()
 
        err := client.Application.Delete(context.TODO(), "valid-app-id")
 
		assert.NoError(t, err)
    }) 

}