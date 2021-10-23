package handlers

import (
	"encoding/json"
	"net/http"
	"protocall/application"
	"protocall/domain/entity"

	"github.com/google/btree"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

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
	})
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

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

	channel, err := apps.Connector.CallAndConnect(user.AsteriskAccount, user.ConferenceID)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

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

	err := apps.Connector.Disconnect(user.ConferenceID, user.Channel)
	if err != nil {
		logrus.Error("Fail to disconnect: ", err)
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
			if err := apps.Connector.Disconnect(participant.ConferenceID, participant.Channel); err != nil {
				logrus.Error("Fail to disconnect: ", err)
			}
			apps.AsteriskAccount.Free(user.AsteriskAccount)
			apps.User.Delete(user.SessionID)
			// TODO: send socket event about end conference
			return true
		})

		apps.Conference.Delete(user.ConferenceID)
		return
	}

	// TODO: send socket event about leave participant
}

func record(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)

	err := apps.Conference.StartRecord(user, user.ConferenceID)
	if err != nil {
		ctx.SetStatusCode(http.StatusForbidden)
		logrus.Error("fail to start record: ", err)
		return
	}
}
