package main

import (
	"net/http"

	"protocall/internal/clerk/stapler"
	"protocall/internal/clerk/translator"
	"protocall/pkg/web"
)

type ClerkHandler struct {
	Stapler    stapler.StaplerHandler
	Translator translator.TranslatorHandler
}

func NewRouter(h *ClerkHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/protocols", web.ApplyMethods(h.Stapler.CreateProtocol, "POST"))
	mux.HandleFunc("/records", web.ApplyMethods(h.Translator.Translate, "POST"))
	return mux
}
