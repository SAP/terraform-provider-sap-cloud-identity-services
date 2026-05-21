package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
)

var corporateIdPPollInterval = 5 * time.Second

type CorporateIdPsCli struct {
	cliClient *Client
}

func NewCorporateIdPCli(cliClient *Client) CorporateIdPsCli {
	return CorporateIdPsCli{cliClient: cliClient}
}

func (c *CorporateIdPsCli) getUrl() string {
	return "IdentityProviders/v1/"
}

func (c *CorporateIdPsCli) Get(ctx context.Context) (corporateidps.IdentityProvidersResponse, string, error) {

	res, _, err := c.cliClient.Execute(ctx, "GET", c.getUrl(), nil, nil, "", RequestHeader, nil)

	if err != nil {
		return corporateidps.IdentityProvidersResponse{}, "", err
	}

	return unMarshalResponse[corporateidps.IdentityProvidersResponse](res, false)
}

func (c *CorporateIdPsCli) GetByIdPId(ctx context.Context, idpId string) (corporateidps.IdentityProvider, string, error) {

	res, _, err := c.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", c.getUrl(), idpId), nil, nil, "", RequestHeader, nil)

	if err != nil {
		return corporateidps.IdentityProvider{}, "", err
	}

	return unMarshalResponse[corporateidps.IdentityProvider](res, false)
}

func (c *CorporateIdPsCli) Create(ctx context.Context, args *corporateidps.IdentityProvider) (corporateidps.IdentityProvider, string, error) {

	_, headers, err := c.cliClient.Execute(ctx, "POST", c.getUrl(), nil, args, "", RequestHeader, []string{
		"location",
	})

	if err != nil {
		return corporateidps.IdentityProvider{}, "", err
	}

	return c.GetByIdPId(ctx, strings.Split(headers["location"], "/")[3])
}

func (c *CorporateIdPsCli) Update(ctx context.Context, args []generic.PatchRequest, idpId string) (corporateidps.IdentityProvider, string, error) {

	reqBody := corporateidps.PatchRequestBody{
		Operations: args,
	}

	_, _, err := c.cliClient.Execute(ctx, "PATCH", fmt.Sprintf("%s%s", c.getUrl(), idpId), nil, reqBody, "", RequestHeader, nil)

	if err != nil {
		return corporateidps.IdentityProvider{}, "", err
	}

	return c.pollUntilUpdated(ctx, idpId, args)
}

func (c *CorporateIdPsCli) Delete(ctx context.Context, idpId string) error {
	_, _, err := c.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", c.getUrl(), idpId), nil, nil, "", RequestHeader, nil)
	return err
}

func (c *CorporateIdPsCli) pollUntilUpdated(ctx context.Context, idpId string, ops []generic.PatchRequest) (corporateidps.IdentityProvider, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return corporateidps.IdentityProvider{}, "", ctx.Err()
		case <-time.After(corporateIdPPollInterval):
		}

		res, _, err := c.GetByIdPId(ctx, idpId)
		if err != nil {
			return corporateidps.IdentityProvider{}, "", err
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			return corporateidps.IdentityProvider{}, "", fmt.Errorf("failed to marshal response: %w", err)
		}
		var resMap map[string]any
		if err := json.Unmarshal(resBytes, &resMap); err != nil {
			return corporateidps.IdentityProvider{}, "", fmt.Errorf("failed to unmarshal response: %w", err)
		}

		if patchOpsReflected(ops, resMap) {
			return res, "", nil
		}
	}
}
