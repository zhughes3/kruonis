package main

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/zhughes3/kruonis/cmd/timelines/models"
)

func TestTimelineGroups(t *testing.T) {
	ctx := context.Background()
	// create a timeline without specifying a group_id
	timeline := createTimeline(ctx, &models.Timeline{
		Title: "TestTimelineGroupsTimelineTitle",
		Tags:  []string{"tag1", "tag2"},
	})

	Convey("Create timeline without group automatically creates a timeline group", t, func() {
		So(timeline.GroupId, ShouldNotBeNil)
	})

	Convey("Test read timeline group", t, func() {
		group, err := testServer.ReadTimelineGroup(ctx, &models.Filter{
			Id: timeline.GroupId,
		})
		So(err, ShouldBeNil)
		So(group, ShouldNotBeNil)
		So(group.Id, ShouldEqual, timeline.GroupId)
		So(group.Private, ShouldBeFalse)
		So(group.UpdatedAt, ShouldNotBeNil)
		So(group.CreatedAt, ShouldNotBeNil)
		So(group.Uuid, ShouldNotBeNil)
		// creating a timeline without a user_id should give default value of user_id, or 0.
		So(group.UserId, ShouldEqual, 0)
	})

	Convey("Add a timeline to an existing group", t, func() {
		timelineWithGroup := createTimeline(ctx, &models.Timeline{
			GroupId: timeline.GroupId,
			Title:   "TestTimelineGroupsTimelineWithGroupIDTitle",
			Tags:    []string{"tag1id", "tag2id"},
		})
		So(timelineWithGroup, ShouldNotBeNil)
		So(timelineWithGroup.GroupId, ShouldNotBeNil)
		So(timelineWithGroup.GroupId, ShouldEqual, timeline.GroupId)
	})

	Convey("Update a timeline group", t, func() {
		updateGroupJSON := &models.UpdateTimelineGroupRequest{
			Id:      timeline.GroupId,
			Title:   "UpdatedTimelineGroup",
			Private: true,
		}
		group, err := testServer.UpdateTimelineGroup(ctx, updateGroupJSON)
		So(err, ShouldBeNil)
		So(group, ShouldNotBeNil)
		So(group.Private, ShouldBeTrue)
		So(group.Title, ShouldEqual, updateGroupJSON.Title)
		So(group.Id, ShouldEqual, timeline.GroupId)

		Convey("Read timeline group shows updates", func() {
			g, err := testServer.ReadTimelineGroup(ctx, &models.Filter{
				Id: group.Id,
			})
			So(err, ShouldBeNil)
			So(g, ShouldNotBeNil)
			So(g.Id, ShouldEqual, group.Id)
			So(g.Title, ShouldEqual, updateGroupJSON.Title)
			So(g.Private, ShouldBeTrue)
			So(g.Uuid, ShouldEqual, group.Uuid)
			So(g.UserId, ShouldEqual, group.UserId)
			So(g.CreatedAt, ShouldResemble, group.CreatedAt)
			So(g.UpdatedAt, ShouldNotBeNil)
		})
	})

	Convey("Delete a timeline group", t, func() {
		err := testServer.db.deleteTimelineGroup(timeline.GroupId)
		So(err, ShouldBeNil)
	})
}

