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

type ArticleUpdateResponse struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

func (s *ArticleService) UpdateArticlesState(articleIDs []string, status ArticleStatus, token string) ([]ArticleUpdateResponse, error) {
	url := fmt.Sprintf("%s/api/protected/articles/status", s.baseURL)
	
	request := StatusUpdateRequest{
		ArticleIDs: articleIDs,
		Status:     string(status),
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var articles []ArticleUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return articles, nil
}
