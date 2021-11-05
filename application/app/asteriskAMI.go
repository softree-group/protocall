package app

import (
	"context"
	"errors"
	"log"
	"protocall/domain/entity"
	"sync"

	"github.com/heltonmarx/goami/ami"
)

type AMIAsterisk struct {
	socket *ami.Socket
	uuid   string

	events chan ami.Response
	stop   chan struct{}
	wg     sync.WaitGroup
}

// NewAsterisk initializes the AMI socket with a login and capturing the events.
func NewAMIAsterisk(ctx context.Context, host string, username string, secret string) (*AMIAsterisk, error) {
	socket, err := ami.NewSocket(ctx, host)
	if err != nil {
		return nil, err
	}
	uuid, err := ami.GetUUID()
	if err != nil {
		return nil, err
	}
	const events = "system,call,all,user"
	err = ami.Login(ctx, socket, username, secret, events, uuid)
	if err != nil {
		return nil, err
	}
	as := &AMIAsterisk{
		socket: socket,
		uuid:   uuid,
		events: make(chan ami.Response),
		stop:   make(chan struct{}),
	}
	as.wg.Add(1)
	go as.run(ctx)
	return as, nil
}

// Logoff closes the current session with AMI.
func (as *AMIAsterisk) Logoff(ctx context.Context) error {
	close(as.stop)
	as.wg.Wait()

	return ami.Logoff(ctx, as.socket, as.uuid)
}

// Events returns an channel with events received from AMI.
func (as *AMIAsterisk) Events() <-chan ami.Response {
	return as.events
}

// SIPPeers fetch the list of SIP peers present on asterisk.
func (as *AMIAsterisk) SIPPeers(ctx context.Context) ([]ami.Response, error) {
	var peers []ami.Response
	resp, err := ami.SIPPeers(ctx, as.socket, as.uuid)
	switch {
	case err != nil:
		return nil, err
	case len(resp) == 0:
		return nil, errors.New("there's no sip peers configured")
	default:
		for _, v := range resp {
			peer, err := ami.SIPShowPeer(ctx, as.socket, as.uuid, v.Get("ObjectName"))
			if err != nil {
				return nil, err
			}
			peers = append(peers, peer)
		}
	}
	return peers, nil
}

func (as *AMIAsterisk) KickUser(ctx context.Context, user *entity.User) error {
	if user.Channel == nil {
		return errors.New("channel is null")
	}
	_, err := ami.ConfbridgeKick(ctx, as.socket, as.uuid, user.ConferenceID, user.Channel.ID)
	if err != nil {
		return err
	}
	return nil
}

func (as *AMIAsterisk) KickAllFromConference(ctx context.Context, conferenceID string) error {
	_, err := ami.ConfbridgeKick(ctx, as.socket, as.uuid, conferenceID, "all")
	return err
}

func (as *AMIAsterisk) run(ctx context.Context) {
	defer as.wg.Done()
	for {
		select {
		case <-as.stop:
			return
		case <-ctx.Done():
			return
		default:
			events, err := ami.Events(ctx, as.socket)
			if err != nil {
				log.Printf("AMI events failed: %v\n", err)
				return
			}
			as.events <- events
		}
	}
}