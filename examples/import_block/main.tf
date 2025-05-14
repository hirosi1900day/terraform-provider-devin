terraform {
  required_providers {
    devin = {
      source = "hirosi1900day/devin"
    }
  }
}

provider "devin" {
  api_key = "test_api_key" # 本番環境では環境変数またはTerraform Cloudの変数として設定
}

# importブロックを使用して既存のナレッジリソースをインポートする
# このブロックはterraform applyを実行すると、リソースをインポートし、その役割を終えます
# 一度インポートが完了したら、このブロックはコメントアウトするか削除しても問題ありません
import {
  to = devin_knowledge.imported_block
  id = "mock-knowledge-1"
}

# インポートされたリソースの定義
resource "devin_knowledge" "imported_block" {
  name                = "モックナレッジ1"
  body                = "これはテスト用のモックナレッジです"
  trigger_description = "テスト用トリガーの説明"
}

# 通常のナレッジリソースの定義
resource "devin_knowledge" "example" {
  name                = "サンプルナレッジ"
  body                = "これはTerraformで作成されたサンプルナレッジです"
  trigger_description = "このナレッジは特定の条件でトリガーされます"
}

# インポートしたリソースのIDを出力
output "imported_id" {
  value = devin_knowledge.imported_block.id
}

# インポートしたリソースの名前を出力
output "imported_name" {
  value = devin_knowledge.imported_block.name
}

# インポートしたリソースの作成日時を出力
output "imported_created_at" {
  value = devin_knowledge.imported_block.created_at
} 