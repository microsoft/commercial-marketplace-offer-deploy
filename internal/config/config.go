package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/log"
)

// The azure settings
type AzureSettings struct {
	ClientId            string `mapstructure:"MODM_AZURE_CLIENT_ID"`
	TenantId            string `mapstructure:"MODM_AZURE_TENANT_ID"`
	SubscriptionId      string `mapstructure:"MODM_AZURE_SUBSCRIPTION_ID"`
	ResourceGroupName   string `mapstructure:"MODM_AZURE_RESOURCE_GROUP"`
	Location            string `mapstructure:"MODM_AZURE_LOCATION"`
	ServiceBusNamespace string `mapstructure:"MODM_AZURE_SERVICEBUS_NAMESPACE"`
}

func (s *AzureSettings) GetFullQualifiedNamespace() string {
	return s.ServiceBusNamespace + ".servicebus.windows.net"
}

func (s *AzureSettings) GetResourceGroupId() string {
	return "/subscriptions/" + s.SubscriptionId + "/resourceGroups/" + s.ResourceGroupName
}

// The database settings
type DatabaseSettings struct {
	Path        string `mapstructure:"MODM_DB_PATH"`
	UseInMemory bool   `mapstructure:"MODM_DB_USE_INMEMEORY"`
}

type LoggingSettings struct {
	DefaultLogLevel string `mapstructure:"MODM_LOG_LEVEL"`
	FilePath        string `mapstructure:"LOG_FILE_PATH"`
}

type HttpSettings struct {
	DomainName string `mapstructure:"MODM_PUBLIC_DOMAIN_NAME"`
	Port       string `mapstructure:"MODM_PUBLIC_PORT"`
}

func (s *AppConfig) GetPublicBaseUrl() string {
	return "https://" + s.Http.DomainName + "/"
}

type AppConfig struct {
	Azure             AzureSettings
	Database          DatabaseSettings
	Http              HttpSettings
	Logging           LoggingSettings
	Environment       string `mapstructure:"GO_ENV"`
	ReadinessFilePath string `mapstructure:"MODM_READINESS_FILE_PATH"`
}

func (c *AppConfig) GetDatabaseOptions() *data.DatabaseOptions {
	dsn := filepath.Join(c.Database.Path, data.DatabaseFileName)
	options := &data.DatabaseOptions{Dsn: dsn, UseInMemory: c.Database.UseInMemory}
	return options
}

func (c *AppConfig) GetReadinessFilePath() string {
	path := "/tmp/ready"
	if len(c.ReadinessFilePath) > 0 {
		path = c.ReadinessFilePath
	}
	return path
}

func (c *AppConfig) GetLoggingOptions(label string) *log.LoggingOptions {
	logfilePath := "/logs"
	if len(c.Logging.FilePath) > 0 {
		logfilePath = c.Logging.FilePath
	}

	logFileName := "log"
	hostname, err := os.Hostname()
	if err != nil {
		hostname = strconv.FormatInt(time.Now().Unix(), 10)
	}
	logFileName = fmt.Sprintf("log-%s.%s.txt", hostname, label)

	return &log.LoggingOptions{
		DefaultLogLevel: c.Logging.DefaultLogLevel,
		FilePath:        filepath.Join(logfilePath, logFileName),
	}
}

func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}
