package provider

import (
	"fmt"
	"sync"
	"time"
)

// IsMockClient checks if the API key is for mock/testing
func IsMockClient(apiKey string) bool {
	return apiKey == "test_api_key"
}

// ===================== Stateful Mock Store =====================

// MockStore holds all mock data in memory with thread-safe access.
// This enables true CRUD lifecycle testing: Create stores data,
// Read retrieves it, Update modifies it, Delete removes it.
type MockStore struct {
	mu           sync.RWMutex
	knowledge    map[string]*KnowledgeNote
	playbooks    map[string]*Playbook
	secrets      map[string]*Secret
	schedules    map[string]*Schedule
	knowledgeSeq int
	playbookSeq  int
	secretSeq    int
	scheduleSeq  int
}

var globalMockStore = &MockStore{
	knowledge: make(map[string]*KnowledgeNote),
	playbooks: make(map[string]*Playbook),
	secrets:   make(map[string]*Secret),
	schedules: make(map[string]*Schedule),
}

// ResetMockStore clears all mock data (call in test setup)
func ResetMockStore() {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()
	globalMockStore.knowledge = make(map[string]*KnowledgeNote)
	globalMockStore.playbooks = make(map[string]*Playbook)
	globalMockStore.secrets = make(map[string]*Secret)
	globalMockStore.schedules = make(map[string]*Schedule)
	globalMockStore.knowledgeSeq = 0
	globalMockStore.playbookSeq = 0
	globalMockStore.secretSeq = 0
	globalMockStore.scheduleSeq = 0
}

// ===================== Knowledge Note Mocks =====================

// GetMockKnowledgeNote returns a mock knowledge note by ID (stateful)
func GetMockKnowledgeNote(noteID string) (*KnowledgeNote, error) {
	globalMockStore.mu.RLock()
	defer globalMockStore.mu.RUnlock()

	if note, ok := globalMockStore.knowledge[noteID]; ok {
		return note, nil
	}
	return nil, fmt.Errorf("ナレッジが見つかりません: ID %s", noteID)
}

// GetMockKnowledgeNoteList returns all knowledge notes from the store
func GetMockKnowledgeNoteList() []KnowledgeNote {
	globalMockStore.mu.RLock()
	defer globalMockStore.mu.RUnlock()

	notes := make([]KnowledgeNote, 0, len(globalMockStore.knowledge))
	for _, note := range globalMockStore.knowledge {
		notes = append(notes, *note)
	}
	return notes
}

// CreateMockKnowledgeNote creates a mock knowledge note (stateful)
func CreateMockKnowledgeNote(req CreateKnowledgeNoteRequest) *KnowledgeNote {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()

	now := float64(time.Now().Unix())
	globalMockStore.knowledgeSeq++
	id := fmt.Sprintf("note-mock-%d", globalMockStore.knowledgeSeq)

	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}
	note := &KnowledgeNote{
		NoteID:     id,
		Name:       req.Name,
		Body:       req.Body,
		Trigger:    req.Trigger,
		FolderID:   req.FolderID,
		FolderPath: "",
		IsEnabled:  isEnabled,
		PinnedRepo: req.PinnedRepo,
		AccessType: "org",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	globalMockStore.knowledge[id] = note
	return note
}

// UpdateMockKnowledgeNote updates a mock knowledge note (stateful)
func UpdateMockKnowledgeNote(noteID string, req UpdateKnowledgeNoteRequest) *KnowledgeNote {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()

	now := float64(time.Now().Unix())
	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}

	existing, ok := globalMockStore.knowledge[noteID]
	createdAt := now - 86400
	if ok {
		createdAt = existing.CreatedAt
	}

	note := &KnowledgeNote{
		NoteID:     noteID,
		Name:       req.Name,
		Body:       req.Body,
		Trigger:    req.Trigger,
		FolderID:   req.FolderID,
		FolderPath: "",
		IsEnabled:  isEnabled,
		PinnedRepo: req.PinnedRepo,
		AccessType: "org",
		CreatedAt:  createdAt,
		UpdatedAt:  now,
	}
	globalMockStore.knowledge[noteID] = note
	return note
}

// DeleteMockKnowledgeNote deletes a mock knowledge note from the store
func DeleteMockKnowledgeNote(noteID string) {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()
	delete(globalMockStore.knowledge, noteID)
}

