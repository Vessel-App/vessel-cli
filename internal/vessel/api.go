package vessel

import (
	"encoding/json"
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/logger"
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
		return nil, fmt.Errorf("http client error")
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

func vesselApiEndpoint() string {
	endpoint := strings.TrimRight(os.Getenv("VESSEL_API_ENDPOINT"), "/")

	if len(endpoint) > 0 {
		return endpoint
	}

	return "http://localhost:8000/api"
}
