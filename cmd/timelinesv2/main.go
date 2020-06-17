package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	initLogging()

	config, err := readConfig("config.env")
	if err != nil {
		panic(err)
	}

	serverConfig := getHttpServerConfig(config)
	databaseConfig := getDatabaseConfig(config)
	imageBlobStoreConfig := getImageBlobStoreConfig(config)

	server := newServer(serverConfig)
	server.WithDB(databaseConfig)
	server.WithImageBlobStoreClient(imageBlobStoreConfig)

	server.Start()
}

func initLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
