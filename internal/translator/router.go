package translator

import (
	"net/http"

	"protocall/pkg/web"
)

func InitRouter(mux *http.ServeMux, h *TranslatorHandler) {
	mux.HandleFunc("/translations", web.ApplyMethods(h.Translate, "POST"))
}
