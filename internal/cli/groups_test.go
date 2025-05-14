package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"
	"testing"

	"github.com/stretchr/testify/assert"
)

var groupsPath = "/scim/Groups/"
var groupsBody = groups.Group{
	DisplayName: "Test Group",
	GroupMembers: []groups.GroupMember{
		{
			Value: "scim-id-member-1",
			Type:  "User",
		},
		{
			Value: "scim-id-member-2",
			Type:  "User",
		},
	},
}

var groupsResponse []byte

func TestGroups_Create(t *testing.T) {

	groupsResponse, _ = json.Marshal(groupsBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(groupsResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[groups.Group](t, r, groupsPath, "POST", groupsBody)
		}))

		defer srv.Close()

		_, _, err := client.Group.Create(context.TODO(), &groupsBody)

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimError{
			Detail: "create failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[groups.Group](t, r, groupsPath, "POST", groupsBody)
		}))

		defer srv.Close()

		res, _, err := client.Group.Create(context.TODO(), &groupsBody)

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "create failed", err.Error())
	})
}

func TestGroups_Get(t *testing.T) {

	allGroups := []groups.Group{
		groupsBody,
		groupsBody,
	}

	t.Run("validate the API request", func(t *testing.T) {

		res, _ := json.Marshal(groups.GroupsResponse{
			Resources: allGroups,
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(res)
			assert.NoError(t, err, "Failed to write response")

			assertCall[groups.Group](t, r, groupsPath, "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.Group.Get(context.TODO())

		assert.NoError(t, err)
	})

	t.Run("validate the API request with error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimError{
			Detail: "get failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[groups.Group](t, r, groupsPath, "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.Group.Get(context.TODO())

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "get failed", err.Error())
	})
}

func TestGroups_GetByGroupId(t *testing.T) {

	groupsResponse, _ = json.Marshal(groupsBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(groupsResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[groups.Group](t, r, fmt.Sprintf("%s%s", groupsPath, "valid-group-id"), "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.Group.GetByGroupId(context.TODO(), "valid-group-id")

		assert.NoError(t, err)
	})

	t.Run("validate the API request with error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimError{
			Detail: "get failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[groups.Group](t, r, fmt.Sprintf("%s%s", groupsPath, "valid-group-id"), "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.Group.GetByGroupId(context.TODO(), "valid-group-id")

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "get failed", err.Error())
	})
}

func TestGroups_Update(t *testing.T) {

	groupsBody.Id = "valid-group-id"
	groupsResponse, _ := json.Marshal(groupsBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(groupsResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[groups.Group](t, r, fmt.Sprintf("%s%s", groupsPath, "valid-group-id"), "PUT", groupsBody)
		}))

		defer srv.Close()

		_, _, err := client.Group.Update(context.TODO(), &groupsBody)

		assert.NoError(t, err)
	})

	t.Run("validate the API request with error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimError{
			Detail: "update failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[groups.Group](t, r, fmt.Sprintf("%s%s", groupsPath, "valid-group-id"), "PUT", groupsBody)
		}))

		defer srv.Close()

		res, _, err := client.Group.Update(context.TODO(), &groupsBody)

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "update failed", err.Error())
	})
}

func TestGroups_Delete(t *testing.T) {

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertCall[groups.Group](t, r, fmt.Sprintf("%s%s", groupsPath, "valid-group-id"), "DELETE", nil)
		}))

		defer srv.Close()

		err := client.Group.Delete(context.TODO(), "valid-group-id")

		assert.NoError(t, err)
	})

	t.Run("validate the API request with error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimError{
			Detail: "delete failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[groups.Group](t, r, fmt.Sprintf("%s%s", groupsPath, "valid-group-id"), "DELETE", groups.Group{})
		}))

		defer srv.Close()

		err := client.Group.Delete(context.TODO(), "valid-group-id")

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
	})
}
