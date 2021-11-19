package translator

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"protocall/pkg/logger"
)

type TranslateRequest struct {
	StartTime time.Time `json:"start" binding:"required"`
	User      struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email"`
		Path     string `json:"path" binding:"required"`
	} `json:"user" binding:"required"`
}

type TranslatorRepository interface {
	Translate(*TranslateRequest)
}

type TranslatorHandler struct {
	app TranslatorRepository
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

	t.app.Translate(&translateRequest)

	res.WriteHeader(http.StatusNoContent)
}
