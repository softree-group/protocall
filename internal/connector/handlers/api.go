package handlers

import (
	"fmt"
	"protocall/internal/connector/application"
	"protocall/internal/connector/config"

	"github.com/fasthttp/router"
	"github.com/mark-by/logutils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func ServeAPI(apps *application.Applications) {
	r := router.New()

	compose := func(
		method func(string, fasthttp.RequestHandler),
		path string,
		handler func(ctx *fasthttp.RequestCtx, applications *application.Applications),
		middleWare func (next fasthttp.RequestHandler) fasthttp.RequestHandler,
	) {
		if middleWare == nil {
			middleWare = fishMiddleware
		}

		method(path, middleWare(func(ctx *fasthttp.RequestCtx) { handler(ctx, apps) }))
	}

	r.GET("/logs", authRequired(logutils.GetLogs))
	r.POST("/logs/changeLevel", authRequired(logutils.ChangeLevel))
	r.POST("/logs/reset", authRequired(logutils.ResetLogs))

	compose(r.GET, "/session", session, nil)
	compose(r.POST, "/conference/start", start, nil)
	compose(r.POST, "/conference/{meetID}/join", join, nil)
	compose(r.POST, "/conference/record", record, nil)
	compose(r.POST, "/conference/leave", leave, nil)
	compose(r.POST, "/conference/ready", ready, nil)
	compose(r.GET, "/conference", info, nil)
	compose(r.POST, "/translates", translate, authRequired)

	startServer(r)
}

func startServer(r *router.Router) {
	logrus.Infof("Запуск сервера на %s:%s ...", viper.Get(config.ServerIP), viper.Get(config.ServerPort))

	err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%s",
		viper.Get(config.ServerIP), viper.Get(config.ServerPort)),
		corsMiddleware().Handler(
			prefixMiddleware("/api/")(
				debugMiddleWare(r.Handler),
			),
		),
	)

	if err != nil {
		logrus.Fatalf("Сервер не запустился с ошибкой: %s", err)
	}
}
