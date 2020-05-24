package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"

	"github.com/ory/dockertest/v3"
	"github.com/pressly/goose"
)

var (
	testServer *server
)

func getTestServer(db *sql.DB) *server {
	scfg := &serverConfig{
		rpcHost:  "localhost",
		rpcPort:  "8081",
		httpHost: "localhost",
		httpPort: "8081",
	}

	return NewServer(scfg, nil, db, nil)
}

func getDatabaseConn(res *dockertest.Resource) *sql.DB {
	data := strings.Split(res.GetHostPort("5432/tcp"), ":")
	if len(data) != 2 {
		log.Fatal("Could not get database host and port")
	}

	host := data[0]
	port := data[1]

	db := NewDB(&dbConfig{
		host:     host,
		port:     port,
		name:     "timelines",
		user:     "timelines",
		password: "timelines"})

	return db
}

func pgdsn(dbname, host, user, password string) string {
	dsnTmpl := "postgres://%s:%s@%s/%s?sslmode=disable"
	return fmt.Sprintf(dsnTmpl, user, password, host, dbname)
}

func runDB(pool *dockertest.Pool) *dockertest.Resource {
	res, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", "timelines"),
			fmt.Sprintf("POSTGRES_USER=%s", "timelines"),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", "timelines"),
		},
	})
	if err != nil {
		log.Fatalf("Could not run resource: %s", err)
	}

	// expire the resource in 3 mins just in case tests fail
	if err := res.Expire(180); err != nil {
		log.Fatal(err)
	}

	var db *sql.DB
	if err = pool.Retry(func() error {
		var err2 error
		db, err2 = sql.Open(
			"postgres",
			pgdsn("timelines",
				res.GetHostPort("5432/tcp"),
				"timelines",
				"timelines"),
		)

		if err2 != nil {
			return err2
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker postgres: %s", err)
	}

	if err := runMigrations(db); err != nil {
		log.Fatalf("Unable to run migrations: %s", err)
	}

	return res
}

func runMigrations(db *sql.DB) error {
	mig := []string{
		"./db/migrations",
	}
	for _, m := range mig {
		if err := goose.Run("up", db, m); err != nil {
			return err
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	var (
		pool  *dockertest.Pool
		dbres *dockertest.Resource
	)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	dbres = runDB(pool)

	db := getDatabaseConn(dbres)

	testServer = getTestServer(db)

	exitVal := m.Run()

	if dbres != nil {
		if err := pool.Purge(dbres); err != nil {
			log.Fatal(err)
		}
	}

	log.Infoln("Ending test run against postgres")

	os.Exit(exitVal)
}
