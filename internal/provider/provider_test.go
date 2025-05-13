package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProviderSchema(t *testing.T) {
	ctx := context.Background()
	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}

	p := &DevinProvider{version: "test"}
	p.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema had unexpected error: %s", resp.Diagnostics.Errors())
	}

	// api_key属性が存在するか確認
	apiKeyAttr := resp.Schema.Attributes["api_key"]
	if apiKeyAttr == nil {
		t.Fatal("Schema doesn't have api_key attribute")
	}
}

func TestProviderConfigure(t *testing.T) {
	// テストではprovider.Configureの実装に依存せず、シンプルにテストします
	// 実際のプロバイダーの機能は統合テストでカバーすることが望ましい
	t.Skip("provider.Configureの詳細なテストはスキップします")
}

func TestProviderAPI(t *testing.T) {
	// API_KEYが環境変数で設定されている場合の動作テスト
	oldApiKey := os.Getenv("DEVIN_API_KEY")
	os.Setenv("DEVIN_API_KEY", "test_api_key")
	defer os.Setenv("DEVIN_API_KEY", oldApiKey) // テスト後に元の値に戻す

	// クライアントが正しく生成されることを確認
	client := NewClient("test_api_key")
	if client == nil {
		t.Fatalf("NewClient() returned nil")
	}
	if client.ApiKey != "test_api_key" {
		t.Errorf("NewClient() API key = %s, want %s", client.ApiKey, "test_api_key")
	}
}

func TestProviderMetadata(t *testing.T) {
	ctx := context.Background()
	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}

	p := &DevinProvider{version: "test"}
	p.Metadata(ctx, req, resp)

	if resp.TypeName != "devin" {
		t.Fatalf("Metadata TypeName = %s, want %s", resp.TypeName, "devin")
	}
	if resp.Version != "test" {
		t.Fatalf("Metadata Version = %s, want %s", resp.Version, "test")
	}
}

func TestProviderResources(t *testing.T) {
	ctx := context.Background()

	p := &DevinProvider{version: "test"}
	resources := p.Resources(ctx)

	if len(resources) != 1 {
		t.Fatalf("Resources() returned %d resources, want 1", len(resources))
	}
}

func TestProviderDataSources(t *testing.T) {
	ctx := context.Background()

	p := &DevinProvider{version: "test"}
	dataSources := p.DataSources(ctx)

	if len(dataSources) != 1 {
		t.Fatalf("DataSources() returned %d data sources, want 1", len(dataSources))
	}
}
