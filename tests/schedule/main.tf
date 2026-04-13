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

# ===================== Schedule Resources =====================

resource "devin_schedule" "daily" {
  name   = "Daily Quality Check"
  prompt = "Run daily code quality checks"
  cron   = "0 9 * * *"
}

resource "devin_schedule" "weekly" {
  name   = "Weekly Dep Review"
  prompt = "Weekly dependency update review"
  cron   = "0 10 * * 1"
}

resource "devin_schedule" "with_playbook" {
  name        = "Playbook Schedule"
  prompt      = "Execute playbook on schedule"
  cron        = "30 14 * * 1-5"
  playbook_id = "playbook-mock-1"
}

# ===================== Outputs =====================

output "daily_id" {
  value = devin_schedule.daily.id
}

output "daily_cron" {
  value = devin_schedule.daily.cron
}

output "weekly_id" {
  value = devin_schedule.weekly.id
}

output "with_playbook_id" {
  value = devin_schedule.with_playbook.id
}

output "with_playbook_playbook_id" {
  value = devin_schedule.with_playbook.playbook_id
}
