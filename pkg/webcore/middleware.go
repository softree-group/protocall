package webcore

import (
	"net/http"
	"os"
	"protocall/pkg/logger"
	"sort"
	"strings"

	"github.com/lab259/cors"
	"github.com/valyala/fasthttp"
)

func ApplyMethods(next http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		methods := sort.StringSlice(methods)
		sort.Sort(methods)
		i := sort.SearchStrings(methods, r.Method)
		if i < len(methods) && methods[i] == r.Method {
			next.ServeHTTP(rw, r)
			return
		}
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"https://protocall.softex-team.ru", "http://localhost:3000", "http://localhost.ru:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Range", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "Content-Range"},
		AllowCredentials: true,
	})
}

var token = os.Getenv("API_TOKEN")

func AuthRequired(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if token == "" {
			next(ctx)
		}

		key := string(ctx.Request.Header.Peek("Authorization"))
		if key == "" {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("Authorization key not specified")
			return
		}

		if key != token {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("Invalid authorization key")
			return
		}

		next(ctx)
	}
}

func Fish(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)
	}
}

func Prefix(prefix string) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			ctx.Request.SetRequestURI(strings.Replace(string(ctx.Request.RequestURI()), prefix, "/", 1))
			next(ctx)
		}
	}
}

func Debug(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logger.L.Debugf("%s %s %s", ctx.Method(), ctx.RequestURI(), ctx.PostBody())
		next(ctx)
	}
}
