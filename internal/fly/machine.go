package fly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RunMachineRequest struct {
	App    string
	Region string
	Image  string
	Env    map[string]string
}

func (m *RunMachineRequest) ToRequest(token string) (*http.Request, error) {
	env := ""
	for k, v := range m.Env {
		env += fmt.Sprintf(`"%s: %s,"`, k, v)
	}
	data := []byte(fmt.Sprintf(`{"region": "%s", "config": {"image": "%s", "env": {%s}}}`, m.Region, m.Image, env))

	// TODO: Decide on url to use (vpn vs proxy)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://_api.internal:4280/v1/apps/%s", m.App), bytes.NewBuffer(data))

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func RunMachine(token, app, region, image, pubKey string) (*Machine, error) {
	e := make(map[string]string)
	e["VESSEL_PUBLIC_KEY"] = pubKey

	req := &RunMachineRequest{
		App:    app,
		Region: region,
		Image:  image,
		Env:    e,
	}

	responseBody, err := DoRequest(token, req)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	m := &Machine{}
	err = json.Unmarshal(responseBody, m)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return m, nil
}

type ListMachinesRequest struct {
	App string
}

type ListMachinesResponse struct {
	Machines []Machine
}

func (m *ListMachinesRequest) ToRequest(token string) (*http.Request, error) {
	// TODO: Decide on url to use (vpn vs proxy)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://_api.internal:4280/v1/apps/%s", m.App), nil)

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func ListMachines(token, app string) (*ListMachinesResponse, error) {

	req := &ListMachinesRequest{
		App: app,
	}

	responseBody, err := DoRequest(token, req)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	m := &ListMachinesResponse{}
	err = json.Unmarshal(responseBody, m)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return m, nil
}

type GetMachineRequest struct {
	App     string
	Machine string
}

func (m *GetMachineRequest) ToRequest(token string) (*http.Request, error) {
	// TODO: Decide on url to use (vpn vs proxy)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://_api.internal:4280/v1/apps/%s/machines/%s", m.App, m.Machine), nil)

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func GetMachine(token, app, machine string) (*Machine, error) {

	req := &GetMachineRequest{
		App:     app,
		Machine: machine,
	}

	responseBody, err := DoRequest(token, req)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	m := &Machine{}
	err = json.Unmarshal(responseBody, m)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return m, nil
}

type StartMachineRequest struct {
	App     string
	Machine string
}

func (m *StartMachineRequest) ToRequest(token string) (*http.Request, error) {
	// TODO: Decide on url to use (vpn vs proxy)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://_api.internal:4280/v1/apps/%s/machines/%s/start", m.App, m.Machine), nil)

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func StartMachine(token, app, machine string) error {
	req := &StartMachineRequest{
		App:     app,
		Machine: machine,
	}

	_, err := DoRequest(token, req)

	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	return nil
}

type StopMachineRequest struct {
	App     string
	Machine string
}

func (m *StopMachineRequest) ToRequest(token string) (*http.Request, error) {
	// TODO: Decide on url to use (vpn vs proxy)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://_api.internal:4280/v1/apps/%s/machines/%s/stop", m.App, m.Machine), nil)

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func StopMachine(token, app, machine string) error {
	req := &StopMachineRequest{
		App:     app,
		Machine: machine,
	}

	_, err := DoRequest(token, req)

	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	return nil
}

type DeleteMachineRequest struct {
	App     string
	Machine string
}

func (m *DeleteMachineRequest) ToRequest(token string) (*http.Request, error) {
	// TODO: Decide on url to use (vpn vs proxy)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://_api.internal:4280/v1/apps/%s/machines/%s", m.App, m.Machine), nil)

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func DeleteMachine(token, app, machine string) error {
	req := &DeleteMachineRequest{
		App:     app,
		Machine: machine,
	}

	_, err := DoRequest(token, req)

	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	return nil
}

func WaitForMachine(token, app, machine string) error {
	// Total wait time ~5 minutes (should only need a minute or 2)
	ticker := time.NewTicker(2 * time.Second)
	totalAttempts := 0
	attemptsAllowed := 150

	for {
		select {
		case <-ticker.C:
			e, err := GetMachine(token, app, machine)
			if err != nil {
				ticker.Stop()
				return fmt.Errorf("could not get machine: %w", err)
			}

			if e.IsInitialized() {
				ticker.Stop()
				return nil
			}

			totalAttempts++
			if totalAttempts >= attemptsAllowed {
				ticker.Stop()
				return fmt.Errorf("too many get machine attempts")
			}
		}
	}
}
