package handlers

import (
	"encoding/json"
	"github.com/hashicorp/go-uuid"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"protocall/application"
	"protocall/config"
	"time"
)

func start(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user, account := createSession(ctx, apps)
	if user == nil {
		return
	}

	conference, err := apps.Conference.StartConference(user)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	_, err = apps.Connector.CreateBridge(conference.ID)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	conference.BridgeID = conference.ID

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
	})
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.Response.SetBody(data)
}

func join(ctx *fasthttp.RequestCtx, apps *application.Applications) {

	meetID := ctx.UserValue("meetID").(string)
	if !apps.Conference.IsExist(meetID) {
		ctx.Error("Conference does not exist", 404)
		return
	}

	user, account := createSession(ctx, apps)
	if user == nil {
		return
	}

	conference, err := apps.Conference.JoinToConference(user, meetID)
	if err != nil {
		ctx.Error(err.Error(), 400)
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
	})
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.Response.SetBody(data)
}

func leave(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	sessionID := ctx.Request.Header.Cookie(sessionCookie)
	if len(sessionID) == 0 {
		ctx.SetStatusCode(400)
		return
	}

	apps.User.Delete(string(sessionID))
	ctx.Response.Header.DelCookie(sessionCookie)
}
