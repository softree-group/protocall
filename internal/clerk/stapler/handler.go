package stapler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"protocall/pkg/logger"
)

type SendProtocolRequest struct {
	ConferenceID string   `json:"conference_id" binding:"required"`
	To           []string `json:"to" binding:"required"`
}

type StaplerRepository interface {
	NewProtocol(
		*SendProtocolRequest,
		...func(context.Context, string, []string),
	)
}

type SenderRepository interface {
	Send(context.Context, string, []string)
}

type StaplerHandler struct {
	app  StaplerRepository
	smtp SenderRepository
}

func (s *StaplerHandler) CreateProtocol(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendRequest := SendProtocolRequest{}

	if err := json.Unmarshal(body, &sendRequest); err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	s.app.NewProtocol(&sendRequest, s.smtp.Send)

	res.WriteHeader(http.StatusNoContent)
}
