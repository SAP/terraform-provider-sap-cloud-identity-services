package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/stretchr/testify/assert"
)

var (
	corporateIdPsPath = "/IdentityProviders/v1/"

	corporateIdPsBody = corporateidps.IdentityProvider{
		DisplayName:           "Test Corporate IdP",
		Name:                  "Test IdP",
		ForwardAllSsoRequests: true,
		IdentityFederation: &corporateidps.IdentityFederation{
			UseLocalUserStore:        true,
			AllowLocalUsersOnly:      true,
			ApplyLocalIdPAuthnChecks: true,
			RequiredGroups: []string{
				"Test Group",
			},
		},
		LoginHintConfiguration: &corporateidps.LoginHintConfiguration{
			LoginHintType: "mail",
			SendMethod:    "urlParam",
		},
		LogoutUrl: "https://example.com/logout",
	}

	oidcConfigBody = corporateidps.OIDCConfiguration{
		DiscoveryUrl:            "https://test.com",
		ClientId:                "test-client-id",
		ClientSecret:            "test-client-secret",
		SubjectNameIdentifier:   "email",
		TokenEndpointAuthMethod: "clientSecretBasic",
		Scopes: []string{
			"test-value-1",
			"openid",
		},
		PkceEnabled: true,
		AdditionalConfig: &corporateidps.OIDCAdditionalConfig{
			OmitIDTokenHintForLogout: true,
			EnforceIssuerCheck:       true,
			EnforceNonce:             true,
		},
	}

	oidcIdP corporateidps.IdentityProvider
)

func TestCorporateIdPs_Create(t *testing.T) {

	oidcIdP = corporateIdPsBody
	oidcIdP.OidcConfiguration = &oidcConfigBody

	t.Run("validate the API request - oidc IdP", func(t *testing.T) {

		oidcIdPResponse, _ := json.Marshal(oidcIdP)

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("location", corporateIdPsPath+"valid-idp-uuid")
			_, err := w.Write(oidcIdPResponse)
			assert.NoError(t, err, "Failed to write response")

			if r.Method != "GET" {
				assertCall[corporateidps.IdentityProvider](t, r, corporateIdPsPath, "POST", oidcIdP)
			}

		}))

		defer srv.Close()

		_, _, err := client.CorporateIdP.Create(context.TODO(), &oidcIdP)

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "create failed",
				Details: []ErrorDetail{
					{
						Message: "server error",
					},
				},
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[corporateidps.IdentityProvider](t, r, corporateIdPsPath, "POST", corporateIdPsBody)
		}))

		defer srv.Close()

		res, _, err := client.CorporateIdP.Create(context.TODO(), &corporateIdPsBody)

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "application error 400 \ncreate failed : server error", err.Error())

	})

}

func TestCorporateIdPs_Get(t *testing.T) {

	allCorporateIdPs := []corporateidps.IdentityProvider{
		oidcIdP,
	}

	t.Run("validate the API request", func(t *testing.T) {

		res, _ := json.Marshal(corporateidps.IdentityProvidersResponse{
			IdentityProviders: allCorporateIdPs,
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(res)
			assert.NoError(t, err, "Failed to write response")

			assertCall[corporateidps.IdentityProvider](t, r, corporateIdPsPath, "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.CorporateIdP.Get(context.TODO())

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "get failed",
				Details: []ErrorDetail{
					{
						Message: "server error",
					},
				},
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[corporateidps.IdentityProvider](t, r, corporateIdPsPath, "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.CorporateIdP.Get(context.TODO())

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "application error 400 \nget failed : server error", err.Error())

	})
}

func TestCorporateIdPs_GetByIdPId(t *testing.T) {

	oidcIdP = corporateIdPsBody
	oidcIdP.OidcConfiguration = &oidcConfigBody

	t.Run("validate the API request - oidc IdP", func(t *testing.T) {

		oidcIdPResponse, _ := json.Marshal(oidcIdP)

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(oidcIdPResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[corporateidps.IdentityProvider](t, r, fmt.Sprintf("%s%s", corporateIdPsPath, "valid-idp-id"), "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.CorporateIdP.GetByIdPId(context.TODO(), "valid-idp-id")
		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "get failed",
				Details: []ErrorDetail{
					{
						Message: "server error",
					},
				},
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[corporateidps.IdentityProvider](t, r, fmt.Sprintf("%s%s", corporateIdPsPath, "valid-idp-id"), "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.CorporateIdP.GetByIdPId(context.TODO(), "valid-idp-id")

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "application error 400 \nget failed : server error", err.Error())
	})
}

func TestCorporateIdPs_Delete(t *testing.T) {

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertCall[corporateidps.IdentityProvider](t, r, fmt.Sprintf("%s%s", corporateIdPsPath, "valid-idp-id"), "DELETE", nil)
		}))

		defer srv.Close()

		err := client.CorporateIdP.Delete(context.TODO(), "valid-idp-id")

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "delete failed",
				Details: []ErrorDetail{
					{
						Message: "server error",
					},
				},
			},
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(resErr)
			assert.NoError(t, err, "Failed to write response")

			assertCall[corporateidps.IdentityProvider](t, r, fmt.Sprintf("%s%s", corporateIdPsPath, "valid-idp-id"), "DELETE", applicationsBody)
		}))

		defer srv.Close()

		err := client.CorporateIdP.Delete(context.TODO(), "valid-idp-id")

		assert.Error(t, err)
		assert.Equal(t, "application error 400 \ndelete failed : server error", err.Error())
	})
}
