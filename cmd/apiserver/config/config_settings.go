package config

// The azure ad settings
type AzureAdSettings struct {
	ClientId       string `mapstructure:"AZURE_CLIENT_ID"`
	TenantId       string `mapstructure:"AZURE_TENANT_ID"`
	SubscriptionId string `mapstructure:"AZURE_SUBSCRIPTION_ID"`
}

// The database settings
type DatabaseSettings struct {
	Path        string `mapstructure:"DB_PATH"`
	UseInMemory bool   `mapstructure:"DB_USE_INMEMEORY"`
}
