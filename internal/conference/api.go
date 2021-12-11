package conference

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"protocall/internal/account"
	"protocall/internal/operator"
	"protocall/internal/socket"
	"protocall/internal/translator"
	"protocall/internal/user"
	"protocall/pkg/bus"
	"protocall/pkg/logger"

	"github.com/golang-jwt/jwt"
	"github.com/hashicorp/go-uuid"
	"github.com/valyala/fasthttp"
)

type API struct {
	domain        string
	centrifugoKey []byte
	bus           bus.Client
	user          *user.Application
	account       *account.Application
	conference    *Application
	socket        *socket.Socket
	manager       *operator.Operator
}

func NewAPI(
	domain string,
	centrifugoKey []byte,
	bus bus.Client,
	user *user.Application,
	account *account.Application,
	conference *Application,
	socket *socket.Socket,
) *API {
	return &API{
		domain:        domain,
		centrifugoKey: centrifugoKey,
		bus:           bus,
		user:          user,
		account:       account,
		conference:    conference,
		socket:        socket,
	}
}

func (a *API) getUser(ctx *fasthttp.RequestCtx) *user.User {
	sessionID := ctx.Request.Header.Cookie(sessionCookie)
	if len(sessionID) == 0 {
		return nil
	}

	return a.user.Find(string(sessionID))
}

const (
	sessionCookie = "session_id"
	day           = 24 * time.Hour
)

func createCookie(domain string) *fasthttp.Cookie {
	token, _ := uuid.GenerateUUID()

	authCookie := fasthttp.AcquireCookie()
	authCookie.SetKey(sessionCookie)
	authCookie.SetValue(token)
	authCookie.SetDomain(domain)
	authCookie.SetPath("/")
	authCookie.SetExpire(time.Now().Add(day))
	authCookie.SetHTTPOnly(true)
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	authCookie.SetSecure(false)
	return authCookie
}

func deleteCookie(ctx *fasthttp.RequestCtx, domain string) {
	cookie := createCookie(domain)
	defer fasthttp.ReleaseCookie(cookie)

	cookie.SetExpire(time.Now().Add(-day * 24))
	ctx.Response.Header.SetCookie(cookie)
}

func createCentToken(id string, centrifugoKey []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
	})

	tokenString, err := token.SignedString(centrifugoKey)
	if err != nil {
		logger.L.Error("fail to generate cent token: ", err)
	}
	return tokenString
}

func (a *API) createSession(ctx *fasthttp.RequestCtx) (*user.User, *account.Account) {
	account := a.account.GetFree()
	if account == nil {
		ctx.Error("Sorry, we are busy ;(", fasthttp.StatusServiceUnavailable)
		return nil, nil
	}

	user := &user.User{}

	err := json.Unmarshal(ctx.PostBody(), user)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return nil, nil
	}

	cookie := createCookie(a.domain)
	defer fasthttp.ReleaseCookie(cookie)

	ctx.Response.Header.SetCookie(cookie)

	user.SessionID = string(cookie.Value())
	user.AsteriskAccount = account.Username
	a.account.Take(account.Username, user.SessionID)

	a.user.Save(user)
	return user, account
}

func (a *API) Start(ctx *fasthttp.RequestCtx) {
	user := a.getUser(ctx)
	if user != nil {
		ctx.Error("You are already signed in", fasthttp.StatusBadRequest)
		return
	}

	user, account := a.createSession(ctx)
	if user == nil {
		return
	}

	conference, err := a.conference.StartConference(user)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	conference.BridgeID = conference.ID

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
		"cent_token": createCentToken(user.AsteriskAccount, a.centrifugoKey),
	})
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.SetBody(data)
	ctx.Response.Header.SetContentType("application/json")
}

