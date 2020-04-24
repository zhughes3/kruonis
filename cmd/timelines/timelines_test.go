package main

import (
	"context"
	"testing"

	"github.com/zhughes3/kruonis/cmd/timelines/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServer_CreateTimeline(t *testing.T) {
	createJSON := &models.Timeline{
		Title: "TestTitle",
		Tags:  []string{"yo", "tag"},
	}

	timeline, err := testServer.CreateTimeline(context.Background(), createJSON)
	So(err, ShouldBeNil)
	So(timeline, ShouldNotBeNil)
	So(timeline.Title, ShouldEqual, "TestTitle")
	So(timeline.Tags, ShouldHaveLength, 2)
	So(timeline.GroupId, ShouldNotBeNil)
	So(timeline.Id, ShouldNotBeNil)
	So(timeline.Events, ShouldHaveLength, 0)

}
