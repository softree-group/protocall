package conference

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"protocall/internal/stapler"
	"protocall/internal/translator"
	"protocall/internal/user"
	"protocall/pkg/bus"
	"protocall/pkg/logger"

	"github.com/hashicorp/go-uuid"

	"github.com/valyala/fasthttp"
)

type ConferenceStorage interface {
	GetConference(conferenceID string) *Conference
	SaveConference(conference *Conference)
	DeleteConference(conferenceID string)
	Store(conferenceID, recName string) error
	DoneJob(conferenceID, recName string) error
	IsDone(conferenceID string) (bool, error)
	GetConferenceInfo(id string) (*ConferenceInfo, error)
}

type Translator interface {
	TranslateRecord(ctx context.Context, data *translator.TranslateRequest) error
	CreateProtocol(ctx context.Context, data *stapler.ProtocolRequest) error
}

type Porter interface {
	UploadRecord(context.Context, string) (string, error)
}

type ARIClient interface {
	PostSnoop(string, string, string) (*fasthttp.Response, error)
}

type Application struct {
	storage    ConferenceStorage
	translator Translator
	porter     Porter
	bus        bus.Client
	ari        ARIClient
	user       *user.Application
}

func NewApplication(
	bus bus.Client,
	storage ConferenceStorage,
	translator Translator,
	porter Porter,
) *Application {
	return &Application{
		bus:        bus,
		storage:    storage,
		translator: translator,
		porter:     porter,
	}
}

func (a *Application) RemoveParticipant(user *user.User, meetID string) {
	conference := a.storage.GetConference(meetID)
	if conference == nil {
		return
	}
}

func (a *Application) StartConference(user *user.User) (*Conference, error) {
	id, _ := uuid.GenerateUUID()
	conference := NewConference(id, user.AsteriskAccount)
	conference.Participants.ReplaceOrInsert(user)
	user.ConferenceID = id
	a.user.SaveUser(user)
	a.storage.SaveConference(conference)
	return conference, nil
}

func (a *Application) JoinToConference(user *user.User, meetID string) (*Conference, error) {
	conference := a.storage.GetConference(meetID)
	if conference == nil {
		return nil, errors.New("no such meeting")
	}
	conference.Participants.ReplaceOrInsert(user)
	user.ConferenceID = meetID
	c.reps.SaveUser(user)
	a.storage.SaveConference(conference)
	return conference, nil
}

func (a *Application) IsExist(meetID string) bool {
	return a.storage.GetConference(meetID) != nil
}

func (a *Application) StartRecordUser(user *user.User, conferenceID string) error {
	resp, err := a.ari.PostSnoop(
		user.Channel.ID,
		fmt.Sprintf("%s_%v_%s", conferenceID, time.Now().UTC().Unix(), user.Username),
		fmt.Sprintf("%v/%v/%v.wav", conferenceID, user.Username, time.Now().UTC().Unix())+","+user.SessionID,
	)
	fasthttp.ReleaseResponse(resp)
	return err
}

func (a *Application) StartRecord(user *user.User, meetID string) error {
	conference := a.storage.GetConference(meetID)
	if conference == nil {
		return errors.New("does not exist")
	}

	if user.AsteriskAccount != conference.HostUserID {
		return errors.New("permissions denied")
	}

	conference.IsRecording = true
	a.storage.SaveConference(conference)

	// conference.Participants.Ascend(func(item btree.Item) bool {
	// 	if item == nil {
	// 		return false
	// 	}
	// 	user := item.(*user.User)
	// 	if user == nil {
	// 		return false
	// 	}
	// 	err := c.StartRecordUser(user, conference.ID)
	// 	if err != nil {
	// 		logger.L.WithField("user", fmt.Sprintf("%+v", user)).Error("Fail to snoop: ", err)
	// 	}
	// 	return true
	// })

	return nil
}

func (a *Application) Store(conferenceID, recName string) error {
	return a.storage.Store(conferenceID, recName)
}

func (a *Application) Get(meetID string) *Conference {
	return a.storage.GetConference(meetID)
}

func (a *Application) Delete(meetID string) {
	a.storage.DeleteConference(meetID)
}

func (a *Application) DoneJob(conferenceID, recName string) error {
	return a.storage.DoneJob(conferenceID, recName)
}

func (a *Application) IsDone(conferenceID string) (bool, error) {
	return a.storage.IsDone(conferenceID)
}

var regexTime = regexp.MustCompile(`(\d+).wav`)

func (a *Application) TranslateRecord(user *user.User, record *translator.Record) error {
	match := regexTime.FindStringSubmatch(record.URI)
	if len(match) < 2 {
		return fmt.Errorf("invalid pattern recordPath: %s", record.URI)
	}
	connTime, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return err
	}

	if err := a.translator.TranslateRecord(context.TODO(), &translator.TranslateRequest{
		User: translator.User{
			Username:    user.Username,
			ConnectTime: time.Unix(connTime, 0).Add(time.Hour * 3),
			SessionID:   user.SessionID,
			Record: translator.Record{
				URI:    record.URI,
				Length: record.Length,
				Path:   record.Path,
			},
			Text: strings.Replace(record.Path, ".wav", ".txt", -1),
		},
	}); err != nil {
		logger.L.Error(err)
		return err
	}
	return nil
}

func (a *Application) UploadRecord(path string) (string, error) {
	url, err := a.porter.UploadRecord(context.TODO(), path)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (a *Application) CreateProtocol(conference *Conference) error {
	users := []stapler.User{}
	// conference.Participants.Ascend(
	// 	func(i btree.Item) bool {
	// 		if i == nil {
	// 			return false
	// 		}

	// 		participant, ok := i.(*user.User)
	// 		if !ok {
	// 			return false
	// 		}
	// 		user := stapler.User{}
	// 		user.Email = participant.Email
	// 		user.NeedProtocol = participant.NeedProtocol
	// 		user.Records = participant.Records
	// 		user.Texts = participant.Texts
	// 		user.Username = participant.Username
	// 		users = append(users, user)
	// 		return true
	// 	},
	// )

	if err := a.translator.CreateProtocol(context.TODO(), &stapler.ProtocolRequest{Users: users}); err != nil {
		logger.L.Error(err)
		return err
	}

	return nil
}

func (a *Application) GetConferenceInfo(id string) (*ConferenceInfo, error) {
	return a.storage.GetConferenceInfo(id)
}
