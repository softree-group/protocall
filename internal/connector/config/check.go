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

func addError(errors *[]confError, level, err string) {
	*errors = append(*errors, confError{
		Level: level,
		Error: err,
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

func configErrorFormatter(key, err string) string {
	return fmt.Sprintf("[%s] %s", key, err)
}

func equal(parameter string, value interface{}, level func(key, err string), message string) {
	if v := viper.Get(parameter); v == value {
		level(parameter, message)
	}
}

func checkConfig() {
	var errors []confError

	pusher := func(level string, err string) {
		addError(&errors, level, err)
	}
	fatal := func(key, err string) {
		pusher(errorFatal, configErrorFormatter(key, err))
	}
	warn := func(key, err string) {
		pusher(errorWarn, configErrorFormatter(key, err))
	}

	equal(LogFile, "-", warn, "Вывод логов в STDOUT")

	equal(ARIUrl, "", fatal, "Не указан адрес астериска")
	equal(ARIApplication, "", fatal, "Не указано название приложения астериска")
	equal(ARISnoopyApplication, "", fatal, "Не указано название приложения астериска для записи")
	equal(ARIWebsocketURL, "", fatal, "Не указан адрес сокета астериска")
	equal(ARIUser, "", fatal, "Не указан пользователь астериска")
	equal(ARIPassword, "", fatal, "Не указан пароль пользователя астериска")
	equal(ARIAccountsFile, "", fatal, "Не указан файл с аккаунтами астериска")

	equal(ServerAPIKey, "", fatal, "Не указан api_key")

	processErrors(errors)
}
