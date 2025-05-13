terraform {
  required_providers {
    devin = {
      source = "hirosi1900day/devin"
    }
  }
}

provider "devin" {
  # 環境変数 DEVIN_API_KEY で設定することも可能
  api_key = var.devin_api_key
}

variable "devin_api_key" {
  description = "Devin API Key"
  sensitive   = true
}

# ナレッジリソースを作成
resource "devin_knowledge" "example" {
  name        = "サンプルナレッジ"
  description = "これはTerraformで作成されたサンプルナレッジです"
}

# ナレッジリソースの情報を取得
data "devin_knowledge" "example" {
  id = devin_knowledge.example.id
}

# 出力
output "knowledge_id" {
  value = devin_knowledge.example.id
}

output "knowledge_name" {
  value = data.devin_knowledge.example.name
}

output "knowledge_description" {
  value = data.devin_knowledge.example.description
}

output "knowledge_created_at" {
  value = devin_knowledge.example.created_at
}
