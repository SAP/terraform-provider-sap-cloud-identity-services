package cli

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
)

type ApplicationSecretsCli struct {
	cliClient *Client
}

func NewApplicationSecretCli(cliClient *Client) ApplicationSecretsCli {
	return ApplicationSecretsCli{cliClient: cliClient}
}

func (a *ApplicationSecretsCli) getUrl(appId string) string {
	return fmt.Sprintf("Applications/v1/%s/apiSecrets", appId)
}

func (a *ApplicationSecretsCli) Create(ctx context.Context, appId string, args applications.ApplicationSecretRequest) (applications.ApplicationSecret, error) {
	res, _, err := a.cliClient.Execute(ctx, "POST", a.getUrl(appId), nil, args, "", RequestHeader, nil)
	if err != nil {
		return applications.ApplicationSecret{}, err
	}

	secret, _, err := unMarshalResponse[applications.ApplicationSecret](res, false)
	return secret, err
}

func (a *ApplicationSecretsCli) Get(ctx context.Context, appId string) (applications.ApplicationSecretsListResponse, error) {
	res, _, err := a.cliClient.Execute(ctx, "GET", a.getUrl(appId), nil, nil, "", RequestHeader, nil)
	if err != nil {
		return applications.ApplicationSecretsListResponse{}, err
	}

	list, _, err := unMarshalResponse[applications.ApplicationSecretsListResponse](res, false)
	return list, err
}

func (a *ApplicationSecretsCli) GetById(ctx context.Context, appId, secretId string) (applications.ApplicationSecret, error) {
	list, err := a.Get(ctx, appId)
	if err != nil {
		return applications.ApplicationSecret{}, err
	}

	for _, secret := range list.Secrets {
		if secret.Id == secretId {
			return secret, nil
		}
	}

	return applications.ApplicationSecret{}, fmt.Errorf("secret with id %s not found for application %s", secretId, appId)
}

func (a *ApplicationSecretsCli) Update(ctx context.Context, appId, secretId string, ops []generic.PatchRequest) (applications.ApplicationSecret, error) {
	reqBody := applications.ApplicationSecretPatchRequestBody{
		Operations: ops,
	}

	_, _, err := a.cliClient.Execute(ctx, "PATCH", fmt.Sprintf("%s/%s", a.getUrl(appId), secretId), nil, reqBody, "", RequestHeader, nil)
	if err != nil {
		return applications.ApplicationSecret{}, err
	}

	return a.GetById(ctx, appId, secretId)
}

func (a *ApplicationSecretsCli) Delete(ctx context.Context, appId, secretId string) error {
	_, _, err := a.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s?id=%s", a.getUrl(appId), secretId), nil, nil, "", RequestHeader, nil)
	return err
}
