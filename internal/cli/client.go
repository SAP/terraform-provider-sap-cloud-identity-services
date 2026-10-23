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
			// Validate customSchemas is a well-formed JSON object (not an array, string, etc.)
			var extra map[string]json.RawMessage
			if err := json.Unmarshal([]byte(customSchemas), &extra); err != nil {
				return nil, fmt.Errorf("custom_schemas is not a valid JSON object: %w", err)
			}
			// Merge by appending extra fields into the base JSON body.
			// We operate on bytes rather than re-encoding through a map so that
			// the field order from the Go struct is preserved (required for VCR cassette replay).
			base := bytes.TrimRight(encodedBody.Bytes(), "\n")
			lastBrace := bytes.LastIndexByte(base, '}')
			if lastBrace < 0 {
				return nil, fmt.Errorf("request body is not a valid JSON object")
			}
			inner := bytes.TrimSpace([]byte(customSchemas))
			inner = bytes.TrimSpace(inner[1 : len(inner)-1]) // strip outer { }
			encodedBody.Reset()
			if len(inner) == 0 {
				encodedBody.Write(base[:lastBrace+1])
			} else {
				encodedBody.Write(base[:lastBrace])
				encodedBody.WriteByte(',')
				encodedBody.Write(inner)
				encodedBody.WriteByte('}')
			}
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

		rawBody, _ := io.ReadAll(io.LimitReader(res.Body, 10*1024*1024))

		// Truncate raw body to avoid echoing large or sensitive server responses verbatim.
		const maxDisplayBytes = 512
		displayBody := string(rawBody)
		if len(displayBody) > maxDisplayBytes {
			displayBody = displayBody[:maxDisplayBytes] + "... [truncated]"
		}

		if strings.Contains(reqHeader, "scim") {

			var responseError ScimResponseError

			if err = json.Unmarshal(rawBody, &responseError); err == nil && responseError.Detail != "" {
				err = fmt.Errorf("SCIM error %s \n%s", responseError.Status, responseError.Detail)
			} else {
				err = fmt.Errorf("SCIM error %d \n%s", res.StatusCode, displayBody)
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
				err = fmt.Errorf("error %d \n%s", res.StatusCode, displayBody)
			}

		}

		return nil, out, err
	}

	for _, header := range headers {
		out[header] = res.Header.Get(header)
	}

	if err = json.NewDecoder(io.LimitReader(res.Body, 50*1024*1024)).Decode(&O); err == nil || err == io.EOF {
		return O, out, nil
	} else {
		return nil, out, err
	}
}
