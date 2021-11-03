package interfaces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"protocall/api"
	"protocall/cmd/clerk/domain"
	"protocall/pkg/logger"
)

type Application struct {
	domain.Sender
	domain.Stapler
	domain.Translator
}

func NewApplication(
	sender domain.Sender,
	gluer domain.Stapler,
	translator domain.Translator,
) *Application {
	return &Application{
		sender,
		gluer,
		translator,
	}
}

func (a *Application) translate(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	translateRequest := api.TranslateRequest{}
	if err := json.Unmarshal(body, &translateRequest); err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	a.TranslateRecord(&translateRequest)

	res.WriteHeader(http.StatusNoContent)
}

func (a *Application) create(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendRequest := api.SendProtocolRequest{}

	if err := json.Unmarshal(body, &sendRequest); err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("DEBUG", sendRequest)

	a.NewProtocol(&sendRequest, a.SendSMTP)

	res.WriteHeader(http.StatusNoContent)
}
