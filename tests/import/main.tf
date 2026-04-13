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

# ===================== Knowledge Import =====================

import {
  to = devin_knowledge.imported_block
  id = "note-mock-1"
}

resource "devin_knowledge" "imported_block" {
  name    = "Mock Knowledge 1"
  body    = "This is a mock knowledge for testing"
  trigger = "Test trigger description"
}

# ===================== Playbook Import =====================

import {
  to = devin_playbook.imported_block
  id = "playbook-mock-1"
}

resource "devin_playbook" "imported_block" {
  title = "Mock Playbook"
  body  = "This is a mock playbook for testing"
}

# ===================== Schedule Import =====================

import {
  to = devin_schedule.imported_block
  id = "schedule-mock-1"
}

resource "devin_schedule" "imported_block" {
  name        = "Imported Schedule"
  prompt      = "Scheduled task prompt"
  cron        = "0 9 * * 1"
  playbook_id = "playbook-mock-1"
}

# ===================== Regular Resources =====================

resource "devin_knowledge" "example" {
  name    = "Sample Knowledge"
  body    = "This is a sample knowledge created with Terraform"
  trigger = "This knowledge is triggered under specific conditions"
}

resource "devin_playbook" "example" {
  title = "Sample Playbook"
  body  = "Playbook content for testing"
}

resource "devin_secret" "example" {
  name  = "IMPORT_TEST_SECRET"
  value = "secret-value"
}

resource "devin_schedule" "example" {
  name   = "Regular Schedule"
  prompt = "Regular schedule task"
  cron   = "0 12 * * *"
}

# ===================== Data Sources =====================

data "devin_knowledge" "lookup" {
  id = "note-mock-2"
}

data "devin_folder" "lookup" {
  id = "folder-mock-2"
}

# ===================== Outputs =====================

output "imported_knowledge_id" {
  value = devin_knowledge.imported_block.id
}

output "imported_playbook_id" {
  value = devin_playbook.imported_block.id
}

output "imported_schedule_id" {
  value = devin_schedule.imported_block.id
}

output "knowledge_id" {
  value = devin_knowledge.example.id
}

output "playbook_id" {
  value = devin_playbook.example.id
}

output "secret_id" {
  value = devin_secret.example.id
}

output "schedule_id" {
  value = devin_schedule.example.id
}

output "datasource_knowledge_name" {
  value = data.devin_knowledge.lookup.name
}

output "datasource_folder_name" {
  value = data.devin_folder.lookup.name
} 