package bridge

import (
	"errors"
	"fmt"
	"net/http"

	"protocall/internal/account"
	"protocall/internal/conference"
	"protocall/internal/socket"
	"protocall/internal/user"
	"protocall/pkg/bus"
	"protocall/pkg/logger"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/valyala/fasthttp"
)

type BridgeStorage interface {
	CreateBridge(hostUsername string, bridgeID string)
	GetForHost(hostUsername string) (string, error)
	DeleteBridge(bridgeID string) error
}

type Application struct {
	config      *Config
	ari         ari.Client
	bus         bus.Client
	bridgeStore BridgeStorage
	account     *account.Application
	user        *user.Application
	conference  *conference.Application
	socket      *socket.Socket
}

func NewApplication(cfg *Config, client ari.Client, bridgeStore BridgeStorage) *Application {
	return &Application{
		config:      cfg,
		ari:         client,
		bridgeStore: bridgeStore,
	}
}

func (a *Application) createCallInternal(account string) (h *ari.ChannelHandle, err error) {
	return a.ari.Channel().Originate(nil, ari.OriginateRequest{
		Endpoint: fmt.Sprintf("PJSIP/%s", account),
		Timeout:  10,
		CallerID: "system",
		App:      a.config.Application,
	})
}

func (a *Application) waitUp(channel *ari.ChannelHandle) error {
	stateChange := channel.Subscribe(ari.Events.ChannelStateChange)
	destroyed := channel.Subscribe(ari.Events.ChannelDestroyed)

	for {
		select {
		case <-stateChange.Events():
			data, err := channel.Data()
			if err != nil {
				logger.L.Error("error to get data from channel: ", err)
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

func (a *Application) CreateBridgeFrom(channel *ari.ChannelHandle) (*ari.BridgeHandle, error) {
	key := channel.Key().New(ari.BridgeKey, channel.ID())

	bridge, err := a.ari.Bridge().Create(key, "video_sfu", key.ID)

	if err != nil {
		return nil, err
	}

	a.bridgeStore.CreateBridge(channel.ID(), bridge.ID())

	return bridge, nil
}

func (a *Application) HasBridge() bool {
	bID, _ := a.bridgeStore.GetForHost("some")
	return bID != ""
}

func (a *Application) GetForHost(hostUsername string) (string, error) {
	return a.bridgeStore.GetForHost(hostUsername)
}

func (a *Application) GetBridge(id string) *ari.BridgeHandle {
	key := &ari.Key{
		Kind:                 ari.BridgeKey,
		ID:                   id,
		Node:                 "",
		Dialog:               "",
		App:                  a.config.Application,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	return a.ari.Bridge().Get(key)
}

func (a *Application) CreateBridge(id string) (*ari.BridgeHandle, error) {
	key := &ari.Key{
		Kind:                 ari.BridgeKey,
		ID:                   id,
		Node:                 "",
		Dialog:               "",
		App:                  a.config.Application,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	return a.ari.Bridge().Create(key, "video_sfu", key.ID)
}

func (a *Application) CallAndConnect(client *user.User) (*ari.Key, error) {
	bridgeID := client.ConferenceID
	account := client.AsteriskAccount
	if client.Channel != nil {
		ch := a.ari.Channel().Get(client.Channel)
		if ch != nil {
			err := ch.Hangup()
			if err != nil {
				logger.L.Error("fail to hangup: ", err)
			}
		}
	}
	bridge := a.GetBridge(client.ConferenceID)
	if bridge == nil {
		return nil, fmt.Errorf("bridge %s does not exist", bridgeID)
	}

	clientChannel, err := a.createCallInternal(account)
	if err != nil {
		return nil, err
	}

	err = a.waitUp(clientChannel)
	if err != nil {
		return nil, err
	}

	err = bridge.AddChannel(clientChannel.ID())
	if err != nil {
		return nil, err
	}

	return clientChannel.Key(), nil
}

func (a *Application) PostSnoop(id, snoopID, appArgs string) (*fasthttp.Response, error) {
	clientt := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI(fmt.Sprintf(
		"%s/channels/%s/snoop?api_key=%s:%s",
		a.config.URL,
		id,
		a.config.User,
		a.config.Password,
	))
	req.SetBodyString(fmt.Sprintf(`{"snoopID": "%s",
"app":  "%s",
"spy":  "%s",
"whisper":  "%s",
"appArgs":  "%s"}`, snoopID, a.config.SnoopyApplication, "in", "both", appArgs))
	req.Header.SetContentType("application/json")
	err := clientt.Do(req, resp)
	if err != nil {
		logger.L.Errorf("Сетевая ошибка по пути")
		return resp, err
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		logger.L.Warnf("Сервер ответил %d", resp.StatusCode())
	}
	return resp, err
}

func (a *Application) Connect(bridge *ari.BridgeHandle, channelID string) error {
	return bridge.AddChannel(channelID)
}

func (a *Application) Disconnect(bridgeID string, channel *ari.Key) error {
	err := a.ari.Channel().Get(channel).Hangup()
	if err != nil {
		logger.L.Error("fail to delete channel: ", err)
	}
	return err
}
