package cli

import (
	"context"
	"encoding/json"
	"terraform-provider-ias/internal/cli/apiObjects/users"
)

type UsersCli struct {
	cliClient *Client
}

func NewUserCli(cliClient *Client) UsersCli {
	return UsersCli{cliClient: cliClient}
}

func (u *UsersCli) getUrl() string {
	return "scim/Users"
}

func (u *UsersCli) Get(ctx context.Context) (users.UsersResponse, error) {
	var users users.UsersResponse

	res, err, _ := u.cliClient.Execute(ctx, "GET", u.getUrl(), nil, DirectoryHeader, nil)

	if err != nil {
		return users, err
	}

	if err = json.Unmarshal(res, &users); err != nil{
		return users, err
	}
	
	return users, nil
}
