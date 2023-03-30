package config

import (
	"path/filepath"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
)

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

type AppSettings struct {
	Azure    AzureSettings
	Database DatabaseSettings
}

func (appSettings *AppSettings) GetDatabaseOptions() *data.DatabaseOptions {
	dsn := filepath.Join(appSettings.Database.Path, data.DatabaseFileName)
	options := &data.DatabaseOptions{Dsn: dsn, UseInMemory: appSettings.Database.UseInMemory}
	return options
}
