package porter

import (
	"context"
	"net/http"
	"path/filepath"

	"protocall/pkg/logger"
)

type Storage interface {
	PutFile(context.Context, string, string) error
}

type PorterHandler struct {
	Storage
	Root string
}

func (ph *PorterHandler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	if from == "" {
		logger.L.Error("empty query parameter: from")
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		logger.L.Error("empty query parameter: to")
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ph.PutFile(context.Background(), filepath.Join(ph.Root, from), filepath.Join(to)); err != nil {
		logger.L.Error(err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wr.WriteHeader(http.StatusNoContent)
}
