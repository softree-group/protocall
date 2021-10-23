package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"protocall/api"
	"protocall/domain/entity"
)

var (
	errTranslate = errors.New("error while send request to translator")
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

func (t *Translator) Translate(u *entity.User, c *entity.Conference, recordPath string) error {
	body, err := json.Marshal(api.TranslatorRequest{
		ConfID:    c.ID,
		StartTime: c.Start,
		User: api.User{
			Username: u.Username,
			Email:    u.Email,
			Path:     recordPath,
		},
	})
	if err != nil {
		return err
	}

	resp, err := t.c.Post(
		fmt.Sprintf("%v/translate", t.addr),
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