func (a *API) Join(ctx *fasthttp.RequestCtx) {
	meetID := ctx.UserValue("meetID").(string)
	if !a.conference.IsExist(meetID) {
		ctx.Error("Conference does not exist", fasthttp.StatusNotFound)
		return
	}

	user, account := a.createSession(ctx)
	if user == nil {
		return
	}

	conference, err := a.conference.JoinToConference(user, meetID)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
		"cent_token": createCentToken(user.AsteriskAccount, a.centrifugoKey),
	})
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	a.socket.PublishConnectionEvent(user)

	ctx.Response.SetBody(data)
	ctx.Response.Header.SetContentType("application/json")
}

func (a *API) Ready(ctx *fasthttp.RequestCtx) {
	user := a.getUser(ctx)
	if user == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	_ = a.socket.PublishConnectedEvent(user)
}

func (a *API) Leave(ctx *fasthttp.RequestCtx) {
	defer deleteCookie(ctx, a.domain)

	user := a.getUser(ctx)
	if user == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	a.socket.PublishLeaveEvent(user)

	if user.Channel != nil {
		if err := a.manager.KickUser(context.Background(), user); err != nil {
			logger.L.Error("Fail to disconnect: ", err)
		}
	}

	a.account.Free(user.AsteriskAccount)

	a.bus.Publish("leave", Event{
		ConferenceID: user.ConferenceID,
		User:         user,
	})
}

func (a *API) Record(ctx *fasthttp.RequestCtx) {
	user := a.getUser(ctx)
	if user == nil {
		ctx.Error("no user", 400)
		return
	}

	err := a.conference.StartRecord(user, user.ConferenceID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		logger.L.Error("fail to start record: ", err)
		return
	}

	a.socket.PublishStartRecordEvent(user.ConferenceID)
}

func (a *API) Info(ctx *fasthttp.RequestCtx) {
	user := a.getUser(ctx)
	if user == nil {
		ctx.Error("no user", 400)
		return
	}

	conference := a.conference.Get(user.ConferenceID)
	if conference == nil {
		ctx.Error("no conference", 400)
		a.account.Free(user.AsteriskAccount)
		a.user.Delete(user.SessionID)
		ctx.Response.Header.DelCookie(sessionCookie)
		return
	}

	confInfo, err := a.conference.GetConferenceInfo(user.ConferenceID)
	if err != nil {
		ctx.Error("", http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(confInfo)
	ctx.SetBody(body)
	ctx.SetContentType("application/json")
}

func (a *API) TranslateDone(ctx *fasthttp.RequestCtx) {
	data := translator.ConnectorRequest{}
	if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	user := a.user.Find(data.SessionID)
	if user == nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}

	a.bus.Publish("translated", Event{
		ConferenceID: user.ConferenceID,
		User:         user,
		Text:         data.Text,
		Record: &translator.Record{
			Path: data.Record.Path,
			URI:  data.Record.URI,
		},
	})
	ctx.Response.SetStatusCode(http.StatusNoContent)
}

func (a *API) Session(ctx *fasthttp.RequestCtx) {
	sessionID := ctx.Request.Header.Cookie(sessionCookie)
	if len(sessionID) == 0 {
		ctx.SetStatusCode(http.StatusNoContent)
		return
	}

	user := a.user.Find(string(sessionID))
	if user == nil {
		deleteCookie(ctx, a.domain)
		ctx.SetStatusCode(http.StatusNoContent)
		return
	}

	account := a.account.Get(user.AsteriskAccount)
	if account == nil {
		account = a.account.GetFree()
		if account == nil {
			deleteCookie(ctx, a.domain)
			ctx.Error("Sorry, we are busy ;(", fasthttp.StatusServiceUnavailable)
			// TODO: wait free account
			return
		}
	}

	conference := a.conference.Get(user.ConferenceID)
	if conference == nil {
		deleteCookie(ctx, a.domain)
		ctx.SetStatusCode(http.StatusNoContent)
		return
	}

	data, _ := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
		"cent_token": createCentToken(account.Username, a.centrifugoKey),
	})

	ctx.Response.SetBody(data)
}
