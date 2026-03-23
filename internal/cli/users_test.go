package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/users"

	"github.com/stretchr/testify/assert"
)

var usersPath = "/scim/Users/"
var usersBody = users.User{
	UserName: "user_for_testing",
	Name: &users.Name{
		GivenName:  "Testing",
		FamilyName: "User",
	},
	Emails: []users.Email{
		{
			Value: "user@testing.com",
			Type:  "work",
		},
	},
}

var usersResponse []byte

func TestUsers_Create(t *testing.T) {

	usersResponse, _ = json.Marshal(usersBody)

	customSchemas, _ := json.Marshal(map[string]any{
		"schema_id": map[string]any{
			"var1": "test",
			"var2": 1,
		},
	})

	incorrectCustomSchemas, _ := json.Marshal(map[string]any{
		"new_schema_id": map[string]any{
			"var1": "test",
			"var2": 1,
		},
	})

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(usersResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[users.User](t, r, usersPath, "POST", usersBody)
		}))

		defer srv.Close()

		_, _, err := client.User.Create(context.TODO(), "", &usersBody)

		assert.NoError(t, err)
	})

	t.Run("validate the API request with custom schemas", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(responseWithCustomSchemas(usersResponse, customSchemas))
			assert.NoError(t, err, "Failed to write response")

			assertCall[users.User](t, r, usersPath, "POST", usersBody)
		}))

		defer srv.Close()

		_, _, err := client.User.Create(context.TODO(), string(customSchemas), &usersBody)

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimResponseError{
			Detail: "create failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[users.User](t, r, usersPath, "POST", users.User{})
		}))

		defer srv.Close()

		res, _, err := client.User.Create(context.TODO(), "", &users.User{})

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "SCIM error 400 \ncreate failed", err.Error())
	})

	t.Run("validate the API request with custom schemas - error", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(responseWithCustomSchemas(usersResponse, incorrectCustomSchemas))
			assert.NoError(t, err, "Failed to write response")

			assertCall[users.User](t, r, usersPath, "POST", users.User{})
		}))

		defer srv.Close()

		res, _, err := client.User.Create(context.TODO(), string(customSchemas), &users.User{})

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "schema_id not found in the returned response", err.Error())
	})
}

func TestUsers_Get(t *testing.T) {

	allUsers := []users.User{
		usersBody,
		usersBody,
	}

	t.Run("validate the API request", func(t *testing.T) {

		res, _ := json.Marshal(users.UsersResponse{
			Resources: allUsers,
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(res)
			assert.NoError(t, err, "Failed to write response")

			assertCall[users.User](t, r, usersPath, "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.User.Get(context.TODO())

		assert.NoError(t, err)
	})

	t.Run("validate the API request with error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimResponseError{
			Detail: "get failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[users.User](t, r, usersPath, "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.User.Get(context.TODO())

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "SCIM error 400 \nget failed", err.Error())
	})
}

func TestUsers_GetByUserId(t *testing.T) {

	usersResponse, _ = json.Marshal(usersBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(usersResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[users.User](t, r, fmt.Sprintf("%s%s", usersPath, "valid-user-id"), "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.User.GetByUserId(context.TODO(), "valid-user-id", false, "")

		assert.NoError(t, err)
	})

	t.Run("validate the API request with error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimResponseError{
			Detail: "get failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[users.User](t, r, fmt.Sprintf("%s%s", usersPath, "valid-user-id"), "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.User.GetByUserId(context.TODO(), "valid-user-id", false, "")

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "SCIM error 400 \nget failed", err.Error())
	})
}

