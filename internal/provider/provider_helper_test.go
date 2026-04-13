package provider

import (
	"testing"
)

func TestClientError(t *testing.T) {
	client := NewClient("invalid_key", "org-test")

	if client.APIKey == "test_api_key" {
		t.Skip("This test requires connection to the actual API")
	}
}

func TestKnowledgeResourceImports(t *testing.T) {
	resource := NewKnowledgeResource()
	if resource == nil {
		t.Fatal("NewKnowledgeResource() returned nil")
	}
}

func TestKnowledgeDataSourceImports(t *testing.T) {
	dataSource := NewKnowledgeDataSource()
	if dataSource == nil {
		t.Fatal("NewKnowledgeDataSource() returned nil")
	}
}

func TestPlaybookResourceImports(t *testing.T) {
	resource := NewPlaybookResource()
	if resource == nil {
		t.Fatal("NewPlaybookResource() returned nil")
	}
}

func TestSecretResourceImports(t *testing.T) {
	resource := NewSecretResource()
	if resource == nil {
		t.Fatal("NewSecretResource() returned nil")
	}
}

func TestScheduleResourceImports(t *testing.T) {
	resource := NewScheduleResource()
	if resource == nil {
		t.Fatal("NewScheduleResource() returned nil")
	}
}

func TestFolderDataSourceImports(t *testing.T) {
	dataSource := NewFolderDataSource()
	if dataSource == nil {
		t.Fatal("NewFolderDataSource() returned nil")
	}
}
