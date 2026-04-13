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

# ===================== Knowledge Resource =====================

resource "devin_knowledge" "example" {
  name    = "Sample Knowledge"
  body    = "This is the content of a sample knowledge created with Terraform"
  trigger = "This knowledge is triggered under specific conditions"
}

resource "devin_knowledge" "with_options" {
  name    = "Knowledge with Options"
  body    = "Knowledge body with optional fields"
  trigger = "Trigger with options"
}

# ===================== Playbook Resource =====================

resource "devin_playbook" "example" {
  title = "Sample Playbook"
  body  = "This playbook does automated code reviews"
}

resource "devin_playbook" "inactive" {
  title = "Inactive Playbook"
  body  = "This playbook is disabled"
}

# ===================== Secret Resource =====================

resource "devin_secret" "example" {
  name  = "TEST_DATABASE_URL"
  value = "postgres://localhost:5432/testdb"
}

resource "devin_secret" "api_token" {
  name  = "EXTERNAL_API_TOKEN"
  value = "token-abc-123"
}

# ===================== Schedule Resource =====================

resource "devin_schedule" "daily" {
  name   = "Daily Maintenance"
  prompt = "Run daily maintenance tasks"
  cron   = "0 9 * * *"
}

resource "devin_schedule" "weekly_with_playbook" {
  name        = "Weekly Code Review"
  prompt      = "Weekly code review"
  cron        = "0 10 * * 1"
  playbook_id = "playbook-mock-1"
}

# ===================== Data Sources =====================

data "devin_knowledge" "lookup" {
  id = "note-mock-1"
}

data "devin_folder" "by_id" {
  id = "folder-mock-1"
}

data "devin_folder" "by_name" {
  name = "モックフォルダ1"
}

# ===================== Outputs =====================

# Knowledge outputs
output "knowledge_id" {
  value = devin_knowledge.example.id
}

output "knowledge_name" {
  value = devin_knowledge.example.name
}

output "knowledge_with_options_enabled" {
  value = devin_knowledge.with_options.is_enabled
}

# Playbook outputs
output "playbook_id" {
  value = devin_playbook.example.id
}

output "playbook_title" {
  value = devin_playbook.example.title
}

# Secret outputs
output "secret_id" {
  value = devin_secret.example.id
}

output "secret_name" {
  value = devin_secret.example.name
}

# Schedule outputs
output "schedule_daily_id" {
  value = devin_schedule.daily.id
}

output "schedule_weekly_cron" {
  value = devin_schedule.weekly_with_playbook.cron
}

# Data source outputs
output "datasource_knowledge_name" {
  value = data.devin_knowledge.lookup.name
}

output "datasource_folder_name" {
  value = data.devin_folder.by_id.name
}

output "datasource_folder_by_name_id" {
  value = data.devin_folder.by_name.id
}
