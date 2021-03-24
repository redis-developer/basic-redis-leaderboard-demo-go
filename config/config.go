package config

import (
	"os"
	"strconv"
)

type Config struct {
	Api    *Api
	Redis  *Redis
	Import *Import
}

func envReadString(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		value = defaultValue
	}
	return value
}

func envReadNumeric(envName string, defaultValue int) int {
	value, _ := strconv.Atoi(os.Getenv(envName))
	if value == 0 {
		value = defaultValue
	}
	return value
}

func envReadBool(envName string) bool {
	value, _ := strconv.ParseBool(os.Getenv(envName))
	return value
}

func NewConfig() *Config {

	apiPort := envReadNumeric(envConfigApiPort, defaultConfigApiPort)

	externalEnvApiPort := os.Getenv(envExternalConfigApiPort)
	if externalEnvApiPort != "" {
		val, err := strconv.Atoi(externalEnvApiPort)
		if err == nil {
			apiPort = val
		}
	}

	config := &Config{
		Api: &Api{
			host:        envReadString(envConfigApiHost, ""),
			port:        apiPort,
			publicPath:  envReadString(envConfigApiPublicPath, defaultConfigApiPublicPath),
			tlsDisabled: envReadBool(envConfigApiTLSDisabled),
		},
		Redis: &Redis{
			host:     envReadString(envConfigRedisHost, ""),
			port:     envReadNumeric(envConfigRedisPort, defaultConfigRedisPort),
			password: envReadString(envConfigRedisPassword, ""),
		},
		Import: &Import{
			path: envReadString(envImportPath, defaultImportPath),
		},
	}

	return config

}
