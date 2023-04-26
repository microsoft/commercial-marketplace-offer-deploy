package config

import (
	"path/filepath"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/log"
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

func (s *AzureSettings) GetFullQualifiedNamespace() string {
	return s.ServiceBusNamespace + ".servicebus.windows.net"
}

func (s *AzureSettings) GetResourceGroupId() string {
	return "/subscriptions/" + s.SubscriptionId + "/resourceGroups/" + s.ResourceGroupName
}

// The database settings
type DatabaseSettings struct {
	Path        string `mapstructure:"DB_PATH"`
	UseInMemory bool   `mapstructure:"DB_USE_INMEMEORY"`
}

type LoggingSettings struct {
	DefaultLogLevel string `mapstructure:"LOG_LEVEL"`
	FilePath        string `mapstructure:"LOG_FILE_PATH"`
}

type HttpSettings struct {
	DomainName string `mapstructure:"PUBLIC_DOMAIN_NAME"`
	Port       string `mapstructure:"PUBLIC_PORT"`
}

func (s *AppConfig) GetPublicBaseUrl() string {
	return "https://" + s.Http.DomainName + "/"
}

type AppConfig struct {
	Azure       AzureSettings
	Database    DatabaseSettings
	Http        HttpSettings
	Logging     LoggingSettings
	Environment string `mapstructure:"GO_ENV"`
}

func (appSettings *AppConfig) GetDatabaseOptions() *data.DatabaseOptions {
	dsn := filepath.Join(appSettings.Database.Path, data.DatabaseFileName)
	options := &data.DatabaseOptions{Dsn: dsn, UseInMemory: appSettings.Database.UseInMemory}
	return options
}

func (appSettings *AppConfig) GetLogOptions() *log.LoggingOptions {
	return &log.LoggingOptions{
		DefaultLogLevel: appSettings.Logging.DefaultLogLevel,
		FilePath:        filepath.Join(appSettings.Logging.FilePath, "modmlog.txt"),
	}
}

func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}
