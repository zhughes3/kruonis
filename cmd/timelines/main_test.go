package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/ory/dockertest/v3"
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	config, err := readConfig("config.env", nil)
	if err != nil {
		panic(err)
	}
	dbCfg := getDBConfig(config)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("kruonis_postgres", "latest", []string{"POSTGRES_USER=" + dbCfg.user, "POSTGRES_PASSWORD=" + dbCfg.password, "POSTGRES_DB=" + dbCfg.name, "DATABASE_HOST=" + dbCfg.host})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(func() error {
		var err error
		database, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), database))
		if err != nil {
			return err
		}
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
