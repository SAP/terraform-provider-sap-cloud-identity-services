package cli

import (
	"context"
	"encoding/json"
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

func (s *SchemasCli) Get(ctx context.Context) (schemas.SchemasResponse, error) {
	var schemas schemas.SchemasResponse

	res, err, _ := s.cliClient.Execute(ctx, "GET", s.getUrl(), nil, DirectoryHeader, nil)

	if err != nil {
		return schemas, err
	}

	if err = json.Unmarshal(res, &schemas); err != nil {
		return schemas, err
	}

	return schemas, nil
}

func (s *SchemasCli) GetBySchemaId(ctx context.Context, schemaId string) (schemas.Schema, error) {
	var schema schemas.Schema

	res, err, _ := s.cliClient.Execute(ctx, "GET", fmt.Sprintf("%s%s", s.getUrl(), schemaId), nil, DirectoryHeader, nil)

	if err != nil {
		return schema, err
	}

	if err = json.Unmarshal(res, &schema); err != nil {
		return schema, err
	}

	return schema, nil
}

func (s *SchemasCli) Create(ctx context.Context, args *schemas.Schema) (schemas.Schema, error) {
	var schema schemas.Schema

	res, err, _ := s.cliClient.Execute(ctx, "POST", s.getUrl(), args, DirectoryHeader, nil)
	if err != nil {
		return schema, err
	}

	if err = json.Unmarshal(res, &schema); err != nil {
		return schema, err
	}

	return schema, nil
}

func (s *SchemasCli) Delete(ctx context.Context, schemaId string) error {

	_, err, _ := s.cliClient.Execute(ctx, "DELETE", fmt.Sprintf("%s%s", s.getUrl(), schemaId), nil, DirectoryHeader, nil)

	return err
}
