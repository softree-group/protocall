package main

import (
	"context"
	"net/http"
	"net/url"
	"path/filepath"

	"protocall/pkg/logger"
)

type storage interface {
	PutFile(context.Context, string, string) error
	GetLink(ctx context.Context, path string) (*url.URL, error)
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

	remote := filepath.Join(to)
	if err := ph.PutFile(r.Context(), filepath.Join(ph.root, from), remote); err != nil {
		logger.L.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	link, err := ph.GetLink(r.Context(), remote)
	if err != nil {
		logger.L.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := link.MarshalBinary()
	if err != nil {
		logger.L.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
