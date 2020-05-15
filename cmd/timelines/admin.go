package main

import (
	"context"

	"github.com/zhughes3/kruonis/cmd/timelines/models"
)

func (s *server) ReadGroups(ctx context.Context, in *models.ReadGroupsRequest) (*models.ReadGroupsResponse, error) {
	resp, err := s.db.readTimelineGroups()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (s *server) ReadTimelines(ctx context.Context, in *models.ReadTimelinesRequest) (*models.ReadTimelinesResponse, error) {
	resp, err := s.db.readTimelines()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (s *server) ReadUsers(ctx context.Context, in *models.ReadUsersRequest) (*models.ReadUsersResponse, error) {
	resp, err := s.db.readUsers()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
