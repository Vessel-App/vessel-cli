package fly

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var flyApiHost string

type GraphResponse struct {
	Data interface{} `json:"data"`
}

type FlyRequest interface {
	ToRequest(token string) (*http.Request, error)
}

func init() {
	// Use "_api.internal" if connected to Fly's VPN
	flyHost := os.Getenv("FLY_HOST")

	if len(flyHost) > 0 {
		flyApiHost = flyHost
		return
	}

	// Else we assume the use of
	// "fly machines api-proxy"
	flyApiHost = "127.0.0.1"
}

func DoRequest(token string, r FlyRequest) ([]byte, error) {
	req, err := r.ToRequest(token)

	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	client := &http.Client{
		Timeout: time.Second * 6,
	}

	result, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http client error: %w", err)
	}

	defer result.Body.Close()

	if result.StatusCode > 299 {
		// todo: Log this raw output
		body, _ := io.ReadAll(result.Body)
		return nil, fmt.Errorf("invalid request: status=%d, body=%s", result.StatusCode, string(body))
	}

	return io.ReadAll(result.Body)
}
