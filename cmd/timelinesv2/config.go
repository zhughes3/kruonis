package main

import "github.com/spf13/viper"

type (
	imageBlobStoreConfig struct {
		acctName, acctKey, containerName string
	}
	httpServerConfig struct {
		host, port, jwtKey, frontendUrl string
	}
	databaseConfig struct {
		name, host, port, user, password string
	}
)

func readConfig(filename string) (*viper.Viper, error) {
	defaults := map[string]interface{}{
		"rpc_host":  "localhost",
		"http_host": "localhost",
		"db_host":   "localhost",
	}

	v := viper.New()
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	v.SetConfigName("") //this will default to config.env
	v.AddConfigPath(".")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func getHttpServerConfig(cfg *viper.Viper) *httpServerConfig {
	return &httpServerConfig{
		host:        cfg.GetString("http_host"),
		port:        cfg.GetString("http_port"),
		jwtKey:      cfg.GetString("jwt_key"),
		frontendUrl: cfg.GetString("frontend_url"),
	}
}

func getDatabaseConfig(cfg *viper.Viper) *databaseConfig {
	return &databaseConfig{
		name:     cfg.GetString("db_name"),
		host:     cfg.GetString("db_host"),
		port:     cfg.GetString("db_port"),
		user:     cfg.GetString("db_user"),
		password: cfg.GetString("db_password"),
	}
}

func getImageBlobStoreConfig(cfg *viper.Viper) *imageBlobStoreConfig {
	return &imageBlobStoreConfig{
		acctName:      cfg.GetString("blob_account"),
		acctKey:       cfg.GetString("blob_key"),
		containerName: cfg.GetString("environment"),
	}
}
