package cli

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"io"

	"encoding/json"
	"net/http"

	"net/url"
)

type ScimError struct {
	Detail  string   `json:"detail"`
	Schemas []string `json:"schemas"`
	Status  string   `json:"status"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type ApplicationError struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
}

const ApplicationHeader = "application/json"
const DirectoryHeader = "application/scim+json"

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

		if len(customSchemas) > 0 {

			// remove the ending characters '}\n' from the encoded buffer
			body := encodedBody.String()[:len(encodedBody.String())-2]
			// remove the beginning character '{' from the custom schemas string
			customSchemas = customSchemas[1:]

			// reset the encoded buffer and concatenate the custom schemas with the rest of the request body
			encodedBody.Reset()
			encodedBody.WriteString(body + "," + customSchemas)
		}

		if err != nil {
			return nil, err
		}
	}

	completeUrl := c.ServerURL.ResolveReference(parsedUrl)

	req, err := http.NewRequestWithContext(ctx, method, completeUrl.String(), &encodedBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+c.AuthorizationToken)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("DataServiceVersion", "2.0")
	req.Header.Set("Content-Type", reqHeader)

	res, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Execute(ctx context.Context, method string, endpoint string, body any, customSchemas string, reqHeader string, headers []string) (interface{}, error, map[string]string) {

	var O interface{}
	out := make(map[string]string, len(headers))

	res, err := c.DoRequest(ctx, method, endpoint, body, customSchemas, reqHeader)

	if err != nil {
		return nil, err, out
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {

		if strings.Contains(reqHeader, "scim") {

			var responseError ScimError

			if err = json.NewDecoder(res.Body).Decode(&responseError); err == nil {
				err = fmt.Errorf("%s", responseError.Detail)
			} else {
				err = fmt.Errorf("responded with unknown error : %s", responseError.Status)
			}

		} else {
			var responseError struct {
				Error ApplicationError `json:"error"`
			}
			if err = json.NewDecoder(res.Body).Decode(&responseError); err == nil {
				err = fmt.Errorf("%s", responseError.Error.Message)

				for _, errMessage := range responseError.Error.Details{
					err = fmt.Errorf("%v : %s", err, errMessage.Message)
				}
			} else {
				err = fmt.Errorf("responded with unknown error : %d", responseError.Error.Code)
			}

		}

		return nil, err, out
	}

	for _, header := range headers {
		out[header] = res.Header.Get(header)
	}

	if err = json.NewDecoder(res.Body).Decode(&O); err == nil || err == io.EOF {
		return O, nil, out
	} else {
		return nil, err, out
	}
}
