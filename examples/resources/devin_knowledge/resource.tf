resource "devin_knowledge" "example" {
  name                = "Example Knowledge"
  body                = "This is an example knowledge resource."
  trigger_description = "Use this knowledge when talking about examples."
  parent_folder_id    = "optional-folder-id"
}
