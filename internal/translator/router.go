package translator

import (
	"net/http"

	"protocall/pkg/webcore"
)

func InitRouter(mux *http.ServeMux, h *TranslatorHandler) {
	mux.HandleFunc("/translations", webcore.ApplyMethods(h.Translate, "POST"))
}
