package infrastructure

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	errUploadFile = errors.New("error while send request to rising")
)

type PorterClientConfig struct {
	Host    string
	Port    string
	Timeout int
}

type PorterClient struct {
	httpClient *http.Client
	addr       string
}

func NewUploader(config *PorterClientConfig) *PorterClient {
	return &PorterClient{
		httpClient: &http.Client{
			Timeout: time.Second * time.Duration(config.Timeout),
		},
		addr: fmt.Sprintf("http://%v:%v", config.Host, config.Port),
	}
}

func (u *PorterClient) UploadConference(path string) error {
	resp, err := u.httpClient.Post(
		fmt.Sprintf("%v/upload?from=%v&to=%v",
			u.addr,
			path,
			path,
		),
		"",
		nil,
	)
	if err != nil {
		return errUploadFile
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errUploadFile
	}

	return nil
}
