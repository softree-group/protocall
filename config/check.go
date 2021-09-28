package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	errorFatal = "fatal"
	errorWarn  = "warn"
)

type confError struct {
	Level string
	Error string
}

func addError(errors *[]confError, level string, error string) {
	*errors = append(*errors, confError{
		Level: level,
		Error: error,
	})
}

func processErrors(errors []confError) {
	fatalErrorExist := false
	for _, err := range errors {
		switch err.Level {
		case errorFatal:
			fatalErrorExist = true
			logrus.Error(err.Error)
		case errorWarn:
			logrus.Warn(err.Error)
		}
	}
	if fatalErrorExist {
		logrus.Fatal("Есть ошибки в конфиге")
	}
}

func configErrorFormatter(key string, error string) string {
	return fmt.Sprintf("[%s] %s", key, error)
}

func equal(parameter string, value interface{}, level func(key, error string), message string) {
	if v := viper.Get(parameter); v == value {
		level(parameter, message)
	}
}

func checkConfig() {
	var errors []confError

	pusher := func(level string, error string) {
		addError(&errors, level, error)
	}
	fatal := func(key, error string) {
		pusher(errorFatal, configErrorFormatter(key, error))
	}
	warn := func(key, error string) {
		pusher(errorWarn, configErrorFormatter(key, error))
	}

	equal(LogFile, "-", warn, "Вывод логов в STDOUT")

	equal(ARIUrl, "", fatal, "Не указан адрес астериска")
	equal(ARIApplication, "", fatal, "Не указано название приложения астериска")
	equal(ARISnoopyApplication, "", fatal, "Не указано название приложения астериска для записи")
	equal(ARIWebsocketUrl, "", fatal, "Не указан адрес сокета астериска")
	equal(ARIUser, "", fatal, "Не указан пользователь астериска")
	equal(ARIPassword, "", fatal, "Не указан пароль пользователя астериска")

	processErrors(errors)
}
