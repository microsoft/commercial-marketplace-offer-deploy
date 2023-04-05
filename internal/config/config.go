package config

import (
	"path/filepath"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

// The azure settings
type AzureSettings struct {
	ClientId            string `mapstructure:"AZURE_CLIENT_ID"`
	TenantId            string `mapstructure:"AZURE_TENANT_ID"`
	SubscriptionId      string `mapstructure:"AZURE_SUBSCRIPTION_ID"`
	ResourceGroupName   string `mapstructure:"AZURE_RESOURCE_GROUP"`
	Location            string `mapstructure:"AZURE_LOCATION"`
	ServiceBusNamespace string `mapstructure:"AZURE_SERVICEBUS_NAMESPACE"`
}

// The database settings
type DatabaseSettings struct {
	Path        string `mapstructure:"DB_PATH"`
	UseInMemory bool   `mapstructure:"DB_USE_INMEMEORY"`
}

type AppConfig struct {
	Azure               AzureSettings
	Database            DatabaseSettings
	ServiceBusNamespace string `mapstructure:"AZURE_SERVICEBUS_NAMESPACE"`
}

func GetAppConfig() *AppConfig {
	appConfig := hosting.GetAppConfig[AppConfig]()
	return &appConfig
}

func (appSettings *AppConfig) GetDatabaseOptions() *data.DatabaseOptions {
	dsn := filepath.Join(appSettings.Database.Path, data.DatabaseFileName)
	options := &data.DatabaseOptions{Dsn: dsn, UseInMemory: appSettings.Database.UseInMemory}
	return options
}
