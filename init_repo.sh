#!/bin/bash

# リポジトリの初期化
echo "Devin Terraform Providerリポジトリを初期化します..."

# 必要なモジュールのダウンロード
go mod tidy

# ビルドとテスト
go build -o terraform-provider-devin
go test -v ./...

# Gitリポジトリの初期化
git init
git add .
git config --local user.name "$(git config --global user.name)"
git config --local user.email "$(git config --global user.email)"
git commit -m "初期コミット: Devin Terraform Provider"

echo ""
echo "以下のコマンドでGitHubリポジトリを作成し、プッシュしてください："
echo ""
echo "  git remote add origin https://github.com/hirosi1900day/terraform-provider-devin-knowledge.git"
echo "  git push -u origin main"
echo ""
echo "初期化が完了しました！"
