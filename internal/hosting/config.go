package hosting

import (
	"log"

	"github.com/spf13/viper"
)

// Loads a configuration instance using values from .env and from environment variables
func LoadConfiguration(path string, name *string, configuration any) error {
	err := readConfig(path, name, configuration)

	if err != nil {
		log.Fatalf("error reading in configuration from '%s'", path)
	}
	return nil
}

// LoadConfig reads configuration from file or environment variables.
func readConfig(path string, name *string, configuration any) error {
	viper.AddConfigPath(path)

	if name == nil {
		viper.SetConfigFile(".env")
	} else {
		viper.SetConfigName(*name)
	}
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	viper.Unmarshal(&configuration)

	if err != nil {
		log.Fatal(err)
	}

	return err
}
