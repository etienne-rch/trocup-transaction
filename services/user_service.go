package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
    baseURL string
}

type TransactionForUserRequest struct {
    UserA         string             `json:"userA"`
    ArticleA      primitive.ObjectID `json:"articleA,omitempty"`
    ArticlePriceA float64           `json:"articlePriceA,omitempty"`
    UserB         string             `json:"userB"`
    ArticleB      primitive.ObjectID `json:"articleB"`
    ArticlePriceB float64           `json:"articlePriceB"`
    ClerkToken    string            `json:"-"`
}

// GetUserService returns a new instance of UserService with the provided base URL
func GetUserService(baseURL string) *UserService {
    return &UserService{
        baseURL: baseURL,
    }
}

// UpdateUsersData sends a request to update users data
func (s *UserService) UpdateUsersData(request TransactionForUserRequest) error {
    payload := map[string]interface{}{
        "userA": request.UserA,
        "userB": request.UserB,
        "articleB": request.ArticleB,
        "articlePriceB": request.ArticlePriceB,
    }

    // Only add articleA and articlePriceA if articleA exists
    if !request.ArticleA.IsZero() {
        payload["articleA"] = request.ArticleA
        payload["articlePriceA"] = request.ArticlePriceA
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("error marshaling payload: %v", err)
    }

    resp, err := http.NewRequest(
        http.MethodPatch,
        fmt.Sprintf("%susers/transactions", s.baseURL),
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return fmt.Errorf("error creating request: %v", err)
    }
    
    resp.Header.Set("Content-Type", "application/json")
    resp.Header.Set("Authorization", fmt.Sprintf("Bearer %s", request.ClerkToken))
    
    client := &http.Client{}
    response, err := client.Do(resp)
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