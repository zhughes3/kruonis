package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Event - a struct representing a timeline event
type Event struct {
	ID          uint64    `json:"id,omitempty"`
	TimelineID  uint64    `json:"timeline_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
	Description string    `json:"description,omitempty"`
	Content     string    `json:"content,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
}

func (s *server) ReadEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	event, err := s.db.readEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(event)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var body CreateTimelineEventRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := s.db.updateEvent(id, body.Title, body.Description, body.Content, body.Timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(event)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.db.deleteEvent(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(booleanResponse{true})

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
