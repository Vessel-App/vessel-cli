package fly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

/*****************
 * CREATE APP
****************/

type CreateAppRequest struct {
	AppName string
	OrgSlug string
}

func (r *CreateAppRequest) ToRequest(token string) (*http.Request, error) {
	data := []byte(fmt.Sprintf(`{"app_name": "%s", "org_slug": "%s"}`, r.AppName, r.OrgSlug))

	// TODO: Decide on url to use (vpn vs proxy)
	req, err := http.NewRequest(http.MethodPost, "http://_api.internal:4280/v1/apps", bytes.NewBuffer(data))

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func CreateApp(token, name, org string) (*App, error) {
	req := &CreateAppRequest{
		AppName: name,
		OrgSlug: org,
	}

	_, err := DoRequest(token, req)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return &App{
		AppName: req.AppName,
		Organization: Organization{
			Slug: req.OrgSlug,
		},
	}, nil
}

/*****************
 * GET APP
****************/

type GetAppRequest struct {
	AppName string
}

func (r *GetAppRequest) ToRequest(token string) (*http.Request, error) {
	// TODO: Decide on url to use (vpn vs proxy)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://_api.internal:4280/v1/apps/%s", r.AppName), nil)

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func GetApp(token, name string) (*App, error) {
	req := &GetAppRequest{
		AppName: name,
	}

	responseBody, err := DoRequest(token, req)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	a := &App{}
	err = json.Unmarshal(responseBody, a)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return a, nil
}

/*****************
 * DELETE APP
****************/

type DeleteAppRequest struct {
	AppName string
}

func (r *DeleteAppRequest) ToRequest(token string) (*http.Request, error) {
	// TODO: Decide on url to use (vpn vs proxy)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://_api.internal:4280/v1/apps/%s", r.AppName), nil)

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func DeleteApp(token, name string) error {
	req := &DeleteAppRequest{
		AppName: name,
	}

	_, err := DoRequest(token, req)

	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	return nil
}
