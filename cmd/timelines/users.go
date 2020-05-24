package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/zhughes3/kruonis/cmd/timelines/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UserWithTimelines struct {
	models.User
	timelines []*models.TimelineGroup
}

type UserWithHash struct {
	models.User
	hash string
}

func (s *server) Signup(ctx context.Context, in *models.SignupRequest) (*models.Error, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Problem generating hash from password: %s", in.GetPassword())
		return nil, err
	}

	_, err = s.db.insertUser(in.GetEmail(), fmt.Sprintf("%s", hash))
	if err != nil {
		return nil, err
	}

	grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "201"))

	return &models.Error{
		Response: true,
	}, nil
}

func (s *server) Login(ctx context.Context, in *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.db.readUserByEmail(in.GetEmail())
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.hash), []byte(in.GetPassword())) != nil {
		//login unsuccessful
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "401"))
		return nil, errBadCredentials
	}

	tokens, err := s.generateTokenPair(user.Id, user.GetEmail(), user.GetIsAdmin())

	if err != nil {
		return nil, err
	}

	if len(tokens) == 2 {
		grpc.SetHeader(ctx, metadata.Pairs("access", tokens["access"]))
		grpc.SetHeader(ctx, metadata.Pairs("refresh", tokens["refresh"]))
	}

	return &models.LoginResponse{}, nil
}

func (s *server) Refresh(ctx context.Context, in *models.RefreshRequest) (*models.RefreshResponse, error) {
	user, err := s.db.readUserByID(UserIDFromContext(ctx))
	accessToken, err := s.generateJWTAccessToken(user.Id, user.GetEmail(), user.GetIsAdmin())

	if err != nil {
		return nil, err
	}

	grpc.SetHeader(ctx, metadata.Pairs("access", accessToken))

	return &models.RefreshResponse{}, nil
}

func (s *server) Ping(ctx context.Context, in *models.PingRequest) (*models.PingResponse, error) {
	userID := UserIDFromContext(ctx)
	if userID == 0 {
		return &models.PingResponse{Response: false}, nil
	}
	return &models.PingResponse{Response: true}, nil
}

func (s *server) Me(ctx context.Context, in *models.MeRequest) (*models.MeResponse, error) {
	user, err := s.db.readUserByID(UserIDFromContext(ctx))
	if err != nil {
		return nil, err
	}

	groups, err := s.db.readUserTimelineGroups(user.Id)
	if err != nil {
		return nil, err
	}
	return &models.MeResponse{
		User:   user,
		Groups: groups,
	}, nil
}
