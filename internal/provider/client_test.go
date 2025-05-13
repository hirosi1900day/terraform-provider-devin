package provider

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")
	if client == nil {
		t.Fatalf("NewClient() returned nil")
	}
	if client.ApiKey != "test-api-key" {
		t.Errorf("NewClient() API key = %s, want %s", client.ApiKey, "test-api-key")
	}
	if client.HttpClient == nil {
		t.Errorf("NewClient() HttpClient is nil")
	}
}

func TestListKnowledge_Mock(t *testing.T) {
	client := NewClient("test_api_key")
	knowledges, err := client.ListKnowledge()
	if err != nil {
		t.Fatalf("ListKnowledge() error = %v", err)
	}

	if len(knowledges) != 2 {
		t.Errorf("ListKnowledge() returned %d items, want 2", len(knowledges))
	}

	// 最初のナレッジの検証
	if knowledges[0].ID != "mock-knowledge-1" {
		t.Errorf("First knowledge ID = %s, want %s", knowledges[0].ID, "mock-knowledge-1")
	}
	if knowledges[0].Name != "モックナレッジ1" {
		t.Errorf("First knowledge Name = %s, want %s", knowledges[0].Name, "モックナレッジ1")
	}

	// 2番目のナレッジの検証
	if knowledges[1].ID != "mock-knowledge-2" {
		t.Errorf("Second knowledge ID = %s, want %s", knowledges[1].ID, "mock-knowledge-2")
	}
	if knowledges[1].Name != "モックナレッジ2" {
		t.Errorf("Second knowledge Name = %s, want %s", knowledges[1].Name, "モックナレッジ2")
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

	// 存在しないIDのテスト
	_, err = client.GetKnowledge("non-existent-id")
	if err == nil {
		t.Errorf("GetKnowledge() with non-existent ID should return error")
	}
}

func TestCreateKnowledge_Mock(t *testing.T) {
	client := NewClient("test_api_key")
	knowledge, err := client.CreateKnowledge("テストナレッジ", "テスト説明")
	if err != nil {
		t.Fatalf("CreateKnowledge() error = %v", err)
	}

	if knowledge.ID != "new-mock-knowledge" {
		t.Errorf("CreateKnowledge() ID = %s, want %s", knowledge.ID, "new-mock-knowledge")
	}
	if knowledge.Name != "テストナレッジ" {
		t.Errorf("CreateKnowledge() Name = %s, want %s", knowledge.Name, "テストナレッジ")
	}
	if knowledge.Description != "テスト説明" {
		t.Errorf("CreateKnowledge() Description = %s, want %s", knowledge.Description, "テスト説明")
	}
}

func TestUpdateKnowledge_Mock(t *testing.T) {
	client := NewClient("test_api_key")
	knowledge, err := client.UpdateKnowledge("mock-knowledge-1", "更新ナレッジ", "更新説明")
	if err != nil {
		t.Fatalf("UpdateKnowledge() error = %v", err)
	}

	if knowledge.ID != "mock-knowledge-1" {
		t.Errorf("UpdateKnowledge() ID = %s, want %s", knowledge.ID, "mock-knowledge-1")
	}
	if knowledge.Name != "更新ナレッジ" {
		t.Errorf("UpdateKnowledge() Name = %s, want %s", knowledge.Name, "更新ナレッジ")
	}
	if knowledge.Description != "更新説明" {
		t.Errorf("UpdateKnowledge() Description = %s, want %s", knowledge.Description, "更新説明")
	}
}

func TestDeleteKnowledge_Mock(t *testing.T) {
	client := NewClient("test_api_key")
	err := client.DeleteKnowledge("mock-knowledge-1")
	if err != nil {
		t.Errorf("DeleteKnowledge() error = %v", err)
	}
}
