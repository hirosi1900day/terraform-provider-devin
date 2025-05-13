package provider

import (
	"fmt"
	"time"
)

// IsMockClient は指定されたAPIキーがモック用かどうかを確認します
func IsMockClient(apiKey string) bool {
	return apiKey == "test_api_key"
}

// MockClient はテスト用のDevinClientの振る舞いを提供します
type MockClient struct {
	// 将来的にはモックのカスタマイズを可能にするフィールドを追加できます
}

// GetMockKnowledge はモックのナレッジリソースを取得します
func GetMockKnowledge(id string) (*Knowledge, error) {
	switch id {
	case "mock-knowledge-1":
		return &Knowledge{
			ID:          "mock-knowledge-1",
			Name:        "モックナレッジ1",
			Description: "これはテスト用のモックナレッジです",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now(),
		}, nil
	case "mock-knowledge-2":
		return &Knowledge{
			ID:          "mock-knowledge-2",
			Name:        "モックナレッジ2",
			Description: "これは別のテスト用のモックナレッジです",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now(),
		}, nil
	case "new-mock-knowledge":
		return &Knowledge{
			ID:          "new-mock-knowledge",
			Name:        "サンプルナレッジ",
			Description: "これはTerraformで作成されたサンプルナレッジです",
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			UpdatedAt:   time.Now(),
		}, nil
	case "":
		return &Knowledge{
			ID:          "mock-knowledge-1",
			Name:        "モックナレッジ1",
			Description: "これはテスト用のモックナレッジです",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now(),
		}, nil
	default:
		return nil, fmt.Errorf("ナレッジが見つかりません: ID %s", id)
	}
}

// GetMockKnowledgeList はモックのナレッジリスト一覧を返します
func GetMockKnowledgeList() []Knowledge {
	return []Knowledge{
		{
			ID:          "mock-knowledge-1",
			Name:        "モックナレッジ1",
			Description: "これはテスト用のモックナレッジです",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "mock-knowledge-2",
			Name:        "モックナレッジ2",
			Description: "これは別のテスト用のモックナレッジです",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}
}

// CreateMockKnowledge は新しいモックナレッジを作成します
func CreateMockKnowledge(name, description string) *Knowledge {
	return &Knowledge{
		ID:          "new-mock-knowledge",
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateMockKnowledge はモックナレッジを更新します
func UpdateMockKnowledge(id, name, description string) *Knowledge {
	return &Knowledge{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}
}
