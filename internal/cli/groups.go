package cli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/groups"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
)

type GroupsCli struct {
	cliClient *Client
}

func NewGroupCli(cliClient *Client) GroupsCli {
	return GroupsCli{cliClient: cliClient}
}

func (g *GroupsCli) getUrl() string {
	return "scim/Groups/"
}

func (g *GroupsCli) Get(ctx context.Context) (groups.GroupsResponse, string, error) {

	var allGroups groups.GroupsResponse
	startId := "initial"
	const maxPages = 10000
	done := false

	for page := 0; page < maxPages; page++ {
		queryStrings := map[string]string{
			"startId": startId,
		}

		res, _, err := g.cliClient.Execute(ctx, "GET", g.getUrl(), queryStrings, nil, "", ScimRequestHeader, nil)
		if err != nil {
			return groups.GroupsResponse{}, "", err
		}

		resp, _, err := unMarshalResponse[groups.GroupsResponse](res, false)
		if err != nil {
			return groups.GroupsResponse{}, "", err
		}

		allGroups.Resources = append(allGroups.Resources, resp.Resources...)

		if resp.NextId == "" || resp.NextId == "end" {
			done = true
			break
		}
		startId = resp.NextId
	}

	if !done {
		return groups.GroupsResponse{}, "", fmt.Errorf("groups pagination exceeded %d pages", 10000)
	}

	return allGroups, "", nil
}

func (g *GroupsCli) GetByGroupId(ctx context.Context, groupId string) (groups.Group, string, error) {

	res, _, err := g.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", g.getUrl(), url.PathEscape(groupId)), nil, nil, "", ScimRequestHeader, nil)

	if err != nil {
		return groups.Group{}, "", err
	}

	return unMarshalResponse[groups.Group](res, false)
}

func (g *GroupsCli) Create(ctx context.Context, args *groups.Group) (groups.Group, string, error) {

	res, _, err := g.cliClient.Execute(ctx, "POST", g.getUrl(), nil, args, "", ScimRequestHeader, nil)

	if err != nil {
		return groups.Group{}, "", err
	}

	return unMarshalResponse[groups.Group](res, false)
}

func (g *GroupsCli) Update(ctx context.Context, args []generic.PatchRequest, groupId string) (groups.Group, string, error) {

	reqBody := groups.PatchRequestBody{
		Schemas:    []string{ScimUpdateSchemas},
		Operations: args,
	}

	_, _, err := g.cliClient.Execute(ctx, "PATCH", fmt.Sprintf("%s%s", g.getUrl(), url.PathEscape(groupId)), nil, reqBody, "", ScimRequestHeader, nil)

	if err != nil {
		return groups.Group{}, "", err
	}

	return g.GetByGroupId(ctx, groupId)
}

func (g *GroupsCli) Delete(ctx context.Context, groupId string) error {

	_, _, err := g.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", g.getUrl(), url.PathEscape(groupId)), nil, nil, "", ScimRequestHeader, nil)

	return err
}
