package stapler

import (
	"net/http"

	"protocall/pkg/webcore"
)

func InitRouter(mux *http.ServeMux, h *StaplerHandler) {
	mux.HandleFunc("/protocols", webcore.ApplyMethods(h.Protocol, "POST"))
}
