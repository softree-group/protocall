package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"protocall/internal/connector/application"
	"protocall/internal/connector/config"
	"protocall/internal/connector/domain/entity"

	"github.com/golang-jwt/jwt"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

const (
	sessionCookie = "session_id"
	day           = 24 * time.Hour
)

func createCookie() *fasthttp.Cookie {
	token, _ := uuid.GenerateUUID()

	authCookie := fasthttp.Cookie{}
	authCookie.SetKey(sessionCookie)
	authCookie.SetValue(token)
	authCookie.SetDomain(viper.GetString(config.ServerDomain))
	authCookie.SetPath("/")
	authCookie.SetExpire(time.Now().Add(day))
	authCookie.SetHTTPOnly(true)
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	authCookie.SetSecure(false)
	return &authCookie
}

func createCentToken(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
	})

	tokenString, err := token.SignedString([]byte(viper.GetString(config.CentrifugoToken)))
	if err != nil {
		logrus.Error("fail to generate cent token: ", err)
	}
	return tokenString
}

func session(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	sessionID := ctx.Request.Header.Cookie(sessionCookie)
	if len(sessionID) == 0 {
		ctx.SetStatusCode(http.StatusNoContent)
		return
	}

	user := apps.User.Find(string(sessionID))
	if user == nil {
		ctx.Response.Header.DelCookie(sessionCookie)
		ctx.SetStatusCode(http.StatusNoContent)
		return
	}

	account := apps.AsteriskAccount.Get(user.AsteriskAccount)
	if account == nil {
		account = apps.AsteriskAccount.GetFree()
		if account == nil {
			ctx.Response.Header.DelCookie(sessionCookie)
			ctx.Error("Sorry, we are busy ;(", fasthttp.StatusServiceUnavailable)
			// TODO: wait free account
			return
		}
	}

	conference := apps.Conference.Get(user.ConferenceID)
	if conference == nil {
		ctx.Response.Header.DelCookie(sessionCookie)
		ctx.SetStatusCode(http.StatusNoContent)
		return
	}

	data, _ := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
		"cent_token": createCentToken(account.Username),
	})

	ctx.Response.SetBody(data)
}

func createSession(ctx *fasthttp.RequestCtx, apps *application.Applications) (*entity.User, *entity.AsteriskAccount) {
	account := apps.AsteriskAccount.GetFree()
	if account == nil {
		ctx.Error("Sorry, we are busy ;(", fasthttp.StatusServiceUnavailable)
		return nil, nil
	}

	user := &entity.User{}

	err := json.Unmarshal(ctx.PostBody(), user)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return nil, nil
	}

	cookie := createCookie()
	ctx.Response.Header.SetCookie(cookie)

	user.SessionID = string(cookie.Value())
	user.AsteriskAccount = account.Username
	apps.AsteriskAccount.Take(account.Username, user.SessionID)

	apps.User.Save(user)
	return user, account
}
