package provider

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")
	if client == nil {
		t.Fatalf("NewClient() returned nil")
	}
	if client.APIKey != "test-api-key" {
		t.Errorf("NewClient() API key = %s, want %s", client.APIKey, "test-api-key")
	}
	if client.HTTPClient == nil {
		t.Errorf("NewClient() HTTPClient is nil")
	}
}

func TestListKnowledge_Mock(t *testing.T) {
	client := NewClient("test_api_key")
	response, err := client.ListKnowledge()
	if err != nil {
		t.Fatalf("ListKnowledge() error = %v", err)
	}

	if len(response.Knowledge) != 2 {
		t.Errorf("ListKnowledge() returned %d items, want 2", len(response.Knowledge))
	}

	if len(response.Folders) != 2 {
		t.Errorf("ListKnowledge() returned %d folders, want 2", len(response.Folders))
	}

	// 最初のナレッジの検証
	if response.Knowledge[0].ID != "mock-knowledge-1" {
		t.Errorf("First knowledge ID = %s, want %s", response.Knowledge[0].ID, "mock-knowledge-1")
	}
	if response.Knowledge[0].Name != "モックナレッジ1" {
		t.Errorf("First knowledge Name = %s, want %s", response.Knowledge[0].Name, "モックナレッジ1")
	}

	// 2番目のナレッジの検証
	if response.Knowledge[1].ID != "mock-knowledge-2" {
		t.Errorf("Second knowledge ID = %s, want %s", response.Knowledge[1].ID, "mock-knowledge-2")
	}
	if response.Knowledge[1].Name != "モックナレッジ2" {
		t.Errorf("Second knowledge Name = %s, want %s", response.Knowledge[1].Name, "モックナレッジ2")
	}

	// フォルダーの検証
	if response.Folders[0].ID != "mock-folder-1" {
		t.Errorf("First folder ID = %s, want %s", response.Folders[0].ID, "mock-folder-1")
	}
	if response.Folders[0].Name != "モックフォルダ1" {
		t.Errorf("First folder Name = %s, want %s", response.Folders[0].Name, "モックフォルダ1")
	}
}

func TestGetKnowledge_Mock(t *testing.T) {
	client := NewClient("test_api_key")

	// 既存のIDの取得テスト
	knowledge, err := client.GetKnowledge("mock-knowledge-1")
	if err != nil {
		t.Fatalf("GetKnowledge() error = %v", err)
	}

	if knowledge.ID != "mock-knowledge-1" {
		t.Errorf("GetKnowledge() ID = %s, want %s", knowledge.ID, "mock-knowledge-1")
	}
	if knowledge.Name != "モックナレッジ1" {
		t.Errorf("GetKnowledge() Name = %s, want %s", knowledge.Name, "モックナレッジ1")
	}
	if knowledge.Body != "これはテスト用のモックナレッジです" {
		t.Errorf("GetKnowledge() Body = %s, want %s", knowledge.Body, "これはテスト用のモックナレッジです")
	}
	if knowledge.TriggerDescription != "テスト用トリガーの説明" {
		t.Errorf("GetKnowledge() TriggerDescription = %s, want %s", knowledge.TriggerDescription, "テスト用トリガーの説明")
	}
	if knowledge.ParentFolderID != "mock-folder-1" {
		t.Errorf("GetKnowledge() ParentFolderID = %s, want %s", knowledge.ParentFolderID, "mock-folder-1")
	}

	// 存在しないIDのテスト
	_, err = client.GetKnowledge("non-existent-id")
	if err == nil {
		t.Errorf("GetKnowledge() with non-existent ID should return error")
	}
}

func TestCreateKnowledge_Mock(t *testing.T) {
	client := NewClient("test_api_key")
	knowledge, err := client.CreateKnowledge("テストナレッジ", "テスト内容", "テストトリガー", "test-folder-id")
	if err != nil {
		t.Fatalf("CreateKnowledge() error = %v", err)
	}

	if knowledge.ID != "new-mock-knowledge" {
		t.Errorf("CreateKnowledge() ID = %s, want %s", knowledge.ID, "new-mock-knowledge")
	}
	if knowledge.Name != "テストナレッジ" {
		t.Errorf("CreateKnowledge() Name = %s, want %s", knowledge.Name, "テストナレッジ")
	}
	if knowledge.Body != "テスト内容" {
		t.Errorf("CreateKnowledge() Body = %s, want %s", knowledge.Body, "テスト内容")
	}
	if knowledge.TriggerDescription != "テストトリガー" {
		t.Errorf("CreateKnowledge() TriggerDescription = %s, want %s", knowledge.TriggerDescription, "テストトリガー")
	}
	if knowledge.ParentFolderID != "test-folder-id" {
		t.Errorf("CreateKnowledge() ParentFolderID = %s, want %s", knowledge.ParentFolderID, "test-folder-id")
	}
}

func TestUpdateKnowledge_Mock(t *testing.T) {
	client := NewClient("test_api_key")
	knowledge, err := client.UpdateKnowledge("mock-knowledge-1", "更新ナレッジ", "更新内容", "更新トリガー", "updated-folder-id")
	if err != nil {
		t.Fatalf("UpdateKnowledge() error = %v", err)
	}

	if knowledge.ID != "mock-knowledge-1" {
		t.Errorf("UpdateKnowledge() ID = %s, want %s", knowledge.ID, "mock-knowledge-1")
	}
	if knowledge.Name != "更新ナレッジ" {
		t.Errorf("UpdateKnowledge() Name = %s, want %s", knowledge.Name, "更新ナレッジ")
	}
	if knowledge.Body != "更新内容" {
		t.Errorf("UpdateKnowledge() Body = %s, want %s", knowledge.Body, "更新内容")
	}
	if knowledge.TriggerDescription != "更新トリガー" {
		t.Errorf("UpdateKnowledge() TriggerDescription = %s, want %s", knowledge.TriggerDescription, "更新トリガー")
	}
	if knowledge.ParentFolderID != "updated-folder-id" {
		t.Errorf("UpdateKnowledge() ParentFolderID = %s, want %s", knowledge.ParentFolderID, "updated-folder-id")
	}
}

func TestDeleteKnowledge_Mock(t *testing.T) {
	client := NewClient("test_api_key")
	err := client.DeleteKnowledge("mock-knowledge-1")
	if err != nil {
		t.Errorf("DeleteKnowledge() error = %v", err)
	}
}
