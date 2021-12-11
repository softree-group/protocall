package operator

import (
	"context"
	"errors"
	"fmt"

	"protocall/internal/user"
	"protocall/pkg/logger"

	"github.com/heltonmarx/goami/ami"
)

type Operator struct {
	host     string
	port     string
	username string
	secret   string
	uuid     string
	socket   *ami.Socket
	events   chan ami.Response
}

func NewOperator(cfg *Config) *Operator {
	return &Operator{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.User,
		secret:   cfg.Password,
		events:   make(chan ami.Response),
	}
}

// Events returns an channel with events received from AMI.
func (as *Operator) Events() <-chan ami.Response {
	return as.events
}

// SIPPeers fetch the list of SIP peers present on asterisk.
func (o *Operator) SIPPeers(ctx context.Context) ([]ami.Response, error) {
	var peers []ami.Response
	resp, err := ami.SIPPeers(ctx, o.socket, o.uuid)
	switch {
	case err != nil:
		return nil, err
	case len(resp) == 0:
		return nil, errors.New("there's no sip peers configured")
	default:
		for _, v := range resp {
			peer, err := ami.SIPShowPeer(ctx, o.socket, o.uuid, v.Get("ObjectName"))
			if err != nil {
				return nil, err
			}
			peers = append(peers, peer)
		}
	}
	return peers, nil
}

func (o *Operator) KickUser(ctx context.Context, client *user.User) error {
	if client.Channel == nil {
		return errors.New("channel is null")
	}
	_, err := ami.ConfbridgeKick(ctx, o.socket, o.uuid, client.ConferenceID, client.Channel.ID)
	if err != nil {
		return err
	}
	return nil
}

func (o *Operator) KickAllFromConference(ctx context.Context, conferenceID string) error {
	_, err := ami.ConfbridgeKick(ctx, o.socket, o.uuid, conferenceID, "all")
	return err
}

const events = "system,call,all,user"

func (o *Operator) Run(ctx context.Context) {
	socket, err := ami.NewSocket(ctx, fmt.Sprintf("%s:%s", o.host, o.port))
	if err != nil {
		logger.L.Error(err)
		return
	}
	o.socket = socket

	uuid, err := ami.GetUUID()
	if err != nil {
		logger.L.Error(err)
		return
	}
	o.uuid = uuid

	err = ami.Login(ctx, socket, o.username, o.secret, events, uuid)
	if err != nil {
		logger.L.Error(err)
		return
	}
	defer func() {
		if err := ami.Logoff(ctx, o.socket, o.uuid); err != nil {
			logger.L.Error(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			events, err := ami.Events(ctx, o.socket)

			if err != nil {
				logger.L.Error("AMI events failed: %v\n", err)
				return
			}
			o.events <- events
		}
	}
}
