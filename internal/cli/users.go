package cli

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
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

	usersList := users.UsersResponse{}
	customSchemas := map[int]string{}
	startId := "initial"

	for {
		queryStrings := map[string]string{
			"startId": startId,
		}

		res, _, err := u.cliClient.Execute(ctx, "GET", u.getUrl(), queryStrings, nil, "", ScimRequestHeader, nil)
		if err != nil {
			return users.UsersResponse{}, map[int]string{}, err
		}

		resBody := res.(map[string]any)
		resMap, _ := resBody["Resources"].([]any)

		for _, r := range resMap {
			// each user is unmarshalled individually and the respective custom schemas are retrieved and added to the map
			var user users.User
			var schema string
			user, schema, err = unMarshalResponse[users.User](r, true)
			if err != nil {
				return users.UsersResponse{}, map[int]string{}, err
			}
			customSchemas[len(usersList.Resources)] = schema
			usersList.Resources = append(usersList.Resources, user)
		}

		nextId, _ := resBody["nextId"].(string)
		if nextId == "" || nextId == "end" {
			break
		}
		startId = nextId
	}

	return usersList, customSchemas, nil
}

func (u *UsersCli) GetByUserId(ctx context.Context, userId string, validateCustomSchemas bool, customSchemas string) (users.User, string, error) {

	res, _, err := u.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", u.getUrl(), userId), nil, nil, "", ScimRequestHeader, nil)

	if err != nil {
		return users.User{}, "", err
	}

	if len(customSchemas) > 0 && validateCustomSchemas {
		if result, err := validateCustomSchemasResponse(res, customSchemas); !result {
			return users.User{}, "", err
		}
	}

	return unMarshalResponse[users.User](res, true)
}

func (u *UsersCli) Create(ctx context.Context, customSchemas string, args *users.User) (users.User, string, error) {

	res, _, err := u.cliClient.Execute(ctx, "POST", u.getUrl(), nil, args, customSchemas, ScimRequestHeader, nil)
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

func (u *UsersCli) Update(ctx context.Context, id string, args []generic.PatchRequest, customSchemas string) (users.User, string, error) {

	reqBody := users.PatchRequestBody{
		Schemas:    []string{ScimUpdateSchemas},
		Operations: args,
	}

	_, _, err := u.cliClient.Execute(ctx, "PATCH", fmt.Sprintf("%s%s", u.getUrl(), id), nil, reqBody, "", ScimRequestHeader, nil)

	if err != nil {
		return users.User{}, "", err
	}

	res, cS, err := u.GetByUserId(ctx, id, true, customSchemas)
	if err != nil {
		return users.User{}, "", err
	}

	return res, cS, nil
}

func (u *UsersCli) Delete(ctx context.Context, userId string) error {

	_, _, err := u.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", u.getUrl(), userId), nil, nil, "", ScimRequestHeader, nil)

	return err
}
