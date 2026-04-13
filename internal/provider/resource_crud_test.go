package provider

import (
	"testing"
)

// ===================== Knowledge CRUD Lifecycle =====================

func TestKnowledgeCRUDLifecycle(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// --- Create ---
	created, err := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name:    "テストナレッジ",
		Body:    "テスト内容",
		Trigger: "テストトリガー",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.NoteID == "" {
		t.Fatal("Created note has empty ID")
	}
	if created.Name != "テストナレッジ" {
		t.Errorf("Name = %q, want %q", created.Name, "テストナレッジ")
	}
	if created.Body != "テスト内容" {
		t.Errorf("Body = %q, want %q", created.Body, "テスト内容")
	}
	if created.Trigger != "テストトリガー" {
		t.Errorf("Trigger = %q, want %q", created.Trigger, "テストトリガー")
	}
	if !created.IsEnabled {
		t.Error("IsEnabled should be true")
	}
	if created.AccessType != "org" {
		t.Errorf("AccessType = %q, want %q", created.AccessType, "org")
	}
	if created.CreatedAt == 0 {
		t.Error("CreatedAt should be set")
	}
	if created.UpdatedAt == 0 {
		t.Error("UpdatedAt should be set")
	}
	noteID := created.NoteID

	// --- Read ---
	read, err := client.GetKnowledgeNote(noteID)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if read.NoteID != noteID {
		t.Errorf("Read ID = %q, want %q", read.NoteID, noteID)
	}
	if read.Name != "テストナレッジ" {
		t.Errorf("Read Name = %q, want %q", read.Name, "テストナレッジ")
	}

	// --- Read via List ---
	list, err := client.ListKnowledgeNotes()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	found := false
	for _, n := range list {
		if n.NoteID == noteID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Created note %q not found in List result", noteID)
	}

	// --- Update ---
	updated, err := client.UpdateKnowledgeNote(noteID, UpdateKnowledgeNoteRequest{
		Name:    "更新ナレッジ",
		Body:    "更新内容",
		Trigger: "更新トリガー",
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.NoteID != noteID {
		t.Errorf("Updated ID = %q, want %q", updated.NoteID, noteID)
	}
	if updated.Name != "更新ナレッジ" {
		t.Errorf("Updated Name = %q, want %q", updated.Name, "更新ナレッジ")
	}
	if updated.Body != "更新内容" {
		t.Errorf("Updated Body = %q, want %q", updated.Body, "更新内容")
	}
	if updated.CreatedAt == 0 {
		t.Error("Updated CreatedAt should be preserved")
	}

	// --- Read after Update ---
	readAfterUpdate, err := client.GetKnowledgeNote(noteID)
	if err != nil {
		t.Fatalf("Read after update failed: %v", err)
	}
	if readAfterUpdate.Name != "更新ナレッジ" {
		t.Errorf("Read after update Name = %q, want %q", readAfterUpdate.Name, "更新ナレッジ")
	}

	// --- Delete ---
	err = client.DeleteKnowledgeNote(noteID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// --- Read after Delete (should fail) ---
	_, err = client.GetKnowledgeNote(noteID)
	if err == nil {
		t.Fatal("Read after delete should fail, but got nil error")
	}
}

func TestKnowledgeCRUD_WithOptionalFields(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	repo := "owner/repo"
	created, err := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name:       "Optional Fields",
		Body:       "Body content",
		Trigger:    "Trigger text",
		PinnedRepo: &repo,
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.PinnedRepo == nil || *created.PinnedRepo != "owner/repo" {
		t.Errorf("PinnedRepo = %v, want %q", created.PinnedRepo, "owner/repo")
	}
}

func TestKnowledgeCRUD_MultipleResources(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	_, err := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "Note 1", Body: "Body 1", Trigger: "Trigger 1",
	})
	if err != nil {
		t.Fatalf("Create 1 failed: %v", err)
	}

	_, err = client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "Note 2", Body: "Body 2", Trigger: "Trigger 2",
	})
	if err != nil {
		t.Fatalf("Create 2 failed: %v", err)
	}

	list, err := client.ListKnowledgeNotes()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("List length = %d, want 2", len(list))
	}
}

