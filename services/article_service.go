package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	origin  string
}

func GetArticleService(baseURL string) *ArticleService {
	return &ArticleService{
		baseURL: baseURL,
		origin:  os.Getenv("SERVICE_ORIGIN"),
	}
}

type ArticleUpdateResponse struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

func (s *ArticleService) UpdateArticlesState(articleIDs []string, status ArticleStatus, token string) ([]ArticleUpdateResponse, error) {
	url := fmt.Sprintf("%sapi/protected/articles/status", s.baseURL)
	log.Printf("Making request to: %s", url)
	log.Printf("Updating articles: %v with status: %s", articleIDs, status)
	
	request := StatusUpdateRequest{
		ArticleIDs: articleIDs,
		Status:     string(status),
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}
	log.Printf("Request payload: %s", string(jsonData))

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	if s.origin != "" {
		req.Header.Set("Origin", s.origin)
	}
	log.Printf("Request headers set - Auth token present: %v", token != "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Response status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Error response body: %s", string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var articles []ArticleUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, fmt.Errorf("error decoding response: %v", err)
	}
	log.Printf("Successfully updated articles: %+v", articles)

	return articles, nil
}
