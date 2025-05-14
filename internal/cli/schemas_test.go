package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/schemas"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var schemasPath = "/scim/Schemas/"
var schemasBody = schemas.Schema{
	Id:   "valid-schema-id",
	Name: "Test Schema",
	Attributes: []schemas.Attribute{
		{
			Name:       "test_attribute_1",
			Type:       "string",
			Mutability: "readWrite",
			Returned:   "never",
			Uniqueness: "none",
		},
		{
			Name:       "test_attribute_2",
			Type:       "bool",
			Mutability: "immutable",
			Returned:   "default",
			Uniqueness: "none",
		},
	},
}

var schemasResponse []byte

func TestSchemas_Create(t *testing.T) {

	schemasResponse, _ = json.Marshal(schemasBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(schemasResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[schemas.Schema](t, r, schemasPath, "POST", schemasBody)
		}))

		defer srv.Close()

		_, _, err := client.Schema.Create(context.TODO(), &schemasBody)

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

			assertCall[schemas.Schema](t, r, schemasPath, "POST", schemasBody)
		}))

		defer srv.Close()

		res, _, err := client.Schema.Create(context.TODO(), &schemasBody)

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "create failed", err.Error())
	})
}

func TestSchemas_Get(t *testing.T) {

	allSchemas := []schemas.Schema{
		schemasBody,
		schemasBody,
	}

	t.Run("validate the API request", func(t *testing.T) {

		res, _ := json.Marshal(schemas.SchemasResponse{
			Resources: allSchemas,
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(res)
			assert.NoError(t, err, "Failed to write response")

			assertCall[schemas.Schema](t, r, schemasPath, "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.Schema.Get(context.TODO())

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimError{
			Detail: "get failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[schemas.Schema](t, r, schemasPath, "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.Schema.Get(context.TODO())

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "get failed", err.Error())
	})
}

func TestSchemas_GetBySchemaId(t *testing.T) {

	schemasResponse, _ = json.Marshal(schemasBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(schemasResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[schemas.Schema](t, r, fmt.Sprintf("%s%s", schemasPath, "valid-schema-id"), "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.Schema.GetBySchemaId(context.TODO(), "valid-schema-id")

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimError{
			Detail: "get failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[schemas.Schema](t, r, fmt.Sprintf("%s%s", schemasPath, "valid-schema-id"), "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.Schema.GetBySchemaId(context.TODO(), "valid-schema-id")

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "get failed", err.Error())
	})
}

func TestSchemas_Delete(t *testing.T) {

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertCall[schemas.Schema](t, r, fmt.Sprintf("%s%s", schemasPath, "valid-schema-id"), "DELETE", schemas.Schema{})
		}))

		defer srv.Close()

		err := client.Schema.Delete(context.TODO(), "valid-schema-id")

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(ScimError{
			Detail: "delete failed",
			Status: "400",
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[schemas.Schema](t, r, fmt.Sprintf("%s%s", schemasPath, "valid-schema-id"), "DELETE", nil)
		}))

		defer srv.Close()

		err := client.Schema.Delete(context.TODO(), "valid-schema-id")

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
	})

}
