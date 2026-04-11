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

# ===================== Data Sources =====================

data "devin_knowledge" "note1" {
  id = "note-mock-1"
}

data "devin_knowledge" "note2" {
  id = "note-mock-2"
}

data "devin_folder" "folder1" {
  id = "folder-mock-1"
}

data "devin_folder" "folder2" {
  id = "folder-mock-2"
}

data "devin_folder" "by_name" {
  name = "モックフォルダ1"
}

# ===================== Outputs =====================

output "knowledge_note1_name" {
  value = data.devin_knowledge.note1.name
}

output "knowledge_note1_body" {
  value = data.devin_knowledge.note1.body
}

output "knowledge_note1_trigger" {
  value = data.devin_knowledge.note1.trigger
}

output "knowledge_note1_enabled" {
  value = data.devin_knowledge.note1.is_enabled
}

output "knowledge_note2_name" {
  value = data.devin_knowledge.note2.name
}

output "folder1_name" {
  value = data.devin_folder.folder1.name
}

output "folder1_path" {
  value = data.devin_folder.folder1.path
}

output "folder1_note_count" {
  value = data.devin_folder.folder1.note_count
}

output "folder2_name" {
  value = data.devin_folder.folder2.name
}

output "folder_by_name_id" {
  value = data.devin_folder.by_name.id
}
