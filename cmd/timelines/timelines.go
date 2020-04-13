package main

import (
	"context"

	"github.com/zhughes3/kruonis/cmd/timelines/models"
)

func (s *server) CreateTimeline(ctx context.Context, t *models.Timeline) (*models.Timeline, error) {
	groupId := t.GetGroupId()
	if groupId == 0 {
		timelineGroup, err := s.db.insertTimelineGroup(t.GetTitle())
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

	return timeline, nil
}

func (s *server) ReadTimeline(ctx context.Context, t *models.Filter) (*models.Timeline, error) {
	timeline, err := s.db.readTimeline(t.GetId())
	if err != nil {
		return nil, err
	}
	return timeline, nil
}

//func (s *server) UpdateTimeline(ctx context.Context, t *models.UpdateTimelineRequest) (*models.Timeline, error) {
//
//}

func (s *server) DeleteTimeline(ctx context.Context, t *models.Filter) (*models.Error, error) {
	if err := s.db.deleteTimeline(t.GetId()); err != nil {
		return &models.Error{Response: false}, err
	}

	return &models.Error{Response: true}, nil
}

func (s *server) ReadTimelineGroup(ctx context.Context, t *models.Filter) (*models.TimelineGroup, error) {
	tg, err := s.db.readTimelineGroup(t.GetId())
	if err != nil {
		return nil, err
	}

	return tg, nil
}

//func (s *server) UpdateTimelineGroup(ctx context.Context, t *models.UpdateTimelineGroupRequest) (*models.TimelineGroup, error) {
//
//}

func (s *server) DeleteTimelineGroup(ctx context.Context, t *models.Filter) (*models.Error, error) {
	if err := s.db.deleteTimelineGroup(t.GetId()); err != nil {
		return &models.Error{Response: false}, err
	}

	return &models.Error{Response: true}, nil
}

func (s *server) CreateTimelineEvent(ctx context.Context, t *models.TimelineEvent) (*models.TimelineEvent, error) {
	timelineEvent, err := s.db.insertTimelineEvent(t.GetTimelineId(), t.GetTitle(), t.GetDescription(), t.GetContent(), t.GetTimestamp())
	if err != nil {
		return nil, err
	}

	return timelineEvent, nil
}

func (s *server) ReadTimelineEvent(ctx context.Context, t *models.Filter) (*models.TimelineEvent, error) {
	tg, err := s.db.readTimelineEvent(t.GetId())
	if err != nil {
		return nil, err
	}

	return tg, nil
}

//func (s *server) UpdateTimelineEvent(ctx context.Context, t *models.UpdateTimelineEventRequest) (*models.TimelineEvent, error) {
//
//}

func (s *server) DeleteTimelineEvent(ctx context.Context, t *models.Filter) (*models.Error, error) {
	if err := s.db.deleteTimelineEvent(t.GetId()); err != nil {
		return &models.Error{Response: false}, err
	}

	return &models.Error{Response: true}, nil
}
