package config

import "github.com/spf13/viper"

func defaults() {
	ariDefaults()
	logDefaults()
	serverDefaults()
	porterDefaults()
	clerkDefaults()
}

func ariDefaults() {
	viper.SetDefault(ARIApplication, "")
	viper.SetDefault(ARISnoopyApplication, "")
	viper.SetDefault(ARIUser, "")
	viper.SetDefault(ARIPassword, "")
	viper.SetDefault(ARIUrl, "")
	viper.SetDefault(ARIWebsocketURL, "")
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
	viper.SetDefault(Participant, 32) //nolint:gomnd
}

func porterDefaults() {
	viper.SetDefault(PorterTimeout, 3) //nolint:gomnd
}

func clerkDefaults() {
	viper.SetDefault(ClerkTimeout, 3) //nolint:gomnd
}
