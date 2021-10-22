package handlers

import (
	"protocall/internal/config"
	"strings"

	"github.com/lab259/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func corsMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"https://protocall.softex-team.ru", "http://localhost:3000", "http://localhost.ru:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Range", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "Content-Range"},
		AllowCredentials: true,
	})
}

func authRequired(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		key := string(ctx.Request.Header.Peek("Authorization"))
		if key == "" {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("Authorization key not specified")
			return
		}

		if key != viper.GetString(config.ServerAPIKey) {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("Invalid authorization key")
			return
		}

		next(ctx)
	}
}

func prefixMiddleware(prefix string) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			uri := string(ctx.Request.RequestURI())
			if strings.Contains(uri, prefix) {
				uri = strings.Replace(uri, prefix, "/", 1)
			}
			ctx.Request.SetRequestURI(uri)

			next(ctx)
		}
	}
}

func debugMiddleWare(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logrus.Debugf("%s %s %s", ctx.Method(), ctx.RequestURI(), ctx.PostBody())

		next(ctx)
	}
}
