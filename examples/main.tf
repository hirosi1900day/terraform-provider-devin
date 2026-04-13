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
  body    = "This is a sample knowledge created with Terraform"
  trigger = "This knowledge is triggered under specific conditions"
  # pinned_repo = "owner/repo"
}

# Retrieve knowledge resource information
data "devin_knowledge" "example" {
  id = devin_knowledge.example.id
}

# Create a playbook
resource "devin_playbook" "example" {
  title = "Sample Playbook"
  body  = "Step-by-step instructions"
}

# Create a secret
resource "devin_secret" "example" {
  name  = "MY_SECRET"
  value = "secret-value"
}

# Create a schedule
resource "devin_schedule" "example" {
  name        = "Weekly Maintenance"
  prompt      = "Run weekly maintenance"
  cron        = "0 9 * * 1"
  playbook_id = devin_playbook.example.id
}

# Outputs
output "knowledge_id" {
  value = devin_knowledge.example.id
}

output "knowledge_name" {
  value = data.devin_knowledge.example.name
}

output "knowledge_trigger" {
  value = data.devin_knowledge.example.trigger
}

output "playbook_id" {
  value = devin_playbook.example.id
}

output "schedule_id" {
  value = devin_schedule.example.id
}
