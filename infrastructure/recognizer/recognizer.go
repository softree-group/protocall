package recognizer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	"github.com/kirito41dd/xslice"

	"protocall/domain/entity"
	"protocall/domain/repository"
)

type Config struct {
	Host string
	Port string
}

type recognizer struct {
	config *Config
}

func NewRecognizer(c *Config) repository.VoiceRecognizer {
	return &SpechToText{
		config: c,
	}
}

const chunkSize = 8000

func (stt *SpechToText) Recognize(ctx context.Context, audio io.ReadCloser) (*entity.Message, error) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%v:%v"), nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	for _, el := range xslice.SplitToChunks(audio, chunkSize).([][]byte) {
		if err := conn.WriteMessage(websocket.BinaryMessage, el); err != nil {
			return nil, err
		}
		if _, _, err = conn.ReadMessage(); err != nil {
			return nil, err
		}
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}")); err != nil {
		return nil, err
	}

	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		return nil, err
	}

	var res entity.Message
	if err = json.Unmarshal(msg, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
