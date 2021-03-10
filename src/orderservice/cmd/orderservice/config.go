package main

import "github.com/kelseyhightower/envconfig"

const appID = "orderservice"

type databaseConfig struct {
	DBName string `envconfig:"db_name"`
	DBHost string `envconfig:"db_host"`
	DBUser string `envconfig:"db_user"`
	DBPass string `envconfig:"db_pass"`
}

type config struct {
	ServeRESTAddress string `envconfig:"serve_rest_address" default:":8080"`
	Database         databaseConfig
}

func parseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, err
	}
	return c, nil
}