func TestTimelines(t *testing.T) {
	ctx := context.Background()
	createTimelineJSON := &models.Timeline{
		Title: "TestTimelinesTimelineTitle",
		Tags:  []string{"tag1timeline", "tag2timeline"},
	}
	timeline, err := testServer.CreateTimeline(ctx, createTimelineJSON)

	Convey("Create timeline", t, func() {
		So(err, ShouldBeNil)
		So(timeline, ShouldNotBeNil)
		So(timeline.Title, ShouldEqual, createTimelineJSON.Title)
		// these tags might not come back in the order specified
		So(timeline.Tags, ShouldHaveLength, 2)
		So(timeline.Tags[0], ShouldEqual, createTimelineJSON.Tags[0])
		So(timeline.Tags[1], ShouldEqual, createTimelineJSON.Tags[1])
		So(timeline.GroupId, ShouldNotBeNil)
		So(timeline.Id, ShouldNotBeNil)
		So(timeline.Events, ShouldHaveLength, 0)
		So(timeline.CreatedAt, ShouldNotBeNil)
		So(timeline.UpdatedAt, ShouldNotBeNil)
	})

	Convey("Test ReadTimeline", t, func() {
		readTimeline, err := testServer.ReadTimeline(ctx, &models.Filter{
			Id: timeline.Id,
		})
		So(err, ShouldBeNil)
		So(timeline, ShouldNotBeNil)
		So(readTimeline.Title, ShouldEqual, timeline.Title)
		// these tags might not come back in the order specified
		So(readTimeline.Tags, ShouldHaveLength, 2)
		So(readTimeline.Tags[0], ShouldEqual, timeline.Tags[0])
		So(readTimeline.Tags[1], ShouldEqual, timeline.Tags[1])
		So(readTimeline.GroupId, ShouldEqual, timeline.GroupId)
		So(readTimeline.Id, ShouldEqual, timeline.Id)
		So(readTimeline.Events, ShouldHaveLength, 0)
		So(readTimeline.CreatedAt, ShouldResemble, timeline.CreatedAt)
		So(readTimeline.UpdatedAt, ShouldResemble, timeline.UpdatedAt)
	})

	Convey("Test updateTimeline", t, func() {
		updateTimelineJSON := &models.UpdateTimelineRequest{
			Id:    timeline.Id,
			Title: "Updated Timeline Title",
			Tags:  []string{"tag1timelineUPDATED", "tag2timelineUPDATED"},
		}
		updatedTimeline, err := testServer.UpdateTimeline(ctx, updateTimelineJSON)
		So(err, ShouldBeNil)
		So(updatedTimeline, ShouldNotBeNil)
		So(updatedTimeline.Id, ShouldEqual, timeline.Id)
		So(updatedTimeline.GroupId, ShouldEqual, timeline.GroupId)
		So(updatedTimeline.Tags, ShouldNotBeNil)
		So(updatedTimeline.Tags, ShouldHaveLength, 2)
		So(updatedTimeline.Tags[0], ShouldEqual, updateTimelineJSON.Tags[0])
		So(updatedTimeline.Tags[1], ShouldEqual, updateTimelineJSON.Tags[1])
		So(updatedTimeline.Title, ShouldEqual, updateTimelineJSON.Title)
		So(updatedTimeline.Events, ShouldHaveLength, 0)
		So(updatedTimeline.CreatedAt, ShouldResemble, timeline.CreatedAt)
		So(updatedTimeline.UpdatedAt, ShouldNotBeNil)
	})

	Convey("Test read timeline events", t, func() {
		createTimelineEvent(ctx, timeline.Id, "ReadTimelineEvents Event")
		events, err := testServer.ReadTimelineEvents(ctx, &models.Filter{Id: timeline.Id})
		So(err, ShouldBeNil)
		So(events, ShouldNotBeNil)
		So(events.Id, ShouldEqual, timeline.Id)
		So(events.Events, ShouldHaveLength, 1)
		So(events.Events[0].Id, ShouldNotBeNil)
	})

	Convey("Read timeline with events", t, func() {
		readTimeline, err := testServer.ReadTimeline(ctx, &models.Filter{
			Id: timeline.Id,
		})
		So(err, ShouldBeNil)
		So(timeline, ShouldNotBeNil)
		So(readTimeline.Events, ShouldNotBeNil)
		So(readTimeline.Events, ShouldHaveLength, 1)
	})

	Convey("Test delete timeline", t, func() {
		ret, err := testServer.DeleteTimeline(ctx, &models.Filter{Id: timeline.Id})
		So(err, ShouldBeNil)
		So(ret, ShouldNotBeNil)
		So(ret.Response, ShouldBeTrue)
	})
}

