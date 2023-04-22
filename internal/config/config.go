package config

import (
	"path/filepath"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
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
	DefaultLogLevel    string `mapstructure:"LOG_LEVEL"`
	InstrumentationKey string `mapstructure:"LOG_KEY"`
}

type HttpSettings struct {
	DomainName string `mapstructure:"PUBLIC_DOMAIN_NAME"`
	HttpPort   string `mapstructure:"PUBLIC_HTTP_PORT"`
	HttpsPort  string `mapstructure:"PUBLIC_HTTPS_PORT"`
	IsSecure   bool   `mapstructure:"HTTPS"`
}

func (s *AppConfig) GetPublicBaseUrl() string {
	protocol := "http"
	if !s.IsDevelopment() {
		protocol = "https"
	}
	return protocol + "://" + s.Http.DomainName + "/"
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

func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}
