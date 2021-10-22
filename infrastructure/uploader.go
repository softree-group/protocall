package infrastructure

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	errUploadFile = errors.New("error while send request to uploader")
)

type UploaderConfig struct {
	Host    string
	Port    string
	Timeout int
}

type Uploader struct {
	httpClient *http.Client
	addr       string
}

func NewUploader(config *UploaderConfig) *Uploader {
	return &Uploader{
		httpClient: &http.Client{
			Timeout: time.Second * time.Duration(config.Timeout),
		},
		addr: fmt.Sprintf("http://%v:%v", config.Host, config.Port),
	}
}

func (u *Uploader) Upload(from, to string) error {
	if resp, err := u.httpClient.Post(
		fmt.Sprintf("%v/upload?from=%v&to=%v",
			u.addr,
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
