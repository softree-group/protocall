package clerk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"protocall/internal/clerk/stapler"
	"protocall/internal/clerk/translator"
	"protocall/internal/connector/domain/entity"
)

var (
	errTranslate = errors.New("error while send request to clerk")
)

type TranslatorConfig struct {
	Host    string
	Port    string
	Timeout int
}

type Translator struct {
	c    *http.Client
	addr string
}

func NewTranslator(config *TranslatorConfig) *Translator {
	return &Translator{
		c: &http.Client{
			Timeout: time.Second * time.Duration(config.Timeout),
		},
		addr: fmt.Sprintf("http://%v:%v", config.Host, config.Port),
	}
}

func (t *Translator) TranslateConference(u *entity.User, c *entity.Conference) error {
	body, err := json.Marshal(translator.TranslateRequest{
		StartTime: c.Start,
		User: translator.User{
			Username: u.Username,
			Email:    u.Email,
			Path:     u.RecordPath,
		},
	})
	if err != nil {
		return err
	}

	resp, err := t.c.Post(
		fmt.Sprintf("%v/records", t.addr),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errTranslate
	}

	return nil
}

func (t *Translator) CreateProtocol(conferenceID string, sendTo []string) error {
	fmt.Println("EMAILS", sendTo)
	body, err := json.Marshal(stapler.SendProtocolRequest{
		ConferenceID: conferenceID,
		To:           sendTo,
	})
	if err != nil {
		return err
	}

	resp, err := t.c.Post(
		fmt.Sprintf("%v/protocols", t.addr),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errTranslate
	}

	return nil
}
