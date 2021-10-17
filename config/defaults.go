package config

import "github.com/spf13/viper"

func defaults() {
	ariDefaults()
	logDefaults()
	serverDefaults()
}

func ariDefaults() {
	viper.SetDefault(ARIApplication, "")
	viper.SetDefault(ARISnoopyApplication, "")
	viper.SetDefault(ARIUser, "")
	viper.SetDefault(ARIPassword, "")
	viper.SetDefault(ARIUrl, "")
	viper.SetDefault(ARIWebsocketUrl, "")
}

func logDefaults() {
	viper.SetDefault(LogFile, "-")
	viper.SetDefault(LogLevel, "debug")
}

func serverDefaults() {
	viper.SetDefault(ServerAPIKey, "")
	viper.SetDefault(ServerIP, "0.0.0.0")
	viper.SetDefault(ServerPort, "9595")
	viper.SetDefault(ServerDomain, "")
}
