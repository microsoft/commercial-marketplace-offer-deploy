package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/spf13/viper"
)

const EnvironmentVariablePrefix = "MODM"

// Loads a configuration instance using values from .env and from environment variables
func LoadConfiguration(path string, name *string, root any) error {
	builder := viper.New()
	builder.AddConfigPath(path)

	if name != nil {
		builder.SetConfigName(*name)
	} else {
		builder.SetConfigFile(filepath.Join(path, ".env"))
	}

	builder.SetConfigType("env")
	builder.SetEnvPrefix(EnvironmentVariablePrefix)
	builder.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	builder.AutomaticEnv()

	errors := []string{}

	err := builder.ReadInConfig()
	if err != nil {
		log.Error(err)
		errors = append(errors, err.Error())
	}

	automaticEnvs(builder)
	err = unmarshal(builder, root)

	if err != nil {
		log.Error(err)
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return utils.NewAggregateError(errors)
	}

	return nil
}

func unmarshal(builder *viper.Viper, root any) error {
	errors := []string{}

	err := builder.Unmarshal(&root)
	if err != nil {
		errors = append(errors, err.Error())
	}

	// unmarshall all the config sections at the root level. We could do this recursively, but we don't need to

	configType := getType(root)
	configValue := GetValue(root)

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)

		// if the field is a struct, we'll unmarshal it
		if field.Type.Kind() == reflect.Struct {
			fieldValue := configValue.FieldByName(field.Name)
			section := reflect.New(field.Type)
			err := builder.Unmarshal(section.Interface())
			fieldValue.Set(reflect.Indirect(section))

			if err != nil {
				errors = append(errors, err.Error())
			}
		}
	}
	return utils.NewAggregateError(errors)
}

// viper doesn't actually load environment variables into the config if they aren't present in the config file
// we still want the env file to take precedence, so we'll load the env variables into the config if they aren't already present
func automaticEnvs(builder *viper.Viper) {
	envKeys := getEnvironmentVariableKeys()
	for _, key := range envKeys {
		err := builder.BindEnv(strings.Replace(key, EnvironmentVariablePrefix+"_", "", 1))
		if err != nil {
			log.Error(err)
		}
	}
}

func getEnvironmentVariableKeys() []string {
	envs := []string{}
	for _, env := range os.Environ() {
		keyValuePair := strings.SplitN(env, "=", 2)
		envs = append(envs, keyValuePair[0])
	}
	return envs
}

func getType(i any) reflect.Type {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func GetValue(i any) reflect.Value {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
