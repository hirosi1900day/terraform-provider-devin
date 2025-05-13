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

# 既存のナレッジリソースをインポートする例
resource "devin_knowledge" "imported" {
  name        = "モックナレッジ1"
  description = "これはテスト用のモックナレッジです"
  # インポート済み: terraform import devin_knowledge.imported mock-knowledge-1
}

# 通常のナレッジリソースの定義
resource "devin_knowledge" "example" {
  name        = "サンプルナレッジ"
  description = "これはTerraformで作成されたサンプルナレッジです"
}
