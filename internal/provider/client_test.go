package provider

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key", "org-test")
	if client == nil {
		t.Fatalf("NewClient() returned nil")
	}
	if client.APIKey != "test-api-key" {
		t.Errorf("NewClient() API key = %s, want %s", client.APIKey, "test-api-key")
	}
	if client.OrgID != "org-test" {
		t.Errorf("NewClient() OrgID = %s, want %s", client.OrgID, "org-test")
	}
	if client.HTTPClient == nil {
		t.Errorf("NewClient() HTTPClient is nil")
	}
}

func TestBaseURL(t *testing.T) {
	client := NewClient("test-api-key", "org-123")
	expected := "https://api.devin.ai/v3/organizations/org-123"
	if client.baseURL() != expected {
		t.Errorf("baseURL() = %s, want %s", client.baseURL(), expected)
	}
}

func TestListKnowledgeNotes_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// Seed 2 notes
	isEnabled := true
	client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "モックナレッジ1", Body: "内容1", Trigger: "トリガー1", IsEnabled: &isEnabled,
	})
	client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "モックナレッジ2", Body: "内容2", Trigger: "トリガー2", IsEnabled: &isEnabled,
	})

	notes, err := client.ListKnowledgeNotes()
	if err != nil {
		t.Fatalf("ListKnowledgeNotes() error = %v", err)
	}

	if len(notes) != 2 {
		t.Errorf("ListKnowledgeNotes() returned %d items, want 2", len(notes))
	}
}

func TestGetKnowledgeNote_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// Seed a note
	isEnabled := true
	created, _ := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name:      "モックナレッジ1",
		Body:      "これはテスト用のモックナレッジです",
		Trigger:   "テスト用トリガーの説明",
		FolderID:  "folder-mock-1",
		IsEnabled: &isEnabled,
	})

	note, err := client.GetKnowledgeNote(created.NoteID)
	if err != nil {
		t.Fatalf("GetKnowledgeNote() error = %v", err)
	}

	if note.NoteID != created.NoteID {
		t.Errorf("GetKnowledgeNote() NoteID = %s, want %s", note.NoteID, created.NoteID)
	}
	if note.Name != "モックナレッジ1" {
		t.Errorf("GetKnowledgeNote() Name = %s, want %s", note.Name, "モックナレッジ1")
	}
	if note.Trigger != "テスト用トリガーの説明" {
		t.Errorf("GetKnowledgeNote() Trigger = %s, want %s", note.Trigger, "テスト用トリガーの説明")
	}
	if note.FolderID != "folder-mock-1" {
		t.Errorf("GetKnowledgeNote() FolderID = %s, want %s", note.FolderID, "folder-mock-1")
	}
	if !note.IsEnabled {
		t.Errorf("GetKnowledgeNote() IsEnabled = false, want true")
	}

	// Non-existent ID test
	_, err = client.GetKnowledgeNote("non-existent-id")
	if err == nil {
		t.Errorf("GetKnowledgeNote() with non-existent ID should return error")
	}
}

func TestCreateKnowledgeNote_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")
	isEnabled := true
	reqBody := CreateKnowledgeNoteRequest{
		Name:      "テストナレッジ",
		Body:      "テスト内容",
		Trigger:   "テストトリガー",
		FolderID:  "test-folder-id",
		IsEnabled: &isEnabled,
	}
	note, err := client.CreateKnowledgeNote(reqBody)
	if err != nil {
		t.Fatalf("CreateKnowledgeNote() error = %v", err)
	}

	if note.NoteID == "" {
		t.Errorf("CreateKnowledgeNote() NoteID is empty")
	}
	if note.Name != "テストナレッジ" {
		t.Errorf("CreateKnowledgeNote() Name = %s, want %s", note.Name, "テストナレッジ")
	}
	if note.Trigger != "テストトリガー" {
		t.Errorf("CreateKnowledgeNote() Trigger = %s, want %s", note.Trigger, "テストトリガー")
	}
}

func TestUpdateKnowledgeNote_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// Seed a note first
	isEnabled := true
	created, _ := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "Original", Body: "Body", Trigger: "Trigger", IsEnabled: &isEnabled,
	})

	reqBody := UpdateKnowledgeNoteRequest{
		Name:      "更新ナレッジ",
		Body:      "更新内容",
		Trigger:   "更新トリガー",
		FolderID:  "updated-folder-id",
		IsEnabled: &isEnabled,
	}
	note, err := client.UpdateKnowledgeNote(created.NoteID, reqBody)
	if err != nil {
		t.Fatalf("UpdateKnowledgeNote() error = %v", err)
	}

	if note.NoteID != created.NoteID {
		t.Errorf("UpdateKnowledgeNote() NoteID = %s, want %s", note.NoteID, created.NoteID)
	}
	if note.Name != "更新ナレッジ" {
		t.Errorf("UpdateKnowledgeNote() Name = %s, want %s", note.Name, "更新ナレッジ")
	}
}

func TestDeleteKnowledgeNote_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	isEnabled := true
	created, _ := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "ToDelete", Body: "Body", Trigger: "Trigger", IsEnabled: &isEnabled,
	})
	err := client.DeleteKnowledgeNote(created.NoteID)
	if err != nil {
		t.Errorf("DeleteKnowledgeNote() error = %v", err)
	}
}

