package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// Devin API のベースURL
	baseURL = "https://api.devin.ai/v1"
)

// DevinClient は Devin API とのやり取りを行うクライアント
type DevinClient struct {
	APIKey     string
	HTTPClient *http.Client
}

// Knowledge はDevinのナレッジリソースを表す構造体
type Knowledge struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Body               string    `json:"body,omitempty"`
	TriggerDescription string    `json:"trigger_description,omitempty"`
	ParentFolderID     string    `json:"parent_folder_id,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

// ListKnowledgeResponse はナレッジリスト取得APIのレスポンス
type ListKnowledgeResponse struct {
	Knowledge []KnowledgeItem `json:"knowledge"`
	Folders   []FolderItem    `json:"folders"`
}

// KnowledgeItem はナレッジ項目を表す構造体
type KnowledgeItem struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Body               string    `json:"body,omitempty"`
	TriggerDescription string    `json:"trigger_description,omitempty"`
	ParentFolderID     string    `json:"parent_folder_id,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

// FolderItem はフォルダ項目を表す構造体
type FolderItem struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateKnowledgeRequest はナレッジ作成APIのリクエスト
type CreateKnowledgeRequest struct {
	Name               string `json:"name"`
	Body               string `json:"body,omitempty"`
	ParentFolderID     string `json:"parent_folder_id,omitempty"`
	TriggerDescription string `json:"trigger_description,omitempty"`
}

// UpdateKnowledgeRequest はナレッジ更新APIのリクエスト
type UpdateKnowledgeRequest struct {
	Name               string `json:"name"`
	Body               string `json:"body,omitempty"`
	ParentFolderID     string `json:"parent_folder_id,omitempty"`
	TriggerDescription string `json:"trigger_description,omitempty"`
}

// ErrorResponse はAPIエラー時のレスポンスを表します
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// NewClient は新しいDevinClientを作成する
func NewClient(apiKey string) *DevinClient {
	return &DevinClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// sendRequest は共通のリクエスト送信関数
func (c *DevinClient) sendRequest(method, path string, body interface{}) ([]byte, error) {
	url := baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("リクエストボディのJSONエンコードに失敗しました: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの作成に失敗しました: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの実行に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み込みに失敗しました: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("API エラー: %s (%s)", errResp.Error.Message, errResp.Error.Type)
		}
		return nil, fmt.Errorf("API エラー: ステータスコード %d", resp.StatusCode)
	}

	return respBody, nil
}

// ListKnowledge はナレッジの一覧を取得する
func (c *DevinClient) ListKnowledge() (*ListKnowledgeResponse, error) {
	// デモ向けにモックデータを返す（開発・テスト用）
	if IsMockClient(c.APIKey) {
		return GetMockKnowledgeList(), nil
	}

	// 通常の処理
	respBody, err := c.sendRequest("GET", "/knowledge", nil)
	if err != nil {
		return nil, err
	}

	var response ListKnowledgeResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("レスポンスのJSONデコードに失敗しました: %w", err)
	}

	return &response, nil
}

// GetKnowledge は特定のIDのナレッジを取得する
// 注意: Devin API には現在、個別のナレッジを取得する専用のエンドポイントが
// 明示的に公開されていないため、List API を使用して特定のIDのナレッジを抽出しています
func (c *DevinClient) GetKnowledge(id string) (*Knowledge, error) {
	// デモ向けにモックデータを返す（開発・テスト用）
	if IsMockClient(c.APIKey) {
		return GetMockKnowledge(id)
	}

	// 通常の処理
	// リスト取得 API を使用して全てのナレッジを取得
	response, err := c.ListKnowledge()
	if err != nil {
		return nil, fmt.Errorf("ナレッジリストの取得中にエラーが発生しました: %w", err)
	}

	// 指定された ID に一致するナレッジを検索
	for _, item := range response.Knowledge {
		if item.ID == id {
			// 見つかったナレッジを返す
			return &Knowledge{
				ID:                 item.ID,
				Name:               item.Name,
				Body:               item.Body,
				TriggerDescription: item.TriggerDescription,
				ParentFolderID:     item.ParentFolderID,
				CreatedAt:          item.CreatedAt,
			}, nil
		}
	}

	return nil, fmt.Errorf("指定されたID '%s' のナレッジが見つかりませんでした", id)
}

// CreateKnowledge は新しいナレッジを作成する
func (c *DevinClient) CreateKnowledge(name, body string, triggerDescription string, parentFolderID string) (*Knowledge, error) {
	// デモ向けにモックデータを返す（開発・テスト用）
	if IsMockClient(c.APIKey) {
		return CreateMockKnowledge(name, body, triggerDescription, parentFolderID), nil
	}

	// 通常の処理
	reqBody := CreateKnowledgeRequest{
		Name:               name,
		Body:               body,
		TriggerDescription: triggerDescription,
		ParentFolderID:     parentFolderID,
	}

	respBody, err := c.sendRequest("POST", "/knowledge", reqBody)
	if err != nil {
		return nil, err
	}

	var knowledge Knowledge
	if err := json.Unmarshal(respBody, &knowledge); err != nil {
		return nil, fmt.Errorf("レスポンスのJSONデコードに失敗しました: %w", err)
	}

	return &knowledge, nil
}

// UpdateKnowledge はナレッジを更新する
func (c *DevinClient) UpdateKnowledge(id, name, body string, triggerDescription string, parentFolderID string) (*Knowledge, error) {
	// デモ向けにモックデータを返す（開発・テスト用）
	if IsMockClient(c.APIKey) {
		return UpdateMockKnowledge(id, name, body, triggerDescription, parentFolderID), nil
	}

	// 通常の処理
	reqBody := UpdateKnowledgeRequest{
		Name:               name,
		Body:               body,
		TriggerDescription: triggerDescription,
		ParentFolderID:     parentFolderID,
	}

	path := fmt.Sprintf("/knowledge/%s", id)
	respBody, err := c.sendRequest("PUT", path, reqBody)
	if err != nil {
		return nil, err
	}

	var knowledge Knowledge
	if err := json.Unmarshal(respBody, &knowledge); err != nil {
		return nil, fmt.Errorf("レスポンスのJSONデコードに失敗しました: %w", err)
	}

	return &knowledge, nil
}

// DeleteKnowledge はナレッジを削除する
func (c *DevinClient) DeleteKnowledge(id string) error {
	// デモ向けにモックデータを返す（開発・テスト用）
	if IsMockClient(c.APIKey) {
		return nil
	}

	// 通常の処理
	path := fmt.Sprintf("/knowledge/%s", id)
	_, err := c.sendRequest("DELETE", path, nil)
	return err
}
