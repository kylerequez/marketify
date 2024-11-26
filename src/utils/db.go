package utils

import "github.com/kylerequez/marketify/src/shared"

func RetrieveSQLConfig() (*shared.SQLConfig, error) {
	uri, err := GetEnv("DB_URI")
	if err != nil {
		return nil, err
	}

	config := shared.SQLConfig{
		URI: *uri,
	}

	return &config, nil
}
