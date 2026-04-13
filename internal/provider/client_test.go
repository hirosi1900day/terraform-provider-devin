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
	client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "モックナレッジ1", Body: "内容1", Trigger: "トリガー1",
	})
	client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "モックナレッジ2", Body: "内容2", Trigger: "トリガー2",
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
	created, _ := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name:    "モックナレッジ1",
		Body:    "これはテスト用のモックナレッジです",
		Trigger: "テスト用トリガーの説明",
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
	reqBody := CreateKnowledgeNoteRequest{
		Name:    "テストナレッジ",
		Body:    "テスト内容",
		Trigger: "テストトリガー",
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
	created, _ := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "Original", Body: "Body", Trigger: "Trigger",
	})

	reqBody := UpdateKnowledgeNoteRequest{
		Name:    "更新ナレッジ",
		Body:    "更新内容",
		Trigger: "更新トリガー",
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

	created, _ := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "ToDelete", Body: "Body", Trigger: "Trigger",
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
		Title: "モックPlaybook", Body: "テスト用Playbookの内容",
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
		Title: "ToDelete", Body: "Body",
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
		Key:         "TEST_SECRET",
		Type:        "key-value",
		Value:       "secret-value",
		IsSensitive: true,
	}
	secret, err := client.CreateSecret(reqBody)
	if err != nil {
		t.Fatalf("CreateSecret() error = %v", err)
	}

	if secret.SecretID == "" {
		t.Errorf("CreateSecret() SecretID is empty")
	}
	if secret.Key != "TEST_SECRET" {
		t.Errorf("CreateSecret() Key = %s, want %s", secret.Key, "TEST_SECRET")
	}
}

func TestDeleteSecret_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	created, _ := client.CreateSecret(CreateSecretRequest{
		Key: "ToDelete", Type: "key-value", Value: "val", IsSensitive: true,
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
		Name: "定期タスク", Prompt: "定期タスクのプロンプト", Frequency: "0 9 * * 1",
	})

	schedule, err := client.GetSchedule(created.ScheduledSessionID)
	if err != nil {
		t.Fatalf("GetSchedule() error = %v", err)
	}

	if schedule.ScheduledSessionID != created.ScheduledSessionID {
		t.Errorf("GetSchedule() ScheduledSessionID = %s, want %s", schedule.ScheduledSessionID, created.ScheduledSessionID)
	}
	if *schedule.Frequency != "0 9 * * 1" {
		t.Errorf("GetSchedule() Frequency = %s, want %s", *schedule.Frequency, "0 9 * * 1")
	}
}

func TestCreateSchedule_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")
	reqBody := CreateScheduleRequest{
		Name:      "テストスケジュール",
		Prompt:    "テストプロンプト",
		Frequency: "0 9 * * 1",
	}
	schedule, err := client.CreateSchedule(reqBody)
	if err != nil {
		t.Fatalf("CreateSchedule() error = %v", err)
	}

	if schedule.ScheduledSessionID == "" {
		t.Errorf("CreateSchedule() ScheduledSessionID is empty")
	}
	if schedule.Prompt != "テストプロンプト" {
		t.Errorf("CreateSchedule() Prompt = %s, want %s", schedule.Prompt, "テストプロンプト")
	}
}

func TestDeleteSchedule_Mock(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	created, _ := client.CreateSchedule(CreateScheduleRequest{
		Name: "ToDelete", Prompt: "ToDelete", Frequency: "0 * * * *",
	})
	err := client.DeleteSchedule(created.ScheduledSessionID)
	if err != nil {
		t.Errorf("DeleteSchedule() error = %v", err)
	}
}
