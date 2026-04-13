package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	// Base URL for Devin API v3
	baseAPIURL = "https://api.devin.ai"
)

// DevinClient is a client for interacting with the Devin API v3
type DevinClient struct {
	APIKey     string
	OrgID      string
	HTTPClient *http.Client

	// Rate limit retry settings
	MaxRetries int

	// Cache for knowledge notes to avoid rate limiting during plan/apply
	knowledgeCache     map[string]*KnowledgeNote
	knowledgeCacheMu   sync.RWMutex
	knowledgeCacheTime time.Time
	CacheTTL           time.Duration
}

// baseURL returns the organization-scoped base URL
func (c *DevinClient) baseURL() string {
	return fmt.Sprintf("%s/v3/organizations/%s", baseAPIURL, c.OrgID)
}

// --- Knowledge types (v3) ---

// KnowledgeNote represents a v3 knowledge note
type KnowledgeNote struct {
	NoteID     string  `json:"note_id"`
	Name       string  `json:"name"`
	Body       string  `json:"body"`
	Trigger    string  `json:"trigger"`
	FolderID   string  `json:"folder_id,omitempty"`
	FolderPath string  `json:"folder_path,omitempty"`
	IsEnabled  bool    `json:"is_enabled"`
	PinnedRepo *string `json:"pinned_repo"`
	Macro      string  `json:"macro,omitempty"`
	AccessType string  `json:"access_type,omitempty"`
	CreatedAt  float64 `json:"created_at"`
	UpdatedAt  float64 `json:"updated_at"`
}

// ListKnowledgeNotesResponse represents the paginated response from v3 knowledge notes API
type ListKnowledgeNotesResponse struct {
	Items       []KnowledgeNote `json:"items"`
	EndCursor   *string         `json:"end_cursor"`
	HasNextPage bool            `json:"has_next_page"`
}

// CreateKnowledgeNoteRequest represents the request for creating a knowledge note
type CreateKnowledgeNoteRequest struct {
	Name       string  `json:"name"`
	Body       string  `json:"body"`
	Trigger    string  `json:"trigger"`
	PinnedRepo *string `json:"pinned_repo,omitempty"`
}

// UpdateKnowledgeNoteRequest represents the request for updating a knowledge note
type UpdateKnowledgeNoteRequest struct {
	Name       string  `json:"name"`
	Body       string  `json:"body"`
	Trigger    string  `json:"trigger"`
	PinnedRepo *string `json:"pinned_repo,omitempty"`
}

// --- Folder types (v3) ---

// FolderItem represents a v3 folder
type FolderItem struct {
	FolderID       string `json:"folder_id"`
	Name           string `json:"name"`
	Path           string `json:"path,omitempty"`
	NoteCount      int    `json:"note_count"`
	ParentFolderID string `json:"parent_folder_id,omitempty"`
}

// ListFoldersResponse represents the response from v3 folders API
type ListFoldersResponse struct {
	Items       []FolderItem `json:"items"`
	EndCursor   *string      `json:"end_cursor"`
	HasNextPage bool         `json:"has_next_page"`
}

// --- Playbook types (v3) ---

// Playbook represents a v3 playbook
type Playbook struct {
	PlaybookID string  `json:"playbook_id"`
	Title      string  `json:"title"`
	Body       string  `json:"body"`
	AccessType string  `json:"access_type,omitempty"`
	Macro      *string `json:"macro"`
	OrgID      string  `json:"org_id,omitempty"`
	CreatedAt  float64 `json:"created_at"`
	UpdatedAt  float64 `json:"updated_at"`
	CreatedBy  string  `json:"created_by,omitempty"`
	UpdatedBy  string  `json:"updated_by,omitempty"`
}

// ListPlaybooksResponse represents the paginated response from playbooks API
type ListPlaybooksResponse struct {
	Items       []Playbook `json:"items"`
	EndCursor   *string    `json:"end_cursor"`
	HasNextPage bool       `json:"has_next_page"`
}

// CreatePlaybookRequest represents the request for creating a playbook
type CreatePlaybookRequest struct {
	Title string  `json:"title"`
	Body  string  `json:"body"`
	Macro *string `json:"macro,omitempty"`
}

// UpdatePlaybookRequest represents the request for updating a playbook
type UpdatePlaybookRequest struct {
	Title string  `json:"title"`
	Body  string  `json:"body"`
	Macro *string `json:"macro,omitempty"`
}

// --- Secret types (v3) ---

