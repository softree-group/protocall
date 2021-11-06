package infrastructure

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"protocall/domain/entity"
	"protocall/domain/services"
	"protocall/internal/config"
)

type Centrifugo struct {
}

func NewCentrifugo() *Centrifugo {
	return &Centrifugo{}
}

type publishData struct {
	Method string `json:"method"`
	Params struct {
		Channel string               `json:"channel"`
		Data    entity.SocketMessage `json:"data"`
	}
}

func (c Centrifugo) Publish(channel string, payload entity.SocketMessage) error {

	publish := publishData{
		Method: "publish",
		Params: struct {
			Channel string               `json:"channel"`
			Data    entity.SocketMessage `json:"data"`
		}{Channel: channel, Data: payload},
	}

	data, err := json.Marshal(publish)
	logrus.Debugf("socket: %s", data)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetBody(data)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", "apikey "+viper.GetString(config.CentrifugoAPIKey))
	req.Header.SetMethod("POST")
	req.SetRequestURI(viper.GetString(config.CentrifugoHost) + "/api")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	return fasthttp.Do(req, resp)
}

var _ services.Socket = &Centrifugo{}
