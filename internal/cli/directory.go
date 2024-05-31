package cli

func NewDirectoryCli(cliClient *Client) DirectoryCli{
	return DirectoryCli{
		User: NewUserCli(cliClient),
	}
}

type DirectoryCli struct{
	User UsersCli
}