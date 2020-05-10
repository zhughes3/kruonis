package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/zhughes3/kruonis/cmd/timelines/models"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	errBadCredentials error          = errors.New("Bad username, password combo")
	errInvalidToken   error          = errors.New("Invalid JSON Web Token")
	errNoBearerToken  error          = errors.New("Must supply bearer token to complete request")
	insecureEndpoints map[string]int = map[string]int{
		"/models.TimelineService/Login":  1,
		"/models.TimelineService/Signup": 1,
	}
)

type claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (s *server) Signup(ctx context.Context, in *models.SignupRequest) (*models.Error, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Problem generating hash from password: %v", in.GetPassword())
		return nil, err
	}

	_, err = s.db.insertUser(in.GetEmail(), fmt.Sprintf("%s", hash))
	if err != nil {
		return nil, err
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "201")); err != nil {
		return nil, err
	}

	return &models.Error{
		Response: true,
	}, nil
}

func (s *server) Login(ctx context.Context, in *models.LoginRequest) (*models.LoginResponse, error) {
	//call user db get hash for username
	user, err := s.db.readUser(in.GetEmail())
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.GetHash()), []byte(in.GetPassword())) != nil {
		//login unsuccessful
		if err := grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "401")); err != nil {
			return nil, err
		}
		return nil, errBadCredentials
	}

	token, err := s.generateJWTToken(in.GetEmail())
	if err != nil {
		return nil, err
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("token", token)); err != nil {
		return nil, err
	}

	return &models.LoginResponse{}, nil
}

func (s *server) authUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := insecureEndpoints[info.FullMethod]; ok {
		// no authentication needed
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok && md.Get("authorization") != nil {
		authHeader := md.Get("authorization")
		jwtoken := stripBearerPrefix(authHeader[0])
		if err := s.validateToken(jwtoken); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}

	return nil, errNoBearerToken
}

func (s *server) validateToken(encodedToken string) error {
	claims := &claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(encodedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		return err
	}
	if !tkn.Valid {
		return errInvalidToken
	}
	return nil
}

func (s *server) generateJWTToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
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

func stripBearerPrefix(in string) string {
	return in[7:]
}
