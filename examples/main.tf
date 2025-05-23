terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "0.0.5"
    }
  }
}

provider "devin" {
  # Set test API key
  api_key = "test_api_key"
}

# Create a knowledge resource
resource "devin_knowledge" "example" {
  name                = "Sample Knowledge"
  body                = "This is a sample knowledge created with Terraform"
  trigger_description = "This knowledge is triggered under specific conditions"
  # parent_folder_id    = "optional-folder-id" # Optional parameter
}

# Retrieve knowledge resource information
data "devin_knowledge" "example" {
  id = devin_knowledge.example.id
}

# Output
output "knowledge_id" {
  value = devin_knowledge.example.id
}

output "knowledge_name" {
  value = data.devin_knowledge.example.name
}

output "knowledge_body" {
  value = data.devin_knowledge.example.body
}

output "knowledge_trigger_description" {
  value = data.devin_knowledge.example.trigger_description
}
