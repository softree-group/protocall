package handlers

import (
	"protocall/internal/connector/config"
	"protocall/internal/connector/memory"
	"protocall/pkg/clerk"
	"protocall/pkg/porter"

	"github.com/spf13/viper"
)

type Handler struct {
	*memory.Bridge
	*memory.AsteriskAccountMemory
	*memory.UserMemory
	*memory.ConferenceMemory
	*porter.PorterClient
	*clerk.ClerkClient
}

func NewHandler() *Handler {
	return &Handler{
		memory.NewBridge(),
		memory.NewAsteriskAccount(),
		memory.NewUser(),
		memory.NewConference(),
		porter.NewPorterClient(&porter.PorterClientConfig{
			Host:    viper.GetString(config.PorterHost),
			Port:    viper.GetString(config.PorterPort),
			Timeout: viper.GetInt(config.PorterTimeout),
		}),
		clerk.NewClerkClient(&clerk.ClerkClientConfig{
			Host:    viper.GetString(config.ClerkHost),
			Port:    viper.GetString(config.ClerkPort),
			Timeout: viper.GetInt(config.ClerkTimeout),
		}),
	}
}
