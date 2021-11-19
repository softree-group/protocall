package stapler

import (
	"net/http"

	"protocall/pkg/web"
)

func InitRouter(mux *http.ServeMux, h *StaplerHandler) {
	mux.HandleFunc("/protocols", web.ApplyMethods(h.Protocol, "POST"))
}
