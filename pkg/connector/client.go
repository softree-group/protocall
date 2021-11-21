package connector

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"protocall/internal/translator"
)

type ConnectorClient struct {
	client *http.Client
	addr   string
	token  string
}

func NewConnectorCLient(client *http.Client, config *ConnectorClientConfig) *ConnectorClient {
	return &ConnectorClient{
		client: client,
		addr:   fmt.Sprintf("http://%v:%v", config.Host, config.Port),
		token:  config.Token,
	}
}

var (
	errResp = errors.New("failed status code from connector")
)

func (c *ConnectorClient) TranslationDone(ctx context.Context, r *translator.TranslateRequest) error {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		fmt.Sprintf("%v/translates", c.addr),
		strings.NewReader(r.User.SessionID),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%w: %d", errResp, resp.StatusCode)
	}
	return nil
}
