package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
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

func (a *ApplicationsCli) Get(ctx context.Context, cursor string) (applications.ApplicationsResponse, string, error) {

	var allApps applications.ApplicationsResponse

	for {
		queryStrings := map[string]string{
			"cursor": cursor,
		}

		res, _, err := a.cliClient.Execute(ctx, "GET", a.getUrl(), queryStrings, nil, "", RequestHeader, nil)
		if err != nil {
			return applications.ApplicationsResponse{}, "", err
		}

		resp, _, err := unMarshalResponse[applications.ApplicationsResponse](res, false)
		if err != nil {
			return applications.ApplicationsResponse{}, "", err
		}

		allApps.Applications = append(allApps.Applications, resp.Applications...)

		if resp.NextCursor == "" {
			break
		}

		cursor = resp.NextCursor
	}

	allApps.NextCursor = cursor

	return allApps, "", nil
}

func (a *ApplicationsCli) GetByAppId(ctx context.Context, appId string) (applications.Application, string, error) {

	res, _, err := a.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", a.getUrl(), appId), nil, nil, "", RequestHeader, nil)

	if err != nil {
		return applications.Application{}, "", err
	}

	return unMarshalResponse[applications.Application](res, false)
}

func (a *ApplicationsCli) Create(ctx context.Context, args *applications.Application) (applications.Application, string, error) {

	// The API returns the unique ID of the created application in the header key "location"
	_, headers, err := a.cliClient.Execute(ctx, "POST", a.getUrl(), nil, args, "", RequestHeader, []string{
		"location",
	})

	if err != nil {
		return applications.Application{}, "", err
	}

	// The retrieved header is returned as a string in the form "/Applications/v1/ID"
	// Hence it is split to retrieve the unique ID which is passed to the GET call
	return a.GetByAppId(ctx, strings.Split(headers["location"], "/")[3])
}

func (a *ApplicationsCli) Update(ctx context.Context, args *applications.Application) (applications.Application, string, error) {

	_, _, err := a.cliClient.Execute(ctx, "PUT", fmt.Sprintf("%s%s", a.getUrl(), args.Id), nil, args, "", RequestHeader, nil)

	if err != nil {
		return applications.Application{}, "", err
	}

	return a.GetByAppId(ctx, args.Id)
}

func (a *ApplicationsCli) Delete(ctx context.Context, appId string) error {
	_, _, err := a.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", a.getUrl(), appId), nil, nil, "", RequestHeader, nil)
	return err
}
