package interfaces

import (
	"encoding/json"
	"io"
	"net/http"

	"protocall/api"
	"protocall/cmd/clerk/application"
	"protocall/cmd/clerk/domain"
	"protocall/pkg/logger"
)

type Application struct {
	domain.Sender
	domain.Gluer
	domain.Translator
}

func NewApplication(
	sender domain.Sender,
	gluer domain.Gluer,
	translator domain.Translator,
) *Application {
	return &Application{
		sender,
		gluer,
		translator,
	}
}

func (a *Application) ProcessRecord(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	translate := api.TranslateRequest{}
	if err := json.Unmarshal(body, &translate); err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	jobID, err := a.CreateJob(&translate)
	if err != nil {
		logger.L.Error(err)
		return
	}

	res.Write([]byte(jobID))
	res.WriteHeader(http.StatusOK)
}

func (a *Application) GetJobStatus(res http.ResponseWriter, req *http.Request) {
	status, err := a.GetStatus(req.URL.Query().Get("id"))
	if err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch status {
	case application.Running:
		res.Write([]byte("RUNNING"))
	case application.Ready:
		res.Write([]byte("READY"))
	case application.Failed:
		res.Write([]byte("FAILED"))
	}
	res.WriteHeader(http.StatusOK)
}

func (a *Application) MakeProtocol(res http.ResponseWriter, req *http.Request) {
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

	protocall, err := a.Merge(req.Context(), sendRequest.ConferenceID)
	if err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.SendProtocol(req.Context(), protocall, sendRequest.To); err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
