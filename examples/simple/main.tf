terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "0.0.3"
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
  body                = "This is the content of a sample knowledge created with Terraform"
  trigger_description = "This knowledge is triggered under specific conditions"
  # parent_folder_id    = "optional-folder-id" # Optional parameter
}

# Output the created knowledge ID
output "knowledge_id" {
  value = devin_knowledge.example.id
}

# Output the created knowledge name
output "knowledge_name" {
  value = devin_knowledge.example.name
}

# 作成したナレッジの内容を出力
output "knowledge_body" {
  value = devin_knowledge.example.body
}

# 作成したナレッジのトリガー説明を出力
output "knowledge_trigger_description" {
  value = devin_knowledge.example.trigger_description
}
