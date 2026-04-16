package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
	"github.com/stretchr/testify/assert"
)

var (
	appSecretPath  = "/Applications/v1/valid-app-id/apiSecrets"
	testAppId      = "valid-app-id"
	testSecretId   = "996cee6e-ca22-4b9e-9e4d-8cad11ad3e10"
	testSecretBody = applications.ApplicationSecret{
		Id:                  "996cee6e-ca22-4b9e-9e4d-8cad11ad3e10",
		ClientId:            "7aef9f3c-23a0-4e2a-8b33-bc45ac13ba6f",
		Secret:              "H6Fd2AYJ[Xm.uI[4EznlbLV5?roC",
		Hint:                "HAC5",
		Description:         "test",
		ValidTo:             "2029-10-12T10:00:00Z",
		AuthorizationScopes: []string{"manageApp", "oAuth"},
		AllApisAccess:       new(bool),
	}
	testSecretRequest = applications.ApplicationSecretRequest{
		Description:         "test",
		ValidTo:             "2029-10-12T10:00:00Z",
		AuthorizationScopes: []string{"manageApp", "oAuth"},
		AllApisAccess:       new(bool),
	}
)

func TestApplicationSecrets_Create(t *testing.T) {

	secretResponse, _ := json.Marshal(testSecretBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(secretResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[applications.ApplicationSecretRequest](t, r, appSecretPath, "POST", testSecretRequest)
		}))
		defer srv.Close()

		res, err := client.ApplicationSecret.Create(context.TODO(), testAppId, testSecretRequest)

		assert.NoError(t, err)
		assert.Equal(t, testSecretBody.Id, res.Id)
		assert.Equal(t, testSecretBody.Secret, res.Secret)
		assert.Equal(t, testSecretBody.Hint, res.Hint)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "create failed",
				Details: []ErrorDetail{{Message: "server error"}},
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")
		}))
		defer srv.Close()

		res, err := client.ApplicationSecret.Create(context.TODO(), testAppId, testSecretRequest)

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "error 400 \ncreate failed : server error", err.Error())
	})
}

func TestApplicationSecrets_Get(t *testing.T) {

	listResponse, _ := json.Marshal(applications.ApplicationSecretsListResponse{
		Secrets: []applications.ApplicationSecret{testSecretBody},
	})

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(listResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[any](t, r, appSecretPath, "GET", nil)
		}))
		defer srv.Close()

		res, err := client.ApplicationSecret.Get(context.TODO(), testAppId)

		assert.NoError(t, err)
		assert.Len(t, res.Secrets, 1)
		assert.Equal(t, testSecretBody.Id, res.Secrets[0].Id)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "get failed",
				Details: []ErrorDetail{{Message: "server error"}},
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")
		}))
		defer srv.Close()

		res, err := client.ApplicationSecret.Get(context.TODO(), testAppId)

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "error 400 \nget failed : server error", err.Error())
	})
}

func TestApplicationSecrets_GetById(t *testing.T) {

	listResponse, _ := json.Marshal(applications.ApplicationSecretsListResponse{
		Secrets: []applications.ApplicationSecret{testSecretBody},
	})

	t.Run("returns the matching secret from the list", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(listResponse)
			assert.NoError(t, err, "Failed to write response")
		}))
		defer srv.Close()

		res, err := client.ApplicationSecret.GetById(context.TODO(), testAppId, testSecretId)

		assert.NoError(t, err)
		assert.Equal(t, testSecretId, res.Id)
		assert.Equal(t, testSecretBody.Description, res.Description)
	})

	t.Run("returns error when secret id is not found in list", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(listResponse)
			assert.NoError(t, err, "Failed to write response")
		}))
		defer srv.Close()

		res, err := client.ApplicationSecret.GetById(context.TODO(), testAppId, "non-existent-id")

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "non-existent-id")
	})

	t.Run("propagates error from Get", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "get failed",
				Details: []ErrorDetail{{Message: "server error"}},
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")
		}))
		defer srv.Close()

		res, err := client.ApplicationSecret.GetById(context.TODO(), testAppId, testSecretId)

		assert.Zero(t, res)
		assert.Error(t, err)
	})
}

func TestApplicationSecrets_Update(t *testing.T) {

	secretResponse, _ := json.Marshal(applications.ApplicationSecretsListResponse{
		Secrets: []applications.ApplicationSecret{testSecretBody},
	})

	patchRequests := []generic.PatchRequest{
		{Op: "replace", Path: "/description", Value: "updated description"},
		{Op: "replace", Path: "/validTo", Value: "2030-01-01T00:00:00Z"},
	}

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PATCH" {
				var actualBody applications.ApplicationSecretPatchRequestBody
				err := json.NewDecoder(r.Body).Decode(&actualBody)
				assert.NoError(t, err)
				assert.Equal(t, 2, len(actualBody.Operations))
				assert.Equal(t, fmt.Sprintf("%s/%s", appSecretPath, testSecretId), r.URL.Path)
			}
			_, err := w.Write(secretResponse)
			assert.NoError(t, err, "Failed to write response")
		}))
		defer srv.Close()

		res, err := client.ApplicationSecret.Update(context.TODO(), testAppId, testSecretId, patchRequests)

		assert.NoError(t, err)
		assert.Equal(t, testSecretId, res.Id)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "update failed",
				Details: []ErrorDetail{{Message: "server error"}},
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PATCH" {
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write(resErr)
				assert.NoError(t, err, "Failed to write response")
			} else {
				_, err := w.Write(secretResponse)
				assert.NoError(t, err, "Failed to write response")
			}
		}))
		defer srv.Close()

		res, err := client.ApplicationSecret.Update(context.TODO(), testAppId, testSecretId, patchRequests)

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "error 400 \nupdate failed : server error", err.Error())
	})
}

func TestApplicationSecrets_Delete(t *testing.T) {

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, appSecretPath, r.URL.Path)
			assert.Equal(t, testSecretId, r.URL.Query().Get("id"))
			assert.Equal(t, "DELETE", r.Method)
		}))
		defer srv.Close()

		err := client.ApplicationSecret.Delete(context.TODO(), testAppId, testSecretId)

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "delete failed",
				Details: []ErrorDetail{{Message: "server error"}},
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")
		}))
		defer srv.Close()

		err := client.ApplicationSecret.Delete(context.TODO(), testAppId, testSecretId)

		assert.Error(t, err)
		assert.Equal(t, "error 400 \ndelete failed : server error", err.Error())
	})
}
