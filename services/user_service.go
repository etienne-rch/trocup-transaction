package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type UserService struct {
    baseURL string
    origin  string
}

type ArticleOwnership struct {
    ArticleID string  `json:"articleId"`
    OwnerID   string  `json:"ownerId"`
    Price     float64 `json:"price"`
}

type TransactionForUserRequest struct {
    UserA          string  `json:"userA"`
    UserB          string  `json:"userB"`
    ArticleA       *string  `json:"articleA,omitempty"`
    ArticleB       string  `json:"articleB"`
    ArticlePriceA  *float64 `json:"articlePriceA,omitempty"`
    ArticlePriceB  float64 `json:"articlePriceB"`
}

// GetUserService returns a new instance of UserService with the provided base URL
func GetUserService(baseURL string) *UserService {
    return &UserService{
        baseURL: baseURL,
        origin:  os.Getenv("SERVICE_ORIGIN"),
    }
}

// UpdateUsersData sends a request to update users data
func (s *UserService) UpdateUsersData(request TransactionForUserRequest, token string) error {
    jsonData, err := json.Marshal(request)
    if err != nil {
        return fmt.Errorf("error marshaling payload: %v", err)
    }

    url := fmt.Sprintf("%susers/transactions", s.baseURL)
    req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", token)
    if s.origin != "" {
        req.Header.Set("Origin", s.origin)
    }

    client := &http.Client{}
    response, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error making request to user service: %v", err)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(response.Body)
        return fmt.Errorf("user service returned status code: %d, body: %s", response.StatusCode, string(body))
    }

    return nil
}

// HealthCheck verifies if the user service is available
func (s *UserService) HealthCheck() bool {
    resp, err := http.Get(fmt.Sprintf("%sapi/health", s.baseURL))
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    return resp.StatusCode == http.StatusOK
} 