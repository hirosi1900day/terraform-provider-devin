terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "1.0.0"
    }
  }
}

provider "devin" {
  # APIキーは環境変数DEVIN_API_KEYで設定するか、下記の行のコメントを解除して設定
  # api_key = "your-api-key-here"
}

# ナレッジリソースを作成
resource "devin_knowledge" "example" {
  name                = "サンプルナレッジ"
  body                = "これはTerraformで作成されたサンプルナレッジの内容です"
  trigger_description = "このナレッジは特定の条件でトリガーされます"
  # parent_folder_id    = "optional-folder-id" # 任意項目
}

# 作成したナレッジのIDを出力
output "knowledge_id" {
  value = devin_knowledge.example.id
}

# 作成したナレッジの名前を出力
output "knowledge_name" {
  value = devin_knowledge.example.name
}

# 作成したナレッジの内容を出力
output "knowledge_body" {
  value = devin_knowledge.example.body
}

# 作成したナレッジのトリガー説明を出力
output "knowledge_trigger_description" {
  value = devin_knowledge.example.trigger_description
}

# 作成したナレッジの作成日時を出力
output "knowledge_created_at" {
  value = devin_knowledge.example.created_at
}
