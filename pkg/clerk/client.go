package clerk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"protocall/internal/stapler"
	"protocall/internal/translator"
)

var (
	errTranslate = errors.New("error while send request to clerk")
)

type ClerkClientConfig struct {
	Host    string
	Port    string
	Timeout int
}

type ClerkClient struct {
	http.Client
	addr string
}

func NewClerkClient(config *ClerkClientConfig) *ClerkClient {
	t := &ClerkClient{
		addr: fmt.Sprintf("http://%v:%v", config.Host, config.Port),
	}
	t.Timeout = time.Duration(config.Timeout) * time.Second
	return t
}

func (c *ClerkClient) TranslateRecord(ctx context.Context, data *translator.TranslateRequest) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%v/translations", c.addr),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errTranslate
	}

	return nil
}

func (c *ClerkClient) CreateProtocol(ctx context.Context, data *stapler.ProtocolRequest) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%v/protocols", c.addr),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errTranslate
	}

	return nil
}
