package utils

import (
	"github.com/kylerequez/marketify/src/shared"
)

func RetrieveServerConfig() (*shared.ServerConfig, error) {
	name, err := GetEnv("APP_NAME")
	if err != nil {
		return nil, err
	}

	hostname, err := GetEnv("SERVER_HOST")
	if err != nil {
		return nil, err
	}

	port, err := GetEnv("SERVER_PORT")
	if err != nil {
		return nil, err
	}

	config := &shared.ServerConfig{
		AppName:  *name,
		Hostname: *hostname,
		Port:     *port,
	}

	return config, nil
}
