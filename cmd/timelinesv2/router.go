package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type (
	route struct {
		path        string
		httpMethods []string
		fn          func(w http.ResponseWriter, r *http.Request)
	}
	routes []route
)

func getRoutes(s *server) routes {
	routes := []route{
		{path: "/v1/groups/trending", httpMethods: []string{http.MethodGet}, fn: s.ListTrendingGroupsHandler},
		{path: "/v1/groups/public", httpMethods: []string{http.MethodGet}, fn: s.ListPublicGroupsHandler},
		{path: "/v1/timelines", httpMethods: []string{http.MethodPost}, fn: s.CreateTimelineHandler},
		{path: "/v1/timelines/{id}", httpMethods: []string{http.MethodGet}, fn: s.ReadTimelineHandler},
		{path: "/v1/timelines/{id}", httpMethods: []string{http.MethodPut}, fn: s.UpdateTimelineHandler},
		{path: "/v1/timelines/{id}", httpMethods: []string{http.MethodDelete}, fn: s.DeleteTimelineHandler},
		{path: "/v1/timelines/{id}/events", httpMethods: []string{http.MethodGet}, fn: s.ReadTimelineEventsHandler},
		{path: "/v1/timelines/{id}/events", httpMethods: []string{http.MethodGet}, fn: s.CreateTimelineEventHandler},
		{path: "/v1/groups/{id}", httpMethods: []string{http.MethodGet}, fn: s.ReadGroupHandler},
		{path: "/v1/groups/{id}", httpMethods: []string{http.MethodPut}, fn: s.UpdateGroupHandler},
		{path: "/v1/groups/{id}", httpMethods: []string{http.MethodDelete}, fn: s.DeleteGroupHandler},
		{path: "/v1/events/{id}", httpMethods: []string{http.MethodGet}, fn: s.ReadEventHandler},
		{path: "/v1/events/{id}", httpMethods: []string{http.MethodPut}, fn: s.UpdateEventHandler},
		{path: "/v1/events/{id}", httpMethods: []string{http.MethodDelete}, fn: s.DeleteEventHandler},
		{path: "/v1/events/{id}/img", httpMethods: []string{http.MethodPost}, fn: s.CreateEventImageHandler},
		{path: "/v1/events/{id}/img", httpMethods: []string{http.MethodPut}, fn: s.UpdateEventImageHandler},
		{path: "/v1/events/{id}/img", httpMethods: []string{http.MethodDelete}, fn: s.DeleteEventImageHandler},
		{path: "/v1/users/login", httpMethods: []string{http.MethodPost}, fn: s.LoginHandler},
		{path: "/v1/users/logout", httpMethods: []string{http.MethodGet}, fn: s.LogoutHandler},
		{path: "/v1/users/signup", httpMethods: []string{http.MethodPost}, fn: s.SignupHandler},
		{path: "/v1/users/ping", httpMethods: []string{http.MethodGet}, fn: s.PingHandler},
		{path: "/v1/users/refresh", httpMethods: []string{http.MethodGet}, fn: s.RefreshHandler},
		{path: "/v1/users/me", httpMethods: []string{http.MethodGet}, fn: s.ReadMeHandler},
		{path: "/v1/admin/users", httpMethods: []string{http.MethodGet}, fn: s.AdminListUsersHandler},
		{path: "/v1/admin/groups", httpMethods: []string{http.MethodGet}, fn: s.AdminListGroupsHandler},
		{path: "/v1/admin/timelines", httpMethods: []string{http.MethodGet}, fn: s.AdminListTimelinesHandler},
	}
	return routes
}

func addRoutes(router *mux.Router, rt routes) {
	for _, r := range rt {
		router.HandleFunc(r.path, r.fn).Methods(r.httpMethods...).Name(r.path)
	}
}

func outputRoutesFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	pathTemplate, _ := route.GetPathTemplate()
	methods, _ := route.GetMethods()
	log.WithFields(log.Fields{
		"Methods":  strings.Join(methods, ","),
		"Endpoint": pathTemplate,
	}).Info("HTTP Route")
	return nil
}
