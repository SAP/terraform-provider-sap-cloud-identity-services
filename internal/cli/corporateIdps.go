package cli

import (
	"context"
	"fmt"
	"strings"

	corporateidps "github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/corporateIdps"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
)

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

	return c.GetByIdPId(ctx, idpId)
}

func (c *CorporateIdPsCli) Delete(ctx context.Context, idpId string) error {
	_, _, err := c.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", c.getUrl(), idpId), nil, nil, "", RequestHeader, nil)
	return err
}
