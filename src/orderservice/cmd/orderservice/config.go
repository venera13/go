package main

import "github.com/kelseyhightower/envconfig"

const appID = "orderservice"

type config struct {
	ServeRESTAddress string `envconfig:"serve_rest_address" default:":8080"`
	DBName           string `envconfig:"mysql_database"`
	DBUser           string `envconfig:"mysql_user"`
	DBPass           string `envconfig:"mysql_password"`
}

func parseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, err
	}
	return c, nil
}
