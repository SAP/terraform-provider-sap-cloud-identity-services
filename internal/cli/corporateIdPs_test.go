package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
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

	saml2ConfigBody = corporateidps.SAML2Configuration{
		SamlMetadataUrl: "https://example.com/saml2-metadata",
		AssertionAttributes: []corporateidps.AssertionAttribute{
			{
				Name:  "attr_name",
				Value: "attr_value",
			},
		},
		DigestAlgorithm:     "sha1",
		IncludeScoping:      true,
		DefaultNameIdFormat: "email",
		AllowCreate:         "true",
		CertificatesForSigning: []corporateidps.SigningCertificateData{
			{
				Base64Certificate: "redacted",
				IsDefault:         true,
				Dn:                "Test",
				ValidFrom:         "1999-01-01T00:00:00Z",
				ValidTo:           "9999-12-31T23:59:59Z",
			},
		},
		SsoEndpoints: []corporateidps.SAML2SSOEndpoint{
			{
				BindingName: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
				Location:    "https://test.com",
				IsDefault:   true,
			},
		},
		SloEndpoints: []corporateidps.SAML2SLOEndpoint{
			{
				BindingName:      "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
				Location:         "https://test.com",
				ResponseLocation: "https://test.com",
				IsDefault:        true,
			},
		},
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

	saml2IdP corporateidps.IdentityProvider
	oidcIdP  corporateidps.IdentityProvider
)

func TestCorporateIdPs_Create(t *testing.T) {

	saml2IdP = corporateIdPsBody
	saml2IdP.Saml2Configuration = &saml2ConfigBody

	oidcIdP = corporateIdPsBody
	oidcIdP.OidcConfiguration = &oidcConfigBody
	t.Run("validate the API request - saml2 IdP", func(t *testing.T) {

		saml2IdPResponse, _ := json.Marshal(saml2IdP)

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("location", corporateIdPsPath+"valid-idp-uuid")
			_, err := w.Write(saml2IdPResponse)
			assert.NoError(t, err, "Failed to write response")

			if r.Method != "GET" {
				assertCall[corporateidps.IdentityProvider](t, r, corporateIdPsPath, "POST", saml2IdP)
			}
		}))

		defer srv.Close()

		_, _, err := client.CorporateIdP.Create(context.TODO(), &saml2IdP)

		assert.NoError(t, err)
	})

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
		assert.Equal(t, "error 400 \ncreate failed : server error", err.Error())

	})

}

func TestCorporateIdPs_Get(t *testing.T) {

	allCorporateIdPs := []corporateidps.IdentityProvider{
		saml2IdP,
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
		assert.Equal(t, "error 400 \nget failed : server error", err.Error())

	})

	t.Run("validate the API - empty response", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    404,
				Message: "Unable to find identity providers.",
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

		assert.NoError(t, err)
		assert.Equal(t, corporateidps.IdentityProvidersResponse{}, res)

	})
}

func TestCorporateIdPs_GetByIdPId(t *testing.T) {

	saml2IdP = corporateIdPsBody
	saml2IdP.Saml2Configuration = &saml2ConfigBody

	oidcIdP = corporateIdPsBody
	oidcIdP.OidcConfiguration = &oidcConfigBody

	t.Run("validate the API request - saml2 IdP", func(t *testing.T) {

		saml2IdPResponse, _ := json.Marshal(saml2IdP)

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(saml2IdPResponse)

			assert.NoError(t, err, "Failed to write response")

			assertCall[corporateidps.IdentityProvider](t, r, fmt.Sprintf("%s%s", corporateIdPsPath, "valid-idp-id"), "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.CorporateIdP.GetByIdPId(context.TODO(), "valid-idp-id")
		assert.NoError(t, err)
	})

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
		assert.Equal(t, "error 400 \nget failed : server error", err.Error())
	})
}

