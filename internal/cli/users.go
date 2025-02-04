package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"terraform-provider-ias/internal/cli/apiObjects/users"
)

type UsersCli struct {
	cliClient *Client
}

func NewUserCli(cliClient *Client) UsersCli {
	return UsersCli{cliClient: cliClient}
}

func (u *UsersCli) getUrl() string {
	return "scim/Users/"
}

func (u *UsersCli) Get(ctx context.Context) (users.UsersResponse, error) {
	var users users.UsersResponse

	res, err, _ := u.cliClient.Execute(ctx, "GET", u.getUrl(), nil, DirectoryHeader, nil)

	if err != nil {
		return users, err
	}

	if err = json.Unmarshal(res, &users); err != nil {
		return users, err
	}

	return users, nil
}

func (u *UsersCli) GetByUserId(ctx context.Context, userId string) (users.User, error) {
	var user users.User

	res, err, _ := u.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", u.getUrl(), userId), nil, DirectoryHeader, nil)

	if err != nil {
		return user, err
	}

	if err = json.Unmarshal(res, &user); err != nil {
		return user, err
	}

	return user, nil
}

func (u *UsersCli) Create(ctx context.Context, args *users.User) (users.User, error) {
	var user users.User

	res, err, _ := u.cliClient.Execute(ctx, "POST", u.getUrl(), args, DirectoryHeader, nil)
	if err != nil {
		return user, err
	}

	if err = json.Unmarshal(res, &user); err != nil {
		return user, err
	}

	return user, nil
}

func (u *UsersCli) Update(ctx context.Context, args *users.User) (users.User, error) {
	var user users.User

	res, err, _ := u.cliClient.Execute(ctx, "PUT", fmt.Sprintf("%s%s", u.getUrl(), args.Id), args, DirectoryHeader, nil)
	if err != nil {
		return user, err
	}

	if err = json.Unmarshal(res, &user); err != nil {
		return user, err
	}

	return user, nil
}

func (u *UsersCli) Delete(ctx context.Context, userId string) error {

	_, err, _ := u.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", u.getUrl(), userId), nil, DirectoryHeader, nil)

	return err
}
