package porter

import (
	"net/http"
	"protocall/pkg/web"
)

func NewRouter(porter *PorterHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/records", web.ApplyMethods(porter.ServeHTTP, "POST"))
	return mux
}
