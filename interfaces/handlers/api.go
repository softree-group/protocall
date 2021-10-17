package handlers

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/mark-by/logutils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"protocall/application"
	"protocall/config"
)

func ServeAPI(apps *application.Applications) {
	r := router.New()

	compose := func(method func(string, fasthttp.RequestHandler), path string, handler func(ctx *fasthttp.RequestCtx, applications *application.Applications)) {
		method(path, func(ctx *fasthttp.RequestCtx) {
			handler(ctx, apps)
		})
	}

	r.GET("/logs", authRequired(logutils.GetLogs))
	r.POST("/logs/changeLevel", authRequired(logutils.ChangeLevel))
	r.POST("/logs/reset", authRequired(logutils.ResetLogs))

	compose(r.GET, "/session", session)
	compose(r.POST, "/conference/start", start)
	compose(r.POST, "/conference/{meetID}/join", join)
	compose(r.POST, "/conference/leave", leave)

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
