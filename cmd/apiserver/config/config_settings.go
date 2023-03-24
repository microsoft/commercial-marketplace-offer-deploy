package config

// The azure settings
type AzureSettings struct {
	ClientId       string `mapstructure:"AZURE_CLIENT_ID"`
	TenantId       string `mapstructure:"AZURE_TENANT_ID"`
	SubscriptionId string `mapstructure:"AZURE_SUBSCRIPTION_ID"`
	ResourceGroup  string `mapstructure:"AZURE_RESOURCE_GROUP"`
	Location       string `mapstructure:"AZURE_LOCATION"`
}

// The database settings
type DatabaseSettings struct {
	Path        string `mapstructure:"DB_PATH"`
	UseInMemory bool   `mapstructure:"DB_USE_INMEMEORY"`
}
