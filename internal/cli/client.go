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
	Target  string `json:"target"`
	Message string `json:"message"`
}

type ResponseError struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
}

const RequestHeader = "application/json"
const ScimRequestHeader = "application/scim+json"
const ScimUpdateSchemas = "urn:ietf:params:scim:api:messages:2.0:PatchOp"

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

func (c *Client) DoRequest(ctx context.Context, method string, endpoint string, queryStrings map[string]string, body any, customSchemas string, reqHeader string) (*http.Response, error) {
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

	if len(queryStrings) > 0 {
		query := req.URL.Query()
		for k, v := range queryStrings {
			query.Set(k, v)
		}
		req.URL.RawQuery = query.Encode()
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

func (c *Client) Execute(ctx context.Context, method string, endpoint string, queryStrings map[string]string, body any, customSchemas string, reqHeader string, headers []string) (any, map[string]string, error) {

	var O any
	out := make(map[string]string, len(headers))

	res, err := c.DoRequest(ctx, method, endpoint, queryStrings, body, customSchemas, reqHeader)

	if err != nil {
		return nil, out, err
	}

	defer func() {
		if tempErr := res.Body.Close(); tempErr != nil {
			err = tempErr
		}
	}()

	if res.StatusCode >= 400 {

		rawBody, _ := io.ReadAll(res.Body)

		if strings.Contains(reqHeader, "scim") {

			var responseError ScimResponseError

			if err = json.Unmarshal(rawBody, &responseError); err == nil && responseError.Detail != "" {
				err = fmt.Errorf("SCIM error %s \n%s", responseError.Status, responseError.Detail)
			} else {
				err = fmt.Errorf("SCIM error %d \n%s", res.StatusCode, string(rawBody))
			}

		} else {
			var responseError struct {
				Error ResponseError `json:"error"`
			}
			if err = json.Unmarshal(rawBody, &responseError); err == nil && responseError.Error.Message != "" {

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
					if errMessage.Target != "" {
						err = fmt.Errorf("%v : %s %s", err, errMessage.Target, errMessage.Message)
					} else {
						err = fmt.Errorf("%v : %s", err, errMessage.Message)
					}
				}
			} else {
				err = fmt.Errorf("error %d \n%s", res.StatusCode, string(rawBody))
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
