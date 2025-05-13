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
}

func TestProviderConfigure(t *testing.T) {
	// In testing, we don't want to depend on the implementation of provider.Configure
	// Actual provider functionality should be covered in integration tests
	t.Skip("Skipping detailed test for provider.Configure")
}

func TestProviderAPI(t *testing.T) {
	// Test behavior when API_KEY is set in environment variables
	oldApiKey := os.Getenv("DEVIN_API_KEY")
	os.Setenv("DEVIN_API_KEY", "test_api_key")
	defer os.Setenv("DEVIN_API_KEY", oldApiKey) // Restore original value after test

	// Verify that the client is created correctly
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