// Secret represents a v3 secret (read response - value is never returned)
type Secret struct {
	SecretID    string  `json:"secret_id"`
	Key         string  `json:"key"`
	SecretType  string  `json:"secret_type,omitempty"`
	IsSensitive bool    `json:"is_sensitive"`
	Note        string  `json:"note,omitempty"`
	AccessType  string  `json:"access_type,omitempty"`
	CreatedAt   float64 `json:"created_at"`
	UpdatedAt   float64 `json:"updated_at"`
	CreatedBy   string  `json:"created_by,omitempty"`
	UpdatedBy   string  `json:"updated_by,omitempty"`
}

// ListSecretsResponse represents the response from secrets API
type ListSecretsResponse struct {
	Items       []Secret `json:"items"`
	EndCursor   *string  `json:"end_cursor"`
	HasNextPage bool     `json:"has_next_page"`
}

// CreateSecretRequest represents the request for creating a secret
type CreateSecretRequest struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	IsSensitive bool   `json:"is_sensitive"`
	Note        string `json:"note,omitempty"`
}

// --- Schedule types (v3) ---

// SchedulePlaybookInfo represents the nested playbook info in a schedule response
type SchedulePlaybookInfo struct {
	PlaybookID string `json:"playbook_id"`
	Title      string `json:"title"`
}

// Schedule represents a v3 schedule
type Schedule struct {
	ScheduledSessionID string                `json:"scheduled_session_id"`
	Name               string                `json:"name"`
	Prompt             string                `json:"prompt"`
	Frequency          *string               `json:"frequency"`
	Enabled            bool                  `json:"enabled"`
	Playbook           *SchedulePlaybookInfo `json:"playbook"`
	Agent              string                `json:"agent,omitempty"`
	OrgID              string                `json:"org_id,omitempty"`
	CreatedAt          string                `json:"created_at"`
	UpdatedAt          string                `json:"updated_at"`
}

// ListSchedulesResponse represents the response from schedules API
type ListSchedulesResponse struct {
	Items       []Schedule `json:"items"`
	EndCursor   *string    `json:"end_cursor"`
	HasNextPage bool       `json:"has_next_page"`
}

// CreateScheduleRequest represents the request for creating a schedule
type CreateScheduleRequest struct {
	Name       string `json:"name"`
	Prompt     string `json:"prompt"`
	Frequency  string `json:"frequency,omitempty"`
	PlaybookID string `json:"playbook_id,omitempty"`
}

// UpdateScheduleRequest represents the request for updating a schedule (PATCH)
type UpdateScheduleRequest struct {
	Name       *string `json:"name,omitempty"`
	Prompt     *string `json:"prompt,omitempty"`
	Frequency  *string `json:"frequency,omitempty"`
	PlaybookID *string `json:"playbook_id,omitempty"`
	Enabled    *bool   `json:"enabled,omitempty"`
}

// --- Error type ---

// ErrorResponse represents the API error response
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// NewClient creates a new DevinClient for v3 API
func NewClient(apiKey, orgID string) *DevinClient {
	return &DevinClient{
		APIKey: apiKey,
		OrgID:  orgID,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		MaxRetries: 5,
		CacheTTL:   5 * time.Minute,
	}
}

// InvalidateKnowledgeCache clears the knowledge notes cache
func (c *DevinClient) InvalidateKnowledgeCache() {
	c.knowledgeCacheMu.Lock()
	defer c.knowledgeCacheMu.Unlock()
	c.knowledgeCache = nil
	c.knowledgeCacheTime = time.Time{}
}

// isKnowledgeCacheValid checks if the knowledge cache is still valid
func (c *DevinClient) isKnowledgeCacheValid() bool {
	if c.knowledgeCache == nil {
		return false
	}
	return time.Since(c.knowledgeCacheTime) < c.CacheTTL
}

// populateKnowledgeCache fetches all knowledge notes and caches them
func (c *DevinClient) populateKnowledgeCache() error {
	notes, err := c.ListKnowledgeNotes()
	if err != nil {
		return err
	}

	cache := make(map[string]*KnowledgeNote, len(notes))
	for i := range notes {
		cache[notes[i].NoteID] = &notes[i]
	}

	c.knowledgeCache = cache
	c.knowledgeCacheTime = time.Now()
	return nil
}

