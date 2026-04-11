package provider

import (
	"fmt"
	"time"
)

// IsMockClient checks if the API key is for mock/testing
func IsMockClient(apiKey string) bool {
	return apiKey == "test_api_key"
}

// ===================== Knowledge Note Mocks =====================

// GetMockKnowledgeNote returns a mock knowledge note by ID
func GetMockKnowledgeNote(noteID string) (*KnowledgeNote, error) {
	now := float64(time.Now().Unix())
	switch noteID {
	case "note-mock-1":
		return &KnowledgeNote{
			NoteID:     "note-mock-1",
			Name:       "モックナレッジ1",
			Body:       "これはテスト用のモックナレッジです",
			Trigger:    "テスト用トリガーの説明",
			FolderID:   "folder-mock-1",
			FolderPath: "/モックフォルダ1",
			IsEnabled:  true,
			PinnedRepo: nil,
			Macro:      "",
			AccessType: "org",
			CreatedAt:  now - 86400,
			UpdatedAt:  now - 3600,
		}, nil
	case "note-mock-2":
		repo := "owner/repo"
		return &KnowledgeNote{
			NoteID:     "note-mock-2",
			Name:       "モックナレッジ2",
			Body:       "これは別のテスト用のモックナレッジです",
			Trigger:    "別のテスト用トリガーの説明",
			FolderID:   "folder-mock-2",
			FolderPath: "/モックフォルダ2",
			IsEnabled:  false,
			PinnedRepo: &repo,
			Macro:      "",
			AccessType: "org",
			CreatedAt:  now - 172800,
			UpdatedAt:  now - 7200,
		}, nil
	default:
		return nil, fmt.Errorf("ナレッジが見つかりません: ID %s", noteID)
	}
}

// GetMockKnowledgeNoteList returns a mock list of knowledge notes
func GetMockKnowledgeNoteList() []KnowledgeNote {
	now := float64(time.Now().Unix())
	return []KnowledgeNote{
		{
			NoteID:     "note-mock-1",
			Name:       "モックナレッジ1",
			Body:       "これはテスト用のモックナレッジの内容です。",
			Trigger:    "テスト用トリガーの説明",
			FolderID:   "folder-mock-1",
			FolderPath: "/モックフォルダ1",
			IsEnabled:  true,
			PinnedRepo: nil,
			AccessType: "org",
			CreatedAt:  now - 86400,
			UpdatedAt:  now - 3600,
		},
		{
			NoteID:     "note-mock-2",
			Name:       "モックナレッジ2",
			Body:       "これは別のテスト用のモックナレッジの内容です。",
			Trigger:    "別のテスト用トリガーの説明",
			FolderID:   "folder-mock-2",
			FolderPath: "/モックフォルダ2",
			IsEnabled:  false,
			PinnedRepo: nil,
			AccessType: "org",
			CreatedAt:  now - 172800,
			UpdatedAt:  now - 7200,
		},
	}
}

// CreateMockKnowledgeNote creates a mock knowledge note
func CreateMockKnowledgeNote(req CreateKnowledgeNoteRequest) *KnowledgeNote {
	now := float64(time.Now().Unix())
	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}
	return &KnowledgeNote{
		NoteID:     "note-new-mock",
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
}

// UpdateMockKnowledgeNote updates a mock knowledge note
func UpdateMockKnowledgeNote(noteID string, req UpdateKnowledgeNoteRequest) *KnowledgeNote {
	now := float64(time.Now().Unix())
	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}
	return &KnowledgeNote{
		NoteID:     noteID,
		Name:       req.Name,
		Body:       req.Body,
		Trigger:    req.Trigger,
		FolderID:   req.FolderID,
		FolderPath: "",
		IsEnabled:  isEnabled,
		PinnedRepo: req.PinnedRepo,
		AccessType: "org",
		CreatedAt:  now - 86400,
		UpdatedAt:  now,
	}
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

// GetMockPlaybook returns a mock playbook by ID
func GetMockPlaybook(playbookID string) (*Playbook, error) {
	now := float64(time.Now().Unix())
	switch playbookID {
	case "playbook-mock-1":
		return &Playbook{
			PlaybookID:        "playbook-mock-1",
			Title:             "モックPlaybook",
			Body:              "テスト用Playbookの内容",
			Status:            "active",
			AccessType:        "org",
			Macro:             nil,
			OrgID:             "org-mock",
			CreatedAt:         now - 86400,
			UpdatedAt:         now - 3600,
			CreatedByUserID:   "user-1",
			CreatedByUserName: "Test User",
			UpdatedByUserID:   "user-1",
			UpdatedByUserName: "Test User",
		}, nil
	default:
		return nil, fmt.Errorf("Playbookが見つかりません: ID %s", playbookID)
	}
}

