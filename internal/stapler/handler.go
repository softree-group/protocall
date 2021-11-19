package stapler

import (
	"encoding/json"
	"io"
	"net/http"

	"protocall/pkg/logger"
)

type StaplerHandler struct {
	App *Stapler
}

func (s *StaplerHandler) Protocol(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendRequest := ProtocolRequest{}
	if err := json.Unmarshal(body, &sendRequest); err != nil {
		logger.L.Errorln("error while collecting records", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.App.Protocol(req.Context(), &sendRequest); err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
