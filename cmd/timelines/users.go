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
	user, err := s.db.readUserByEmail(in.GetEmail())
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

	tokens, err := s.generateTokenPair(user.Id, user.GetEmail(), user.GetIsAdmin())

	if err != nil {
		return nil, err
	}

	if len(tokens) == 2 {
		if err := grpc.SetHeader(ctx, metadata.Pairs("access", tokens["access"])); err != nil {
			return nil, err
		}
		if err := grpc.SetHeader(ctx, metadata.Pairs("refresh", tokens["refresh"])); err != nil {
			return nil, err
		}
	}

	return &models.LoginResponse{}, nil
}

func (s *server) Refresh(ctx context.Context, in *models.RefreshRequest) (*models.RefreshResponse, error) {
	user, err := s.db.readUserByID(UserIDFromContext(ctx))
	accessToken, err := s.generateJWTAccessToken(user.Id, user.GetEmail(), user.GetIsAdmin())

	if err != nil {
		return nil, err
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("access", accessToken)); err != nil {
		return nil, err
	}

	return &models.RefreshResponse{}, nil
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