// ===================== Folder Mocks =====================

// GetMockFolderList returns a mock list of folders
func GetMockFolderList() []FolderItem {
	return []FolderItem{
		{
			FolderID:       "folder-mock-1",
			Name:           "モックフォルダ1",
			Path:           "/モックフォルダ1",
			NoteCount:      5,
			ParentFolderID: "",
		},
		{
			FolderID:       "folder-mock-2",
			Name:           "モックフォルダ2",
			Path:           "/モックフォルダ2",
			NoteCount:      3,
			ParentFolderID: "",
		},
	}
}

// GetMockFolderByID returns a mock folder by ID
func GetMockFolderByID(id string) (*FolderItem, error) {
	switch id {
	case "folder-mock-1":
		return &FolderItem{
			FolderID:       "folder-mock-1",
			Name:           "モックフォルダ1",
			Path:           "/モックフォルダ1",
			NoteCount:      5,
			ParentFolderID: "",
		}, nil
	case "folder-mock-2":
		return &FolderItem{
			FolderID:       "folder-mock-2",
			Name:           "モックフォルダ2",
			Path:           "/モックフォルダ2",
			NoteCount:      3,
			ParentFolderID: "",
		}, nil
	default:
		return nil, fmt.Errorf("フォルダが見つかりません: ID %s", id)
	}
}

// GetMockFolderByName returns a mock folder by name
func GetMockFolderByName(name string) (*FolderItem, error) {
	switch name {
	case "モックフォルダ1":
		return &FolderItem{
			FolderID:       "folder-mock-1",
			Name:           "モックフォルダ1",
			Path:           "/モックフォルダ1",
			NoteCount:      5,
			ParentFolderID: "",
		}, nil
	case "モックフォルダ2":
		return &FolderItem{
			FolderID:       "folder-mock-2",
			Name:           "モックフォルダ2",
			Path:           "/モックフォルダ2",
			NoteCount:      3,
			ParentFolderID: "",
		}, nil
	default:
		return nil, fmt.Errorf("フォルダが見つかりません: 名前 %s", name)
	}
}

// ===================== Playbook Mocks =====================

// GetMockPlaybook returns a mock playbook by ID (stateful)
func GetMockPlaybook(playbookID string) (*Playbook, error) {
	globalMockStore.mu.RLock()
	defer globalMockStore.mu.RUnlock()

	if pb, ok := globalMockStore.playbooks[playbookID]; ok {
		return pb, nil
	}
	return nil, fmt.Errorf("Playbookが見つかりません: ID %s", playbookID)
}

// CreateMockPlaybook creates a mock playbook (stateful)
func CreateMockPlaybook(req CreatePlaybookRequest) *Playbook {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()

	now := float64(time.Now().Unix())
	globalMockStore.playbookSeq++
	id := fmt.Sprintf("playbook-mock-%d", globalMockStore.playbookSeq)

	status := "active"
	if req.Status != "" {
		status = req.Status
	}
	pb := &Playbook{
		PlaybookID:        id,
		Title:             req.Title,
		Body:              req.Body,
		Status:            status,
		AccessType:        "org",
		Macro:             nil,
		OrgID:             "org-mock",
		CreatedAt:         now,
		UpdatedAt:         now,
		CreatedByUserID:   "user-1",
		CreatedByUserName: "Test User",
		UpdatedByUserID:   "user-1",
		UpdatedByUserName: "Test User",
	}
	globalMockStore.playbooks[id] = pb
	return pb
}

// UpdateMockPlaybook updates a mock playbook (stateful)
func UpdateMockPlaybook(playbookID string, req UpdatePlaybookRequest) *Playbook {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()

	now := float64(time.Now().Unix())
	status := "active"
	if req.Status != "" {
		status = req.Status
	}

	existing, ok := globalMockStore.playbooks[playbookID]
	createdAt := now - 86400
	if ok {
		createdAt = existing.CreatedAt
	}

	pb := &Playbook{
		PlaybookID:        playbookID,
		Title:             req.Title,
		Body:              req.Body,
		Status:            status,
		AccessType:        "org",
		Macro:             nil,
		OrgID:             "org-mock",
		CreatedAt:         createdAt,
		UpdatedAt:         now,
		CreatedByUserID:   "user-1",
		CreatedByUserName: "Test User",
		UpdatedByUserID:   "user-1",
		UpdatedByUserName: "Test User",
	}
	globalMockStore.playbooks[playbookID] = pb
	return pb
}

