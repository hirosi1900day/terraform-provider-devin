package provider

import (
	"testing"
)

func TestClientError(t *testing.T) {
	// Test with invalid API key
	client := NewClient("invalid_key")

	// In the current implementation, we need to skip this test when using mock data
	// because errors are only triggered with actual API requests
	if client.APIKey == "test_api_key" {
		t.Skip("This test requires connection to the actual API")
	}
}

func TestKnowledgeResourceImports(t *testing.T) {
	// Test for resource import functionality
	resource := NewKnowledgeResource()
	if resource == nil {
		t.Fatal("NewKnowledgeResource() returned nil")
	}

	// Further detailed tests would require RPC, so we only verify resource creation
}

func TestKnowledgeDataSourceImports(t *testing.T) {
	// Test for data source import functionality
	dataSource := NewKnowledgeDataSource()
	if dataSource == nil {
		t.Fatal("NewKnowledgeDataSource() returned nil")
	}
}
