package app

import (
	"errors"
	"fmt"

	"protocall/internal/config"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (c *Connector) createCallInternal(account string) (h *ari.ChannelHandle, err error) {
	return c.ari.Channel().Originate(nil, ari.OriginateRequest{
		Endpoint: fmt.Sprintf("PJSIP/%s", account),
		Timeout:  10,
		CallerID: "system",
		App:      viper.GetString(config.ARIApplication),
	})
}

func (c *Connector) waitUp(channel *ari.ChannelHandle) error {
	stateChange := channel.Subscribe(ari.Events.ChannelStateChange)
	destroyed := channel.Subscribe(ari.Events.ChannelDestroyed)

	for {
		select {
		case <-stateChange.Events():
			data, err := channel.Data()
			if err != nil {
				logrus.Error("error to get data from channel: ", err)
				continue
			}

			if data.State == "Up" {
				return nil
			}
		case <-destroyed.Events():
			return errors.New("channel destroyed")
		}
	}
}
