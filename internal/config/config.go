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

type HttpSettings struct {
	FQDN      string `mapstructure:"PUBLIC_FQDN"`
	HttpPort  string `mapstructure:"PUBLIC_HTTP_PORT"`
	HttpsPort string   `mapstructure:"PUBLIC_HTTPS_PORT"`
}

func (s *HttpSettings) GetBaseUrl(secure bool) string {
	protocol := "http"
	if secure {
		protocol = "https"
	}
	return protocol + "://" + s.FQDN + "/"
}

type AppConfig struct {
	Azure    AzureSettings
	Database DatabaseSettings
	Http     HttpSettings
}

func GetAppConfig() *AppConfig {
	appConfig := hosting.GetAppConfig[*AppConfig]()
	return appConfig
}

func (appSettings *AppConfig) GetDatabaseOptions() *data.DatabaseOptions {
	dsn := filepath.Join(appSettings.Database.Path, data.DatabaseFileName)
	options := &data.DatabaseOptions{Dsn: dsn, UseInMemory: appSettings.Database.UseInMemory}
	return options
}
