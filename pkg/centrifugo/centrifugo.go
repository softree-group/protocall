package centrifugo

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type Centrifugo struct {
	addr  string
	token string
}

func NewCentrifugo(cfg *Config) *Centrifugo {
	applySecrets(cfg)

	return &Centrifugo{
		addr:  cfg.Addr,
		token: cfg.Token,
	}
}

type SocketMessage map[string]interface{}

func (c *Centrifugo) Publish(channel string, payload SocketMessage) error {
	publish := struct {
		Method string `json:"method"`
		Params struct {
			Channel string        `json:"channel"`
			Data    SocketMessage `json:"data"`
		}
	}{
		Method: "publish",
		Params: struct {
			Channel string        `json:"channel"`
			Data    SocketMessage `json:"data"`
		}{Channel: channel, Data: payload},
	}

	data, err := json.Marshal(publish)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetBody(data)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", "apikey "+c.token)
	req.Header.SetMethod("POST")
	req.SetRequestURI(c.addr + "/api")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	return fasthttp.Do(req, resp)
}
