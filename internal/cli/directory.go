package cli

func NewDirectoryCli(cliClient *Client) DirectoryCli{
	return DirectoryCli{
		User: NewUserCli(cliClient),
		Schema: NewSchemaCli(cliClient),
		Group: NewGroupCli(cliClient),
	}
}

type DirectoryCli struct{
	User 	UsersCli
	Schema 	SchemasCli
	Group 	GroupsCli	
}