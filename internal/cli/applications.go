package cli

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/applications"
	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
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

func (a *ApplicationsCli) Get(ctx context.Context) (applications.ApplicationsResponse, string, error) {

	res, _, err := a.cliClient.Execute(ctx, "GET", a.getUrl(), nil, "", RequestHeader, nil)

	if err != nil {
		return applications.ApplicationsResponse{}, "", err
	}

	return unMarshalResponse[applications.ApplicationsResponse](res, false)
}

func (a *ApplicationsCli) GetByAppId(ctx context.Context, appId string) (applications.Application, string, error) {

	res, _, err := a.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", a.getUrl(), appId), nil, "", RequestHeader, nil)

	if err != nil {
		return applications.Application{}, "", err
	}

	return unMarshalResponse[applications.Application](res, false)
}

func (a *ApplicationsCli) Create(ctx context.Context, args *applications.Application) (applications.Application, string, error) {

	// The API returns the unique ID of the created application in the header key "location"
	_, headers, err := a.cliClient.Execute(ctx, "POST", a.getUrl(), args, "", RequestHeader, []string{
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

	req := getPatchRequestBody(*args, "")
	reqBody := generic.PatchRequestBody{
		Operations: req,
	}

	// _, _, err := a.cliClient.Execute(ctx, "PUT", fmt.Sprintf("%s%s", a.getUrl(), args.Id), args, "", RequestHeader, nil)
	_, _, err := a.cliClient.Execute(ctx, "PATCH", fmt.Sprintf("%s%s", a.getUrl(), args.Id), reqBody, "", RequestHeader, nil)

	if err != nil {
		return applications.Application{}, "", err
	}

	return a.GetByAppId(ctx, args.Id)
}

func (a *ApplicationsCli) Delete(ctx context.Context, appId string) error {
	_, _, err := a.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", a.getUrl(), appId), nil, "", RequestHeader, nil)
	return err
}

func getPatchRequestBody(args any, tag string) ([]generic.PatchRequest) {
	var patchRequests, reqs []generic.PatchRequest

	argsVal := reflect.ValueOf(args)
	argsType := reflect.TypeOf(args)

	if argsType.Kind().String() == "ptr" {
		argsVal = argsVal.Elem()
		argsType = argsType.Elem()
	}

	for i:=0; i<argsType.NumField(); i++ {

		field := argsType.Field(i)
	
		fieldName := field.Name
		fieldType := field.Type
		
		fieldTag := fmt.Sprintf("/%s", strings.Split(field.Tag.Get("json"), ",")[0])

		if tag != "" {
			fieldTag = fmt.Sprintf("%s%s", tag, fieldTag)
		}

		fieldValue := argsVal.FieldByName(fieldName)

		var val any
		var valSet bool
		t := fieldType.Kind().String()

		switch t {
			case "struct":
			case "ptr":
				if fieldValue.IsNil() { 
					val = nil
				} else {
					reqs = getPatchRequestBody(fieldValue.Interface(), fieldTag)
					valSet = false
				}
			case "string":
				val = fieldValue.String()
				valSet = true

				if val == "" {
					remove := validate(fieldTag)

					if remove {
						continue
					}
				}
			case "bool":
				val = fieldValue.Bool()
				valSet = true
			case "int":
				val = fieldValue.Int()
				valSet = true
			case "slice":
				if fieldValue.IsNil() {
					val = []any{}
					valSet = true
				} else {
					

					for i := range fieldValue.Len() {
						obj := fieldValue.Index(i).Interface()
						var req []generic.PatchRequest

						objType := reflect.TypeOf(obj)
						if objType.Kind().String() != "ptr" && objType.Kind().String() != "struct" {
							req = []generic.PatchRequest{ 
								{
									Op:    "replace",
									Path:  fmt.Sprintf("%s/%d", fieldTag, i),
									Value: obj,
								},
							}
						} else {
							req = getPatchRequestBody(obj, fmt.Sprintf("%s/%d", fieldTag, i))
						}
						
						reqs = append(reqs, req...)
					}
					valSet = false
				}
		}

		if !valSet {
			patchRequests = append(patchRequests, reqs...)
		} else {
			patchRequest := generic.PatchRequest{
				Op:   "replace",
				Path: fieldTag,
				Value: val,
			}

			patchRequests = append(patchRequests, patchRequest)
		}

	}

	return patchRequests
}

func validate(fieldTag string) bool {

	remove := false

	switch fieldTag {
	case "/urn:sap:identity:application:schemas:extension:sci:1.0:Authentication/saml2Configuration/samlMetadataUrl":
		fallthrough
	case "/urn:sap:identity:application:schemas:extension:sci:1.0:Authentication/subjectNameIdentifierFunction":
		remove = true
	}

	return remove
}