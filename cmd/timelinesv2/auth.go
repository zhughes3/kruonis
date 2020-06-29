package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"

	log "github.com/sirupsen/logrus"
)

var (
	errBadAccessToken  error           = errors.New("Invalid access token.")
	errBadRefreshToken error           = errors.New("Invalid refresh token. Must login.")
	errForbiddenRoute  error           = errors.New("Access forbidden.")
	insecureEndpoints  map[string]bool = map[string]bool{
		"/v1/Login":              true,
		"/v1/Signup":             true,
		"/v1/ListTrendingGroups": true,
	}
	adminEndpoints map[string]bool = map[string]bool{
		"/v1/AdminListGroups":    true,
		"/v1/AdminListTimelines": true,
		"/v1/AdminListUsers":     true,
	}
	accessEndpoints map[string]bool = map[string]bool{
		"/v1/Ping":             true,
		"/v1/Me":               true,
		"/v1/CreateEventImage": true,
		"/v1/UpdateEventImage": true,
		"/v1/DeleteEventImage": true,
	}
)

type (
	accessTokenClaims struct {
		UserID  uint64 `json:"id"`
		Email   string `json:"email"`
		IsAdmin bool   `json:"is_admin"`
		jwt.StandardClaims
	}
	refreshTokenClaims struct {
		UserID uint64 `json:"id"`
		jwt.StandardClaims
	}
)

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentRoute := mux.CurrentRoute(r)
		routeName := currentRoute.GetName()

		log.Println(routeName)

		if _, ok := insecureEndpoints[routeName]; ok {
			// no authentication needed
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()

		access, _ := getCookieByName(r, "access")
		refresh, _ := getCookieByName(r, "refresh")

		// handle refresh token
		if routeName == "/v1/Refresh" {
			if len(refresh) == 0 {
				http.Error(w, errBadAccessToken.Error(), http.StatusUnauthorized)
				return
			}
			claims, err := s.validateRefreshToken(refresh)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r = r.WithContext(NewContextWithRefreshTokenClaims(ctx, *claims))
			next.ServeHTTP(w, r)
			return
		}

		if len(access) > 0 {
			// handle access token
			claims, err := s.validateAccessToken(access)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if _, ok := adminEndpoints[routeName]; ok {
				if claims.IsAdmin {
					next.ServeHTTP(w, r.WithContext(NewContextWithAccessTokenClaims(ctx, *claims)))
					return
				}
				http.Error(w, errForbiddenRoute.Error(), http.StatusForbidden)
				return
			}

			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r.WithContext(NewContextWithAccessTokenClaims(ctx, *claims)))
			return
		}

		next.ServeHTTP(w, r)
		return
	})
}

func getCookieByName(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (s *server) validateAccessToken(encodedToken string) (*accessTokenClaims, error) {
	claims := &accessTokenClaims{}

	// Parse the JWT string and store the result in `accessTokenClaims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(encodedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errBadAccessToken
	}

	return claims, nil
}

func (s *server) validateRefreshToken(encodedToken string) (*refreshTokenClaims, error) {
	claims := &refreshTokenClaims{}

	tkn, err := jwt.ParseWithClaims(encodedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errBadRefreshToken
	}

	return claims, nil
}

func (s *server) generateTokenPair(id uint64, email string, isAdmin bool) (map[string]string, error) {
	access, err := s.generateJWTAccessToken(id, email, isAdmin)
	if err != nil {
		return nil, err
	}
	refresh, err := s.generateJWTRefreshToken(id)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"access":  access,
		"refresh": refresh,
	}, nil
}

func (s *server) generateJWTAccessToken(id uint64, email string, isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &accessTokenClaims{
		IsAdmin: isAdmin,
		Email:   email,
		UserID:  id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *server) generateJWTRefreshToken(id uint64) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &refreshTokenClaims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
