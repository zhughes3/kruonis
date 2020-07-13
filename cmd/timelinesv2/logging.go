package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

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
			"Method":   routeName,
		}).Info("request info")
		return
	})
}
