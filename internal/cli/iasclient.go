package cli

func NewSciClient(cliClient *Client) *SciClient {
	return &SciClient{
		Client:            cliClient,
		Application:       NewApplicationCli(cliClient),
		ApplicationSecret: NewApplicationSecretCli(cliClient),
		User:              NewUserCli(cliClient),
		Schema:            NewSchemaCli(cliClient),
		Group:             NewGroupCli(cliClient),
		CorporateIdP:      NewCorporateIdPCli(cliClient),
	}
}

type SciClient struct {
	*Client
	Application       ApplicationsCli
	ApplicationSecret ApplicationSecretsCli
	User              UsersCli
	Schema            SchemasCli
	Group             GroupsCli
	CorporateIdP      CorporateIdPsCli
}