// CreateMockPlaybook creates a mock playbook
func CreateMockPlaybook(req CreatePlaybookRequest) *Playbook {
	now := float64(time.Now().Unix())
	status := "active"
	if req.Status != "" {
		status = req.Status
	}
	return &Playbook{
		PlaybookID:        "playbook-new-mock",
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
}

// UpdateMockPlaybook updates a mock playbook
func UpdateMockPlaybook(playbookID string, req UpdatePlaybookRequest) *Playbook {
	now := float64(time.Now().Unix())
	status := "active"
	if req.Status != "" {
		status = req.Status
	}
	return &Playbook{
		PlaybookID:        playbookID,
		Title:             req.Title,
		Body:              req.Body,
		Status:            status,
		AccessType:        "org",
		Macro:             nil,
		OrgID:             "org-mock",
		CreatedAt:         now - 86400,
		UpdatedAt:         now,
		CreatedByUserID:   "user-1",
		CreatedByUserName: "Test User",
		UpdatedByUserID:   "user-1",
		UpdatedByUserName: "Test User",
	}
}

// ===================== Secret Mocks =====================

// GetMockSecretList returns a mock list of secrets
func GetMockSecretList() []Secret {
	now := float64(time.Now().Unix())
	return []Secret{
		{
			SecretID:  "secret-mock-1",
			Name:      "DATABASE_URL",
			CreatedAt: now - 86400,
			UpdatedAt: now - 3600,
		},
		{
			SecretID:  "secret-mock-2",
			Name:      "API_TOKEN",
			CreatedAt: now - 172800,
			UpdatedAt: now - 7200,
		},
	}
}

// GetMockSecretByID returns a mock secret by ID
func GetMockSecretByID(secretID string) (*Secret, error) {
	now := float64(time.Now().Unix())
	switch secretID {
	case "secret-mock-1":
		return &Secret{
			SecretID:  "secret-mock-1",
			Name:      "DATABASE_URL",
			CreatedAt: now - 86400,
			UpdatedAt: now - 3600,
		}, nil
	case "secret-mock-2":
		return &Secret{
			SecretID:  "secret-mock-2",
			Name:      "API_TOKEN",
			CreatedAt: now - 172800,
			UpdatedAt: now - 7200,
		}, nil
	default:
		return nil, fmt.Errorf("Secretが見つかりません: ID %s", secretID)
	}
}

// CreateMockSecret creates a mock secret
func CreateMockSecret(req CreateSecretRequest) *Secret {
	now := float64(time.Now().Unix())
	return &Secret{
		SecretID:  "secret-new-mock",
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ===================== Schedule Mocks =====================

// GetMockSchedule returns a mock schedule by ID
func GetMockSchedule(scheduleID string) (*Schedule, error) {
	now := float64(time.Now().Unix())
	switch scheduleID {
	case "schedule-mock-1":
		return &Schedule{
			ScheduleID: "schedule-mock-1",
			Prompt:     "定期タスクのプロンプト",
			Cron:       "0 9 * * 1",
			PlaybookID: "playbook-mock-1",
			Status:     "active",
			CreatedAt:  now - 86400,
			UpdatedAt:  now - 3600,
		}, nil
	default:
		return nil, fmt.Errorf("Scheduleが見つかりません: ID %s", scheduleID)
	}
}

// CreateMockSchedule creates a mock schedule
func CreateMockSchedule(req CreateScheduleRequest) *Schedule {
	now := float64(time.Now().Unix())
	return &Schedule{
		ScheduleID: "schedule-new-mock",
		Prompt:     req.Prompt,
		Cron:       req.Cron,
		PlaybookID: req.PlaybookID,
		Status:     "active",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// UpdateMockSchedule updates a mock schedule
func UpdateMockSchedule(scheduleID string, req UpdateScheduleRequest) *Schedule {
	now := float64(time.Now().Unix())
	schedule := &Schedule{
		ScheduleID: scheduleID,
		Prompt:     "既存のプロンプト",
		Cron:       "0 9 * * 1",
		Status:     "active",
		CreatedAt:  now - 86400,
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

	return schedule
}
