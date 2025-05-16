data "devin_knowledge" "example" {
  id = "knowledge-resource-id"
}

output "knowledge_name" {
  value = data.devin_knowledge.example.name
}

output "knowledge_body" {
  value = data.devin_knowledge.example.body
}