// sendRequest is a common function for sending requests with rate limit retry
func (c *DevinClient) sendRequest(method, path string, body interface{}) ([]byte, error) {
	url := c.baseURL() + path

	var jsonData []byte
	if body != nil {
		var err error
		jsonData, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON encode request body: %w", err)
		}
	}

	var lastErr error
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		var reqBody io.Reader
		if jsonData != nil {
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

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Rate limit handling: retry with exponential backoff
		if resp.StatusCode == 429 {
			if attempt >= c.MaxRetries {
				lastErr = fmt.Errorf("API rate limit exceeded after %d retries", c.MaxRetries)
				break
			}

			// Use Retry-After header if available, otherwise exponential backoff
			wait := time.Duration(1<<uint(attempt)) * time.Second
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					wait = time.Duration(seconds) * time.Second
				}
			}

			time.Sleep(wait)
			continue
		}

		if resp.StatusCode >= 400 {
			var errResp ErrorResponse
			if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error.Message != "" {
				return nil, fmt.Errorf("API error (status %d): %s (%s)", resp.StatusCode, errResp.Error.Message, errResp.Error.Type)
			}
			return nil, fmt.Errorf("API error: status code %d, body: %s", resp.StatusCode, string(respBody))
		}

		return respBody, nil
	}

	return nil, lastErr
}

// ===================== Knowledge Notes =====================

// GetKnowledgeNote retrieves a single knowledge note by ID
// Uses cache when available to avoid rate limiting during terraform plan/apply
func (c *DevinClient) GetKnowledgeNote(noteID string) (*KnowledgeNote, error) {
	if IsMockClient(c.APIKey) {
		return GetMockKnowledgeNote(noteID)
	}

	// Try cache first
	c.knowledgeCacheMu.RLock()
	if c.isKnowledgeCacheValid() {
		if note, ok := c.knowledgeCache[noteID]; ok {
			c.knowledgeCacheMu.RUnlock()
			return note, nil
		}
	}
	c.knowledgeCacheMu.RUnlock()

	// Cache miss: populate cache from list API
	c.knowledgeCacheMu.Lock()
	// Double-check after acquiring write lock
	if !c.isKnowledgeCacheValid() {
		if err := c.populateKnowledgeCache(); err != nil {
			c.knowledgeCacheMu.Unlock()
			// Fallback to individual API call
			return c.getKnowledgeNoteDirect(noteID)
		}
	}
	if note, ok := c.knowledgeCache[noteID]; ok {
		c.knowledgeCacheMu.Unlock()
		return note, nil
	}
	c.knowledgeCacheMu.Unlock()

	// Not found in cache, try direct API (might be newly created)
	return c.getKnowledgeNoteDirect(noteID)
}

// getKnowledgeNoteDirect retrieves a single knowledge note directly from the API
func (c *DevinClient) getKnowledgeNoteDirect(noteID string) (*KnowledgeNote, error) {
	path := fmt.Sprintf("/knowledge/notes/%s", noteID)
	respBody, err := c.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var note KnowledgeNote
	if err := json.Unmarshal(respBody, &note); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &note, nil
}

// ListKnowledgeNotes retrieves all knowledge notes with pagination
func (c *DevinClient) ListKnowledgeNotes() ([]KnowledgeNote, error) {
	if IsMockClient(c.APIKey) {
		return GetMockKnowledgeNoteList(), nil
	}

	var allNotes []KnowledgeNote
	after := ""

	for {
		path := "/knowledge/notes?first=200"
		if after != "" {
			path += "&after=" + after
		}

		respBody, err := c.sendRequest("GET", path, nil)
		if err != nil {
			return nil, err
		}

		var response ListKnowledgeNotesResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to decode JSON response: %w", err)
		}

		allNotes = append(allNotes, response.Items...)

		if !response.HasNextPage {
			break
		}
		if response.EndCursor != nil {
			after = *response.EndCursor
		} else {
			break
		}
	}

	return allNotes, nil
}

// CreateKnowledgeNote creates a new knowledge note
func (c *DevinClient) CreateKnowledgeNote(reqBody CreateKnowledgeNoteRequest) (*KnowledgeNote, error) {
	if IsMockClient(c.APIKey) {
		return CreateMockKnowledgeNote(reqBody), nil
	}

	respBody, err := c.sendRequest("POST", "/knowledge/notes", reqBody)
	if err != nil {
		return nil, err
	}

	var note KnowledgeNote
	if err := json.Unmarshal(respBody, &note); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	c.InvalidateKnowledgeCache()

	return &note, nil
}

