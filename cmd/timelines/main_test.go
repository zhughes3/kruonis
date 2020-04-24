package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/ory/dockertest/v3"
)

var (
	testServer *server
)

func TestMain(m *testing.M) {
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
	cfg := getDBConfig(config)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("kruonis_postgres", "latest", []string{"POSTGRES_USER=" + cfg.user, "POSTGRES_PASSWORD=" + cfg.password, "POSTGRES_DB=" + cfg.name, "DATABASE_HOST=" + cfg.host})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	conn, err := net.Listen("tcp", ":"+sCfg.httpPort)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(func() error {
		var err error
		database, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.user, cfg.password, cfg.host, cfg.port, cfg.name, "disable"))
		if err != nil {
			return err
		}
		testServer = NewServer(sCfg, conn, database)

		//testServer.Start()
		return database.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	exitVal := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(exitVal)
}
