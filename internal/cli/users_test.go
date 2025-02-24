package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-ias/internal/cli/apiObjects/users"
	"testing"

	"github.com/stretchr/testify/assert"
)

var usersPath = "/scim/Users/"
var usersBody = users.User{
	UserName: "user_for_testing",
	Name: users.Name{
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

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(usersResponse)
			assertCall[users.User](t, r, usersPath, "POST", usersBody)
		}))

		defer srv.Close()

		_, _, err := client.User.Create(context.TODO(), "", &usersBody)

		assert.NoError(t, err)
	})

	t.Run("validate the API request with custom schemas", func(t *testing.T) {

		customSchemas, _ := json.Marshal(map[string]interface{}{
			"schema_id": map[string]interface{}{
				"var1": "test",
				"var2": 1,
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(responseWithCustomSchemas(usersResponse, customSchemas))
			assertCall[users.User](t, r, usersPath, "POST", usersBody)
		}))

		defer srv.Close()

		_, _, err := client.User.Create(context.TODO(), string(customSchemas), &usersBody)

		assert.NoError(t, err)
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
			w.Write(res)
			assertCall[users.User](t, r, usersPath, "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.User.Get(context.TODO())

		assert.NoError(t, err)
	})
}

func TestUsers_GetByUserId(t *testing.T) {

	usersResponse, _ = json.Marshal(usersBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(usersResponse)
			assertCall[users.User](t, r, fmt.Sprintf("%s%s", usersPath, "valid-user-id"), "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.User.GetByUserId(context.TODO(), "valid-user-id")

		assert.NoError(t, err)
	})
}

func TestUsers_Update(t *testing.T) {

	usersBody.Id = "valid-user-id"
	usersResponse, _ := json.Marshal(usersBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(usersResponse)
			assertCall[users.User](t, r, fmt.Sprintf("%s%s", usersPath, "valid-user-id"), "PUT", usersBody)
		}))

		defer srv.Close()

		_, _, err := client.User.Update(context.TODO(), "", &usersBody)

		assert.NoError(t, err)
	})

	t.Run("validate the API request with custom schemas", func(t *testing.T) {

		customSchemas, _ := json.Marshal(map[string]interface{}{
			"schema_id": map[string]interface{}{
				"var1": "test",
				"var2": 1,
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(responseWithCustomSchemas(usersResponse, customSchemas))
			assertCall[users.User](t, r, fmt.Sprintf("%s%s", usersPath, "valid-user-id"), "PUT", usersBody)
		}))

		defer srv.Close()

		_, _, err := client.User.Update(context.TODO(), string(customSchemas), &usersBody)

		assert.NoError(t, err)
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
