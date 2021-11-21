package app

import (
	"context"
	"protocall/internal/connector/application/applications"
	"protocall/internal/connector/domain/entity"
	"protocall/internal/connector/domain/repository"
	"protocall/internal/connector/domain/services"

	"github.com/google/btree"
	"github.com/sirupsen/logrus"
)

type ApplicationEventListener struct {
	reps          repository.Repositories
	bus           services.Bus
	conferenceApp applications.Conference
	ami           *AMIAsterisk
	socket        *Socket
}

func NewApplicationEventListener(
	reps repository.Repositories,
	bus services.Bus,
	conference applications.Conference,
	ami *AMIAsterisk,
	socket *Socket,
) *ApplicationEventListener {
	return &ApplicationEventListener{
		reps:          reps,
		bus:           bus,
		conferenceApp: conference,
		ami:           ami,
		socket:        socket,
	}
}

func (a *ApplicationEventListener) handleStartRecordEvent(event interface{}) {
	data, ok := event.(entity.EventDefault)
	if !ok {
		logrus.Error("StartRecord event invalid event type")
		return
	}
	user := a.reps.FindUser(data.User.SessionID)
	if user == nil {
		logrus.Error("no user in handle startRecord event")
		a.bus.Publish("fail", event)
		return
	}

	data.User = user
	data.ConferenceID = user.ConferenceID
	a.reps.Store(user.ConferenceID, data.RecName)
}

func (a *ApplicationEventListener) handleSavedEvent(event interface{}) {
	data, ok := event.(entity.EventDefault)
	if !ok {
		logrus.Error("Saved event invalid event type")
		return
	}
	user := a.reps.FindUser(data.User.SessionID)
	if user == nil {
		logrus.Error("no user in handle saved event")
		a.bus.Publish("fail", event)
		return
	}

	data.User = user
	data.ConferenceID = user.ConferenceID

	err := a.conferenceApp.UploadRecord(data.RecName)
	if err != nil {
		logrus.Error("fail to upload: ", err)
		a.bus.Publish("fail", data)
		return
	}
	a.bus.Publish("uploaded", data)
}

func (a *ApplicationEventListener) handleUploadedEvent(event interface{}) {
	data, ok := event.(entity.EventDefault)
	if !ok {
		logrus.Error("uploaded event invalid event type")
		return
	}

	err := a.conferenceApp.TranslateRecord(data.User, data.RecName)
	if err != nil {
		logrus.Error("fail to translate: ", err)
		a.bus.Publish("fail", event)
	}
	a.bus.Publish("translated", event)
}

func (a *ApplicationEventListener) handleTranslatedEvent(event interface{}) {
	data, ok := event.(entity.EventDefault)
	if !ok {
		logrus.Error("Translated event invalid event type")
		return
	}
	a.reps.DoneJob(data.ConferenceID, data.RecName)

	if isDone, _ := a.reps.IsDone(data.ConferenceID); isDone {
		a.bus.Publish("conferenceTranslated", event)
	}
}

func (a *ApplicationEventListener) handleConferenceTranslatedEvent(event interface{}) {
	data, ok := event.(entity.EventDefault)
	if !ok {
		logrus.Error("Conference translated event invalid event type")
		return
	}

	conference := a.reps.GetConference(data.ConferenceID)
	err := a.conferenceApp.CreateProtocol(conference)
	if err != nil {
		logrus.Error("fail to create protocol: ", err)
		a.bus.Publish("fail", event)
	}

	conference.Participants.Ascend(func(i btree.Item) bool {
		if i == nil {
			return false
		}

		participant, ok := i.(*entity.User)
		if !ok {
			return false
		}
		a.reps.DeleteUser(participant.SessionID)
		return true
	})

	a.reps.DeleteConference(data.ConferenceID)
}

func (a *ApplicationEventListener) handleFailEvent(event interface{}) {
	data, ok := event.(entity.EventDefault)
	if !ok {
		logrus.Error("Translated event invalid event type")
		return
	}
	logrus.Errorln("error while process events")

	a.reps.DoneJob(data.ConferenceID, data.RecName)
	a.reps.DeleteUser(data.User.SessionID)
}

func (a *ApplicationEventListener) handleLeaveUser(event interface{}) {
	data, ok := event.(entity.EventDefault)
	if !ok {
		logrus.Error("Translated event invalid event type")
		return
	}
	logrus.Errorln("error while process events")

	conf := a.reps.GetConference(data.ConferenceID)

	if conf.HostUserID == data.User.AsteriskAccount {
		conf.Participants.Ascend(func(item btree.Item) bool {
			if item == nil {
				return false
			}

			participant := item.(*entity.User)
			if participant == nil {
				return false
			}

			a.bus.Publish("leave/"+participant.SessionID, "")
			a.reps.FreeAccount(participant.AsteriskAccount)
			return true
		})

		err := a.ami.KickAllFromConference(context.Background(), data.ConferenceID)
		if err != nil {
			logrus.Error("fail to kick all users: ", err)
		}
		a.socket.PublishEndConference(data.ConferenceID)
	}
}

func (a *ApplicationEventListener) Listen() {
	startRecordEvent := a.bus.Subscribe("startRecord")
	savedEvent := a.bus.Subscribe("saved")
	uploadedEvent := a.bus.Subscribe("uploaded")
	translatedEvent := a.bus.Subscribe("translated")
	conferenceTranslatedEvent := a.bus.Subscribe("conferenceTranslated")
	failEvent := a.bus.Subscribe("fail")
	deleteEvent := a.bus.Subscribe("leave")

	for {
		select {
		case event := <-startRecordEvent.Channel():
			a.handleStartRecordEvent(event)
		case event := <-savedEvent.Channel():
			a.handleSavedEvent(event)
		case event := <-uploadedEvent.Channel():
			a.handleUploadedEvent(event)
		case event := <-translatedEvent.Channel():
			a.handleTranslatedEvent(event)
		case event := <-conferenceTranslatedEvent.Channel():
			a.handleConferenceTranslatedEvent(event)
		case event := <-failEvent.Channel():
			a.handleFailEvent(event)
		case event := <-deleteEvent.Channel():
			a.handleLeaveUser(event)
		}
	}
}