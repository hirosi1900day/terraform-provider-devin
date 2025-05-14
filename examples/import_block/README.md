# importブロックの使用例

このディレクトリには、Terraform 1.5以降で導入された「importブロック」を使って既存のDevin APIナレッジリソースをインポートする例が含まれています。

## 準備

1. まず、このディレクトリで初期化を行います：

```bash
terraform init
```

## importブロックの仕組み

従来のコマンドラインからのインポート（`terraform import`）と異なり、importブロックは設定ファイル内に宣言的に記述します。これにより：

- インポート操作がコードとして管理できる
- CIパイプラインでの自動化が容易になる
- インポート処理が再現可能になる

## 使用方法

1. `main.tf`ファイルに、以下のようなimportブロックが含まれています：

```hcl
import {
  to = devin_knowledge.imported_block
  id = "mock-knowledge-1"
}

resource "devin_knowledge" "imported_block" {
  name                = "モックナレッジ1"
  body                = "これはテスト用のモックナレッジです"
  trigger_description = "テスト用トリガーの説明"
}
```

2. 以下のコマンドを実行して、importブロックを適用します：

```bash
terraform plan
terraform apply
```

3. インポート完了後は、通常のリソースとして管理できます。

## 再実行方法

インポートをやり直したい場合は、以下の手順で操作します：

1. 状態ファイルからリソースを削除します：

```bash
terraform state rm devin_knowledge.imported_block
```

2. 再度 `terraform plan` と `terraform apply` を実行します。

## 複数リソースのインポート

`main.tf`ファイルにはコメントアウトされた状態で複数リソースのインポート例も含まれています。複数のリソースをインポートするには、コメントを解除して必要に応じて修正してください。

## 注意点

- importブロックは一度適用されると、その役割を終えます。適用後のファイルからimportブロックを削除しても問題ありません。
- 適切なリソース属性を事前に指定しておくことで、インポート後の変更を最小限に抑えられます。
- テスト用のAPI Keyを使用する場合は、`test_api_key`を使用することでモックサーバーが利用されます。
- 実際の環境では、APIキーを環境変数（`DEVIN_API_KEY`）として設定するか、`api_key`パラメータに直接指定します。

## その他の情報

詳細は[Terraform公式ドキュメント](https://developer.hashicorp.com/terraform/language/import)を参照してください。 