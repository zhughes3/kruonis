package main

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type imageServerConfig struct {
	acctName, acctKey, containerName string
}

type serverConfig struct {
	rpcHost, rpcPort, httpHost, httpPort, jwtKey, frontend string
}

type dbConfig struct {
	name, host, port, user, password string
}

func main() {
	initLogging()
	config_defaults := map[string]interface{}{
		"rpc_host":  "localhost",
		"http_host": "localhost",
		"db_host":   "localhost",
	}
	config, err := readConfig("config.env", config_defaults)
	if err != nil {
		panic(err)
	}
	sCfg := getServerConfig(config)
	dbCfg := getDBConfig(config)
	imageServerConfig := getImageServerConfig(config)
	conn, err := net.Listen("tcp", ":"+sCfg.httpPort)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	db := NewDB(dbCfg)
	defer db.Close()

	imageClient := NewImageClient(imageServerConfig)
	s := NewServer(sCfg, conn, db, imageClient)
	s.Start()
}

func initLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func getServerConfig(cfg *viper.Viper) *serverConfig {
	return &serverConfig{
		rpcHost:  cfg.GetString("rpc_host"),
		rpcPort:  cfg.GetString("rpc_port"),
		httpHost: cfg.GetString("http_host"),
		httpPort: cfg.GetString("http_port"),
		jwtKey:   cfg.GetString("jwt_key"),
		frontend: cfg.GetString("frontend_url"),
	}
}

func getDBConfig(cfg *viper.Viper) *dbConfig {
	return &dbConfig{
		name:     cfg.GetString("db_name"),
		host:     cfg.GetString("db_host"),
		port:     cfg.GetString("db_port"),
		user:     cfg.GetString("db_user"),
		password: cfg.GetString("db_password"),
	}
}

func getImageServerConfig(cfg *viper.Viper) *imageServerConfig {
	return &imageServerConfig{
		acctName:      cfg.GetString("blob_account"),
		acctKey:       cfg.GetString("blob_key"),
		containerName: cfg.GetString("environment"),
	}
}
