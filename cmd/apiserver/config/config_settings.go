package config

// The azure ad settings
type AzureAdSettings struct {
	ClientId string `mapstructure:"ClientId"`
	TenantId string `mapstructure:"TenantId"`
}

// The database settings
type DatabaseSettings struct {
	Path        string `mapstructure:"Path"`
	UseInMemory bool   `mapstructure:"UseInMemory"`
}
