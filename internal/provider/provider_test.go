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

	// Check if api_key attribute exists
	apiKeyAttr := resp.Schema.Attributes["api_key"]
	if apiKeyAttr == nil {
		t.Fatal("Schema doesn't have api_key attribute")
	}

	// Check if org_id attribute exists
	orgIDAttr := resp.Schema.Attributes["org_id"]
	if orgIDAttr == nil {
		t.Fatal("Schema doesn't have org_id attribute")
	}
}

func TestProviderConfigure(t *testing.T) {
	t.Skip("Skipping detailed test for provider.Configure")
}

func TestProviderAPI(t *testing.T) {
	oldAPIKey := os.Getenv("DEVIN_API_KEY")
	oldOrgID := os.Getenv("DEVIN_ORG_ID")
	os.Setenv("DEVIN_API_KEY", "test_api_key")
	os.Setenv("DEVIN_ORG_ID", "org-test")
	defer func() {
		os.Setenv("DEVIN_API_KEY", oldAPIKey)
		os.Setenv("DEVIN_ORG_ID", oldOrgID)
	}()

	client := NewClient("test_api_key", "org-test")
	if client == nil {
		t.Fatalf("NewClient() returned nil")
	}
	if client.APIKey != "test_api_key" {
		t.Errorf("NewClient() API key = %s, want %s", client.APIKey, "test_api_key")
	}
	if client.OrgID != "org-test" {
		t.Errorf("NewClient() OrgID = %s, want %s", client.OrgID, "org-test")
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

	if len(resources) != 4 {
		t.Fatalf("Resources() returned %d resources, want 4", len(resources))
	}
}

func TestProviderDataSources(t *testing.T) {
	ctx := context.Background()

	p := &DevinProvider{version: "test"}
	dataSources := p.DataSources(ctx)

	if len(dataSources) != 2 {
		t.Fatalf("DataSources() returned %d data sources, want 2", len(dataSources))
	}
}
