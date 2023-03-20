package config

import (
	"errors"
	"log"
	"sync"

	"github.com/spf13/viper"
)

var lock = &sync.Mutex{}

type Configuration struct {
	Azure AzureAdSettings `mapstructure:"Azure"`
}

var configuration *Configuration

// Configures a new app configuration instance with the ability to configure it
// if default values are desired, pass nil as configure
func LoadConfiguration(path string) (*Configuration, error) {
	config := getConfiguration(path)

	if config == nil {
		return nil, errors.New("Configuration failed to load")
	}
	return config, nil
}

func getConfiguration(path string) *Configuration {
	if configuration == nil {
		lock.Lock()
		defer lock.Unlock()
		if configuration == nil {
			var err error
			configuration, err = readConfig(path)

			if err != nil {
				log.Printf("error reading in configuration from '%s'", path)
			}
		}
	}
	return configuration
}

// LoadConfig reads configuration from file or environment variables.
func readConfig(path string) (config *Configuration, err error) {
	viper.AddConfigPath(path)

	// TODO: incorporate ability to load based on APP_ENV
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		log.Fatal(err)
	}

	err = viper.UnmarshalKey("Azure", &config.Azure)
	return
}
