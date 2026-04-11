terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "1.0.0"
    }
  }
}

provider "devin" {
  api_key = "test_api_key"
  org_id  = "org-test"
}

# Use an import block to import existing knowledge resource
import {
  to = devin_knowledge.imported_block
  id = "note-mock-1"
}

# Definition of the imported resource
resource "devin_knowledge" "imported_block" {
  name    = "Mock Knowledge 1"
  body    = "This is a mock knowledge for testing"
  trigger = "Test trigger description"
}

# Definition of a regular knowledge resource
resource "devin_knowledge" "example" {
  name    = "Sample Knowledge"
  body    = "This is a sample knowledge created with Terraform"
  trigger = "This knowledge is triggered under specific conditions"
}

# Output the imported resource ID
output "imported_id" {
  value = devin_knowledge.imported_block.id
}

# Output the imported resource name
output "imported_name" {
  value = devin_knowledge.imported_block.name
} 