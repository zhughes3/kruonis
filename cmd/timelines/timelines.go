package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zhughes3/kruonis/cmd/timelines/models"
)

func (s *server) CreateTimeline(ctx context.Context, t *models.Timeline) (*models.Timeline, error) {
	groupId := t.GetGroupId()
	if groupId == 0 {
		timelineGroup, err := s.db.insertTimelineGroup(t.GetTitle(), UserIDFromContext(ctx), false)
		if err != nil {
			return nil, err
		}
		groupId = timelineGroup.GetId()
	}

	timeline, err := s.db.insertTimeline(groupId, t.GetTitle())
	if err != nil {
		return nil, err
	}

	for _, tag := range t.GetTags() {
		_, err := s.db.insertTag(tag, timeline.GetId())
		if err != nil {
			return nil, err
		}
		timeline.Tags = append(timeline.Tags, tag)
	}

	grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "201"))

	return timeline, nil
}

func (s *server) ReadTimeline(ctx context.Context, t *models.Filter) (*models.Timeline, error) {
	timeline, err := s.db.readTimeline(t.GetId())
	if err != nil {
		return nil, err
	}
	return timeline, nil
}

func (s *server) UpdateTimeline(ctx context.Context, in *models.UpdateTimelineRequest) (*models.Timeline, error) {
	timeline, err := s.db.updateTimeline(in.GetId(), in.GetTitle(), in.GetTags())
	if err != nil {
		return nil, err
	}

	return timeline, nil
}

func (s *server) DeleteTimeline(ctx context.Context, t *models.Filter) (*models.Error, error) {
	if err := s.db.deleteTimeline(t.GetId()); err != nil {
		return &models.Error{Response: false}, err
	}

	return &models.Error{Response: true}, nil
}

func (s *server) ReadTimelineGroup(ctx context.Context, t *models.Filter) (*models.TimelineGroup, error) {
	tg, err := s.db.readTimelineGroup(ctx, t.GetId())
	if err != nil {
		if err == errUnauthorized {
			grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "401"))
		}
		return nil, err
	}

	return tg, nil
}

func (s *server) UpdateTimelineGroup(ctx context.Context, in *models.UpdateTimelineGroupRequest) (*models.TimelineGroup, error) {
	timelineGroup, err := s.db.updateTimelineGroup(in.GetId(), in.GetTitle(), in.GetPrivate())
	if err != nil {
		return nil, err
	}

	return timelineGroup, nil
}

func (s *server) DeleteTimelineGroup(ctx context.Context, t *models.Filter) (*models.Error, error) {
	if err := s.db.deleteTimelineGroup(t.GetId()); err != nil {
		return &models.Error{Response: false}, err
	}

	return &models.Error{Response: true}, nil
}

func (s *server) CreateTimelineEvent(ctx context.Context, t *models.TimelineEvent) (*models.TimelineEvent, error) {
	timelineEvent, err := s.db.insertTimelineEvent(t.GetId(), t.GetTitle(), t.GetDescription(), t.GetContent(), t.GetTimestamp())
	if err != nil {
		return nil, err
	}

	grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "201"))

	return timelineEvent, nil
}

func (s *server) ReadTimelineEvents(ctx context.Context, in *models.Filter) (*models.ReadTimelineEventsResponse, error) {
	events, err := s.db.readTimelineEvents(in.GetId())
	if err != nil {
		return nil, err
	}

	return &models.ReadTimelineEventsResponse{
		Id:     in.Id,
		Events: events,
	}, nil
}

func (s *server) ReadTimelineEvent(ctx context.Context, t *models.Filter) (*models.TimelineEvent, error) {
	tg, err := s.db.readTimelineEvent(t.GetId())
	if err != nil {
		return nil, err
	}

	return tg, nil
}

func (s *server) UpdateTimelineEvent(ctx context.Context, t *models.UpdateTimelineEventRequest) (*models.TimelineEvent, error) {
	timelineEvent, err := s.db.updateTimelineEvent(t.GetId(), t.GetTitle(), t.GetDescription(), t.GetContent(), t.Timestamp)
	if err != nil {
		return nil, err
	}

	return timelineEvent, nil
}

func (s *server) DeleteTimelineEvent(ctx context.Context, t *models.Filter) (*models.Error, error) {
	if err := s.db.deleteTimelineEvent(t.GetId()); err != nil {
		return &models.Error{Response: false}, err
	}

	return &models.Error{Response: true}, nil
}

func (s *server) ListTrendingTimelineGroups(ctx context.Context, in *models.TrendingTimelineGroupsRequest) (*models.TrendingTimelineGroupsResponse, error) {
	// TODO this needs to get the last timeline groups, sorted by hit count
	resp, err := s.db.readTimelineGroups()
	if err != nil {
		return nil, err
	}
	return &models.TrendingTimelineGroupsResponse{
		Groups: resp.Groups,
	}, nil
}
