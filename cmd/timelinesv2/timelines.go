package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type (
	Timeline struct {
		Id        uint64    `json:"id,omitempty"`
		GroupId   uint64    `json:"group_id,omitempty"`
		Title     string    `json:"title,omitempty"`
		Tags      []string  `json:"tags,omitempty"`
		Events    []*Event  `json:"events,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}
	CreateTimelineRequest struct {
		GroupId uint64   `json:"group_id"`
		Title   string   `json:"title,omitempty"`
		Tags    []string `json:"tags,omitempty"`
	}
	CreateTimelineEventRequest struct {
		Title       string    `json:"title,omitempty"`
		Timestamp   time.Time `json:"timestamp,omitempty"`
		Description string    `json:"description,omitempty"`
		Content     string    `json:"content,omitempty"`
	}
	TimelineIDWithEvents struct {
		Id     string   `json:"id,omitempty"`
		Events []*Event `json:"events,omitempty"`
	}
	UpdateTimelineRequest struct {
		Title string   `json:"title,omitempty"`
		Tags  []string `json:"tags,omitempty"`
	}
)

func (s *server) CreateTimelineHandler(w http.ResponseWriter, r *http.Request) {
	var body CreateTimelineRequest
	claims := AccessTokenClaimsFromContext(r.Context())

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.GroupId == 0 {
		group, err := s.db.insertGroup(body.Title, claims.UserID, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body.GroupId = group.Id
	}

	timeline, err := s.db.insertTimeline(body.GroupId, body.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, tag := range body.Tags {
		_, err := s.db.insertTag(tag, timeline.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	respJSON, _ := json.Marshal(timeline)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) CreateTimelineEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var body CreateTimelineEventRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := s.db.insertTimelineEvent(id, body.Title, body.Description, body.Content, body.Timestamp)

	w.WriteHeader(http.StatusCreated)
	respJSON, _ := json.Marshal(event)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) ReadTimelineHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	timeline, err := s.db.readTimeline(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(timeline)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) ReadTimelineEventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	events, err := s.db.readTimelineEvents(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(TimelineIDWithEvents{
		Id:     id,
		Events: events,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) UpdateTimelineHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var body UpdateTimelineRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	timeline, err := s.db.updateTimeline(id, body.Title, body.Tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(timeline)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) DeleteTimelineHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.db.deleteTimeline(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(booleanResponse{true})

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
