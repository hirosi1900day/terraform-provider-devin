package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	// Base URL for Devin API
	baseURL = "https://api.devin.ai/v1"
)

// DevinClient is a client for interacting with the Devin API
type DevinClient struct {
	APIKey     string
	HTTPClient *http.Client

	// Cache for knowledge list to avoid rate limiting
	knowledgeCache     *ListKnowledgeResponse
	knowledgeCacheMu   sync.RWMutex
	knowledgeCacheTime time.Time
	// Cache TTL (default: 5 minutes for terraform plan duration)
	CacheTTL time.Duration
}

// Knowledge represents a Devin knowledge resource
type Knowledge struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Body               string    `json:"body"`                       // Required
	TriggerDescription string    `json:"trigger_description"`        // Required
	ParentFolderID     string    `json:"parent_folder_id,omitempty"` // Optional
	CreatedAt          time.Time `json:"created_at"`
}

// ListKnowledgeResponse represents the response from the knowledge list API
type ListKnowledgeResponse struct {
	Knowledge []KnowledgeItem `json:"knowledge"`
	Folders   []FolderItem    `json:"folders"`
}

// KnowledgeItem represents a knowledge item
type KnowledgeItem struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Body               string    `json:"body"`                       // Required
	TriggerDescription string    `json:"trigger_description"`        // Required
	ParentFolderID     string    `json:"parent_folder_id,omitempty"` // Optional
	CreatedAt          time.Time `json:"created_at"`
}

// FolderItem represents a folder item
type FolderItem struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateKnowledgeRequest represents the request for knowledge creation API
type CreateKnowledgeRequest struct {
	Name               string `json:"name"`                       // Required
	Body               string `json:"body"`                       // Required
	ParentFolderID     string `json:"parent_folder_id,omitempty"` // Optional
	TriggerDescription string `json:"trigger_description"`        // Required
}

// UpdateKnowledgeRequest represents the request for knowledge update API
type UpdateKnowledgeRequest struct {
	Name               string `json:"name"`                       // Required
	Body               string `json:"body"`                       // Required
	ParentFolderID     string `json:"parent_folder_id,omitempty"` // Optional
	TriggerDescription string `json:"trigger_description"`        // Required
}

// ErrorResponse represents the API error response
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// NewClient creates a new DevinClient
func NewClient(apiKey string) *DevinClient {
	return &DevinClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		CacheTTL: 5 * time.Minute, // Default cache TTL
	}
}

// InvalidateCache clears the knowledge cache
func (c *DevinClient) InvalidateCache() {
	c.knowledgeCacheMu.Lock()
	defer c.knowledgeCacheMu.Unlock()
	c.knowledgeCache = nil
	c.knowledgeCacheTime = time.Time{}
}

// isCacheValid checks if the cache is still valid
func (c *DevinClient) isCacheValid() bool {
	if c.knowledgeCache == nil {
		return false
	}
	return time.Since(c.knowledgeCacheTime) < c.CacheTTL
}

// sendRequest is a common function for sending requests
func (c *DevinClient) sendRequest(method, path string, body interface{}) ([]byte, error) {
	url := baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON encode request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("API error: %s (%s)", errResp.Error.Message, errResp.Error.Type)
		}
		return nil, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	return respBody, nil
}

// ListKnowledge retrieves a list of knowledge resources
// Results are cached to avoid rate limiting during terraform plan/apply
func (c *DevinClient) ListKnowledge() (*ListKnowledgeResponse, error) {
	// Return mock data for demo (development/testing)
	if IsMockClient(c.APIKey) {
		return GetMockKnowledgeList(), nil
	}

	// Check cache first
	c.knowledgeCacheMu.RLock()
	if c.isCacheValid() {
		cached := c.knowledgeCache
		c.knowledgeCacheMu.RUnlock()
		return cached, nil
	}
	c.knowledgeCacheMu.RUnlock()

	// Cache miss or expired, fetch from API
	c.knowledgeCacheMu.Lock()
	defer c.knowledgeCacheMu.Unlock()

	// Double-check after acquiring write lock (another goroutine might have updated)
	if c.isCacheValid() {
		return c.knowledgeCache, nil
	}

	// Normal processing
	respBody, err := c.sendRequest("GET", "/knowledge", nil)
	if err != nil {
		return nil, err
	}

	var response ListKnowledgeResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Update cache
	c.knowledgeCache = &response
	c.knowledgeCacheTime = time.Now()

	return &response, nil
}

