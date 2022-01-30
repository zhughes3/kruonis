package main

import (
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type (
	server struct {
		host, port  string
		lis         net.Listener
		server      *http.Server
		imageClient *imageBlobStoreClient
		db          *db
		jwtKey      []byte
		cors        *cors.Cors
	}
)

func newServer(cfg *httpServerConfig) *server {
	return &server{
		host:   cfg.host,
		port:   cfg.port,
		jwtKey: []byte(cfg.jwtKey),
		cors: cors.New(cors.Options{
			AllowedOrigins: []string{cfg.frontendURL},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}),
	}
}

func (s *server) WithDB(cfg *databaseConfig) {
	s.db = newDB(cfg)
}

func (s *server) WithImageBlobStoreClient(cfg *imageBlobStoreConfig) {
	s.imageClient = newImageBlobStoreClient(cfg)
}

func (s *server) Start() error {
	var err error
	s.server, err = s.prepareServer()
	if err != nil {
		panic(err)
	}

	conn, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log.WithFields(log.Fields{
		"Port": s.port,
		"Host": s.host,
	}).Info("API Server started")

	return s.server.Serve(conn)
}

func (s *server) prepareServer() (*http.Server, error) {
	r := mux.NewRouter()
	routes := getRoutes(s)
	addRoutes(r, routes)

	r.Use(s.loggingMiddleware, s.authMiddleware)

	err := r.Walk(outputRoutesFn)
	if err != nil {
		log.Error(err)
	}
	handler := s.cors.Handler(r)

	return &http.Server{
		Handler:      handler,
		Addr:         s.host + ":" + s.port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}, nil
}
