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
			ID:                 "mock-knowledge-1",
			Name:               "モックナレッジ1",
			Body:               "これはテスト用のモックナレッジです",
			TriggerDescription: "テスト用トリガーの説明",
			ParentFolderID:     "mock-folder-1",
			CreatedAt:          time.Now().Add(-24 * time.Hour),
		}, nil
	case "mock-knowledge-2":
		return &Knowledge{
			ID:                 "mock-knowledge-2",
			Name:               "モックナレッジ2",
			Body:               "これは別のテスト用のモックナレッジです",
			TriggerDescription: "別のテスト用トリガーの説明",
			ParentFolderID:     "mock-folder-2",
			CreatedAt:          time.Now().Add(-48 * time.Hour),
		}, nil
	case "new-mock-knowledge":
		return &Knowledge{
			ID:                 "new-mock-knowledge",
			Name:               "サンプルナレッジ",
			Body:               "これはTerraformで作成されたサンプルナレッジです",
			TriggerDescription: "Terraformサンプルトリガー",
			ParentFolderID:     "mock-folder-1",
			CreatedAt:          time.Now().Add(-1 * time.Hour),
		}, nil
	case "":
		return &Knowledge{
			ID:                 "mock-knowledge-1",
			Name:               "モックナレッジ1",
			Body:               "これはテスト用のモックナレッジです",
			TriggerDescription: "テスト用トリガーの説明",
			ParentFolderID:     "mock-folder-1",
			CreatedAt:          time.Now().Add(-24 * time.Hour),
		}, nil
	default:
		return nil, fmt.Errorf("ナレッジが見つかりません: ID %s", id)
	}
}

// GetMockKnowledgeList はモックのナレッジリスト一覧を返します
func GetMockKnowledgeList() *ListKnowledgeResponse {
	return &ListKnowledgeResponse{
		Knowledge: []KnowledgeItem{
			{
				ID:                 "mock-knowledge-1",
				Name:               "モックナレッジ1",
				Body:               "これはテスト用のモックナレッジの内容です。",
				TriggerDescription: "テスト用トリガーの説明",
				ParentFolderID:     "mock-folder-1",
				CreatedAt:          time.Now().Add(-24 * time.Hour),
			},
			{
				ID:                 "mock-knowledge-2",
				Name:               "モックナレッジ2",
				Body:               "これは別のテスト用のモックナレッジの内容です。",
				TriggerDescription: "別のテスト用トリガーの説明",
				ParentFolderID:     "mock-folder-2",
				CreatedAt:          time.Now().Add(-48 * time.Hour),
			},
		},
		Folders: []FolderItem{
			{
				ID:          "mock-folder-1",
				Name:        "モックフォルダ1",
				Description: "これはテスト用のモックフォルダです",
				CreatedAt:   time.Now().Add(-72 * time.Hour),
			},
			{
				ID:          "mock-folder-2",
				Name:        "モックフォルダ2",
				Description: "これは別のテスト用のモックフォルダです",
				CreatedAt:   time.Now().Add(-96 * time.Hour),
			},
		},
	}
}

// CreateMockKnowledge は新しいモックナレッジを作成します
func CreateMockKnowledge(name, body string, triggerDescription, parentFolderID string) *Knowledge {
	return &Knowledge{
		ID:                 "new-mock-knowledge",
		Name:               name,
		Body:               body,
		TriggerDescription: triggerDescription,
		ParentFolderID:     parentFolderID,
		CreatedAt:          time.Now(),
	}
}

// UpdateMockKnowledge はモックナレッジを更新します
func UpdateMockKnowledge(id, name, body string, triggerDescription, parentFolderID string) *Knowledge {
	return &Knowledge{
		ID:                 id,
		Name:               name,
		Body:               body,
		TriggerDescription: triggerDescription,
		ParentFolderID:     parentFolderID,
		CreatedAt:          time.Now().Add(-24 * time.Hour),
	}
}

// GetMockFolderByID はモックフォルダリソースをIDで取得します
func GetMockFolderByID(id string) (*FolderItem, error) {
	switch id {
	case "mock-folder-1":
		return &FolderItem{
			ID:          "mock-folder-1",
			Name:        "モックフォルダ1",
			Description: "これはテスト用のモックフォルダです",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
		}, nil
	case "mock-folder-2":
		return &FolderItem{
			ID:          "mock-folder-2",
			Name:        "モックフォルダ2",
			Description: "これは別のテスト用のモックフォルダです",
			CreatedAt:   time.Now().Add(-96 * time.Hour),
		}, nil
	default:
		return nil, fmt.Errorf("フォルダが見つかりません: ID %s", id)
	}
}

// GetMockFolderByName はモックフォルダリソースを名前で取得します
func GetMockFolderByName(name string) (*FolderItem, error) {
	switch name {
	case "モックフォルダ1":
		return &FolderItem{
			ID:          "mock-folder-1",
			Name:        "モックフォルダ1",
			Description: "これはテスト用のモックフォルダです",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
		}, nil
	case "モックフォルダ2":
		return &FolderItem{
			ID:          "mock-folder-2",
			Name:        "モックフォルダ2",
			Description: "これは別のテスト用のモックフォルダです",
			CreatedAt:   time.Now().Add(-96 * time.Hour),
		}, nil
	default:
		return nil, fmt.Errorf("フォルダが見つかりません: 名前 %s", name)
	}
}