// UpdateKnowledgeNote updates a knowledge note
func (c *DevinClient) UpdateKnowledgeNote(noteID string, reqBody UpdateKnowledgeNoteRequest) (*KnowledgeNote, error) {
	if IsMockClient(c.APIKey) {
		return UpdateMockKnowledgeNote(noteID, reqBody), nil
	}

	path := fmt.Sprintf("/knowledge/notes/%s", noteID)
	respBody, err := c.sendRequest("PUT", path, reqBody)
	if err != nil {
		return nil, err
	}

	var note KnowledgeNote
	if err := json.Unmarshal(respBody, &note); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	c.InvalidateKnowledgeCache()

	return &note, nil
}

// DeleteKnowledgeNote deletes a knowledge note
func (c *DevinClient) DeleteKnowledgeNote(noteID string) error {
	if IsMockClient(c.APIKey) {
		DeleteMockKnowledgeNote(noteID)
		return nil
	}

	path := fmt.Sprintf("/knowledge/notes/%s", noteID)
	_, err := c.sendRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	c.InvalidateKnowledgeCache()

	return nil
}

// ===================== Folders =====================

// ListFolders retrieves all folders with pagination
func (c *DevinClient) ListFolders() ([]FolderItem, error) {
	if IsMockClient(c.APIKey) {
		return GetMockFolderList(), nil
	}

	var allFolders []FolderItem
	after := ""

	for {
		path := "/knowledge/folders?first=200"
		if after != "" {
			path += "&after=" + after
		}

		respBody, err := c.sendRequest("GET", path, nil)
		if err != nil {
			return nil, err
		}

		var response ListFoldersResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to decode JSON response: %w", err)
		}

		allFolders = append(allFolders, response.Items...)

		if !response.HasNextPage {
			break
		}
		if response.EndCursor != nil {
			after = *response.EndCursor
		} else {
			break
		}
	}

	return allFolders, nil
}

// GetFolderByID retrieves a folder by ID from the list
func (c *DevinClient) GetFolderByID(id string) (*FolderItem, error) {
	if IsMockClient(c.APIKey) {
		return GetMockFolderByID(id)
	}

	folders, err := c.ListFolders()
	if err != nil {
		return nil, fmt.Errorf("error occurred while retrieving folder list: %w", err)
	}

	for _, folder := range folders {
		if folder.FolderID == id {
			return &folder, nil
		}
	}

	return nil, fmt.Errorf("folder resource with ID '%s' not found", id)
}

// GetFolderByName retrieves a folder by name from the list
func (c *DevinClient) GetFolderByName(name string) (*FolderItem, error) {
	if IsMockClient(c.APIKey) {
		return GetMockFolderByName(name)
	}

	folders, err := c.ListFolders()
	if err != nil {
		return nil, fmt.Errorf("error occurred while retrieving folder list: %w", err)
	}

	for _, folder := range folders {
		if folder.Name == name {
			return &folder, nil
		}
	}

	return nil, fmt.Errorf("folder resource with name '%s' not found", name)
}

// ===================== Playbooks =====================

// GetPlaybook retrieves a single playbook by ID
func (c *DevinClient) GetPlaybook(playbookID string) (*Playbook, error) {
	if IsMockClient(c.APIKey) {
		return GetMockPlaybook(playbookID)
	}

	path := fmt.Sprintf("/playbooks/%s", playbookID)
	respBody, err := c.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var playbook Playbook
	if err := json.Unmarshal(respBody, &playbook); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &playbook, nil
}

// CreatePlaybook creates a new playbook
func (c *DevinClient) CreatePlaybook(reqBody CreatePlaybookRequest) (*Playbook, error) {
	if IsMockClient(c.APIKey) {
		return CreateMockPlaybook(reqBody), nil
	}

	respBody, err := c.sendRequest("POST", "/playbooks", reqBody)
	if err != nil {
		return nil, err
	}

	var playbook Playbook
	if err := json.Unmarshal(respBody, &playbook); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &playbook, nil
}

// UpdatePlaybook updates a playbook
func (c *DevinClient) UpdatePlaybook(playbookID string, reqBody UpdatePlaybookRequest) (*Playbook, error) {
	if IsMockClient(c.APIKey) {
		return UpdateMockPlaybook(playbookID, reqBody), nil
	}

	path := fmt.Sprintf("/playbooks/%s", playbookID)
	respBody, err := c.sendRequest("PUT", path, reqBody)
	if err != nil {
		return nil, err
	}

	var playbook Playbook
	if err := json.Unmarshal(respBody, &playbook); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &playbook, nil
}

