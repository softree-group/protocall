package snoopy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"protocall/application/app"
	"protocall/config"
	"protocall/domain/entity"
	"protocall/domain/repository"
	"strconv"
	"time"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/CyCoreSystems/ari/v5/ext/record"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Snoopy struct {
	ari           ari.Client
	reps          *repository.Repositories
	eventListener *app.EventListener
}

func New(reps *repository.Repositories, el *app.EventListener) *Snoopy {
	ariClient, err := native.Connect(&native.Options{
		Application:  viper.GetString(config.ARISnoopyApplication),
		URL:          viper.GetString(config.ARIUrl),
		WebsocketURL: viper.GetString(config.ARIWebsocketUrl),
		Username:     viper.GetString(config.ARIUser),
		Password:     viper.GetString(config.ARIPassword),
	})
	if err != nil {
		logrus.Fatal("Fail to connect snoopy app")
	}

	return &Snoopy{
		ari:           ariClient,
		reps:          reps,
		eventListener: el,
	}
}

const timeout = 3

var (
	errUploadFile       = errors.New("cannot upload file to S3")
	errNotifyTranslator = errors.New("error while send request to translator")
)

func uploadFile(from, to string) error {
	c := http.Client{
		Timeout: timeout * time.Second,
	}
	if resp, err := c.Post(
		fmt.Sprintf("http://%v:%v/upload?from=%v&to=%v",
			viper.GetString(config.UploaderHostConf),
			viper.GetString(config.UploaderPortConf),
			from,
			to,
		),
		"",
		nil,
	); resp.StatusCode != http.StatusNoContent || err != nil {
		errUploadFile = fmt.Errorf("%w: %v", errUploadFile, from)
		if err != nil {
			errUploadFile = fmt.Errorf("%w: %v", errUploadFile, err)

		}
		return errUploadFile
	}
	return nil
}

type user struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Path     string `json:"path" binding:"required"`
}

type translatorReq struct {
	ConfID    string    `json:"conf_id" binding:"required"`
	StartTime time.Time `json:"start" binding:"required"`
	User      user
}

func notifyTranslator(u *entity.User, conference *entity.Conference, audioPath string) error {
	body, err := json.Marshal(translatorReq{
		ConfID:    conference.ID,
		StartTime: conference.Start,
		User: user{
			u.Username,
			u.Email,
			audioPath,
		},
	})
	if err != nil {
		return err
	}
	c := http.Client{
		Timeout: timeout * time.Second,
	}
	if resp, err := c.Post(
		fmt.Sprintf("http://%v:%v/translate",
			viper.GetString(config.TranslatorHostConf),
			viper.GetString(config.TranslatorPortConf),
		),
		"application/json",
		bytes.NewReader(body),
	); resp.StatusCode != http.StatusOK || err != nil {
		if err != nil {
			errNotifyTranslator = fmt.Errorf("%w: %v", errNotifyTranslator, err)
		}
		return errNotifyTranslator
	}
	return nil
}

func (s Snoopy) channelHandler(channel *ari.ChannelHandle, sessionID string, timestamp *time.Time) {
	sub := channel.Subscribe(ari.Events.All)
	end := channel.Subscribe(ari.Events.StasisEnd)

	defer sub.Cancel()
	defer end.Cancel()

	ctx := context.Background()
	rec := record.Record(ctx, channel)

	for {
		select {
		case event := <-sub.Events():
			logrus.Info("In SPY: ", event.GetType())
		case <-end.Events():
			logrus.Info("saving record for ", channel.ID())
			res, err := rec.Result()
			if err != nil {
				logrus.Error("Fail to get result from record for channel ", channel.ID(), ". Error: ", err)
				return
			}

			user := s.reps.User.Find(sessionID)
			if user == nil {
				logrus.Error("cannot find user")
				return
			}

			location := fmt.Sprintf("%v/%v/%v.wav", user.ConferenceID, user.Username, strconv.FormatInt(timestamp.Unix(), 10))
			err = res.Save(location)
			if err != nil {
				logrus.Error("fail to save result record for channel ", channel.ID(), ". Error: ", err)
				return
			}

			if err := uploadFile(location, location); err != nil {
				logrus.Error(err)
				return
			}

			logrus.Info("saved record for ", channel.ID())

			fmt.Println("CONFID", user.ConferenceID, "||||")

			conference := s.reps.Conference.Get(user.ConferenceID)
			if conference == nil {
				logrus.Error("cannot find conference")
				return
			}
			if err := notifyTranslator(user, conference, fmt.Sprintf("%v/%v", user.ConferenceID, user.Username)); err != nil {
				logrus.Error(err)
				return
			}

			logrus.Info("successfully notified translator", channel.ID())
			s.eventListener.FreeResouces()
			return
		}
	}
}

func (s Snoopy) listen() {
	start := s.ari.Bus().Subscribe(nil, ari.Events.StasisStart)
	for {
		select {
		case event := <-start.Events():
			value := event.(*ari.StasisStart)
			channel := s.ari.Channel().Get(value.Key(ari.ChannelKey, value.Channel.ID))
			timestamp := time.Time(value.Timestamp)

			logrus.Info("snoop channel: ", channel.ID())
			go s.channelHandler(channel, value.Args[0], &timestamp)
		}
	}
}

func (s Snoopy) Snoop() {
	logrus.Info("Start snooping...")
	s.listen()
}