// DeleteMockPlaybook deletes a mock playbook from the store
func DeleteMockPlaybook(playbookID string) {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()
	delete(globalMockStore.playbooks, playbookID)
}

// ===================== Secret Mocks =====================

// GetMockSecretList returns all secrets from the store
func GetMockSecretList() []Secret {
	globalMockStore.mu.RLock()
	defer globalMockStore.mu.RUnlock()

	secrets := make([]Secret, 0, len(globalMockStore.secrets))
	for _, s := range globalMockStore.secrets {
		secrets = append(secrets, *s)
	}
	return secrets
}

// GetMockSecretByID returns a mock secret by ID (stateful)
func GetMockSecretByID(secretID string) (*Secret, error) {
	globalMockStore.mu.RLock()
	defer globalMockStore.mu.RUnlock()

	if s, ok := globalMockStore.secrets[secretID]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("Secretが見つかりません: ID %s", secretID)
}

// CreateMockSecret creates a mock secret (stateful)
func CreateMockSecret(req CreateSecretRequest) *Secret {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()

	now := float64(time.Now().Unix())
	globalMockStore.secretSeq++
	id := fmt.Sprintf("secret-mock-%d", globalMockStore.secretSeq)

	s := &Secret{
		SecretID:  id,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	globalMockStore.secrets[id] = s
	return s
}

// DeleteMockSecret deletes a mock secret from the store
func DeleteMockSecret(secretID string) {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()
	delete(globalMockStore.secrets, secretID)
}

// ===================== Schedule Mocks =====================

// GetMockSchedule returns a mock schedule by ID (stateful)
func GetMockSchedule(scheduleID string) (*Schedule, error) {
	globalMockStore.mu.RLock()
	defer globalMockStore.mu.RUnlock()

	if s, ok := globalMockStore.schedules[scheduleID]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("Scheduleが見つかりません: ID %s", scheduleID)
}

// CreateMockSchedule creates a mock schedule (stateful)
func CreateMockSchedule(req CreateScheduleRequest) *Schedule {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()

	now := float64(time.Now().Unix())
	globalMockStore.scheduleSeq++
	id := fmt.Sprintf("schedule-mock-%d", globalMockStore.scheduleSeq)

	s := &Schedule{
		ScheduleID: id,
		Prompt:     req.Prompt,
		Cron:       req.Cron,
		PlaybookID: req.PlaybookID,
		Status:     "active",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	globalMockStore.schedules[id] = s
	return s
}

// UpdateMockSchedule updates a mock schedule (stateful)
func UpdateMockSchedule(scheduleID string, req UpdateScheduleRequest) *Schedule {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()

	now := float64(time.Now().Unix())

	existing, ok := globalMockStore.schedules[scheduleID]
	if !ok {
		existing = &Schedule{
			ScheduleID: scheduleID,
			Prompt:     "既存のプロンプト",
			Cron:       "0 9 * * 1",
			Status:     "active",
			CreatedAt:  now - 86400,
		}
	}

	schedule := &Schedule{
		ScheduleID: scheduleID,
		Prompt:     existing.Prompt,
		Cron:       existing.Cron,
		PlaybookID: existing.PlaybookID,
		Status:     existing.Status,
		CreatedAt:  existing.CreatedAt,
		UpdatedAt:  now,
	}

	if req.Prompt != nil {
		schedule.Prompt = *req.Prompt
	}
	if req.Cron != nil {
		schedule.Cron = *req.Cron
	}
	if req.PlaybookID != nil {
		schedule.PlaybookID = *req.PlaybookID
	}
	if req.Status != nil {
		schedule.Status = *req.Status
	}

	globalMockStore.schedules[scheduleID] = schedule
	return schedule
}

// DeleteMockSchedule deletes a mock schedule from the store
func DeleteMockSchedule(scheduleID string) {
	globalMockStore.mu.Lock()
	defer globalMockStore.mu.Unlock()
	delete(globalMockStore.schedules, scheduleID)
}