// DeletePlaybook deletes a playbook
func (c *DevinClient) DeletePlaybook(playbookID string) error {
	if IsMockClient(c.APIKey) {
		DeleteMockPlaybook(playbookID)
		return nil
	}

	path := fmt.Sprintf("/playbooks/%s", playbookID)
	_, err := c.sendRequest("DELETE", path, nil)
	return err
}

// ===================== Secrets =====================

// ListSecrets retrieves all secrets
func (c *DevinClient) ListSecrets() ([]Secret, error) {
	if IsMockClient(c.APIKey) {
		return GetMockSecretList(), nil
	}

	var allSecrets []Secret
	after := ""

	for {
		path := "/secrets?first=200"
		if after != "" {
			path += "&after=" + after
		}

		respBody, err := c.sendRequest("GET", path, nil)
		if err != nil {
			return nil, err
		}

		var response ListSecretsResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to decode JSON response: %w", err)
		}

		allSecrets = append(allSecrets, response.Items...)

		if !response.HasNextPage {
			break
		}
		if response.EndCursor != nil {
			after = *response.EndCursor
		} else {
			break
		}
	}

	return allSecrets, nil
}

// GetSecretByID retrieves a secret by ID from the list
func (c *DevinClient) GetSecretByID(secretID string) (*Secret, error) {
	if IsMockClient(c.APIKey) {
		return GetMockSecretByID(secretID)
	}

	secrets, err := c.ListSecrets()
	if err != nil {
		return nil, fmt.Errorf("error occurred while retrieving secret list: %w", err)
	}

	for _, secret := range secrets {
		if secret.SecretID == secretID {
			return &secret, nil
		}
	}

	return nil, fmt.Errorf("secret with ID '%s' not found", secretID)
}

// CreateSecret creates a new secret
func (c *DevinClient) CreateSecret(reqBody CreateSecretRequest) (*Secret, error) {
	if IsMockClient(c.APIKey) {
		return CreateMockSecret(reqBody), nil
	}

	respBody, err := c.sendRequest("POST", "/secrets", reqBody)
	if err != nil {
		return nil, err
	}

	var secret Secret
	if err := json.Unmarshal(respBody, &secret); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &secret, nil
}

// DeleteSecret deletes a secret
func (c *DevinClient) DeleteSecret(secretID string) error {
	if IsMockClient(c.APIKey) {
		DeleteMockSecret(secretID)
		return nil
	}

	path := fmt.Sprintf("/secrets/%s", secretID)
	_, err := c.sendRequest("DELETE", path, nil)
	return err
}

// ===================== Schedules =====================

// GetSchedule retrieves a single schedule by ID
func (c *DevinClient) GetSchedule(scheduleID string) (*Schedule, error) {
	if IsMockClient(c.APIKey) {
		return GetMockSchedule(scheduleID)
	}

	path := fmt.Sprintf("/schedules/%s", scheduleID)
	respBody, err := c.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var schedule Schedule
	if err := json.Unmarshal(respBody, &schedule); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &schedule, nil
}

// CreateSchedule creates a new schedule
func (c *DevinClient) CreateSchedule(reqBody CreateScheduleRequest) (*Schedule, error) {
	if IsMockClient(c.APIKey) {
		return CreateMockSchedule(reqBody), nil
	}

	respBody, err := c.sendRequest("POST", "/schedules", reqBody)
	if err != nil {
		return nil, err
	}

	var schedule Schedule
	if err := json.Unmarshal(respBody, &schedule); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &schedule, nil
}

// UpdateSchedule updates a schedule (PATCH - partial update)
func (c *DevinClient) UpdateSchedule(scheduleID string, reqBody UpdateScheduleRequest) (*Schedule, error) {
	if IsMockClient(c.APIKey) {
		return UpdateMockSchedule(scheduleID, reqBody), nil
	}

	path := fmt.Sprintf("/schedules/%s", scheduleID)
	respBody, err := c.sendRequest("PATCH", path, reqBody)
	if err != nil {
		return nil, err
	}

	var schedule Schedule
	if err := json.Unmarshal(respBody, &schedule); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &schedule, nil
}

// DeleteSchedule deletes a schedule
func (c *DevinClient) DeleteSchedule(scheduleID string) error {
	if IsMockClient(c.APIKey) {
		DeleteMockSchedule(scheduleID)
		return nil
	}

	path := fmt.Sprintf("/schedules/%s", scheduleID)
	_, err := c.sendRequest("DELETE", path, nil)
	return err
}
