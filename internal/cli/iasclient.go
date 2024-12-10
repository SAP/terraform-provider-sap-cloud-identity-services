package cli

func NewIasClient(cliClient *Client) *IasClient {
	return &IasClient{
		Client: cliClient,
		Directory: NewDirectoryCli(cliClient),
		ApplicationConfiguration: NewApplicationConfigurationCli(cliClient),
	}
}

type IasClient struct{
	*Client
	Directory DirectoryCli
	ApplicationConfiguration ApplicationConfigurationCli
}