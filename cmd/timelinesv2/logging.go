package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (s *server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentRoute := mux.CurrentRoute(r)
		routeName := currentRoute.GetName()
		start := time.Now()
		next.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"Duration": time.Since(start).String(),
		}).Info(routeName)
		return
	})
}