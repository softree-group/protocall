package porter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	errUploadFile = errors.New("error while send request to porter")
)

type PorterClientConfig struct {
	Host    string
	Port    string
	Timeout int
}

type PorterClient struct {
	http.Client
	addr string
}

func NewPorterClient(config *PorterClientConfig) *PorterClient {
	t := &PorterClient{
		addr: fmt.Sprintf("http://%v:%v", config.Host, config.Port),
	}
	t.Timeout = time.Duration(config.Timeout) * time.Second
	return t
}

func (p *PorterClient) UploadRecord(ctx context.Context, path string) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%v/records?from=%v&to=%v",
			p.addr,
			path,
			path,
		),
		http.NoBody,
	)
	if err != nil {
		return "", err
	}

	resp, err := p.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status code %v", errUploadFile, resp.StatusCode)
	}

	url, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(url), nil
}
