package cli

import (
	"bytes"
	"context"
	"io"

	// "context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"net/url"
	"os"
)

func NewClient(h *http.Client, u *url.URL) *Client{
	return &Client{
		HttpClient: h,
		ServerURL: u,
	}
}

type Client struct{
	HttpClient *http.Client
	ServerURL  *url.URL
}


func (c *Client) DoRequest(ctx context.Context, method string, endpoint string, body any) (*http.Response, error) {
	parsedUrl, err := url.Parse(endpoint)

	if err != nil {
		return nil, err
	}

	var username, password string

	username = os.Getenv("ias_username")
	password = os.Getenv("ias_password")

	base64Encoded := base64.StdEncoding.EncodeToString([]byte (username+":"+password))

	var encodedBody bytes.Buffer
	if body != nil {
		encoder := json.NewEncoder(&encodedBody)
		err := encoder.Encode(body)
		if err != nil{
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
	req.Header.Set("DataServiceVersion","2.0")
	req.Header.Set("Content-Type", "application/scim+json")

	res, err := c.HttpClient.Do(req)

	if err != nil{
		return nil, err 
	}

	return res, nil
}

func (c *Client) Execute (ctx context.Context, method string, endpoint string, body any) ([]byte, error) {

	var O interface{}

	res, err := c.DoRequest(ctx, method, endpoint, body)

	if err!=nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&O); err == nil || err == io.EOF {
		encodedRes, err := json.Marshal(O)

		if err!=nil{
			return nil,err
		}

		return encodedRes, nil
	} else {
		return nil, err
	}
}
