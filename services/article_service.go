package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ArticleStatus string

const (
	ArticleStatusAvailable   ArticleStatus = "AVAILABLE"
	ArticleStatusUnavailable ArticleStatus = "UNAVAILABLE"
)

type StatusUpdateRequest struct {
	ArticleIDs []string `json:"articleIds"`
	Status     string   `json:"status"`
}

type ArticleService struct {
	baseURL string
}

func GetArticleService(baseURL string) *ArticleService {
	return &ArticleService{
		baseURL: baseURL,
	}
}

func (s *ArticleService) UpdateArticlesState(articleIDs []string, status ArticleStatus, clerkToken string) error {
	url := fmt.Sprintf("%s/api/protected/articles/batch", s.baseURL)
	
	// Build request matching the expected structure
	request := StatusUpdateRequest{
		ArticleIDs: articleIDs,
		Status:     string(status),  // Convert ArticleStatus to string
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", clerkToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
} 