func TestCorporateIdPs_Update(t *testing.T) {

	originalInterval := corporateIdPPollInterval
	defer func() { corporateIdPPollInterval = originalInterval }()
	corporateIdPPollInterval = time.Millisecond

	saml2IdP = corporateIdPsBody
	saml2IdP.Saml2Configuration = &saml2ConfigBody

	patchOps := []generic.PatchRequest{
		{
			Op:    "replace",
			Path:  "/displayName",
			Value: "Updated Corporate IdP",
		},
	}

	expectedPatchBody := corporateidps.PatchRequestBody{
		Operations: patchOps,
	}

	updatedIdP := saml2IdP
	updatedIdP.DisplayName = "Updated Corporate IdP"
	updatedIdPResponse, _ := json.Marshal(updatedIdP)
	oldIdPResponse, _ := json.Marshal(saml2IdP)

	t.Run("validate the API request - success", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PATCH" {
				assertCall[corporateidps.PatchRequestBody](t, r, fmt.Sprintf("%s%s", corporateIdPsPath, "valid-idp-id"), "PATCH", expectedPatchBody)
				w.WriteHeader(http.StatusOK)
			} else {
				_, err := w.Write(updatedIdPResponse)
				assert.NoError(t, err, "Failed to write response")
			}
		}))

		defer srv.Close()

		res, _, err := client.CorporateIdP.Update(context.TODO(), patchOps, "valid-idp-id")

		assert.NoError(t, err)
		assert.Equal(t, "Updated Corporate IdP", res.DisplayName)
	})

	t.Run("validate polling - GET returns updated value on retry", func(t *testing.T) {

		callCount := 0
		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PATCH" {
				w.WriteHeader(http.StatusOK)
			} else {
				callCount++
				if callCount == 1 {
					_, err := w.Write(oldIdPResponse)
					assert.NoError(t, err, "Failed to write response")
				} else {
					_, err := w.Write(updatedIdPResponse)
					assert.NoError(t, err, "Failed to write response")
				}
			}
		}))

		defer srv.Close()

		res, _, err := client.CorporateIdP.Update(context.TODO(), patchOps, "valid-idp-id")

		assert.NoError(t, err)
		assert.Equal(t, "Updated Corporate IdP", res.DisplayName)
		assert.Equal(t, 2, callCount)
	})

	t.Run("validate the API request - error on PATCH", func(t *testing.T) {

		resErr, _ := json.Marshal(struct {
			Error ResponseError `json:"error"`
		}{
			Error: ResponseError{
				Code:    400,
				Message: "update failed",
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

			assertCall[corporateidps.PatchRequestBody](t, r, fmt.Sprintf("%s%s", corporateIdPsPath, "valid-idp-id"), "PATCH", expectedPatchBody)
		}))

		defer srv.Close()

		res, _, err := client.CorporateIdP.Update(context.TODO(), patchOps, "valid-idp-id")

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "error 400 \nupdate failed : server error", err.Error())
	})

	t.Run("validate the API request - error on GET after PATCH", func(t *testing.T) {

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

		callCount := 0
		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PATCH" {
				w.WriteHeader(http.StatusOK)
			} else {
				callCount++
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write(resErr)
				assert.NoError(t, err, "Failed to write response")
			}
		}))

		defer srv.Close()

		res, _, err := client.CorporateIdP.Update(context.TODO(), patchOps, "valid-idp-id")

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "error 400 \nget failed : server error", err.Error())
		assert.Equal(t, 1, callCount)
	})

	t.Run("validate write-only paths are skipped during polling", func(t *testing.T) {

		secretPatchOps := []generic.PatchRequest{
			{
				Op:    "replace",
				Path:  "/oidcConfiguration/clientSecret",
				Value: "new-secret",
			},
		}

		oidcIdP = corporateIdPsBody
		oidcIdP.OidcConfiguration = &oidcConfigBody
		oidcIdPResponse, _ := json.Marshal(oidcIdP)

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(oidcIdPResponse)
			assert.NoError(t, err, "Failed to write response")
		}))

		defer srv.Close()

		_, _, err := client.CorporateIdP.Update(context.TODO(), secretPatchOps, "valid-idp-id")
		assert.NoError(t, err)
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
		assert.Equal(t, "error 400 \ndelete failed : server error", err.Error())
	})
}
