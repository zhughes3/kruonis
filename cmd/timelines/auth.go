package main

import (
	"context"
	"errors"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/metadata"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
)

var (
	errBadCredentials  error          = errors.New("Bad username, password combo")
	errBadAccessToken  error          = errors.New("Invalid access token.")
	errBadRefreshToken error          = errors.New("Invalid refresh token. Must login.")
	errNoBearerToken   error          = errors.New("Must supply bearer token to complete request")
	insecureEndpoints  map[string]int = map[string]int{
		"/models.TimelineService/Login":  1,
		"/models.TimelineService/Signup": 1,
	}
	adminEndpoints map[string]int = map[string]int{
		"/models.TimelineService/ReadGroups":    1,
		"/models.TimelineService/ReadTimelines": 1,
		"/models.TimelineService/ReadUsers":     1,
	}
)

type accessTokenClaims struct {
	UserID  uint64 `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json: "is_admin"`
	jwt.StandardClaims
}

type refreshTokenClaims struct {
	UserID uint64 `json:"id"`
	jwt.StandardClaims
}

func (s *server) authUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := insecureEndpoints[info.FullMethod]; ok {
		// no authentication needed
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok && md.Get("grpcgateway-cookie") != nil {
		cookieHeader := md.Get("grpcgateway-cookie")
		tokens := tokensFromCookieHeader(cookieHeader[0])

		// handle refresh bearer token
		if info.FullMethod == "/models.TimelineService/Refresh" {
			claims, err := s.validateRefreshToken(tokens.Refresh)
			if err != nil {
				return nil, err
			}

			return handler(NewContextWithUserID(ctx, claims.UserID), req)
		}

		// handle access bearer token
		claims, err := s.validateAccessToken(tokens.Access)
		if err != nil {
			return nil, err
		}

		if _, ok := adminEndpoints[info.FullMethod]; ok {
			if claims.IsAdmin {
				return handler(NewContextWithUserID(ctx, claims.UserID), req)
			}
			return nil, err
		}

		return handler(NewContextWithUserID(ctx, claims.UserID), req)
	}

	// no auth
	return handler(ctx, req)
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
	claims := &accessTokenClaims{
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

type Tokens struct {
	Refresh, Access string
}

func tokensFromCookieHeader(in string) *Tokens {
	tokens := strings.Split(in, ";")
	if len(tokens) != 2 {
		return nil
	}

	ret := &Tokens{}
	for _, token := range tokens {
		kv := strings.Split(token, "=")
		if len(tokens) != 2 {
			return nil
		}
		switch strings.TrimSpace(kv[0]) {
		case "refresh":
			ret.Refresh = kv[1]
		case "access":
			ret.Access = kv[1]
		default:
			log.Error("Unkown key in cookie header: ", kv[0], " with value: ", kv[1])

		}
	}

	return ret
}
