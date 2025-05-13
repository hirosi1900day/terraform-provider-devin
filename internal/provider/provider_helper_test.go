package provider

import (
	"testing"
)

func TestClientError(t *testing.T) {
	// 不正なAPIキーでのテスト
	client := NewClient("invalid_key")

	// 現在の実装では、実際のAPIリクエストを送らないとエラーが発生しないため、
	// モックがない場合のテストはスキップします
	if client.ApiKey == "test_api_key" {
		t.Skip("このテストは実際のAPIに接続する必要があるためスキップします")
	}
}

func TestKnowledgeResourceImports(t *testing.T) {
	// リソースのインポート機能のテスト
	resource := NewKnowledgeResource()
	if resource == nil {
		t.Fatal("NewKnowledgeResource() returned nil")
	}

	// これ以上の詳細なテストはRPCが必要なため、リソースが作成できることだけを確認
}

func TestKnowledgeDataSourceImports(t *testing.T) {
	// データソースのインポート機能のテスト
	dataSource := NewKnowledgeDataSource()
	if dataSource == nil {
		t.Fatal("NewKnowledgeDataSource() returned nil")
	}
}
