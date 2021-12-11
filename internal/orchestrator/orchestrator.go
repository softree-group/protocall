package orchestrator

import (
	"context"
	"fmt"

	"protocall/internal/account"
	"protocall/internal/conference"
	"protocall/internal/operator"
	"protocall/internal/socket"
	"protocall/internal/user"
	"protocall/pkg/bus"
	"protocall/pkg/logger"
)

type Orchestrator struct {
	bus        bus.Client
	conference *conference.Application
	account    *account.Application
	user       *user.Application
	ami        *operator.Operator
	socket     *socket.Socket
}

func NewOrchestrator(
	bus bus.Client,
	conference *conference.Application,
	account *account.Application,
	user *user.Application,
	ami *operator.Operator,
	socket *socket.Socket,
) *Orchestrator {
	return &Orchestrator{
		bus:        bus,
		user:       user,
		conference: conference,
		account:    account,
		ami:        ami,
		socket:     socket,
	}
}

func (o *Orchestrator) handleStartRecordEvent(event interface{}) {
	data, ok := event.(conference.Event)
	if !ok {
		logger.L.Error("StartRecord event invalid event type")
		return
	}
	logger.L.WithField("user", data.User.SessionID).Debug("handle start record")

	user := o.user.Find(data.User.SessionID)
	if user == nil {
		logger.L.Error("no user in handle startRecord event")
		o.bus.Publish("fail", event)
		return
	}

	data.User = user
	data.ConferenceID = user.ConferenceID

	o.conference.Store(user.ConferenceID, data.Record.Path)
}

func (a *Orchestrator) handleSavedEvent(event interface{}) {
	data, ok := event.(conference.Event)
	if !ok {
		logger.L.Error("Saved event invalid event type")
		return
	}
	logger.L.WithField("user", data.User.SessionID).Debug("handle saved")

	user := a.user.Find(data.User.SessionID)
	if user == nil {
		logger.L.Error("no user in handle saved event")
		a.bus.Publish("fail", event)
		return
	}

	data.User = user
	data.ConferenceID = user.ConferenceID

	var err error
	data.Record.URI, err = a.conference.UploadRecord(data.Record.Path)
	if err != nil {
		logger.L.Error("fail to upload: ", err)
		a.bus.Publish("fail", data)
		return
	}

	user.Records = append(user.Records, data.Record.Path)
	a.user.Save(user)

	a.bus.Publish("uploaded", data)
}

func (a *Orchestrator) handleUploadedEvent(event interface{}) {
	data, ok := event.(conference.Event)
	if !ok {
		logger.L.Error("uploaded event invalid event type")
		return
	}
	logger.L.WithField("user", data.User.SessionID).Debug("handle uploaded")

	err := a.conference.TranslateRecord(data.User, data.Record)
	if err != nil {
		logger.L.Error("fail to translate: ", err)
		a.bus.Publish("fail", event)
	}
}

func (a *Orchestrator) handleTranslatedEvent(event interface{}) {
	data, ok := event.(conference.Event)
	if !ok {
		logger.L.Error("Translated event invalid event type")
		return
	}
	logger.L.WithField("user", data.User.SessionID).Debug("handle translated")

	user := a.user.Find(data.User.SessionID)
	if user == nil {
		logger.L.Error("no user in handle saved event")
		a.bus.Publish("fail", event)
		return
	}
	user.Texts = append(user.Texts, data.Text)

	fmt.Printf("123 %+v", user)

	a.user.Save(user)

	a.conference.DoneJob(data.ConferenceID, data.Record.Path)
	if isDone, _ := a.conference.IsDone(data.ConferenceID); isDone {
		a.bus.Publish("conferenceTranslated", event)
	}
}

func (a *Orchestrator) handleConferenceTranslatedEvent(event interface{}) {
	data, ok := event.(conference.Event)
	if !ok {
		logger.L.Error("Conference translated event invalid event type")
		return
	}
	logger.L.WithField("user", data.User.SessionID).Debug("handle conference translated")

	conference := a.conference.Get(data.ConferenceID)
	err := a.conference.CreateProtocol(conference)
	if err != nil {
		logger.L.Error("fail to create protocol: ", err)
		a.bus.Publish("fail", event)
	}

	// conference.Participants.Ascend(func(i btree.Item) bool {
	// 	if i == nil {
	// 		return false
	// 	}

	// 	participant, ok := i.(*user.User)
	// 	if !ok {
	// 		return false
	// 	}
	// 	a.reps.DeleteUser(participant.SessionID)
	// 	return true
	// })

	a.conference.Delete(data.ConferenceID)
}

func (a *Orchestrator) handleFailEvent(event interface{}) {
	data, ok := event.(conference.Event)
	if !ok {
		logger.L.Error("Translated event invalid event type")
		return
	}
	logger.L.WithField("user", data.User.SessionID).Debug("handle fail")

	a.conference.DoneJob(data.ConferenceID, data.Record.Path)
	a.user.Delete(data.User.SessionID)
}

func (o *Orchestrator) handleLeaveUser(event interface{}) {
	data, ok := event.(conference.Event)
	if !ok {
		logger.L.Error("Translated event invalid event type")
		return
	}
	logger.L.WithField("user", data.User.SessionID).Debug("handle leave")

	conf := o.conference.Get(data.ConferenceID)

	if conf.HostUserID == data.User.AsteriskAccount {
		// conf.Participants.Ascend(func(item btree.Item) bool {
		// 	if item == nil {
		// 		return false
		// 	}

		// 	participant := item.(*user.User)
		// 	if participant == nil {
		// 		return false
		// 	}

		// 	a.bus.Publish("leave/"+participant.SessionID, "")
		// 	a.account.FreeAccount(participant.AsteriskAccount)
		// 	return true
		// })

		err := o.ami.KickAllFromConference(context.Background(), data.ConferenceID)
		if err != nil {
			logger.L.Error("fail to kick all users: ", err)
		}
		o.socket.PublishEndConference(data.ConferenceID)
		return
	}

	o.bus.Publish("leave/"+data.User.SessionID, "")
}

func (o *Orchestrator) Run() {
	startRecordEvent := o.bus.Subscribe("startRecord")
	savedEvent := o.bus.Subscribe("saved")
	uploadedEvent := o.bus.Subscribe("uploaded")
	translatedEvent := o.bus.Subscribe("translated")
	conferenceTranslatedEvent := o.bus.Subscribe("conferenceTranslated")
	failEvent := o.bus.Subscribe("fail")
	deleteEvent := o.bus.Subscribe("leave")

	for {
		select {
		case event := <-startRecordEvent.Channel():
			go o.handleStartRecordEvent(event)
		case event := <-savedEvent.Channel():
			go o.handleSavedEvent(event)
		case event := <-uploadedEvent.Channel():
			go o.handleUploadedEvent(event)
		case event := <-translatedEvent.Channel():
			go o.handleTranslatedEvent(event)
		case event := <-conferenceTranslatedEvent.Channel():
			go o.handleConferenceTranslatedEvent(event)
		case event := <-failEvent.Channel():
			go o.handleFailEvent(event)
		case event := <-deleteEvent.Channel():
			go o.handleLeaveUser(event)
		}
	}
}