// GetKnowledge retrieves a knowledge resource by ID
// Note: Currently, the Devin API does not explicitly expose a dedicated endpoint
// for retrieving individual knowledge resources, so we use the List API to extract
// a specific knowledge resource by ID
func (c *DevinClient) GetKnowledge(id string) (*Knowledge, error) {
	// Return mock data for demo (development/testing)
	if IsMockClient(c.APIKey) {
		return GetMockKnowledge(id)
	}

	// Normal processing
	// Use the list API to get all knowledge resources
	response, err := c.ListKnowledge()
	if err != nil {
		return nil, fmt.Errorf("error occurred while retrieving knowledge list: %w", err)
	}

	// Search for knowledge resource with matching ID
	for _, item := range response.Knowledge {
		if item.ID == id {
			return &Knowledge{
				ID:                 item.ID,
				Name:               item.Name,
				Body:               item.Body,
				TriggerDescription: item.TriggerDescription,
				ParentFolderID:     item.ParentFolderID,
				CreatedAt:          item.CreatedAt,
			}, nil
		}
	}

	return nil, fmt.Errorf("knowledge resource with ID '%s' not found", id)
}

// CreateKnowledge creates a new knowledge resource
func (c *DevinClient) CreateKnowledge(name, body string, triggerDescription string, parentFolderID string) (*Knowledge, error) {
	// Return mock data for demo (development/testing)
	if IsMockClient(c.APIKey) {
		return CreateMockKnowledge(name, body, triggerDescription, parentFolderID), nil
	}

	// Normal processing
	reqBody := CreateKnowledgeRequest{
		Name:               name,
		Body:               body,
		TriggerDescription: triggerDescription,
		ParentFolderID:     parentFolderID,
	}

	respBody, err := c.sendRequest("POST", "/knowledge", reqBody)
	if err != nil {
		return nil, err
	}

	var knowledge Knowledge
	if err := json.Unmarshal(respBody, &knowledge); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Invalidate cache after creating a new knowledge
	c.InvalidateCache()

	return &knowledge, nil
}

// UpdateKnowledge updates a knowledge resource
func (c *DevinClient) UpdateKnowledge(id, name, body string, triggerDescription string, parentFolderID string) (*Knowledge, error) {
	// Return mock data for demo (development/testing)
	if IsMockClient(c.APIKey) {
		return UpdateMockKnowledge(id, name, body, triggerDescription, parentFolderID), nil
	}

	// Normal processing
	reqBody := UpdateKnowledgeRequest{
		Name:               name,
		Body:               body,
		TriggerDescription: triggerDescription,
		ParentFolderID:     parentFolderID,
	}

	path := fmt.Sprintf("/knowledge/%s", id)
	respBody, err := c.sendRequest("PUT", path, reqBody)
	if err != nil {
		return nil, err
	}

	var knowledge Knowledge
	if err := json.Unmarshal(respBody, &knowledge); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Invalidate cache after updating knowledge
	c.InvalidateCache()

	return &knowledge, nil
}

// DeleteKnowledge deletes a knowledge resource
func (c *DevinClient) DeleteKnowledge(id string) error {
	// Return mock data for demo (development/testing)
	if IsMockClient(c.APIKey) {
		return nil
	}

	// Normal processing
	path := fmt.Sprintf("/knowledge/%s", id)
	_, err := c.sendRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	// Invalidate cache after deleting knowledge
	c.InvalidateCache()

	return nil
}

// GetFolderByID retrieves a folder resource by ID
// Note: Currently, the Devin API does not explicitly expose a dedicated endpoint
// for retrieving individual folder resources, so we use the List API to extract
// a specific folder resource by ID
func (c *DevinClient) GetFolderByID(id string) (*FolderItem, error) {
	// Return mock data for demo (development/testing)
	if IsMockClient(c.APIKey) {
		return GetMockFolderByID(id)
	}

	// Normal processing
	// Use the list API to get all knowledge resources (which includes folders)
	response, err := c.ListKnowledge()
	if err != nil {
		return nil, fmt.Errorf("error occurred while retrieving folder list: %w", err)
	}

	// Search for folder resource with matching ID
	for _, item := range response.Folders {
		if item.ID == id {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("folder resource with ID '%s' not found", id)
}

// GetFolderByName retrieves a folder resource by name
// Note: Currently, the Devin API does not explicitly expose a dedicated endpoint
// for retrieving individual folder resources, so we use the List API to extract
// a specific folder resource by name
func (c *DevinClient) GetFolderByName(name string) (*FolderItem, error) {
	// Return mock data for demo (development/testing)
	if IsMockClient(c.APIKey) {
		return GetMockFolderByName(name)
	}

	// Normal processing
	// Use the list API to get all knowledge resources (which includes folders)
	response, err := c.ListKnowledge()
	if err != nil {
		return nil, fmt.Errorf("error occurred while retrieving folder list: %w", err)
	}

	// Search for folder resource with matching name
	for _, item := range response.Folders {
		if item.Name == name {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("folder resource with name '%s' not found", name)
}
