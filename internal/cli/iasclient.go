package cli

func NewIasClient(cliClient *Client) *IasClient {
	return &IasClient{
		Client:      cliClient,
		Application: NewApplicationCli(cliClient),
		User:        NewUserCli(cliClient),
		Schema:      NewSchemaCli(cliClient),
		Group:       NewGroupCli(cliClient),
	}
}

type IasClient struct {
	*Client
	Application ApplicationsCli
	User        UsersCli
	Schema      SchemasCli
	Group       GroupsCli
}
