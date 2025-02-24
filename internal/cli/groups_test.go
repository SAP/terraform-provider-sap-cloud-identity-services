package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-ias/internal/cli/apiObjects/groups"
	"testing"

	"github.com/stretchr/testify/assert"
)
 
var groupsPath = "/scim/Groups/"
var groupsBody = groups.Group{
	DisplayName: "Test Group",
	GroupMembers: []groups.GroupMember{
		{
			Value: "scim-id-member-1",
			Type: "User",
		},
		{
			Value: "scim-id-member-2",
			Type: "User",
		},
	},
}

var groupsResponse []byte
 
func TestGroups_Create(t *testing.T) {

	groupsResponse, _ = json.Marshal(groupsBody)

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(groupsResponse)
            assertCall[groups.Group](t, r, groupsPath, "POST", groupsBody)
        }))
 
        defer srv.Close()
 
        _, _, err := client.Group.Create(context.TODO(), &groupsBody)
 
		assert.NoError(t, err)
    }) 
}

func TestGroups_Get(t *testing.T){

    allGroups := []groups.Group{
        groupsBody,
        groupsBody,
    }

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        res, _ := json.Marshal(groups.GroupsResponse{
            Resources: allGroups,
        })

        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(res)
            assertCall[groups.Group](t, r, groupsPath, "GET", nil)
        }))
 
        defer srv.Close()
 
        _, _, err := client.Group.Get(context.TODO())
 
		assert.NoError(t, err)
    }) 
}

func TestGroups_GetByGroupId(t *testing.T){

    groupsResponse, _ = json.Marshal(groupsBody)

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(groupsResponse)
            assertCall[groups.Group](t, r, fmt.Sprintf("%s%s", groupsPath, "valid-group-id"), "GET", nil)
        }))
 
        defer srv.Close()
 
        _, _, err := client.Group.GetByGroupId(context.TODO(), "valid-group-id")
 
		assert.NoError(t, err)
    }) 
}

func TestGroups_Update(t *testing.T) {

	groupsBody.Id = "valid-group-id"
    groupsResponse, _ := json.Marshal(groupsBody)

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(groupsResponse)
            assertCall[groups.Group](t, r, fmt.Sprintf("%s%s", groupsPath, "valid-group-id"), "PUT", groupsBody)
        }))
 
        defer srv.Close()
 
        _, _, err := client.Group.Update(context.TODO(), &groupsBody)
 
		assert.NoError(t, err)
    }) 
}

func TestGroups_Delete(t *testing.T) {

    t.Run("construct the request body correctly", func(t *testing.T) {
 
        client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            assertCall[groups.Group](t, r, fmt.Sprintf("%s%s", groupsPath, "valid-group-id"), "DELETE", nil)
        }))
 
        defer srv.Close()
 
        err := client.Group.Delete(context.TODO(), "valid-group-id")
 
		assert.NoError(t, err)
    }) 

}