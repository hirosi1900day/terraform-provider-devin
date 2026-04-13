resource "devin_knowledge" "example" {
  name       = "Example Knowledge"
  body       = "This is an example knowledge resource."
  trigger    = "Use this knowledge when talking about examples."
  folder_id  = "optional-folder-id"  # Optional
  is_enabled = true                   # Optional, default: true
  # pinned_repo = "owner/repo"       # Optional
}
