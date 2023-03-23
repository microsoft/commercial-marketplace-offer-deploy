package config

import (
	"errors"
	"log"
	"sync"

	"github.com/spf13/viper"
)

var lock = &sync.Mutex{}

type Configuration struct {
	Azure    AzureAdSettings
	Database DatabaseSettings
}

var configuration *Configuration

func GetConfiguration() *Configuration {
	return getConfiguration(".", nil)
}

// Configures a new app configuration instance with the ability to configure it
// if default values are desired, pass nil as configure
func LoadConfiguration(path string, name *string) (*Configuration, error) {
	config := getConfiguration(path, name)

	if config == nil {
		return nil, errors.New("Configuration failed to load")
	}
	return config, nil
}

// singelton implementation for configuration
func getConfiguration(path string, name *string) *Configuration {
	if configuration == nil {
		lock.Lock()
		defer lock.Unlock()
		if configuration == nil {
			var err error
			configuration, err = readConfig(path, name)

			if err != nil {
				log.Printf("error reading in configuration from '%s'", path)
			}
		}
	}
	return configuration
}

// LoadConfig reads configuration from file or environment variables.
func readConfig(path string, name *string) (config *Configuration, err error) {
	viper.AddConfigPath(path)

	if name == nil {
		viper.SetConfigFile(".env")
	} else {
		viper.SetConfigName(*name)
	}
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	c := &Configuration{}

	err = viper.Unmarshal(&c.Azure)

	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&c.Database)

	if err != nil {
		log.Fatal(err)
	}

	return c, err
}
