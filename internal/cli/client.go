package cli

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"io"

	"encoding/json"
	"net/http"

	"net/url"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
)

// Error message for cases when the GET call for listing resources fails with a 404
// This occurs when there are no resources to be listed
var emptyResponseError, _ = regexp.Compile("Unable to find (.+)")

type ScimResponseError struct {
	Detail  string   `json:"detail"`
	Schemas []string `json:"schemas"`
	Status  string   `json:"status"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
}

const RequestHeader = "application/json"
const ScimRequestHeader = "application/scim+json"

func NewClient(h *http.Client, u *url.URL) *Client {
	return &Client{
		HttpClient: h,
		ServerURL:  u,
	}
}

type Client struct {
	HttpClient         *http.Client
	ServerURL          *url.URL
	AuthorizationToken string
}

func (c *Client) DoRequest(ctx context.Context, method string, endpoint string, body any, customSchemas string, reqHeader string) (*http.Response, error) {
	parsedUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	var encodedBody bytes.Buffer
	if body != nil {
		encoder := json.NewEncoder(&encodedBody)
		err := encoder.Encode(body)
		if err != nil {
			return nil, err
		}

		if len(customSchemas) > 0 {
			// remove the ending characters '}\n' from the encoded buffer
			body := encodedBody.String()[:len(encodedBody.String())-2]
			// remove the beginning character '{' from the custom schemas string
			customSchemas = customSchemas[1:]

			// reset the encoded buffer and concatenate the custom schemas with the rest of the request body
			encodedBody.Reset()
			encodedBody.WriteString(body + "," + customSchemas)
		}
	}

	completeUrl := c.ServerURL.ResolveReference(parsedUrl)

	req, err := http.NewRequestWithContext(ctx, method, completeUrl.String(), &encodedBody)
	if err != nil {
		return nil, err
	}

	// Only set Authorization header if it's not empty
	if c.AuthorizationToken != "" {
		req.Header.Set("Authorization", c.AuthorizationToken)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("DataServiceVersion", "2.0")
	req.Header.Set("Content-Type", reqHeader)

	return c.HttpClient.Do(req)
}

func (c *Client) Execute(ctx context.Context, method string, endpoint string, body any, customSchemas string, reqHeader string, headers []string) (interface{}, map[string]string, error) {

	var O interface{}
	out := make(map[string]string, len(headers))

	res, err := c.DoRequest(ctx, method, endpoint, body, customSchemas, reqHeader)

	if err != nil {
		return nil, out, err
	}

	defer func() {
		if tempErr := res.Body.Close(); tempErr != nil {
			err = tempErr
		}
	}()

	if res.StatusCode >= 400 {

		if strings.Contains(reqHeader, "scim") {

			var responseError ScimResponseError

			if err = json.NewDecoder(res.Body).Decode(&responseError); err == nil {
				// For  users and schemas can never be empty, hence no explicit handling is needed
				// For groups, the GET call does not throw an error in case of an empty resource list, hence no explicit handling is needed
				err = fmt.Errorf("SCIM error %s \n%s", responseError.Status, responseError.Detail)
			} else {
				err = fmt.Errorf("responded with unknown error : %s", responseError.Status)
			}

		} else {
			var responseError struct {
				Error ResponseError `json:"error"`
			}
			if err = json.NewDecoder(res.Body).Decode(&responseError); err == nil {

				// check the error message and status code to handle the situation when
				// an error is thrown on a GET call for returning an empty list of resources
				if emptyResponseError.MatchString(responseError.Error.Message) && responseError.Error.Code == 404 {

					// fetch the type of the resource
					val := emptyResponseError.FindStringSubmatch(responseError.Error.Message)

					// check the resource and return the appropriate empty response object
					// Applications can never be empty, hence there is no check
					// TODO add cases for other non-scim resources once part of the provider
					switch val[1] {
					case "identity providers.":
						return corporateidps.IdentityProvidersResponse{}, out, nil
					}

				}

				err = fmt.Errorf("error %d \n%s", responseError.Error.Code, responseError.Error.Message)

				for _, errMessage := range responseError.Error.Details {
					err = fmt.Errorf("%v : %s", err, errMessage.Message)
				}
			} else {
				err = fmt.Errorf("responded with unknown error : %d", responseError.Error.Code)
			}

		}

		return nil, out, err
	}

	for _, header := range headers {
		out[header] = res.Header.Get(header)
	}

	if err = json.NewDecoder(res.Body).Decode(&O); err == nil || err == io.EOF {
		return O, out, nil
	} else {
		return nil, out, err
	}
}
