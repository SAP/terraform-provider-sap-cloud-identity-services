package cli

import (
	"context"
	"fmt"
	"terraform-provider-ias/internal/cli/apiObjects/groups"
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

	res, err, _ := g.cliClient.Execute(ctx, "GET", g.getUrl(), nil, "", DirectoryHeader, nil)

	if err != nil {
		return groups.GroupsResponse{}, "", err
	}

	return unMarshalResponse[groups.GroupsResponse](res, false)
}

func (g *GroupsCli) GetByGroupId(ctx context.Context, groupId string) (groups.Group, string, error) {

	res, err, _ := g.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", g.getUrl(), groupId), nil, "", DirectoryHeader, nil)

	if err != nil {
		return groups.Group{}, "", err
	}

	return unMarshalResponse[groups.Group](res, false)
}

func (g *GroupsCli) Create(ctx context.Context, args *groups.Group) (groups.Group, string, error) {

	res, err, _ := g.cliClient.Execute(ctx, "POST", g.getUrl(), args, "", DirectoryHeader, nil)

	if err != nil {
		return groups.Group{}, "", err
	}

	return unMarshalResponse[groups.Group](res, false)
}

func (g *GroupsCli) Update(ctx context.Context, args *groups.Group) (groups.Group, string, error) {

	res, err, _ := g.cliClient.Execute(ctx, "PUT", fmt.Sprintf("%s%s", g.getUrl(), args.Id), args, "", DirectoryHeader, nil)

	if err != nil {
		return groups.Group{}, "", err
	}

	return unMarshalResponse[groups.Group](res, false)
}

func (g *GroupsCli) Delete(ctx context.Context, groupId string) error {

	_, err, _ := g.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", g.getUrl(), groupId), nil, "", DirectoryHeader, nil)

	return err
}
