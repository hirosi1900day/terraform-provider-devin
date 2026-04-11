data "devin_knowledge" "example" {
  id = "note-xxxx"  # note_id
}

output "knowledge_name" {
  value = data.devin_knowledge.example.name
}

output "knowledge_body" {
  value = data.devin_knowledge.example.body
}

output "knowledge_trigger" {
  value = data.devin_knowledge.example.trigger
}
