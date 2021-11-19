package main

import (
	"context"
	"net/http"
	"path/filepath"

	"protocall/pkg/logger"
)

type storage interface {
	PutFile(context.Context, string, string) error
}

type porterHandler struct {
	storage
	root string
}

func (ph *porterHandler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	if from == "" {
		logger.L.Error("empty query parameter: from")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		logger.L.Error("empty query parameter: to")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ph.PutFile(r.Context(), filepath.Join(ph.root, from), filepath.Join(to)); err != nil {
		logger.L.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
