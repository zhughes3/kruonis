package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		Id        uint64    `json:"id,omitempty"`
		Email     string    `json:"email,omitempty"`
		IsAdmin   bool      `json:"is_admin,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}
	UserWithTimelines struct {
		User   User     `json:"user,omitempty"`
		Groups []*Group `json:"groups,omitempty"`
	}
	accountCredentials struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}
	booleanResponse struct {
		Response bool `json:"response, omitempty"`
	}
)

func (s *server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var body accountCredentials

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.db.readUserByEmail(body.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.hash), []byte(body.Password)) != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tokens, err := s.generateTokenPair(user.Id, user.Email, user.IsAdmin)
	http.SetCookie(w, &http.Cookie{
		Name:     "access",
		Value:    tokens["access"],
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    tokens["refresh"],
		HttpOnly: true,
		Path:     "/",
	})
}

func (s *server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	setLogoutCookies(w)
}

func (s *server) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var body accountCredentials

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = s.db.insertUser(body.Email, fmt.Sprintf("%s", hash))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := &booleanResponse{Response: true}
	respJSON, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	claims := RefreshTokenClaimsFromContext(r.Context())
	user, err := s.db.readUserByID(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respJSON, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) PingHandler(w http.ResponseWriter, r *http.Request) {
	var resp bool
	claims := AccessTokenClaimsFromContext(r.Context())
	if claims.UserID != 0 {
		resp = true
	} else {
		resp = false
	}
	respJSON, _ := json.Marshal(booleanResponse{Response: resp})
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) ReadMeHandler(w http.ResponseWriter, r *http.Request) {
	claims := AccessTokenClaimsFromContext(r.Context())
	groups, err := s.db.readUserTimelineGroups(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respJSON, _ := json.Marshal(UserWithTimelines{
		User: User{
			Id:      claims.UserID,
			Email:   claims.Email,
			IsAdmin: claims.IsAdmin,
		},
		Groups: groups,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func setLogoutCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})
}
