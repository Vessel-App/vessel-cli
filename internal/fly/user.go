package fly

import (
	"bytes"
	"fmt"
	"net/http"
)

type GetUserRequest struct{}

func (r *GetUserRequest) ToRequest(token string) (*http.Request, error) {
	query := []byte(`
{
	"query": "query {
		currentUser {email}
		personalOrganization {id slug name type viewerRole}
		organizations {id slug name type viewerRole}
	}"
}
`)
	req, err := http.NewRequest("POST", "https://api.fly.io/graphql", bytes.NewBuffer(query))

	if err != nil {
		return nil, fmt.Errorf("could get create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}
