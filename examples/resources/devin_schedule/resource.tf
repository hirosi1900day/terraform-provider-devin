resource "devin_schedule" "example" {
  name        = "Weekly Dependency Check"
  prompt      = "Check repository dependencies and create PRs"
  cron        = "0 9 * * 1"  # Every Monday at 9:00 AM
  playbook_id = devin_playbook.example.id  # Optional
}
