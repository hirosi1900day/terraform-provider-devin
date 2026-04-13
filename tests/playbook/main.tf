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

# ===================== Playbook Resources =====================

resource "devin_playbook" "basic" {
  title = "Basic Playbook"
  body  = "A simple playbook for code review automation"
}

resource "devin_playbook" "inactive" {
  title = "Inactive Playbook"
  body  = "This playbook is paused"
}

resource "devin_playbook" "default_status" {
  title = "Default Status Playbook"
  body  = "Status defaults to active"
}

# ===================== Outputs =====================

output "basic_id" {
  value = devin_playbook.basic.id
}

output "basic_title" {
  value = devin_playbook.basic.title
}
