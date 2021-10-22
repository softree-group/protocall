package config

import (
	"flag"
	"fmt"
	"github.com/mark-by/logutils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var FilePath string

func initFlags() {
	flag.StringVar(&FilePath, "s", "config.yaml", "config file")

	flag.Parse()
}

func InitConfig() {
	initFlags()

	viper.SetConfigFile(FilePath)

	defaults()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			logrus.Fatalf("Невалидный синтаксис")
		} else {
			_ = viper.WriteConfig()
			fmt.Printf("Файл не найден. По пути '%s' записан шаблон\n", "config.yaml")
		}
	}

	viper.AllKeys()
	checkConfig()
	logutils.InitLogrus(viper.GetString(LogFile), viper.GetString(LogLevel))
	viper.WatchConfig()
}
