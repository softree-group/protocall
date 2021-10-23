package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/google/btree"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"

	"protocall/application/applications"
	"protocall/domain/entity"
	"protocall/domain/repository"
	"protocall/internal/config"
)

type Conference struct {
	reps repository.Repositories
	ari  ari.Client
}

func NewConference(reps repository.Repositories, ariClient ari.Client) *Conference {
	return &Conference{reps: reps, ari: ariClient}
}

func (c *Conference) RemoveParticipant(user *entity.User, meetID string) {
	conference := c.reps.GetConference(meetID)
	if conference == nil {
		return
	}
}

func (c *Conference) StartConference(user *entity.User) (*entity.Conference, error) {
	id, _ := uuid.GenerateUUID()
	conference := entity.NewConference(id, user.AsteriskAccount)
	conference.Participants.ReplaceOrInsert(user)
	user.ConferenceID = id
	c.reps.SaveUser(user)
	c.reps.SaveConference(conference)
	return conference, nil
}

func (c *Conference) JoinToConference(user *entity.User, meetID string) (*entity.Conference, error) {
	conference := c.reps.GetConference(meetID)
	if conference == nil {
		return nil, errors.New("no such meeting")
	}
	conference.Participants.ReplaceOrInsert(user)
	user.ConferenceID = meetID
	c.reps.SaveUser(user)
	c.reps.SaveConference(conference)
	return conference, nil
}

func (c *Conference) IsExist(meetID string) bool {
	return c.reps.GetConference(meetID) != nil
}

func postSnoop(id, snoopID, appArgs, app, spy, whisper string) (*fasthttp.Response, error) {
	clientt := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI(
		"http://pbx.softex-team.ru:10088/ari/channels/" +
			id +
			"/snoop?api_key=" +
			viper.GetString(config.ARIUser) +
			":" +
			viper.GetString(config.ARIPassword),
	)
	req.SetBodyString(fmt.Sprintf(`{"snoopID": "%s",
"app":  "%s",
"spy":  "%s",
"whisper":  "%s",
"appArgs":  "%s"}`, snoopID, app, spy, whisper, appArgs))
	req.Header.SetContentType("application/json")
	err := clientt.Do(req, resp)
	logrus.Info("REQ: ", req.String())
	if err != nil {
		logrus.Errorf("Сетевая ошибка по пути")
		return resp, err
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		logrus.Warnf("Сервер ответил %d", resp.StatusCode())
	}
	return resp, err
}

func (c *Conference) StartRecordUser(user *entity.User, conferenceID string) error {
	user.RecordPath = fmt.Sprintf("%v/%v/%v.wav", conferenceID, user.Username, time.Now().UTC().Unix())

	fmt.Println("HHHHHHH", user.RecordPath)

	resp, err := postSnoop(
		user.Channel.ID,
		fmt.Sprintf("%s_%v_%s", conferenceID, time.Now().UTC().Unix(), user.Username),
		user.RecordPath,
		viper.GetString(config.ARISnoopyApplication),
		"in",
		"both",
	)
	fasthttp.ReleaseResponse(resp)
	return err
}

func (c *Conference) StartRecord(user *entity.User, meetID string) error {
	conference := c.reps.GetConference(meetID)
	if conference == nil {
		return errors.New("does not exist")
	}

	if user.AsteriskAccount != conference.HostUserID {
		return errors.New("permissions denied")
	}

	conference.IsRecording = true
	c.reps.SaveConference(conference)

	conference.Participants.Ascend(func(item btree.Item) bool {
		if item == nil {
			return false
		}
		user := item.(*entity.User)
		if user == nil {
			return false
		}
		err := c.StartRecordUser(user, conference.ID)
		if err != nil {
			logrus.WithField("user", fmt.Sprintf("%+v", user)).Error("Fail to snoop: ", err)
		}
		return true
	})

	return nil
}

func (c *Conference) Get(meetID string) *entity.Conference {
	return c.reps.GetConference(meetID)
}

func (c *Conference) Delete(meetID string) {
	c.reps.DeleteConference(meetID)
	err := c.ari.Bridge().Delete(&ari.Key{
		Kind: ari.BridgeKey,
		ID:   meetID,
		App:  viper.GetString(config.ARIApplication),
	})
	if err != nil {
		logrus.Error("fail to delete bridge: ", err)
	}
}

func (c *Conference) TranslateRecord(string) (*entity.Message, error) {
	return nil, nil
}

func (c *Conference) UploadRecord(user *entity.User, meetID string) error {
	if err := c.reps.Upload(user.RecordPath, user.RecordPath); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (c *Conference) SendConference() {}

var _ applications.Conference = &Conference{}
