package infrastructure

import (
	"protocall/infrastructure/memory"
	"protocall/internal/config"

	"github.com/spf13/viper"
)

type Repositories struct {
	*memory.Bridge
	*memory.AsteriskAccountMemory
	*memory.UserMemory
	*memory.ConferenceMemory
	*Uploader
	*Translator
	*Sender
}

func New() *Repositories {
	return &Repositories{
		memory.NewBridge(),
		memory.NewAsteriskAccount(),
		memory.NewUser(),
		memory.NewConference(),
		NewUploader(&UploaderConfig{
			Host:    viper.GetString(config.UploaderHost),
			Port:    viper.GetString(config.UploaderPort),
			Timeout: viper.GetInt(config.UploaderTimeout),
		}),
		NewTranslator(&TranslatorConfig{
			Host:    viper.GetString(config.TranslatorHost),
			Port:    viper.GetString(config.UploaderPort),
			Timeout: viper.GetInt(config.UploaderTimeout),
		}),
		NewSender(),
	}
}
