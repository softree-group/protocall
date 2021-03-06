package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"protocall/internal/connector/application"
	"protocall/internal/connector/domain/entity"
	"protocall/internal/translator"

	"github.com/google/btree"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func start(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user != nil {
		ctx.Error("You are already signed in", fasthttp.StatusBadRequest)
		return
	}

	user, account := createSession(ctx, apps)
	if user == nil {
		return
	}

	requestData := map[string]interface{}{}
	json.Unmarshal(ctx.Request.Body(), &requestData)
	title, ok := requestData["title"]
	if !ok {
		ctx.Error("No title specified", fasthttp.StatusBadRequest)
		return
	}

	conference, err := apps.Conference.StartConference(user, title.(string))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	conference.BridgeID = conference.ID

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
		"cent_token": createCentToken(user.AsteriskAccount),
	})
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.SetBody(data)
	ctx.Response.Header.SetContentType("application/json")
}

func join(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	meetID := ctx.UserValue("meetID").(string)
	if !apps.Conference.IsExist(meetID) {
		ctx.Error("Conference does not exist", fasthttp.StatusNotFound)
		return
	}

	user, account := createSession(ctx, apps)
	if user == nil {
		return
	}

	conference, err := apps.Conference.JoinToConference(user, meetID)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
		"cent_token": createCentToken(user.AsteriskAccount),
	})
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	_ = apps.Socket.PublishConnectionEvent(user)

	ctx.Response.SetBody(data)
	ctx.Response.Header.SetContentType("application/json")
}

func getUser(ctx *fasthttp.RequestCtx, apps *application.Applications) *entity.User {
	sessionID := ctx.Request.Header.Cookie(sessionCookie)
	if len(sessionID) == 0 {
		return nil
	}

	return apps.User.Find(string(sessionID))
}

func ready(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	_ = apps.Socket.PublishConnectedEvent(user)
}

func leave(ctx *fasthttp.RequestCtx, apps *application.Applications) {

	defer deleteCookie(ctx)

	user := getUser(ctx, apps)
	if user == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	apps.Socket.PublishLeaveEvent(user)

	if user.Channel != nil {
		if err := apps.AMI.KickUser(context.Background(), user); err != nil {
			logrus.Error("Fail to disconnect: ", err)
		}
	}

	apps.AsteriskAccount.Free(user.AsteriskAccount)

	apps.Bus.Publish("leave", entity.EventDefault{
		ConferenceID: user.ConferenceID,
		User:         user,
	})
}

type mediaRequest struct {
	AudioMuted bool `json:"audioMuted"`
	VideoMuted bool `json:"videoMuted"`
}

func media(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	var request mediaRequest
	_ = json.Unmarshal(ctx.PostBody(), &request)

	user.VideoMuted = request.VideoMuted
	user.AudioMuted = request.AudioMuted

	apps.User.Save(user)
}

func record(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user == nil {
		ctx.Error("no user", 400)
		return
	}

	err := apps.Conference.StartRecord(user, user.ConferenceID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		logrus.Error("fail to start record: ", err)
		return
	}

	apps.Socket.PublishStartRecordEvent(user.ConferenceID)
}

type UserInfo struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Channel    string `json:"channel"`
	AudioMuted bool   `json:"audioMuted"`
	VideoMuted bool   `json:"videoMuted"`
}

type ConferenceInfo struct {
	ID           string     `json:"id"`
	HostID       string     `json:"host_id"`
	Participants []UserInfo `json:"participants"`
	IsRecording  bool       `json:"is_recording"`
	StartedAt    int64      `json:"started_at"`
}

func info(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user == nil {
		ctx.Error("no user", 400)
		return
	}

	conference := apps.Conference.Get(user.ConferenceID)
	if conference == nil {
		ctx.Error("no conference", 400)
		apps.AsteriskAccount.Free(user.AsteriskAccount)
		apps.User.Delete(user.SessionID)
		ctx.Response.Header.DelCookie(sessionCookie)
		return
	}

	conferenceInfo := ConferenceInfo{
		ID:           conference.ID,
		HostID:       conference.HostUserID,
		Participants: nil,
		IsRecording:  conference.IsRecording,
		StartedAt:    conference.Start.Unix(),
	}

	participants := make([]UserInfo, 0, conference.Participants.Len())

	conference.Participants.Ascend(func(item btree.Item) bool {
		if item == nil {
			return false
		}
		user := item.(*entity.User)
		if user == nil {
			return false
		}
		channel := ""
		if user.Channel != nil {
			channel = user.Channel.ID
		}
		participants = append(participants, UserInfo{
			Name:       user.Username,
			ID:         user.AsteriskAccount,
			Channel:    channel,
			AudioMuted: user.AudioMuted,
			VideoMuted: user.VideoMuted,
		})
		return true
	})

	conferenceInfo.Participants = participants

	body, _ := json.Marshal(conferenceInfo)
	ctx.SetBody(body)
	ctx.SetContentType("application/json")
}

func translate(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	data := translator.ConnectorRequest{}
	if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	user := apps.User.Find(data.SessionID)
	if user == nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}

	apps.Bus.Publish("translated", entity.EventDefault{
		ConferenceID: user.ConferenceID,
		User:         user,
		Text:         data.Text,
		Record: &entity.Record{
			Path: data.Record.Path,
			URI:  data.Record.URI,
		},
	})
	ctx.Response.SetStatusCode(http.StatusNoContent)
}
