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
  name        = "更新されたナレッジ"
  description = "これはTerraformで更新されたサンプルナレッジです"
}

# 作成したナレッジのIDを出力
output "knowledge_id" {
  value = devin_knowledge.example.id
}

# 作成したナレッジの名前を出力
output "knowledge_name" {
  value = devin_knowledge.example.name
}

# 作成したナレッジの説明を出力
output "knowledge_description" {
  value = devin_knowledge.example.description
}

# 作成したナレッジの作成日時を出力
output "knowledge_created_at" {
  value = devin_knowledge.example.created_at
}
