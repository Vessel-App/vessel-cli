package fly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*****************
 * ALLOCATE IP
****************/

type IpAddressAllocationResponse struct {
	Allocation IpAddressAllocation `json:"allocateIpAddress"`
}

type GetAppIpResponse struct {
	App App `json:"app"`
}

type AllocateIpRequest struct {
	App string
	V6  bool
}

func (i *AllocateIpRequest) ToRequest(token string) (*http.Request, error) {
	ipType := "v6"
	if i.V6 == false {
		ipType = "v4"
	}

	query := []byte(fmt.Sprintf(strings.Replace(`{
		"query": "mutation($input: AllocateIPAddressInput!) { allocateIpAddress(input: $input) { ipAddress { id address type region createdAt } } }",
		"variables": { "input": { "appId": "%s", "type": "%s" } }
	}`, "\n", "", -1), i.App, ipType))

	req, err := http.NewRequest(http.MethodPost, "https://api.fly.io/graphql", bytes.NewBuffer(query))

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func AllocateIp(token, app string, useV6 bool) (*IpAddressAllocation, error) {
	req := &AllocateIpRequest{
		App: app,
		V6:  useV6,
	}

	responseBody, err := DoRequest(token, req)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	i := &IpAddressAllocationResponse{}
	err = json.Unmarshal(responseBody, i)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return &i.Allocation, nil
}

/*****************
 * GET APP IP
****************/

type GetAppIpRequest struct {
	App string
}

func (i *GetAppIpRequest) ToRequest(token string) (*http.Request, error) {
	query := []byte(fmt.Sprintf(strings.Replace(`{
		"query": "query ($appName: String!) { app(name: $appName) { ipAddresses { nodes {id address type region createdAt } } } }",
		"variables": { "appName": "%s" }
	}`, "\n", "", -1), i.App))

	req, err := http.NewRequest(http.MethodPost, "https://api.fly.io/graphql", bytes.NewBuffer(query))

	if err != nil {
		return nil, fmt.Errorf("could not create http request object: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func GetAppIp(token, app string) (*IpAddress, error) {
	req := &GetAppIpRequest{
		App: app,
	}

	responseBody, err := DoRequest(token, req)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	a := &GetAppIpResponse{}
	err = json.Unmarshal(responseBody, a)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	if len(a.App.IpAddresses.Nodes) == 0 {
		return nil, fmt.Errorf("no IP addresses allocated to app: %s", app)
	}

	return &a.App.IpAddresses.Nodes[0], nil
}
