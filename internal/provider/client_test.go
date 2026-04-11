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
	client := NewClient("test_api_key", "org-mock")
	notes, err := client.ListKnowledgeNotes()
	if err != nil {
		t.Fatalf("ListKnowledgeNotes() error = %v", err)
	}

	if len(notes) != 2 {
		t.Errorf("ListKnowledgeNotes() returned %d items, want 2", len(notes))
	}

	if notes[0].NoteID != "note-mock-1" {
		t.Errorf("First note ID = %s, want %s", notes[0].NoteID, "note-mock-1")
	}
	if notes[0].Name != "モックナレッジ1" {
		t.Errorf("First note Name = %s, want %s", notes[0].Name, "モックナレッジ1")
	}
}

func TestGetKnowledgeNote_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")

	note, err := client.GetKnowledgeNote("note-mock-1")
	if err != nil {
		t.Fatalf("GetKnowledgeNote() error = %v", err)
	}

	if note.NoteID != "note-mock-1" {
		t.Errorf("GetKnowledgeNote() NoteID = %s, want %s", note.NoteID, "note-mock-1")
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

	if note.NoteID != "note-new-mock" {
		t.Errorf("CreateKnowledgeNote() NoteID = %s, want %s", note.NoteID, "note-new-mock")
	}
	if note.Name != "テストナレッジ" {
		t.Errorf("CreateKnowledgeNote() Name = %s, want %s", note.Name, "テストナレッジ")
	}
	if note.Trigger != "テストトリガー" {
		t.Errorf("CreateKnowledgeNote() Trigger = %s, want %s", note.Trigger, "テストトリガー")
	}
}

func TestUpdateKnowledgeNote_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	isEnabled := true
	reqBody := UpdateKnowledgeNoteRequest{
		Name:      "更新ナレッジ",
		Body:      "更新内容",
		Trigger:   "更新トリガー",
		FolderID:  "updated-folder-id",
		IsEnabled: &isEnabled,
	}
	note, err := client.UpdateKnowledgeNote("note-mock-1", reqBody)
	if err != nil {
		t.Fatalf("UpdateKnowledgeNote() error = %v", err)
	}

	if note.NoteID != "note-mock-1" {
		t.Errorf("UpdateKnowledgeNote() NoteID = %s, want %s", note.NoteID, "note-mock-1")
	}
	if note.Name != "更新ナレッジ" {
		t.Errorf("UpdateKnowledgeNote() Name = %s, want %s", note.Name, "更新ナレッジ")
	}
}

func TestDeleteKnowledgeNote_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	err := client.DeleteKnowledgeNote("note-mock-1")
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
	client := NewClient("test_api_key", "org-mock")
	playbook, err := client.GetPlaybook("playbook-mock-1")
	if err != nil {
		t.Fatalf("GetPlaybook() error = %v", err)
	}

	if playbook.PlaybookID != "playbook-mock-1" {
		t.Errorf("GetPlaybook() PlaybookID = %s, want %s", playbook.PlaybookID, "playbook-mock-1")
	}
	if playbook.Title != "モックPlaybook" {
		t.Errorf("GetPlaybook() Title = %s, want %s", playbook.Title, "モックPlaybook")
	}
	if playbook.Status != "active" {
		t.Errorf("GetPlaybook() Status = %s, want %s", playbook.Status, "active")
	}
}

func TestCreatePlaybook_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	reqBody := CreatePlaybookRequest{
		Title: "新しいPlaybook",
		Body:  "Playbookの内容",
	}
	playbook, err := client.CreatePlaybook(reqBody)
	if err != nil {
		t.Fatalf("CreatePlaybook() error = %v", err)
	}

	if playbook.PlaybookID != "playbook-new-mock" {
		t.Errorf("CreatePlaybook() PlaybookID = %s, want %s", playbook.PlaybookID, "playbook-new-mock")
	}
	if playbook.Title != "新しいPlaybook" {
		t.Errorf("CreatePlaybook() Title = %s, want %s", playbook.Title, "新しいPlaybook")
	}
}

func TestDeletePlaybook_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	err := client.DeletePlaybook("playbook-mock-1")
	if err != nil {
		t.Errorf("DeletePlaybook() error = %v", err)
	}
}

func TestCreateSecret_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	reqBody := CreateSecretRequest{
		Name:  "TEST_SECRET",
		Value: "secret-value",
	}
	secret, err := client.CreateSecret(reqBody)
	if err != nil {
		t.Fatalf("CreateSecret() error = %v", err)
	}

	if secret.SecretID != "secret-new-mock" {
		t.Errorf("CreateSecret() SecretID = %s, want %s", secret.SecretID, "secret-new-mock")
	}
	if secret.Name != "TEST_SECRET" {
		t.Errorf("CreateSecret() Name = %s, want %s", secret.Name, "TEST_SECRET")
	}
}

func TestDeleteSecret_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	err := client.DeleteSecret("secret-mock-1")
	if err != nil {
		t.Errorf("DeleteSecret() error = %v", err)
	}
}

func TestGetSchedule_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	schedule, err := client.GetSchedule("schedule-mock-1")
	if err != nil {
		t.Fatalf("GetSchedule() error = %v", err)
	}

	if schedule.ScheduleID != "schedule-mock-1" {
		t.Errorf("GetSchedule() ScheduleID = %s, want %s", schedule.ScheduleID, "schedule-mock-1")
	}
	if schedule.Cron != "0 9 * * 1" {
		t.Errorf("GetSchedule() Cron = %s, want %s", schedule.Cron, "0 9 * * 1")
	}
}

func TestCreateSchedule_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	reqBody := CreateScheduleRequest{
		Prompt: "テストプロンプト",
		Cron:   "0 9 * * 1",
	}
	schedule, err := client.CreateSchedule(reqBody)
	if err != nil {
		t.Fatalf("CreateSchedule() error = %v", err)
	}

	if schedule.ScheduleID != "schedule-new-mock" {
		t.Errorf("CreateSchedule() ScheduleID = %s, want %s", schedule.ScheduleID, "schedule-new-mock")
	}
	if schedule.Prompt != "テストプロンプト" {
		t.Errorf("CreateSchedule() Prompt = %s, want %s", schedule.Prompt, "テストプロンプト")
	}
}

func TestDeleteSchedule_Mock(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")
	err := client.DeleteSchedule("schedule-mock-1")
	if err != nil {
		t.Errorf("DeleteSchedule() error = %v", err)
	}
}