// ===================== Playbook CRUD Lifecycle =====================

func TestPlaybookCRUDLifecycle(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// --- Create ---
	created, err := client.CreatePlaybook(CreatePlaybookRequest{
		Title: "テストPlaybook",
		Body:  "Playbookの内容",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.PlaybookID == "" {
		t.Fatal("Created playbook has empty ID")
	}
	if created.Title != "テストPlaybook" {
		t.Errorf("Title = %q, want %q", created.Title, "テストPlaybook")
	}
	if created.Body != "Playbookの内容" {
		t.Errorf("Body = %q, want %q", created.Body, "Playbookの内容")
	}
	if created.OrgID != "org-mock" {
		t.Errorf("OrgID = %q, want %q", created.OrgID, "org-mock")
	}
	if created.AccessType != "org" {
		t.Errorf("AccessType = %q, want %q", created.AccessType, "org")
	}
	if created.CreatedAt == 0 {
		t.Error("CreatedAt should be set")
	}
	if created.CreatedBy != "user-1" {
		t.Errorf("CreatedBy = %q, want %q", created.CreatedBy, "user-1")
	}
	playbookID := created.PlaybookID

	// --- Read ---
	read, err := client.GetPlaybook(playbookID)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if read.PlaybookID != playbookID {
		t.Errorf("Read ID = %q, want %q", read.PlaybookID, playbookID)
	}
	if read.Title != "テストPlaybook" {
		t.Errorf("Read Title = %q, want %q", read.Title, "テストPlaybook")
	}

	// --- Update ---
	updated, err := client.UpdatePlaybook(playbookID, UpdatePlaybookRequest{
		Title: "更新Playbook",
		Body:  "更新内容",
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.PlaybookID != playbookID {
		t.Errorf("Updated ID = %q, want %q", updated.PlaybookID, playbookID)
	}
	if updated.Title != "更新Playbook" {
		t.Errorf("Updated Title = %q, want %q", updated.Title, "更新Playbook")
	}
	if updated.Body != "更新内容" {
		t.Errorf("Updated Body = %q, want %q", updated.Body, "更新内容")
	}

	// --- Read after Update ---
	readAfterUpdate, err := client.GetPlaybook(playbookID)
	if err != nil {
		t.Fatalf("Read after update failed: %v", err)
	}
	if readAfterUpdate.Title != "更新Playbook" {
		t.Errorf("Read after update Title = %q, want %q", readAfterUpdate.Title, "更新Playbook")
	}

	// --- Delete ---
	err = client.DeletePlaybook(playbookID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// --- Read after Delete (should fail) ---
	_, err = client.GetPlaybook(playbookID)
	if err == nil {
		t.Fatal("Read after delete should fail, but got nil error")
	}
}

// ===================== Secret Create/Delete Lifecycle =====================

func TestSecretCreateDeleteLifecycle(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// --- Create ---
	created, err := client.CreateSecret(CreateSecretRequest{
		Key:         "TEST_SECRET",
		Type:        "key-value",
		Value:       "secret-value-123",
		IsSensitive: true,
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.SecretID == "" {
		t.Fatal("Created secret has empty ID")
	}
	if created.Key != "TEST_SECRET" {
		t.Errorf("Key = %q, want %q", created.Key, "TEST_SECRET")
	}
	if created.CreatedAt == 0 {
		t.Error("CreatedAt should be set")
	}
	secretID := created.SecretID

	// --- Read by ID ---
	read, err := client.GetSecretByID(secretID)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if read.SecretID != secretID {
		t.Errorf("Read ID = %q, want %q", read.SecretID, secretID)
	}
	if read.Key != "TEST_SECRET" {
		t.Errorf("Read Key = %q, want %q", read.Key, "TEST_SECRET")
	}

	// --- Read via List ---
	list, err := client.ListSecrets()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	found := false
	for _, s := range list {
		if s.SecretID == secretID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Created secret %q not found in List result", secretID)
	}

	// --- Delete ---
	err = client.DeleteSecret(secretID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// --- Read after Delete (should fail) ---
	_, err = client.GetSecretByID(secretID)
	if err == nil {
		t.Fatal("Read after delete should fail, but got nil error")
	}
}

func TestSecretCreateMultiple(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	_, err := client.CreateSecret(CreateSecretRequest{Key: "SECRET_A", Type: "key-value", Value: "aaa", IsSensitive: true})
	if err != nil {
		t.Fatalf("Create A failed: %v", err)
	}
	_, err = client.CreateSecret(CreateSecretRequest{Key: "SECRET_B", Type: "key-value", Value: "bbb", IsSensitive: true})
	if err != nil {
		t.Fatalf("Create B failed: %v", err)
	}

	list, err := client.ListSecrets()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("List length = %d, want 2", len(list))
	}
}

// ===================== Schedule CRUD Lifecycle =====================

func TestScheduleCRUDLifecycle(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// --- Create ---
	created, err := client.CreateSchedule(CreateScheduleRequest{
		Name:      "テストスケジュール",
		Prompt:    "テストプロンプト",
		Frequency: "0 9 * * 1",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.ScheduledSessionID == "" {
		t.Fatal("Created schedule has empty ID")
	}
	if created.Prompt != "テストプロンプト" {
		t.Errorf("Prompt = %q, want %q", created.Prompt, "テストプロンプト")
	}
	if *created.Frequency != "0 9 * * 1" {
		t.Errorf("Frequency = %q, want %q", *created.Frequency, "0 9 * * 1")
	}
	if !created.Enabled {
		t.Error("Enabled should be true")
	}
	if created.CreatedAt == "" {
		t.Error("CreatedAt should be set")
	}
	scheduleID := created.ScheduledSessionID

	// --- Read ---
	read, err := client.GetSchedule(scheduleID)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if read.ScheduledSessionID != scheduleID {
		t.Errorf("Read ID = %q, want %q", read.ScheduledSessionID, scheduleID)
	}
	if read.Prompt != "テストプロンプト" {
		t.Errorf("Read Prompt = %q, want %q", read.Prompt, "テストプロンプト")
	}

	// --- Update (PATCH) ---
	newPrompt := "更新プロンプト"
	newFrequency := "0 10 * * *"
	updated, err := client.UpdateSchedule(scheduleID, UpdateScheduleRequest{
		Prompt:    &newPrompt,
		Frequency: &newFrequency,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.ScheduledSessionID != scheduleID {
		t.Errorf("Updated ID = %q, want %q", updated.ScheduledSessionID, scheduleID)
	}
	if updated.Prompt != "更新プロンプト" {
		t.Errorf("Updated Prompt = %q, want %q", updated.Prompt, "更新プロンプト")
	}
	if *updated.Frequency != "0 10 * * *" {
		t.Errorf("Updated Frequency = %q, want %q", *updated.Frequency, "0 10 * * *")
	}

	// --- Read after Update ---
	readAfterUpdate, err := client.GetSchedule(scheduleID)
	if err != nil {
		t.Fatalf("Read after update failed: %v", err)
	}
	if readAfterUpdate.Prompt != "更新プロンプト" {
		t.Errorf("Read after update Prompt = %q, want %q", readAfterUpdate.Prompt, "更新プロンプト")
	}
	if *readAfterUpdate.Frequency != "0 10 * * *" {
		t.Errorf("Read after update Frequency = %q, want %q", *readAfterUpdate.Frequency, "0 10 * * *")
	}

	// --- Delete ---
	err = client.DeleteSchedule(scheduleID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// --- Read after Delete (should fail) ---
	_, err = client.GetSchedule(scheduleID)
	if err == nil {
		t.Fatal("Read after delete should fail, but got nil error")
	}
}

func TestScheduleCRUD_WithPlaybook(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// Create a playbook first
	pb, err := client.CreatePlaybook(CreatePlaybookRequest{
		Title: "Linked Playbook",
		Body:  "Body",
	})
	if err != nil {
		t.Fatalf("CreatePlaybook failed: %v", err)
	}

	// Create schedule with playbook_id
	created, err := client.CreateSchedule(CreateScheduleRequest{
		Name:       "Playbook付きスケジュール",
		Prompt:     "Playbook付きスケジュール",
		Frequency:  "0 9 * * 1",
		PlaybookID: pb.PlaybookID,
	})
	if err != nil {
		t.Fatalf("Create schedule failed: %v", err)
	}
	if created.Playbook == nil || created.Playbook.PlaybookID != pb.PlaybookID {
		t.Errorf("PlaybookID mismatch")
	}

	// Read and verify
	read, err := client.GetSchedule(created.ScheduledSessionID)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if read.Playbook == nil || read.Playbook.PlaybookID != pb.PlaybookID {
		t.Errorf("Read PlaybookID mismatch")
	}

	// Update: remove playbook_id via PATCH
	newPlaybookID := ""
	updated, err := client.UpdateSchedule(created.ScheduledSessionID, UpdateScheduleRequest{
		PlaybookID: &newPlaybookID,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Playbook != nil {
		t.Errorf("Updated Playbook should be nil")
	}
}

func TestScheduleCRUD_PartialUpdate(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	created, err := client.CreateSchedule(CreateScheduleRequest{
		Name:      "Partial Update Test",
		Prompt:    "Original Prompt",
		Frequency: "0 9 * * 1",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Only update prompt, cron should remain unchanged
	newPrompt := "Updated Prompt Only"
	updated, err := client.UpdateSchedule(created.ScheduledSessionID, UpdateScheduleRequest{
		Prompt: &newPrompt,
	})
	if err != nil {
		t.Fatalf("Partial update failed: %v", err)
	}
	if updated.Prompt != "Updated Prompt Only" {
		t.Errorf("Prompt = %q, want %q", updated.Prompt, "Updated Prompt Only")
	}
	if *updated.Frequency != "0 9 * * 1" {
		t.Errorf("Frequency should be unchanged: got %q, want %q", *updated.Frequency, "0 9 * * 1")
	}
}

// ===================== Folder DataSource =====================

func TestFolderDataSource_ByID_CRUD(t *testing.T) {
	folder, err := GetMockFolderByID("folder-mock-1")
	if err != nil {
		t.Fatalf("GetMockFolderByID failed: %v", err)
	}
	if folder.FolderID != "folder-mock-1" {
		t.Errorf("FolderID = %q, want %q", folder.FolderID, "folder-mock-1")
	}
	if folder.Name != "モックフォルダ1" {
		t.Errorf("Name = %q, want %q", folder.Name, "モックフォルダ1")
	}
	if folder.NoteCount != 5 {
		t.Errorf("NoteCount = %d, want %d", folder.NoteCount, 5)
	}
}

func TestFolderDataSource_ByName_CRUD(t *testing.T) {
	folder, err := GetMockFolderByName("モックフォルダ2")
	if err != nil {
		t.Fatalf("GetMockFolderByName failed: %v", err)
	}
	if folder.FolderID != "folder-mock-2" {
		t.Errorf("FolderID = %q, want %q", folder.FolderID, "folder-mock-2")
	}
}

func TestFolderDataSource_NotFound(t *testing.T) {
	_, err := GetMockFolderByID("nonexistent")
	if err == nil {
		t.Fatal("GetMockFolderByID should fail for nonexistent folder")
	}
	_, err = GetMockFolderByName("nonexistent")
	if err == nil {
		t.Fatal("GetMockFolderByName should fail for nonexistent folder")
	}
}

func TestFolderDataSource_ViaClient(t *testing.T) {
	client := NewClient("test_api_key", "org-mock")

	folder, err := client.GetFolderByID("folder-mock-1")
	if err != nil {
		t.Fatalf("GetFolderByID failed: %v", err)
	}
	if folder.Name != "モックフォルダ1" {
		t.Errorf("Name = %q, want %q", folder.Name, "モックフォルダ1")
	}

	folder, err = client.GetFolderByName("モックフォルダ2")
	if err != nil {
		t.Fatalf("GetFolderByName failed: %v", err)
	}
	if folder.FolderID != "folder-mock-2" {
		t.Errorf("FolderID = %q, want %q", folder.FolderID, "folder-mock-2")
	}

	folders, err := client.ListFolders()
	if err != nil {
		t.Fatalf("ListFolders failed: %v", err)
	}
	if len(folders) != 2 {
		t.Errorf("ListFolders length = %d, want 2", len(folders))
	}
}

// ===================== Knowledge Cache Lifecycle =====================

func TestKnowledgeCache_CreateInvalidates(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	created, err := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "Cached", Body: "Body", Trigger: "Trigger",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Read to verify Create stored data properly
	_, err = client.GetKnowledgeNote(created.NoteID)
	if err != nil {
		t.Fatalf("GetKnowledgeNote after Create failed: %v", err)
	}
}

func TestKnowledgeCache_DeleteInvalidates(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	created, err := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "ToDelete", Body: "Body", Trigger: "Trigger",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = client.DeleteKnowledgeNote(created.NoteID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = client.GetKnowledgeNote(created.NoteID)
	if err == nil {
		t.Error("Should not be able to read deleted note")
	}
}

// ===================== Cross-resource Isolation =====================

func TestMockStoreIsolation(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	_, err := client.CreateKnowledgeNote(CreateKnowledgeNoteRequest{
		Name: "Note", Body: "B", Trigger: "T",
	})
	if err != nil {
		t.Fatalf("Create knowledge failed: %v", err)
	}

	_, err = client.CreatePlaybook(CreatePlaybookRequest{
		Title: "PB", Body: "B",
	})
	if err != nil {
		t.Fatalf("Create playbook failed: %v", err)
	}

	_, err = client.CreateSecret(CreateSecretRequest{
		Key: "SEC", Type: "key-value", Value: "val", IsSensitive: true,
	})
	if err != nil {
		t.Fatalf("Create secret failed: %v", err)
	}

	_, err = client.CreateSchedule(CreateScheduleRequest{
		Name: "Sched", Prompt: "P", Frequency: "0 * * * *",
	})
	if err != nil {
		t.Fatalf("Create schedule failed: %v", err)
	}

	// Verify each resource type has exactly 1 item
	notes, _ := client.ListKnowledgeNotes()
	if len(notes) != 1 {
		t.Errorf("Knowledge count = %d, want 1", len(notes))
	}

	secrets, _ := client.ListSecrets()
	if len(secrets) != 1 {
		t.Errorf("Secret count = %d, want 1", len(secrets))
	}

	// ResetMockStore clears everything
	ResetMockStore()
	notes, _ = client.ListKnowledgeNotes()
	if len(notes) != 0 {
		t.Errorf("After reset, knowledge count = %d, want 0", len(notes))
	}
	secrets, _ = client.ListSecrets()
	if len(secrets) != 0 {
		t.Errorf("After reset, secret count = %d, want 0", len(secrets))
	}
}

// ===================== Error Handling =====================

func TestReadNonexistentResources(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	_, err := client.GetKnowledgeNote("nonexistent")
	if err == nil {
		t.Error("GetKnowledgeNote should fail for nonexistent")
	}

	_, err = client.GetPlaybook("nonexistent")
	if err == nil {
		t.Error("GetPlaybook should fail for nonexistent")
	}

	_, err = client.GetSecretByID("nonexistent")
	if err == nil {
		t.Error("GetSecretByID should fail for nonexistent")
	}

	_, err = client.GetSchedule("nonexistent")
	if err == nil {
		t.Error("GetSchedule should fail for nonexistent")
	}
}

func TestDeleteIdempotent(t *testing.T) {
	ResetMockStore()
	client := NewClient("test_api_key", "org-mock")

	// Delete nonexistent resources should not panic
	client.DeleteKnowledgeNote("nonexistent")
	client.DeletePlaybook("nonexistent")
	client.DeleteSecret("nonexistent")
	client.DeleteSchedule("nonexistent")
}
