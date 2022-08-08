package fly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NearestRegion struct {
	NearestRegion Region `json:"nearestRegion"`
}

type GetNearestRegionRequest struct{}

func (r *GetNearestRegionRequest) ToRequest(token string) (*http.Request, error) {
	query := []byte(`{"query": "query { nearestRegion { code name gatewayAvailable } }"}`)

	req, err := http.NewRequest(http.MethodPost, "https://api.fly.io/graphql", bytes.NewBuffer(query))

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func GetNearestRegion(token string) (*NearestRegion, error) {
	req := &GetNearestRegionRequest{}

	responseBody, err := DoRequest("", req)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	r := &NearestRegion{}
	gr := &GraphResponse{
		Data: r,
	}

	err = json.Unmarshal(responseBody, gr)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return r, nil
}
