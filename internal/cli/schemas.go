package cli

import (
	"context"
	"fmt"

	"terraform-provider-ias/internal/cli/apiObjects/schemas"
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

	res, err, _ := s.cliClient.Execute(ctx, "GET", s.getUrl(), nil, "", DirectoryHeader, nil)

	if err != nil {
		return schemas.SchemasResponse{}, "", err
	}

	return unMarshalResponse[schemas.SchemasResponse](res, false)
}

func (s *SchemasCli) GetBySchemaId(ctx context.Context, schemaId string) (schemas.Schema, string, error) {

	res, err, _ := s.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", s.getUrl(), schemaId), nil, "", DirectoryHeader, nil)

	if err != nil {
		return schemas.Schema{}, "", err
	}

	return unMarshalResponse[schemas.Schema](res, false)
}

func (s *SchemasCli) Create(ctx context.Context, args *schemas.Schema) (schemas.Schema, string, error) {

	res, err, _ := s.cliClient.Execute(ctx, "POST", s.getUrl(), args, "", DirectoryHeader, nil)
	if err != nil {
		return schemas.Schema{}, "", err
	}

	return unMarshalResponse[schemas.Schema](res, false)
}

func (s *SchemasCli) Delete(ctx context.Context, schemaId string) error {

	_, err, _ := s.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", s.getUrl(), schemaId), nil, "", DirectoryHeader, nil)

	return err
}
