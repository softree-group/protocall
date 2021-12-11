package main

import (
	"protocall/internal/conference"
	"protocall/pkg/webcore"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func compose(
	method func(string, fasthttp.RequestHandler),
	path string,
	handler func(ctx *fasthttp.RequestCtx),
	middleware func(next fasthttp.RequestHandler) fasthttp.RequestHandler,
) {
	if middleware == nil {
		middleware = webcore.Fish
	}
	method(path, middleware(func(ctx *fasthttp.RequestCtx) { handler(ctx) }))
}

func NewRouter(conference *conference.API) *router.Router {
	r := router.New()
	compose(r.GET, "/session", conference.Session, nil)
	compose(r.POST, "/conference/start", conference.Start, nil)
	compose(r.POST, "/conference/{meetID}/join", conference.Join, nil)
	compose(r.POST, "/conference/record", conference.Record, nil)
	compose(r.POST, "/conference/leave", conference.Leave, nil)
	compose(r.POST, "/conference/ready", conference.Ready, nil)
	compose(r.GET, "/conference", conference.Info, nil)
	compose(r.POST, "/conference/translate", conference.TranslateDone, webcore.AuthRequired)

	return r
}

// func startServer(r *router.Router) {
// 	logrus.Infof("Запуск сервера на %s:%s ...", viper.Get(config.ServerIP), viper.Get(config.ServerPort))

// 	err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%s",
// 		viper.Get(config.ServerIP), viper.Get(config.ServerPort)),
// 		corsMiddleware().Handler(
// 			prefixMiddleware("/api/")(
// 				debugMiddleWare(r.Handler),
// 			),
// 		),
// 	)

// 	if err != nil {
// 		logrus.Fatalf("Сервер не запустился с ошибкой: %s", err)
// 	}
// }
