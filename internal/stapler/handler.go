package stapler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"protocall/pkg/logger"
)

type Notifier interface {
	Send(context.Context, []Phrase, []User)
}

type StaplerHandler struct {
	*Stapler
	Notifier
}

func (s *StaplerHandler) Protocol(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	pRec := ProtocolRequest{}
	if err := json.Unmarshal(body, &pRec); err != nil {
		logger.L.Errorln("error while collecting records", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := s.Make(req.Context(), &pRec)
	if err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.Send(req.Context(), data, pRec.Users)

	res.WriteHeader(http.StatusNoContent)
}
