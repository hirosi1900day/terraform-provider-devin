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

# Create a knowledge resource
resource "devin_knowledge" "example" {
  name    = "Sample Knowledge"
  body    = "This is the content of a sample knowledge created with Terraform"
  trigger = "This knowledge is triggered under specific conditions"
  # folder_id = "optional-folder-id"
}

# Output the created knowledge ID
output "knowledge_id" {
  value = devin_knowledge.example.id
}

# Output the created knowledge name
output "knowledge_name" {
  value = devin_knowledge.example.name
}

# Output the content of the created knowledge
output "knowledge_body" {
  value = devin_knowledge.example.body
}

# Output the trigger of the created knowledge
output "knowledge_trigger" {
  value = devin_knowledge.example.trigger
}
