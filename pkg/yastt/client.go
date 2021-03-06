package yastt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Yastt struct {
	*http.Client
	*YasttConfig
}

func NewYastt(client *http.Client, config *YasttConfig) *Yastt {
	return &Yastt{client, config}
}

func (y *Yastt) updateJob(ctx context.Context, id string) (*RecognizerResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%v/operations/%v", y.OperationAddr, id),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Api-Key %v", y.SecretKey))

	resp, err := y.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status: %d %v", resp.StatusCode, string(body))
	}

	result := &RecognizerResponse{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (y *Yastt) pool(ctx context.Context, id string, length time.Duration) (*RecognizerResponse, error) {
	interval := time.Duration(float64(time.Second) * length.Seconds() * y.PoolCoefficient)
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			resp, err := y.updateJob(ctx, id)
			if err != nil {
				return nil, err
			}
			if !resp.Done {
				time.Sleep(interval)
				interval /= 2
				continue
			}

			return resp, nil
		}
	}
}

func (y *Yastt) createJob(ctx context.Context, data *RecognizerRequest) (*RecognizerResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%v/speech/stt/v2/longRunningRecognize", y.TranscribeAddr),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Api-Key %v", y.SecretKey))

	resp, err := y.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status: %d %v", resp.StatusCode, string(body))
	}

	result := &RecognizerResponse{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result, nil
}

// Recognize is sending requests to yandex cloud sst api.
// Pool interval calculate from track length and coefficient from config.
func (y *Yastt) Recognize(ctx context.Context, filename string, length time.Duration) (<-chan Chunk, <-chan error) {
	out := make(chan Chunk)
	outErr := make(chan error)

	go func() {
		defer close(out)
		defer close(outErr)

		resp, err := y.createJob(ctx, &RecognizerRequest{
			Config: Config{y.Specification},
			Audio:  Audio{filename},
		})
		if err != nil {
			outErr <- err
			return
		}

		if !resp.Done {
			resp, err = y.pool(ctx, resp.ID, length)
			if err != nil {
				outErr <- err
				return
			}
		}

		for _, chunk := range resp.Response.Chunks {
			out <- chunk
		}
		outErr <- nil
	}()
	return out, outErr
}
