terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "0.0.2"
    }
  }
}

provider "devin" {
  # テスト用APIキーを設定
  api_key = "test_api_key"
}

# ナレッジリソースを作成
resource "devin_knowledge" "example" {
  name                = "サンプルナレッジ"
  body                = "これはTerraformで作成されたサンプルナレッジです"
  trigger_description = "このナレッジは特定の条件でトリガーされます"
  # parent_folder_id    = "optional-folder-id" # 任意項目
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

output "knowledge_body" {
  value = data.devin_knowledge.example.body
}

output "knowledge_trigger_description" {
  value = data.devin_knowledge.example.trigger_description
}
