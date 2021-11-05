package handlers

import (
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"protocall/application"
	"protocall/domain/entity"
	"protocall/internal/config"

	"github.com/google/btree"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Publish struct {
	Method string `json:"method"`
	Params struct {
		Channel string      `json:"channel"`
		Data    PublishData `json:"data"`
	}
}

type PublishData struct {
	Event string `json:"event"`
	User  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}

func publish(user *entity.User, event string) {
	logrus.Info("publish: ", user.ConferenceID, event)
	publishData := PublishData{
		Event: event,
		User: struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{ID: user.AsteriskAccount, Name: user.Username},
	}

	publishMethod := Publish{
		Method: "publish",
		Params: struct {
			Channel string      `json:"channel"`
			Data    PublishData `json:"data"`
		}{Channel: "conference~" + user.ConferenceID, Data: publishData},
	}

	data, err := json.Marshal(publishMethod)
	if err != nil {
		logrus.Error("fail to marshal: ", err)
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetBody(data)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", "apikey "+viper.GetString(config.CentrifugoAPIKey))
	req.Header.SetMethod("POST")
	req.SetRequestURI(viper.GetString(config.CentrifugoHost) + "/api")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = fasthttp.Do(req, resp)

	logrus.Info(resp.StatusCode())
	if err != nil {
		logrus.Error("fail to publish: ", err)
	}
}

func start(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user != nil {
		ctx.Error("You are already signed in", http.StatusBadRequest)
		return
	}

	user, account := createSession(ctx, apps)
	if user == nil {
		return
	}

	conference, err := apps.Conference.StartConference(user)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = apps.Connector.CreateBridge(conference.ID)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	conference.BridgeID = conference.ID

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
		"cent_token": createCentToken(user.SessionID),
	})
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Response.SetBody(data)
	ctx.Response.Header.SetContentType("application/json")
}

func join(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	meetID := ctx.UserValue("meetID").(string)
	if !apps.Conference.IsExist(meetID) {
		ctx.Error("Conference does not exist", http.StatusNotFound)
		return
	}

	user, account := createSession(ctx, apps)
	if user == nil {
		return
	}

	conference, err := apps.Conference.JoinToConference(user, meetID)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
		"cent_token": createCentToken(user.SessionID),
	})
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	publish(user, "connection")

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
		ctx.Error("no session", http.StatusBadRequest)
		return
	}

	logrus.Info("Calling ", user.AsteriskAccount)
	channel, err := apps.Connector.CallAndConnect(user)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	defer func() {
		if err != nil {
			publish(user, "fail_connection")
			return
		}
		publish(user, "connected")
	}()

	user.Channel = channel
	apps.User.Save(user)

	conference := apps.Conference.Get(user.ConferenceID)
	if conference == nil {
		logrus.Error("fail to get conference ", user.ConferenceID)
		return
	}

	if conference.IsRecording {
		err = apps.Conference.StartRecordUser(user, conference.ID)
		if err != nil {
			logrus.Error("fail to start record user: ", err)
			return
		}
	}
}

func leave(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user == nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	publish(user, "leave")

	if user.Channel != nil {
		err := apps.Connector.Disconnect(user.ConferenceID, user.Channel)
		if err != nil {
			logrus.Error("Fail to disconnect: ", err)
		}
	}

	apps.AsteriskAccount.Free(user.AsteriskAccount)

	defer apps.User.Delete(user.SessionID)
	defer ctx.Response.Header.DelCookie(sessionCookie)

	conference := apps.Conference.Get(user.ConferenceID)
	conference.Participants.Delete(user)

	if conference.HostUserID == user.AsteriskAccount {
		conference.Participants.Ascend(func(item btree.Item) bool {
			if item == nil {
				return false
			}
			participant := item.(*entity.User)
			if participant == nil {
				return false
			}
			logrus.Info("Disconnect: ", participant.Username)
			if err := apps.Connector.Disconnect(participant.ConferenceID, participant.Channel); err != nil {
				logrus.Error("Fail to disconnect: ", err)
			}
			apps.AsteriskAccount.Free(user.AsteriskAccount)
			apps.User.Delete(user.SessionID)
			publish(user, "end")
			return true
		})

		apps.Conference.Delete(user.ConferenceID)
	}

	if conference.IsRecording {
		if err := apps.Conference.UploadRecord(user, user.ConferenceID); err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

}

func record(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user == nil {
		ctx.Error("no user", 400)
		return
	}

	err := apps.Conference.StartRecord(user, user.ConferenceID)
	if err != nil {
		ctx.SetStatusCode(http.StatusForbidden)
		logrus.Error("fail to start record: ", err)
		return
	}

	publish(user, "start_record")
}

type UserInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
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
		participants = append(participants, UserInfo{
			Name: user.Username,
			ID:   user.AsteriskAccount,
		})
		return true
	})

	conferenceInfo.Participants = participants

	body, _ := json.Marshal(conferenceInfo)
	ctx.SetBody(body)
	ctx.SetContentType("application/json")
}
