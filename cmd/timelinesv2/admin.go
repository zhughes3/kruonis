package main

import (
	"encoding/json"
	"net/http"
)

func (s *server) AdminListUsersHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.db.readUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) AdminListGroupsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.db.readGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) AdminListTimelinesHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.db.readTimelines()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
