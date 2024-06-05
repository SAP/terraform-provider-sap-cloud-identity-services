package cli

import (
	"bytes"
	"context"
	"fmt"

	"io"

	// "context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"net/url"
	"os"
)

const ApplicationHeader = "application/json"
const DirectoryHeader = "application/scim+json"

func NewClient(h *http.Client, u *url.URL) *Client {
	return &Client{
		HttpClient: h,
		ServerURL:  u,
	}
}

type Client struct {
	HttpClient *http.Client
	ServerURL  *url.URL
}

func (c *Client) DoRequest(ctx context.Context, method string, endpoint string, body any, reqHeader string) (*http.Response, error) {
	parsedUrl, err := url.Parse(endpoint)

	if err != nil {
		return nil, err
	}

	var username, password string

	username = os.Getenv("ias_username")
	password = os.Getenv("ias_password")

	base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	var encodedBody bytes.Buffer
	if body != nil {
		encoder := json.NewEncoder(&encodedBody)
		err := encoder.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	completeUrl := c.ServerURL.ResolveReference(parsedUrl)

	req, err := http.NewRequestWithContext(ctx, method, completeUrl.String(), &encodedBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+base64Encoded)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("DataServiceVersion", "2.0")
	req.Header.Set("Content-Type", reqHeader)

	res, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Execute(ctx context.Context, method string, endpoint string, body any, reqHeader string, headers []string) ([]byte, error, map[string]string) {

	var O interface{}
	out := make(map[string]string, len(headers))

	res, err := c.DoRequest(ctx, method, endpoint, body, reqHeader)

	if err != nil {
		return nil, err, out
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {

		type ErrorDetail struct {
			Target  string `json:"target"`
			Message string `json:"message"`
		}

		type Error struct {
			Code    int           `json:"code"`
			Message string        `json:"message"`
			Details []ErrorDetail `json:"details"`
		}

		var responseError struct {
			Error Error `json:"error"`
		}

		// _ = json.NewDecoder(res.Body).Decode(&O)

		// body, _ := io.ReadAll(res.Body)
		// bodyString := string(body)

		// fmt.Println(bodyString)

		if err = json.NewDecoder(res.Body).Decode(&responseError); err == nil {
			err = fmt.Errorf(fmt.Sprintf("%d : %s", responseError.Error.Code, responseError.Error.Message))
		} else {
			err = fmt.Errorf("responded with unknown error : %d", responseError.Error.Code)
		}

		return nil, err, out
	}

	for i := 0; i < len(headers); i++ {
		out[headers[i]] = res.Header.Get(headers[i])
	}

	if err = json.NewDecoder(res.Body).Decode(&O); err == nil || err == io.EOF {
		encodedRes, err := json.Marshal(O)

		if err != nil {
			return nil, err, out
		}

		return encodedRes, nil, out
	} else {
		return nil, err, out
	}
}
