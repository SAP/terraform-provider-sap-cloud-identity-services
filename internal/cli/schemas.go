package cli

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/schemas"
)

type SchemasCli struct {
	cliClient *Client
}

func NewSchemaCli(cliClient *Client) SchemasCli {
	return SchemasCli{cliClient: cliClient}
}

func (s *SchemasCli) getUrl() string {
	return "scim/Schemas/"
}

func (s *SchemasCli) Get(ctx context.Context) (schemas.SchemasResponse, string, error) {

	res, _, err := s.cliClient.Execute(ctx, "GET", s.getUrl(), nil, "", ScimRequestHeader, nil)

	if err != nil {
		return schemas.SchemasResponse{}, "", err
	}

	return unMarshalResponse[schemas.SchemasResponse](res, false)
}

func (s *SchemasCli) GetBySchemaId(ctx context.Context, schemaId string) (schemas.Schema, string, error) {

	res, _, err := s.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", s.getUrl(), schemaId), nil, "", ScimRequestHeader, nil)

	if err != nil {
		return schemas.Schema{}, "", err
	}

	return unMarshalResponse[schemas.Schema](res, false)
}

func (s *SchemasCli) Create(ctx context.Context, args *schemas.Schema) (schemas.Schema, string, error) {

	res, _, err := s.cliClient.Execute(ctx, "POST", s.getUrl(), args, "", ScimRequestHeader, nil)
	if err != nil {
		return schemas.Schema{}, "", err
	}

	return unMarshalResponse[schemas.Schema](res, false)
}

func (s *SchemasCli) Delete(ctx context.Context, schemaId string) error {

	_, _, err := s.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", s.getUrl(), schemaId), nil, "", ScimRequestHeader, nil)

	return err
}
