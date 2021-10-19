package app

import (
	"errors"
	"fmt"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"protocall/application/applications"
	"protocall/config"
	"protocall/domain/entity"
	"protocall/domain/repository"
	"time"
)

type Conference struct {
	reps *repository.Repositories
	ari  ari.Client
}

func NewConference(reps *repository.Repositories, ari ari.Client) *Conference {
	return &Conference{reps: reps, ari: ari}
}

func (c Conference) StartConference(user *entity.User) (*entity.Conference, error) {
	id, _ := uuid.GenerateUUID()
	conference := &entity.Conference{
		ID:           id,
		Participants: []*entity.User{user},
		HostUserID:   user.AsteriskAccount,
		BridgeID:     "",
	}
	user.ConferenceID = id
	c.reps.User.Save(user)
	c.reps.Conference.Save(conference)
	return conference, nil
}

func (c Conference) JoinToConference(user *entity.User, meetID string) (*entity.Conference, error) {
	conference := c.reps.Conference.Get(meetID)
	if conference == nil {
		return nil, errors.New("no such meeting")
	}
	conference.Participants = append(conference.Participants, user)
	user.ConferenceID = meetID
	c.reps.User.Save(user)
	c.reps.Conference.Save(conference)
	return conference, nil
}

func (c Conference) IsExist(meetID string) bool {
	return c.reps.Conference.Get(meetID) != nil
}

func postSnoop(id, snoopId, appArgs, app, spy, whisper string) (*fasthttp.Response, error) {
	clientt := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI("http://pbx.softex-team.ru:10088/ari/channels/" + id + "/snoop?api_key=" + viper.GetString(config.ARIUser) + ":" + viper.GetString(config.ARIPassword))
	req.SetBodyString(fmt.Sprintf(`{"snoopId": "%s",
"app":  "%s",
"spy":  "%s",
"whisper":  "%s",
"appArgs":  "%s"}`, snoopId, app, spy, whisper, appArgs))
	req.Header.SetContentType("application/json")
	err := clientt.Do(req, resp)
	logrus.Info("REQ: ", req.String())
	if err != nil {
		logrus.Errorf("Сетевая ошибка по пути")
		return resp, err
	}
	if resp.StatusCode() >= 400 {
		logrus.Warnf("Сервер ответил %d", resp.StatusCode())
	}
	return resp, err
}

func (c Conference) StartRecordUser(user *entity.User, conferenceID string) error {
	resp, err := postSnoop(user.Channel.ID, fmt.Sprintf("%s_%v_%s", conferenceID, time.Now().UTC().Unix(), user.Username), "some", viper.GetString(config.ARISnoopyApplication), "in", "both")
	fasthttp.ReleaseResponse(resp)
	//logrus.Info("Channel: ", user.Channel)
	//_, err := c.ari.Channel().Snoop(user.Channel, fmt.Sprintf("%s/%v_%s", conferenceID, time.Now().UTC().Unix(), "some"), &ari.SnoopOptions{
	//	App:     viper.GetString(config.ARISnoopyApplication),
	//	AppArgs: user.Channel.ID,
	//	Spy:     "in",
	//	Whisper: "both",
	//})
	return err
}

func (c Conference) StartRecord(user *entity.User, meetID string) error {
	conference := c.reps.Conference.Get(meetID)
	if conference == nil {
		return errors.New("does not exist")
	}

	if user.AsteriskAccount != conference.HostUserID {
		return errors.New("permissions denied")
	}

	conference.IsRecording = true
	c.reps.Conference.Save(conference)

	for _, user := range conference.Participants {

		err := c.StartRecordUser(user, conference.ID)

		if err != nil {
			logrus.WithField("user", fmt.Sprintf("%+v", user)).Error("Fail to snoop: ", err)
		}
	}

	return nil
}

func (c Conference) Get(meetID string) *entity.Conference {
	return c.reps.Conference.Get(meetID)
}

var _ applications.Conference = Conference{}
