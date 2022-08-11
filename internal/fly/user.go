package fly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetUserRequest struct{}

func (r *GetUserRequest) ToRequest(token string) (*http.Request, error) {
	query := []byte(`{"query": "query {currentUser {email} organizations {nodes{id slug name type viewerRole}}}"}`)
	req, err := http.NewRequest(http.MethodPost, "https://api.fly.io/graphql", bytes.NewBuffer(query))

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func GetUser(token string) (*User, error) {
	req := &GetUserRequest{}

	responseBody, err := DoRequest(token, req)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	u := &User{}
	gr := &GraphResponse{
		Data: u,
	}

	err = json.Unmarshal(responseBody, gr)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return u, nil
}
