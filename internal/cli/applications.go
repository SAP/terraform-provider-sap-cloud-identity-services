package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"terraform-provider-ias/internal/cli/apiObjects/applications"
)

type ApplicationsCli struct {
	cliClient *Client
}

func NewApplicationCli(cliClient *Client) ApplicationsCli {
	return ApplicationsCli{cliClient: cliClient}
}

func (a *ApplicationsCli) getUrl() string {
	return "Applications/v1/"
}

func (a *ApplicationsCli) Get(ctx context.Context) (applications.ApplicationsResponse, error) {
	var app applications.ApplicationsResponse

	res, err, _ := a.cliClient.Execute(ctx, "GET", a.getUrl(), nil, ApplicationHeader, nil)

	if err != nil {
		return app, err
	}

	if err = json.Unmarshal(res, &app); err != nil {
		return app, err
	}
	return app, nil
}

func (a *ApplicationsCli) GetByAppId(ctx context.Context, appId string) (applications.Application, error) {
	var app applications.Application

	res, err, _ := a.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", a.getUrl(), appId), nil, ApplicationHeader, nil)

	if err != nil {
		return app, err
	}

	if err = json.Unmarshal(res, &app); err != nil {
		return app, err
	}
	return app, nil
}

func (a *ApplicationsCli) Create(ctx context.Context, args *applications.Application) (string, error) {

	var appId string

	_, err, res := a.cliClient.Execute(ctx, "POST", a.getUrl(), args, ApplicationHeader, []string{
		"location",
	})

	if err != nil {
		return appId, err
	}

	appId = res["location"]

	return appId, nil
}

func (a *ApplicationsCli) Update(ctx context.Context, args *applications.Application) error {
	_, err, _ := a.cliClient.Execute(ctx, "PUT", fmt.Sprintf("%s%s", a.getUrl(), args.Id), args, ApplicationHeader, nil)
	return err
}

func (a *ApplicationsCli) Delete(ctx context.Context, id string) error {
	_, err, _ := a.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", a.getUrl(), id), nil, ApplicationHeader, nil)
	return err
}
