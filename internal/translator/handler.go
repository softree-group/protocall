package translator

import (
	"encoding/json"
	"io"
	"net/http"

	"protocall/pkg/logger"
)

type TranslatorHandler struct {
	App *Translator
}

func (t *TranslatorHandler) Translate(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	translateRequest := TranslateRequest{}
	if err := json.Unmarshal(body, &translateRequest); err != nil {
		logger.L.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	t.App.Translate(&translateRequest)

	res.WriteHeader(http.StatusNoContent)
}
