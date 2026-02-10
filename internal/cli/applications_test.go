package cli

import (
	"context"
	"encoding/json"

	"fmt"
	"net/http"
	"testing"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"

	"github.com/stretchr/testify/assert"
)

var applicationsPath = "/Applications/v1/"
var applicationsBody = applications.Application{
	Name:        "Test Application",
	Description: "test app",
	AuthenticationSchema: &applications.AuthenticationSchema{
		SsoType:               "saml",
		SubjectNameIdentifier: "mail",
		AssertionAttributes: []applications.AssertionAttribute{
			{
				AssertionAttributeName: "attr_name_1",
				UserAttributeName:      "user_name_1",
			},
			{
				AssertionAttributeName: "attr_name_2",
				UserAttributeName:      "user_name_2",
			},
		},
		AdvancedAssertionAttributes: []applications.AdvancedAssertionAttribute{
			{
				AttributeName:  "attr_name_1",
				AttributeValue: "attr_value_1",
			},
			{
				AttributeName:  "attr_name_2",
				AttributeValue: "attr_value_2",
			},
		},
		DefaultAuthenticatingIdpId: "valid-idp-uuid",
		ConditionalAuthentication: []applications.AuthenicationRule{
			{
				UserType:           "employee",
				IdentityProviderId: "valid-idp-uuid",
			},
			{
				UserGroup:       "valid-group-id",
				UserEmailDomain: "test.com",
			},
		},
		OidcConfig: &applications.OidcConfig{
			RedirectUris: []string{
				"https:redirectUris.com",
			},
			PostLogoutRedirectUris: []string{
				"https:postlogoutRedirectUris.com",
			},
			FrontChannelLogoutUris: []string{
				"https:frontChannelLogoutUris.com",
			},
			BackChannelLogoutUris: []string{
				"https:backChannelLogoutUris.com",
			},
			TokenPolicy: &applications.TokenPolicy{
				JwtValidity:                  3600,
				RefreshValidity:              43200,
				RefreshParallel:              1,
				MaxExchangePeriod:            "unlimited",
				RefreshTokenRotationScenario: "off",
				AccessTokenFormat:            "default",
			},
			RestrictedGrantTypes: []applications.GrantType{
				"clientCredentials",
				"authorizationCode",
				"password",
				"jwtBearer",
			},
			ProxyConfig: &applications.OidcProxyConfig{
				Acrs: []string{
					"acrs-1",
				},
			},
		},
		Saml2Configuration: &applications.SamlConfiguration{
			SamlMetadataUrl: "https://test.com",
			AcsEndpoints: []applications.Saml2AcsEndpoint{
				{
					BindingName: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
					Location:    "https://test.1.com",
					Index:       1,
					IsDefault:   true,
				},
			},
			SloEndpoints: []applications.Saml2SLOEndpoint{
				{
					BindingName:      "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
					Location:         "https://logout.1.com",
					ResponseLocation: "https://logout-response.1.com",
				},
			},
			CertificatesForSigning: []corporateidps.SigningCertificateData{
				{
					Base64Certificate: "-----BEGIN CERTIFICATE-----\\n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\\n-----END CERTIFICATE-----",
					Dn:                "CN=Test Cert",
					IsDefault:         true,
				},
			},
			CertificateForEncryption: &applications.EncryptionCertificateData{
				Base64Certificate: "-----BEGIN CERTIFICATE-----\\n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\\n-----END CERTIFICATE-----",
			},
			ResponseElementsToEncrypt: "wholeAssertion",
			DefaultNameIdFormat:       "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
			SignSLOMessages:           true,
			RequireSignedSLOMessages:  true,
			RequireSignedAuthnRequest: true,
			SignAssertions:            true,
			SignAuthnResponses:        true,
			DigestAlgorithm:           "sha1",
		},
	},
}

var applicationsResponse []byte

func TestApplications_Create(t *testing.T) {

	applicationsResponse, _ = json.Marshal(applicationsBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("location", "/Applications/v1/valid-app-id")
			_, err := w.Write(applicationsResponse)
			assert.NoError(t, err, "Failed to write response")

			if r.Method != "GET" {
				assertCall[applications.Application](t, r, applicationsPath, "POST", applicationsBody)
			}

		}))

		defer srv.Close()

		_, _, err := client.Application.Create(context.TODO(), &applicationsBody)

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

			assertCall[applications.Application](t, r, applicationsPath, "POST", applicationsBody)
		}))

		defer srv.Close()

		res, _, err := client.Application.Create(context.TODO(), &applicationsBody)

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "error 400 \ncreate failed : server error", err.Error())
	})
}

