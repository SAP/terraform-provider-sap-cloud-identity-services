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
		assert.Equal(t, "application error 400 \ncreate failed : server error", err.Error())

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
		assert.Equal(t, "application error 400 \nget failed : server error", err.Error())

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
