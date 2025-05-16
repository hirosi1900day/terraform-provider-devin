terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "0.0.3"
    }
  }
}

provider "devin" {
  api_key = "test_api_key" # In production, set this as an environment variable or Terraform Cloud variable
}

# Use an import block to import existing knowledge resource
# This block imports the resource when terraform apply is executed, and then its job is done
# Once the import is complete, you can comment out or remove this block
import {
  to = devin_knowledge.imported_block
  id = "mock-knowledge-1"
}

# Definition of the imported resource
resource "devin_knowledge" "imported_block" {
  name                = "Mock Knowledge 1"
  body                = "This is a mock knowledge for testing"
  trigger_description = "Test trigger description"
}

# Definition of a regular knowledge resource
resource "devin_knowledge" "example" {
  name                = "Sample Knowledge"
  body                = "This is a sample knowledge created with Terraform"
  trigger_description = "This knowledge is triggered under specific conditions"
}

# Output the imported resource ID
output "imported_id" {
  value = devin_knowledge.imported_block.id
}

# Output the imported resource name
output "imported_name" {
  value = devin_knowledge.imported_block.name
} 