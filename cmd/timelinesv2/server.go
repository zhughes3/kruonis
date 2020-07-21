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

	s.addHTTPRoutes(r)

	r.Use(s.loggingMiddleware, s.authMiddleware)
	handler := s.cors.Handler(r)

	return &http.Server{
		Handler:      handler,
		Addr:         s.host + ":" + s.port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}, nil
}

func (s *server) addHTTPRoutes(r *mux.Router) {
	r.HandleFunc("/v1/groups/trending", s.ListTrendingGroupsHandler).Methods(http.MethodGet).Name("/v1/ListTrendingGroups")
	r.HandleFunc("/v1/groups/public", s.ListPublicGroupsHandler).Methods(http.MethodGet).Name("/v1/ListPublicGroups")
	r.HandleFunc("/v1/timelines", s.CreateTimelineHandler).Methods(http.MethodPost).Name("/v1/CreateTimeline")
	r.HandleFunc("/v1/timelines/{id}", s.ReadTimelineHandler).Methods(http.MethodGet).Name("/v1/ReadTimeline")
	r.HandleFunc("/v1/timelines/{id}", s.UpdateTimelineHandler).Methods(http.MethodPut).Name("/v1/UpdateTimeline")
	r.HandleFunc("/v1/timelines/{id}", s.DeleteTimelineHandler).Methods(http.MethodDelete).Name("/v1/DeleteTimeline")
	r.HandleFunc("/v1/timelines/{id}/events", s.ReadTimelineEventsHandler).Methods(http.MethodGet).Name("/v1/ReadTimelineEvents")
	r.HandleFunc("/v1/timelines/{id}/events", s.CreateTimelineEventHandler).Methods(http.MethodPost).Name("/v1/CreateTimelineEvent")
	r.HandleFunc("/v1/groups/{id}", s.ReadGroupHandler).Methods(http.MethodGet).Name("/v1/ReadGroup")
	r.HandleFunc("/v1/groups/{id}", s.UpdateGroupHandler).Methods(http.MethodPut).Name("/v1/UpdateGroup")
	r.HandleFunc("/v1/groups/{id}", s.DeleteGroupHandler).Methods(http.MethodDelete).Name("/v1/DeleteGroup")
	r.HandleFunc("/v1/events/{id}", s.ReadEventHandler).Methods(http.MethodGet).Name("/v1/ReadEvent")
	r.HandleFunc("/v1/events/{id}", s.UpdateEventHandler).Methods(http.MethodPut).Name("/v1/UpdateEvent")
	r.HandleFunc("/v1/events/{id}", s.DeleteEventHandler).Methods(http.MethodDelete).Name("/v1/DeleteEvent")
	r.HandleFunc("/v1/events/{id}/img", s.CreateEventImageHandler).Methods(http.MethodPost).Name("/v1/CreateEventImage")
	r.HandleFunc("/v1/events/{id}/img", s.UpdateEventImageHandler).Methods(http.MethodPut).Name("/v1/UpdateEventImage")
	r.HandleFunc("/v1/events/{id}/img", s.DeleteEventImageHandler).Methods(http.MethodDelete).Name("/v1/DeleteEventImage")
	r.HandleFunc("/v1/users/login", s.LoginHandler).Methods(http.MethodPost).Name("/v1/Login")
	r.HandleFunc("/v1/users/logout", s.LogoutHandler).Methods(http.MethodGet).Name("/v1/Logout")
	r.HandleFunc("/v1/users/signup", s.SignupHandler).Methods(http.MethodPost).Name("/v1/Signup")
	r.HandleFunc("/v1/users/ping", s.PingHandler).Methods(http.MethodGet).Name("/v1/Ping")
	r.HandleFunc("/v1/users/refresh", s.RefreshHandler).Methods(http.MethodGet).Name("/v1/Refresh")
	r.HandleFunc("/v1/users/me", s.ReadMeHandler).Methods(http.MethodGet).Name("/v1/ReadMe")
	r.HandleFunc("/v1/admin/users", s.AdminListUsersHandler).Methods(http.MethodGet).Name("/v1/AdminListUsers")
	r.HandleFunc("/v1/admin/groups", s.AdminListGroupsHandler).Methods(http.MethodGet).Name("/v1/AdminListGroups")
	r.HandleFunc("/v1/admin/timelines", s.AdminListTimelinesHandler).Methods(http.MethodGet).Name("/v1/AdminListTimelines")
}
