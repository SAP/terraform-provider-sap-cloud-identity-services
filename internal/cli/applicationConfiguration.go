package cli

func NewApplicationConfigurationCli(cliclient *Client) ApplicationConfigurationCli{
	return ApplicationConfigurationCli{
		Application: NewApplicationCli(cliclient),
	}
}

type ApplicationConfigurationCli struct{
	Application ApplicationsCli
}