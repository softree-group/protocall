package config

import "github.com/spf13/viper"

func defaults() {
	ariDefaults()
	logDefaults()
	serverDefaults()
	translatorDefaults()
	uploaderDefaults()
}

func ariDefaults() {
	viper.SetDefault(ARIApplication, "")
	viper.SetDefault(ARISnoopyApplication, "")
	viper.SetDefault(ARIUser, "")
	viper.SetDefault(ARIPassword, "")
	viper.SetDefault(ARIUrl, "")
	viper.SetDefault(ARIWebsocketUrl, "")
	viper.SetDefault(ARIAccountsFile, "")
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

func uploaderDefaults() {
	viper.SetDefault(UploaderTimeout, 3)
}

func translatorDefaults() {
	viper.SetDefault(TranslatorTimeout, 3)
}
