package main

import (
	"net/http"

	"protocall/pkg/webcore"
)

func initRouter(mux *http.ServeMux, porter *porterHandler) {
	mux.HandleFunc("/records", webcore.ApplyMethods(porter.serveHTTP, "POST"))
}
