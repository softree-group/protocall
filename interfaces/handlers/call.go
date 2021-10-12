package handlers

import (
	"github.com/hashicorp/go-uuid"
	"github.com/valyala/fasthttp"
	"protocall/application"
)

//CreateCookie Method that return a cookie valorized as input (GoLog-Token as key)
func createCookie() *fasthttp.Cookie {
	token, _ := uuid.GenerateUUID()

	authCookie := fasthttp.Cookie{}
	authCookie.SetKey("sessionID")
	authCookie.SetValue(token)
	authCookie.SetMaxAge(120)
	authCookie.SetHTTPOnly(true)
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	return &authCookie
}

func start(ctx *fasthttp.RequestCtx, apps *application.Applications) {


	cookie := createCookie()
	ctx.Response.Header.SetCookie(cookie)
}

func join(ctx *fasthttp.RequestCtx, apps *application.Applications) {

	cookie := createCookie()
	ctx.Response.Header.SetCookie(cookie)
}

func leave(ctx *fasthttp.RequestCtx, apps *application.Applications) {

}