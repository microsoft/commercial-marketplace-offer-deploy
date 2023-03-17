package config

// The azure ad settings
type AzureAdSettings struct {
	ClientId string `mapstructure:"ClientId"`
	TenantId string `mapstructure:"TenantId"`
}
