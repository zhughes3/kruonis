package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

var uuidRegex = regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$")
var errInvalidGroupOperation = errors.New("Cannot update a group that you didn't create.")

type (
	Group struct {
		Id        uint64      `json:"id,omitempty"`
		Uuid      string      `json:"uuid,omitempty"`
		Title     string      `json:"title,omitempty"`
		Timelines []*Timeline `json:"timelines,omitempty"`
		CreatedAt time.Time   `json:"created_at,omitempty"`
		UpdatedAt time.Time   `json:"updated_at,omitempty"`
		Private   bool        `json:"private,omitempty"`
		UserId    uint64      `json:"user_id,omitempty"`
		Views     uint64 	  `json:"views,omitempty"`
	}

	UpdateGroupRequest struct {
		Title   string `json:"title,omitempty"`
		Private bool   `json:"private,omitempty"`
	}
)

func (s *server) ReadGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	claims := AccessTokenClaimsFromContext(r.Context())

	group, err := s.db.readGroup(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if group.Private && group.UserId != claims.UserID {
		http.Error(w, errForbiddenRoute.Error(), http.StatusForbidden)
		return
	}

	go func(id string) {
		err := s.db.incrementGroupViews(id)
		if err != nil {
			log.Errorf("Error incrementing group views: %s", err)
		}
	}(id)

	go func(id string) {
		err := s.db.insertGroupView(id)
		if err != nil {
			log.Errorf("Error inserting group view: %s", err)
		}
	}(id)

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(group)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	claims := AccessTokenClaimsFromContext(r.Context())

	currentGroup, err := s.db.readGroup(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if claims.UserID != currentGroup.UserId && !claims.IsAdmin {
		http.Error(w, errInvalidGroupOperation.Error(), http.StatusUnauthorized)
		return
	}

	var body UpdateGroupRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group, err := s.db.updateGroup(id, body.Title, body.Private)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(group)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	claims := AccessTokenClaimsFromContext(r.Context())

	currentGroup, err := s.db.readGroup(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if claims.UserID != currentGroup.UserId && !claims.IsAdmin {
		http.Error(w, errInvalidGroupOperation.Error(), http.StatusUnauthorized)
		return
	}

	if err := s.db.deleteGroup(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(booleanResponse{true})

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) ListTrendingGroupsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO this should be returned in popularity order
	s.AdminListGroupsHandler(w, r)
}

func (s *server) ListPublicGroupsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO 
	s.AdminListGroupsHandler(w, r)
}
