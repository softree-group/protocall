package recognizer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"protocall/pkg/logger"

	"github.com/gorilla/websocket"
)

var (
	errRecv         = errors.New("error while receiving response from recognizer")
	errReadFromFile = errors.New("error while reading from file")
)

type RecognizerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Recognizer struct {
	config *RecognizerConfig
	conn   *websocket.Conn
}

func NewRecognizer(ctx context.Context, c *RecognizerConfig) (*Recognizer, error) {
	return &Recognizer{
		config: c,
	}, nil
}

func (r *Recognizer) open(ctx context.Context) (err error) {
	r.conn, _, err = websocket.DefaultDialer.DialContext(
		ctx,
		fmt.Sprintf("ws://%v:%v", r.config.Host, r.config.Port),
		nil,
	)

	return
}

// TODO better error handling
func (r *Recognizer) close() {
	r.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			"",
		),
	)
}

var EOF = `{"eof" : 1}`

type TextRespone struct {
	Result []struct {
		Conf  float64 `json:"conf" binding:"required"`
		Start float64 `json:"start" binding:"required"`
		End   float64 `json:"end" binding:"required"`
		Word  string  `json:"word" binding:"required"`
	} `json:"result" binding:"required"`
	Text string `json:"text" binding:"required"`
}

const (
	chunkSize = 8000
)

func (r *Recognizer) Recognize(ctx context.Context, input io.Reader) <-chan TextRespone {
	output := make(chan TextRespone)

	go func() {
		if err := r.open(ctx); err != nil {
			logger.L.Error(err)
			return
		}
		defer r.close()
		defer close(output)

		for {
			select {
			case <-ctx.Done():
				r.conn.WriteMessage(websocket.TextMessage, []byte(EOF))
				return

			default:
				buf := make([]byte, chunkSize)
				if n, err := input.Read(buf); err != nil {
					switch {
					case err == io.EOF && n == 0:
						r.conn.WriteMessage(websocket.TextMessage, []byte(EOF))
						return
					case err != io.EOF:
						logger.L.Errorf("%w: %v", errReadFromFile, err)
						return
					}
				}

				if err := r.conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {
					logger.L.Errorln(err)
					return
				}

				_, data, err := r.conn.ReadMessage()
				if err != nil {
					logger.L.Errorf("%w: %v", errRecv, err)
					return
				}

				content := TextRespone{}
				if err := json.Unmarshal(data, &content); err != nil {
					logger.L.Errorln(err)
					return
				}

				output <- content
			}
		}
	}()

	return output
}