func TestUsers_Update(t *testing.T) {

	usersBody.Id = "valid-user-id"
	usersResponse, _ := json.Marshal(usersBody)

	patchRequests := []generic.PatchRequest{
		{
			Op:    "replace",
			Path:  "userName",
			Value: "updated-user-name",
		},
		{
			Op:    "replace",
			Path:  "schemas",
			Value: []string{"urn:ietf:params:scim:schemas:core:2.0:User", "urn:ietf:params:scim:CustomSchema"},
		},
		{
			Op:    "replace",
			Path:  "password",
			Value: "updated-password",
		},
		{
			Op:    "replace",
			Path:  "name",
			Value: "updated-name",
		},
		{
			Op:    "replace",
			Path:  "displayName",
			Value: "updated-display-name",
		},
		{
			Op:    "replace",
			Path:  "userType",
			Value: "external",
		},
		{
			Op:    "replace",
			Path:  "active",
			Value: true,
		},
		{
			Op:   "replace",
			Path: "emails",
			Value: []users.Email{
				{
					Value:   "test1@gmail.com",
					Type:    "work",
					Primary: true,
				},
			},
		},
		{
			Op:    "replace",
			Path:  "urn:ietf:params:scim:schemas:extension:sap:2.0:User:sendMail",
			Value: true,
		},
		{
			Op:    "replace",
			Path:  "urn:ietf:params:scim:schemas:extension:sap:2.0:User:mailVerified",
			Value: true,
		},
		{
			Op:    "replace",
			Path:  "urn:ietf:params:scim:schemas:extension:sap:2.0:User:status",
			Value: "active",
		},
		{
			Op:    "replace",
			Path:  "urn:ietf:params:scim:CustomSchema:test1",
			Value: "test-val",
		},
	}

	customSchemas, _ := json.Marshal(map[string]any{
		"urn:ietf:params:scim:CustomSchema": map[string]any{
			"test1": "test-val",
		},
	})

	incorrectCustomSchemas, _ := json.Marshal(map[string]any{
		"new_schema_id": map[string]any{
			"var1": "test",
			"var2": 1,
		},
	})

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PATCH" {

				var actualBody users.PatchRequestBody
				err := json.NewDecoder(r.Body).Decode(&actualBody)

				assert.NoError(t, err)
				assert.Equal(t, 12, len(actualBody.Operations))
			}
			_, err := w.Write(usersResponse)
			assert.NoError(t, err, "Failed to write response")
		}))

		defer srv.Close()

		_, _, err := client.User.Update(context.TODO(), "valid-user-id", patchRequests, "")

		assert.NoError(t, err)
	})

	t.Run("validate the API request with custom schemas", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PATCH" {

				var actualBody users.PatchRequestBody
				err := json.NewDecoder(r.Body).Decode(&actualBody)

				assert.NoError(t, err)
				assert.Equal(t, 12, len(actualBody.Operations))
			}
			_, err := w.Write(responseWithCustomSchemas(usersResponse, customSchemas))
			assert.NoError(t, err, "Failed to write response")
		}))

		defer srv.Close()

		_, _, err := client.User.Update(context.TODO(), "valid-user-id", patchRequests, string(customSchemas))

		assert.NoError(t, err)
	})

	t.Run("validate the API request with error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimResponseError{
			Detail: "update failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")
		}))

		defer srv.Close()

		res, _, err := client.User.Update(context.TODO(), "valid-user-id", patchRequests, "")

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "SCIM error 400 \nupdate failed", err.Error())
	})

	t.Run("validate the API request with custom schemas - error", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(responseWithCustomSchemas(usersResponse, incorrectCustomSchemas))
			assert.NoError(t, err, "Failed to write response")
		}))

		defer srv.Close()

		res, _, err := client.User.Update(context.TODO(), "valid-user-id", patchRequests, string(customSchemas))

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "urn:ietf:params:scim:CustomSchema not found in the returned response", err.Error())
	})
}

func TestUsers_Delete(t *testing.T) {

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertCall[users.User](t, r, fmt.Sprintf("%s%s", usersPath, "valid-user-id"), "DELETE", users.User{})
		}))

		defer srv.Close()

		err := client.User.Delete(context.TODO(), "valid-user-id")

		assert.NoError(t, err)
	})

	t.Run("validate the API request with error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimResponseError{
			Detail: "delete failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[users.User](t, r, fmt.Sprintf("%s%s", usersPath, "valid-user-id"), "DELETE", users.User{})
		}))

		defer srv.Close()

		err := client.User.Delete(context.TODO(), "valid-user-id")

		assert.Error(t, err)
		assert.Equal(t, "SCIM error 400 \ndelete failed", err.Error())
	})
}

func responseWithCustomSchemas(userRes []byte, customSchemas []byte) []byte {
	var res bytes.Buffer
	// remove closing brace : '}'
	modifiedUsersRes := string(userRes)[:len(string(userRes))-1]
	// remove opening brace : '{'
	modifiedCustomSchemas := string(customSchemas)[1:]
	// append with a comma between the two
	res.WriteString(modifiedUsersRes + "," + modifiedCustomSchemas)

	return res.Bytes()
}
