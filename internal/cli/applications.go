package cli

import (
	"context"
	"encoding/json"
	// "strconv"
	"terraform-provider-cloudidentityservices/internal/cli/apiObjects/applications"
)

type ApplicationsCli struct {
	cliClient *Client
}

func NewApplicationCli (cliClient *Client) ApplicationsCli{
	return ApplicationsCli{cliClient: cliClient}
}

func (a *ApplicationsCli) getUrl() string{
	return "Applications/v1//"
}

func (a *ApplicationsCli) Get(ctx context.Context) (applications.ApplicationsResponse, error){
	var app applications.ApplicationsResponse

	res, err := a.cliClient.Execute(ctx, "GET", a.getUrl(), nil)

	if err!=nil {
		return app, err
	}

	if err = json.Unmarshal(res, &app); err != nil {
		return app, err
	}
	return app, nil
}

func (a *ApplicationsCli) GetByAppId(ctx context.Context, appId string) (applications.ApplicationResponse, error){
	var app applications.ApplicationResponse

	// id, _ := strconv.Atoi(appId)

	res, err := a.cliClient.Execute(ctx, "GET", a.getUrl()+appId[ 1: len(appId)-1], nil)

	if err!=nil {
		return app, err
	}

	if err = json.Unmarshal(res, &app); err != nil {
		return app, err
	}
	return app, nil
}
