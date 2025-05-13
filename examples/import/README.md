# インポート機能の使用方法

このディレクトリには、既存のDevin APIナレッジリソースをTerraformの管理下に置くための例が含まれています。

## 準備

1. まず、このディレクトリで初期化を行います：

```bash
terraform init
```

## 既存リソースのインポート方法

1. main.tfファイルを確認し、インポートしたいリソースのテンプレートが定義されていることを確認します
2. 以下のコマンドを実行して、既存のナレッジリソースをインポートします：

```bash
terraform import devin_knowledge.imported mock-knowledge-1
```

`mock-knowledge-1`の部分は、インポートしたい実際のナレッジIDに置き換えてください。

3. インポートが成功すると、Terraformはリソースの状態を読み込み、状態ファイルに記録します
4. 実際のリソースと設定を同期させるため、インポート後に`terraform plan`を実行し、必要な差分を確認してください：

```bash
terraform plan
```

5. 問題がなければ変更を適用します：

```bash
terraform apply
```

## 注意事項

- インポートコマンドはリソースの状態のみをインポートします。設定ファイル（main.tf）は自動的に更新されません
- インポート後に`terraform plan`を実行して、必要な属性を設定ファイルに追加する必要がある場合があります
- テスト用のAPI Keyを使用する場合は、`test_api_key`を使用することでモックサーバーが利用されます

## その他の情報

詳細は[Terraform公式ドキュメント](https://www.terraform.io/cli/import)を参照してください。
