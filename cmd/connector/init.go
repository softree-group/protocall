package main

import (
	"protocall/internal/connector/application/app"
	"protocall/internal/connector/application/applications"
	"protocall/internal/connector/centrifugo"
	"protocall/internal/connector/domain/services"
	"protocall/pkg/bus"
)

// type application struct {
// 	Listener                 applications.EventListener
// 	Snoopy                   applications.Snoopy
// 	User                     applications.User
// 	AsteriskAccount          applications.AsteriskAccount
// 	Conference               applications.Conference
// 	Connector                applications.Connector
// 	AMI                      *app.AMIAsterisk
// 	Socket                   applications.Socket
// 	Bus                      services.Bus
// 	ApplicationEventListener applications.EventListener
// }

func app() error {
	// logger.L.Info("start connect")
	// ariClient, err := native.Connect(&native.Options{
	// 	Application:  viper.GetString(config.ARIApplication),
	// 	URL:          viper.GetString(config.ARIUrl),
	// 	WebsocketURL: viper.GetString(config.ARIWebsocketURL),
	// 	Username:     viper.GetString(config.ARIUser),
	// 	Password:     viper.GetString(config.ARIPassword),
	// })
	// logger.L.Info("end connect")
	// if err != nil {
	// 	logger.L.Fatal("cannot connect to asterisk: ", err)
	// }

	// ariClient, err := native.Connect(&native.Options{
	// 	Application:  viper.GetString(config.ARISnoopyApplication),
	// 	URL:          viper.GetString(config.ARIUrl),
	// 	WebsocketURL: viper.GetString(config.ARIWebsocketURL),
	// 	Username:     viper.GetString(config.ARIUser),
	// 	Password:     viper.GetString(config.ARIPassword),
	// })
	// if err != nil {
	// 	logrus.Fatal("Fail to connect snoopy app")
	// }

	socket := app.NewSocket(centrifugo.NewCentrifugo())

	busService := bus.New()

	// userApp := app.NewUser(reps)
	// conferenceApp := app.NewConference(reps, ariClient, busService)
	// asteriskApp := app.NewAsteriskAccount(reps, viper.GetString(config.ARIAccountsFile))
	// ami, err := app.NewAMIAsterisk(context.Background(), viper.GetString(config.AMIHost)+":"+viper.GetString(config.AMIPort), viper.GetString(config.AMIUser), viper.GetString(config.AMIPassword))
	// if err != nil {
	// 	logger.L.Fatal("Fail connect to ami: ", err)
	// }

	return &Applications{
		Listener: app.NewListener(
			reps,
			ariClient,
			app.NewHandler(ariClient, reps, connector),
			userApp,
			conferenceApp,
			asteriskApp,
			socketApp,
			busService,
		),
		Snoopy:          app.NewSnoopy(busService),
		Conference:      conferenceApp,
		AsteriskAccount: asteriskApp,
		User:            userApp,
		Connector:       connector,
		AMI:             ami,
		Socket:          socketApp,
		Bus:             busService,
		ApplicationEventListener: app.NewApplicationEventListener(
			reps,
			busService,
			conferenceApp,
			ami,
			socketApp,
		),
	}

	go a.Snoopy.Snoop()
	go a.ApplicationEventListener.Listen()
	go a.Listener.Listen()

	// NewRouter(&api{
	// 	session.,
	// })
	// fasthttp.ListenAndServe(fmt.Sprintf("%s:%s",
	// 	viper.Get(config.ServerIP), viper.Get(config.ServerPort)),
	// 	corsMiddleware().Handler(
	// 		prefixMiddleware("/api/")(
	// 			debugMiddleWare(r.Handler),
	// 		),
	// 	),
	// )
	handlers.ServeAPI(app)
}

