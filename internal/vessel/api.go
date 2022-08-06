package vessel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Teams []Team `json:"all_teams"`
}

type Team struct {
	Id   int    `json:"id"`
	Guid string `json:"guid"`
	Name string `json:"name"`
}

type CreateEnvironmentRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
	Region    string `json:"region"`
}

type Environment struct {
	Id          uint64 `json:"id"`
	ProviderId  string `json:"provider_id,omitempty"`
	Name        string `json:"name"`
	Size        string `json:"size"`
	PublicKey   string `json:"public_key"`
	Region      string `json:"region"`
	IpAddress   string `json:"ip_address,omitempty"`
	Status      string `json:"status"`
	Initialized bool   `json:"initialized"`
}

func GetUser(token string) (*User, error) {
	url := fmt.Sprintf("%s/user", vesselApiEndpoint())

	client := &http.Client{
		Timeout: time.Second * 2,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	logger.GetLogger().Debug("caller", "api::GetUser", "msg", "about to call Vessel API", "url", url, "http_method", "GET")

	r, err := client.Do(req)

	if err != nil {
		// 500 errors go here
		return nil, fmt.Errorf("http client error: %w", err)
	}

	defer r.Body.Close()

	if r.StatusCode > 299 {
		return nil, fmt.Errorf("invalid user request: %w", err)
	}

	user := &User{}
	err = json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return user, nil
}

func CreateEnvironment(team, name, publicKey, region, token string) (*Environment, error) {
	url := fmt.Sprintf("%s/team/%s/environment", vesselApiEndpoint(), team)

	client := &http.Client{
		//Timeout: time.Second * 2,
		Timeout: time.Minute * 5, // For local dev while queue is "sync"
	}

	environmentRequest := CreateEnvironmentRequest{
		Name:      name,
		PublicKey: publicKey,
		Region:    region,
	}
	body, err := json.Marshal(environmentRequest)

	if err != nil {
		return nil, fmt.Errorf("create environment json error")
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logger.GetLogger().Debug("caller", "api::CreateEnvironment", "msg", "about to call Vessel API", "url", url, "http_method", "POST", "body", string(body))

	r, err := client.Do(req)

	if err != nil {
		// 500 errors go here
		return nil, fmt.Errorf("http client error: %w", err)
	}

	defer r.Body.Close()

	if r.StatusCode > 299 {
		b, _ := ioutil.ReadAll(r.Body)
		logger.GetLogger().Debug("caller", "api::CreateEnvironment", "msg", "http request error", "status", r.StatusCode, "body", string(b))
		return nil, fmt.Errorf("invalid create environment request: %w", err)
	}

	env := &Environment{}
	err = json.NewDecoder(r.Body).Decode(env)

	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return env, nil
}

// WaitForEnvironment polls the API and waits for the environment to be ready to use
//  before finally returning it
func WaitForEnvironment(team string, machine uint64, token string) (*Environment, error) {
	// Total wait time ~5 minutes (should only need a minute or 2)
	ticker := time.NewTicker(2 * time.Second)
	totalAttempts := 0
	attemptsAllowed := 150

	for {
		select {
		case <-ticker.C:
			e, err := GetEnvironment(team, machine, token)
			if err != nil {
				ticker.Stop()
				return nil, fmt.Errorf("could not get environment: %w", err)
			}

			if e.Initialized {
				ticker.Stop()
				return e, nil
			}

			totalAttempts++
			if totalAttempts >= attemptsAllowed {
				ticker.Stop()
				return nil, fmt.Errorf("too many get environment attempts")
			}
		}
	}
}

// GetEnvironment retrieves a development environment
func GetEnvironment(team string, machine uint64, token string) (*Environment, error) {
	url := fmt.Sprintf("%s/team/%s/environment/%d", vesselApiEndpoint(), team, machine)

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	logger.GetLogger().Debug("caller", "api::GetUser", "msg", "about to call Vessel API", "url", url, "http_method", "GET")

	r, err := client.Do(req)

	if err != nil {
		// 500 errors go here
		return nil, fmt.Errorf("http client error: %w", err)
	}

	defer r.Body.Close()

	if r.StatusCode > 299 {
		return nil, fmt.Errorf("invalid user request: %w", err)
	}

	environment := &Environment{}
	err = json.NewDecoder(r.Body).Decode(environment)

	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return environment, nil
}

func vesselApiEndpoint() string {
	endpoint := strings.TrimRight(os.Getenv("VESSEL_API_ENDPOINT"), "/")

	if len(endpoint) > 0 {
		return endpoint
	}

	return "http://localhost:8888/api"
}
