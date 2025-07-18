package cli

import (
	"context"

	"fmt"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/users"
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

func (u *UsersCli) Get(ctx context.Context) (users.UsersResponse, map[int]string, error) {

	res, _, err := u.cliClient.Execute(ctx, "GET", u.getUrl(), nil, "", ScimRequestHeader, nil)
	if err != nil {
		return users.UsersResponse{}, map[int]string{}, err
	}

	usersList := users.UsersResponse{}
	resMap := res.(map[string]interface{})["Resources"].([]interface{})
	customSchemas := map[int]string{}

	for i, r := range resMap {

		// each user is unmarshalled individually and the respective custom schemas are retrieved and added to the map
		var user users.User
		user, customSchemas[i], err = unMarshalResponse[users.User](r, true)

		if err != nil {
			return users.UsersResponse{}, map[int]string{}, err
		}
		usersList.Resources = append(usersList.Resources, user)

	}

	return usersList, customSchemas, err
}

func (u *UsersCli) GetByUserId(ctx context.Context, userId string) (users.User, string, error) {

	res, _, err := u.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", u.getUrl(), userId), nil, "", ScimRequestHeader, nil)

	if err != nil {
		return users.User{}, "", err
	}

	return unMarshalResponse[users.User](res, true)
}

func (u *UsersCli) Create(ctx context.Context, customSchemas string, args *users.User) (users.User, string, error) {

	res, _, err := u.cliClient.Execute(ctx, "POST", u.getUrl(), args, customSchemas, ScimRequestHeader, nil)
	if err != nil {
		return users.User{}, "", err
	}

	if len(customSchemas) > 0 {
		if result, err := validateCustomSchemasResponse(res, customSchemas); !result {
			return users.User{}, "", err
		}
	}

	return unMarshalResponse[users.User](res, false)
}

func (u *UsersCli) Update(ctx context.Context, customSchemas string, args *users.User) (users.User, string, error) {

	res, _, err := u.cliClient.Execute(ctx, "PUT", fmt.Sprintf("%s%s", u.getUrl(), args.Id), args, customSchemas, ScimRequestHeader, nil)

	if err != nil {
		return users.User{}, "", err
	}

	if len(customSchemas) > 0 {
		if result, err := validateCustomSchemasResponse(res, customSchemas); !result {
			return users.User{}, "", err
		}
	}

	return unMarshalResponse[users.User](res, false)
}

func (u *UsersCli) Delete(ctx context.Context, userId string) error {

	_, _, err := u.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", u.getUrl(), userId), nil, "", ScimRequestHeader, nil)

	return err
}
