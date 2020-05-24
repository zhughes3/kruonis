package main

import (
	"context"

	"testing"

	"github.com/zhughes3/kruonis/cmd/timelines/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServer_Login(t *testing.T) {
	req := &models.SignupRequest{
		Email:    "somemeial@email.com",
		Password: "SomeGoodPassword",
	}
	Convey("Testing user signup and login", t, func() {
		ctx := context.Background()
		resp, err := testServer.Signup(ctx, req)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.Response, ShouldBeTrue)

		loginResp, err := testServer.Login(ctx, &models.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		So(err, ShouldBeNil)
		So(loginResp, ShouldNotBeNil)
	})

}
