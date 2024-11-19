package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserService struct {
    baseURL string
}

type ArticleOwnership struct {
    ArticleID string  `json:"articleId"`
    OwnerID   string  `json:"ownerId"`
    Price     float64 `json:"price"`
}

type TransactionForUserRequest struct {
    UserA          string  `json:"userA"`
    UserB          string  `json:"userB"`
    ArticleA       string  `json:"articleA,omitempty"`
    ArticleB       string  `json:"articleB"`
    ArticlePriceA  float64 `json:"articlePriceA,omitempty"`
    ArticlePriceB  float64 `json:"articlePriceB"`
}

// GetUserService returns a new instance of UserService with the provided base URL
func GetUserService(baseURL string) *UserService {
    return &UserService{
        baseURL: baseURL,
    }
}

// UpdateUsersData sends a request to update users data
func (s *UserService) UpdateUsersData(request TransactionForUserRequest, token string) error {
    jsonData, err := json.Marshal(request)
    if err != nil {
        return fmt.Errorf("error marshaling payload: %v", err)
    }

    req, err := http.NewRequest(
        http.MethodPatch,
        fmt.Sprintf("%susers/transactions", s.baseURL),
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return fmt.Errorf("error creating request: %v", err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

    client := &http.Client{}
    response, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error making request to user service: %v", err)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return fmt.Errorf("user service returned status code: %d", response.StatusCode)
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