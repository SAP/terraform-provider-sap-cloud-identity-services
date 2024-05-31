package cli

import (
	"context"
	"encoding/json"
	"terraform-provider-cloudidentityservices/internal/cli/apiObjects/users"
)

type UsersCli struct {
	cliClient *Client
}

func NewUserCli (cliClient *Client) UsersCli{
	return UsersCli{cliClient: cliClient}
}

func (u *UsersCli) getUrl() string{
	return "scim/Users"
}

func (u *UsersCli) Get(ctx context.Context) (users.UserGet, error){
	var user users.UserGet

	res, err := u.cliClient.Execute(ctx, "GET", u.getUrl(), nil)

	if err!=nil {
		return user, err
	}

	encodedRes, _ := json.Marshal(res)
	err = json.Unmarshal(encodedRes, &user)

	if err!=nil{
		return user, err
	}

	return user, nil
}
