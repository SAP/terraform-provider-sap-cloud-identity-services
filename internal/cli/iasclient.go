package cli

func NewSciClient(cliClient *Client) *SciClient {
	return &SciClient{
		Client:      cliClient,
		Application: NewApplicationCli(cliClient),
		User:        NewUserCli(cliClient),
		Schema:      NewSchemaCli(cliClient),
		Group:       NewGroupCli(cliClient),
	}
}

type SciClient struct {
	*Client
	Application ApplicationsCli
	User        UsersCli
	Schema      SchemasCli
	Group       GroupsCli
}