func TestApplications_Get(t *testing.T) {

	allApplications := []applications.Application{
		applicationsBody,
		applicationsBody,
	}

	t.Run("validate the API request", func(t *testing.T) {

		res, _ := json.Marshal(applications.ApplicationsResponse{
			Applications: allApplications,
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(res)
			assert.NoError(t, err, "Failed to write response")

			assertCall[applications.Application](t, r, applicationsPath, "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.Application.Get("", context.TODO())

		assert.NoError(t, err)
	})

	t.Run("validate the API request - with cursor", func(t *testing.T) {

		resWithCursor, _ := json.Marshal(applications.ApplicationsResponse{
			NextCursor: "test",
			Applications: allApplications,
		})

		resWithoutCursor, _ := json.Marshal(applications.ApplicationsResponse{
			Applications: allApplications,
		})

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			if r.URL.Query().Get("cursor") == "" {
				_, err := w.Write(resWithCursor)
				assert.NoError(t, err, "Failed to write response")
			} else {
				_, err := w.Write(resWithoutCursor)
				assert.NoError(t, err, "Failed to write response")
			}
			

			assertCall[applications.Application](t, r, applicationsPath, "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.Application.Get("", context.TODO())

		assert.NoError(t, err)
		assert.Len(t, res.Applications, 4)
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

			assertCall[applications.Application](t, r, applicationsPath, "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.Application.Get("", context.TODO())

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "error 400 \nget failed : server error", err.Error())
	})
}

func TestApplications_GetByAppId(t *testing.T) {

	applicationsResponse, _ = json.Marshal(applicationsBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(applicationsResponse)
			assert.NoError(t, err, "Failed to write response")

			assertCall[applications.Application](t, r, fmt.Sprintf("%s%s", applicationsPath, "valid-app-id"), "GET", nil)
		}))

		defer srv.Close()

		_, _, err := client.Application.GetByAppId(context.TODO(), "valid-app-id")
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

			assertCall[applications.Application](t, r, fmt.Sprintf("%s%s", applicationsPath, "valid-app-id"), "GET", nil)
		}))

		defer srv.Close()

		res, _, err := client.Application.GetByAppId(context.TODO(), "valid-app-id")

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "error 400 \nget failed : server error", err.Error())
	})
}

func TestApplications_Update(t *testing.T) {

	applicationsBody.Id = "valid-app-id"
	applicationsResponse, _ := json.Marshal(applicationsBody)

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(applicationsResponse)
			assert.NoError(t, err, "Failed to write response")

			if r.Method != "GET" {
				assertCall[applications.Application](t, r, fmt.Sprintf("%s%s", applicationsPath, "valid-app-id"), "PUT", applicationsBody)
			}
		}))

		defer srv.Close()

		_, _, err := client.Application.Update(context.TODO(), &applicationsBody)

		assert.NoError(t, err)
	})

	t.Run("validate the API request - error", func(t *testing.T) {

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

			assertCall[applications.Application](t, r, fmt.Sprintf("%s%s", applicationsPath, "valid-app-id"), "PUT", applicationsBody)
		}))

		defer srv.Close()

		res, _, err := client.Application.Update(context.TODO(), &applicationsBody)

		assert.Zero(t, res)
		assert.Error(t, err)
		assert.Equal(t, "error 400 \nupdate failed : server error", err.Error())
	})
}

func TestApplications_Delete(t *testing.T) {

	t.Run("validate the API request", func(t *testing.T) {

		client, srv := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertCall[applications.Application](t, r, fmt.Sprintf("%s%s", applicationsPath, "valid-app-id"), "DELETE", nil)
		}))

		defer srv.Close()

		err := client.Application.Delete(context.TODO(), "valid-app-id")

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

			assertCall[applications.Application](t, r, fmt.Sprintf("%s%s", applicationsPath, "valid-app-id"), "DELETE", applicationsBody)
		}))

		defer srv.Close()

		err := client.Application.Delete(context.TODO(), "valid-app-id")

		assert.Error(t, err)
		assert.Equal(t, "error 400 \ndelete failed : server error", err.Error())
	})
}