func TestGetFolderByID_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	folder, err := client.GetFolderByID("folder-mock-1")
	if err != nil {
		t.Fatalf("GetFolderByID() error = %v", err)
	}

	if folder.FolderID != "folder-mock-1" {
		t.Errorf("GetFolderByID() FolderID = %s, want %s", folder.FolderID, "folder-mock-1")
	}
	if folder.Name != "モックフォルダ1" {
		t.Errorf("GetFolderByID() Name = %s, want %s", folder.Name, "モックフォルダ1")
	}
}

func TestGetFolderByName_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	folder, err := client.GetFolderByName("モックフォルダ1")
	if err != nil {
		t.Fatalf("GetFolderByName() error = %v", err)
	}

	if folder.FolderID != "folder-mock-1" {
		t.Errorf("GetFolderByName() FolderID = %s, want %s", folder.FolderID, "folder-mock-1")
	}
}

func TestGetPlaybook_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// Seed a playbook
	created, _ := client.CreatePlaybook(CreatePlaybookRequest{
		Title: "モックPlaybook", Body: "テスト用Playbookの内容", Status: "active",
	})

	playbook, err := client.GetPlaybook(created.PlaybookID)
	if err != nil {
		t.Fatalf("GetPlaybook() error = %v", err)
	}

	if playbook.PlaybookID != created.PlaybookID {
		t.Errorf("GetPlaybook() PlaybookID = %s, want %s", playbook.PlaybookID, created.PlaybookID)
	}
	if playbook.Title != "モックPlaybook" {
		t.Errorf("GetPlaybook() Title = %s, want %s", playbook.Title, "モックPlaybook")
	}
	if playbook.Status != "active" {
		t.Errorf("GetPlaybook() Status = %s, want %s", playbook.Status, "active")
	}
}

func TestCreatePlaybook_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")
	reqBody := CreatePlaybookRequest{
		Title: "新しいPlaybook",
		Body:  "Playbookの内容",
	}
	playbook, err := client.CreatePlaybook(reqBody)
	if err != nil {
		t.Fatalf("CreatePlaybook() error = %v", err)
	}

	if playbook.PlaybookID == "" {
		t.Errorf("CreatePlaybook() PlaybookID is empty")
	}
	if playbook.Title != "新しいPlaybook" {
		t.Errorf("CreatePlaybook() Title = %s, want %s", playbook.Title, "新しいPlaybook")
	}
}

func TestDeletePlaybook_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	created, _ := client.CreatePlaybook(CreatePlaybookRequest{
		Title: "ToDelete", Body: "Body", Status: "active",
	})
	err := client.DeletePlaybook(created.PlaybookID)
	if err != nil {
		t.Errorf("DeletePlaybook() error = %v", err)
	}
}

func TestCreateSecret_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")
	reqBody := CreateSecretRequest{
		Name:  "TEST_SECRET",
		Value: "secret-value",
	}
	secret, err := client.CreateSecret(reqBody)
	if err != nil {
		t.Fatalf("CreateSecret() error = %v", err)
	}

	if secret.SecretID == "" {
		t.Errorf("CreateSecret() SecretID is empty")
	}
	if secret.Name != "TEST_SECRET" {
		t.Errorf("CreateSecret() Name = %s, want %s", secret.Name, "TEST_SECRET")
	}
}

func TestDeleteSecret_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	created, _ := client.CreateSecret(CreateSecretRequest{
		Name: "ToDelete", Value: "val",
	})
	err := client.DeleteSecret(created.SecretID)
	if err != nil {
		t.Errorf("DeleteSecret() error = %v", err)
	}
}

func TestGetSchedule_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	created, _ := client.CreateSchedule(CreateScheduleRequest{
		Prompt: "定期タスクのプロンプト", Cron: "0 9 * * 1",
	})

	schedule, err := client.GetSchedule(created.ScheduleID)
	if err != nil {
		t.Fatalf("GetSchedule() error = %v", err)
	}

	if schedule.ScheduleID != created.ScheduleID {
		t.Errorf("GetSchedule() ScheduleID = %s, want %s", schedule.ScheduleID, created.ScheduleID)
	}
	if schedule.Cron != "0 9 * * 1" {
		t.Errorf("GetSchedule() Cron = %s, want %s", schedule.Cron, "0 9 * * 1")
	}
}

func TestCreateSchedule_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")
	reqBody := CreateScheduleRequest{
		Prompt: "テストプロンプト",
		Cron:   "0 9 * * 1",
	}
	schedule, err := client.CreateSchedule(reqBody)
	if err != nil {
		t.Fatalf("CreateSchedule() error = %v", err)
	}

	if schedule.ScheduleID == "" {
		t.Errorf("CreateSchedule() ScheduleID is empty")
	}
	if schedule.Prompt != "テストプロンプト" {
		t.Errorf("CreateSchedule() Prompt = %s, want %s", schedule.Prompt, "テストプロンプト")
	}
}

func TestDeleteSchedule_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	created, _ := client.CreateSchedule(CreateScheduleRequest{
		Prompt: "ToDelete", Cron: "0 * * * *",
	})
	err := client.DeleteSchedule(created.ScheduleID)
	if err != nil {
		t.Errorf("DeleteSchedule() error = %v", err)
	}
}
