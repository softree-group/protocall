package main

import (
	"net/http"
	"protocall/pkg/web"
)

func initRouter(mux *http.ServeMux, porter *porterHandler) {
	mux.HandleFunc("/records", web.ApplyMethods(porter.serveHTTP, "POST"))
}
