package config

import "github.com/spf13/viper"

func defaults() {
	ariDefaults()
	logDefaults()
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