func TestTimelineEvents(t *testing.T) {
	ctx := context.Background()
	timeline := createTimeline(ctx, &models.Timeline{
		Title: "TestTimelineEventsTimelineTitle",
		Tags:  []string{"tag1", "tag2"},
	})

	Convey("Create a new event", t, func() {
		timestamp, _ := convertTime(time.Now())
		createEventJSON := &models.TimelineEvent{
			Id:          timeline.Id,
			Title:       "Timeline Event Title",
			Timestamp:   timestamp,
			Description: "A description of a timeline event.",
			Content:     "The actual content of a timeline event. This will usually be paragraph(s).",
		}
		event, err := testServer.CreateTimelineEvent(ctx, createEventJSON)
		So(err, ShouldBeNil)
		So(event, ShouldNotBeNil)
		So(event.Id, ShouldEqual, createEventJSON.Id)
		So(event.Title, ShouldEqual, createEventJSON.Title)
		So(event.Timestamp, ShouldEqual, createEventJSON.Timestamp)
		So(event.Description, ShouldEqual, createEventJSON.Description)
		So(event.Content, ShouldEqual, createEventJSON.Content)
		So(event.CreatedAt, ShouldNotBeNil)
		So(event.UpdatedAt, ShouldNotBeNil)
		So(event.ImageUrl, ShouldHaveLength, 0)
		So(event.EventId, ShouldNotBeNil)
		defer testServer.db.deleteTimelineEvent(event.EventId)

		Convey("Read newly created event", func() {
			evt, err := testServer.ReadTimelineEvent(ctx, &models.Filter{
				Id: event.EventId,
			})
			So(err, ShouldBeNil)
			So(evt, ShouldNotBeNil)
			So(evt.EventId, ShouldEqual, event.EventId)
			So(evt.Id, ShouldEqual, event.Id)
			So(evt.Title, ShouldEqual, event.Title)
			So(evt.Timestamp, ShouldResemble, event.Timestamp)
			So(evt.Description, ShouldEqual, event.Description)
			So(evt.Content, ShouldEqual, event.Content)
			So(evt.CreatedAt, ShouldResemble, event.CreatedAt)
			So(evt.UpdatedAt, ShouldResemble, event.UpdatedAt)
			So(evt.ImageUrl, ShouldEqual, event.ImageUrl)
		})

		Convey("Update event", func() {
			updatedTimestamp, _ := convertTime(time.Now())
			updateTimelineEventJSON := &models.UpdateTimelineEventRequest{
				Id:          event.EventId,
				Title:       "Updated Event Title",
				Timestamp:   updatedTimestamp,
				Description: "Updated Description",
				Content:     "Updated content",
			}
			evt, err := testServer.UpdateTimelineEvent(ctx, updateTimelineEventJSON)
			So(err, ShouldBeNil)
			So(evt, ShouldNotBeNil)
			So(evt.EventId, ShouldEqual, event.EventId)
			So(evt.Id, ShouldNotBeNil)
			So(evt.Title, ShouldEqual, updateTimelineEventJSON.Title)
			So(evt.Timestamp, ShouldResemble, updateTimelineEventJSON.Timestamp)
			So(evt.Description, ShouldEqual, updateTimelineEventJSON.Description)
			So(evt.Content, ShouldEqual, updateTimelineEventJSON.Content)
			So(evt.CreatedAt, ShouldResemble, event.CreatedAt)
			So(evt.UpdatedAt, ShouldNotBeNil)
			So(evt.ImageUrl, ShouldEqual, event.ImageUrl)
		})

		Convey("Delete timeline event", func() {
			err := testServer.db.deleteTimelineEvent(event.EventId)
			So(err, ShouldBeNil)
		})
	})
}

// helper functions
func createTimeline(ctx context.Context, in *models.Timeline) *models.Timeline {
	timeline, _ := testServer.CreateTimeline(ctx, in)
	return timeline
}

func createTimelineEvent(ctx context.Context, timelineID uint64, eventName string) *models.TimelineEvent {
	timestamp, _ := convertTime(time.Now())
	eventJSON := &models.TimelineEvent{
		Id:          timelineID,
		Title:       eventName,
		Timestamp:   timestamp,
		Description: "random description",
		Content:     "random content",
	}
	event, _ := testServer.CreateTimelineEvent(ctx, eventJSON)
	return event
}
