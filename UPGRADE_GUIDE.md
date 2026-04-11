# v1.0.0 アップグレードガイド (v0.0.7 → v1.0.0)

このガイドは `hirosi1900day/devin` Terraform Provider を v0.0.7 から v1.0.0 へアップグレードする際の手順と破壊的変更をまとめたものです。

v1.0.0 では Devin API v1 から v3 への移行に伴い、認証方式・リソーススキーマ・内部実装が大幅に変更されています。

---

## 目次

1. [破壊的変更の一覧](#1-破壊的変更の一覧)
2. [Provider 設定の変更](#2-provider-設定の変更)
3. [devin_knowledge リソースの変更](#3-devin_knowledge-リソースの変更)
4. [devin_knowledge データソースの変更](#4-devin_knowledge-データソースの変更)
5. [devin_folder データソースの変更](#5-devin_folder-データソースの変更)
6. [新規リソース](#6-新規リソース)
7. [State 移行手順](#7-state-移行手順)
8. [自動生成スクリプトの更新](#8-自動生成スクリプトの更新)

---

## 1. 破壊的変更の一覧

| 対象 | 変更内容 |
|------|----------|
| Provider | `org_id` が**必須**に。API Key は `cog_*` 形式に変更 |
| `devin_knowledge` | `trigger_description` → `trigger` に**リネーム** |
| `devin_knowledge` | `parent_folder_id` → `folder_id` に**リネーム** |
| `devin_knowledge` | `id` の値が `note-xxxx` 形式に変更 |
| `devin_knowledge` | `created_at` の型が ISO 8601 文字列 → UNIX timestamp (float64) に変更 |
| `data.devin_knowledge` | 属性名のリネーム（`trigger_description` → `trigger` 等） |
| `data.devin_folder` | `description` 属性が**廃止**、`path` / `note_count` / `parent_folder_id` が追加 |
| `data.devin_folder` | `id` の値が `folder-xxxx` 形式に変更 |

---

## 2. Provider 設定の変更

### Before (v0.0.7)

```hcl
provider "devin" {
  api_key = "apk_..."  # Personal API Key
}
```

### After (v1.0.0)

```hcl
provider "devin" {
  api_key = "cog_..."  # Service User credential（必須）
  org_id  = "org_..."  # Organization ID（必須）
}
```

### 環境変数

| 環境変数 | 用途 |
|----------|------|
| `DEVIN_API_KEY` | API Key（変更なし、値を `cog_*` に更新が必要） |
| `DEVIN_ORG_ID` | Organization ID（**新規**） |

### 取得方法

1. Devin の Settings → Service Users で Service User を作成
2. 発行された `cog_*` トークンを `api_key` に設定
3. Settings → Organization から `org_*` で始まる Organization ID を取得し `org_id` に設定

---

## 3. devin_knowledge リソースの変更

### 属性のリネーム

```hcl
# Before (v0.0.7)
resource "devin_knowledge" "example" {
  name                = "rule name"
  body                = "content"
  trigger_description = "when working in repo ..."
  parent_folder_id    = data.devin_folder.xxx.id
}

# After (v1.0.0)
resource "devin_knowledge" "example" {
  name       = "rule name"
  body       = "content"
  trigger    = "when working in repo ..."    # リネーム
  folder_id  = data.devin_folder.xxx.id      # リネーム
  is_enabled = true                           # 新規（Optional, default: true）
  # pinned_repo = "owner/repo"               # 新規（Optional）
}
```

### 変更された属性の対応表

| v0.0.7 | v1.0.0 | 変更 |
|--------|--------|------|
| `id` | `id` | 値が `note-xxxx` 形式に |
| `name` | `name` | 変更なし |
| `body` | `body` | 変更なし |
| `trigger_description` | `trigger` | **リネーム** |
| `parent_folder_id` | `folder_id` | **リネーム** |
| — | `is_enabled` | **新規** (bool, default: true) |
| — | `pinned_repo` | **新規** (string, optional) |
| — | `folder_path` | **新規** (computed) |
| — | `macro` | **新規** (computed) |
| — | `access_type` | **新規** (computed) |
| — | `created_at` | **型変更** ISO 8601 → float64 (UNIX timestamp) |
| — | `updated_at` | **新規** (computed) |

---

## 4. devin_knowledge データソースの変更

```hcl
# Before (v0.0.7)
data "devin_knowledge" "example" {
  id = "knowledge-id"
}
# 参照: data.devin_knowledge.example.trigger_description

# After (v1.0.0)
data "devin_knowledge" "example" {
  id = "note-xxxx"  # note_id を指定
}
# 参照: data.devin_knowledge.example.trigger
```

属性のリネームはリソースと同様です。

---

## 5. devin_folder データソースの変更

```hcl
# Before (v0.0.7)
data "devin_folder" "example" {
  name = "My Folder"
}
# 参照: data.devin_folder.example.description  # ← v1.0.0 で廃止

# After (v1.0.0)
data "devin_folder" "example" {
  name = "My Folder"
}
# 参照:
#   data.devin_folder.example.id               # folder_id
#   data.devin_folder.example.path             # 新規: 階層パス
#   data.devin_folder.example.note_count       # 新規: ノート数
#   data.devin_folder.example.parent_folder_id # 新規: 親フォルダID
```

| v0.0.7 | v1.0.0 | 変更 |
|--------|--------|------|
| `id` | `id` | 値が `folder-xxxx` 形式に |
| `name` | `name` | 変更なし |
| `description` | — | **廃止** |
| — | `path` | **新規** |
| — | `note_count` | **新規** |
| — | `parent_folder_id` | **新規** |

---

## 6. 新規リソース

v1.0.0 で以下のリソースが追加されました。

### devin_playbook

セッション開始時にユーザーが選択する手順書を管理します。

```hcl
resource "devin_playbook" "example" {
  title  = "データ分析フロー"    # Required
  body   = "手順の内容..."        # Required
  status = "active"               # Optional, default: "active"
}
```

### devin_secret

Devin が利用する環境変数（シークレット）を管理します。

```hcl
resource "devin_secret" "example" {
  name  = "DATABASE_URL"          # Required, ForceNew
  value = var.database_url        # Required, Sensitive, ForceNew
}
```

> **注意**: Secret には Update API がないため、`name` または `value` を変更するとリソースが再作成されます。`value` は API から取得できない Write-Only 属性です。

### devin_schedule

定期実行のスケジュールを管理します。

```hcl
resource "devin_schedule" "example" {
  prompt      = "リポジトリの依存更新を確認して PR を作成"  # Required
  cron        = "0 9 * * 1"                                # Required
  playbook_id = devin_playbook.example.id                  # Optional
}
```

> **注意**: Schedule の更新は PATCH（部分更新）で行われます。

---

## 7. State 移行手順

v1.0.0 はスキーマが大きく変わっているため、**既存 State からのリソース除去 → 再インポート** が推奨です。

### 手順

```bash
# 1. Provider バージョンを更新
#    required_providers の version を "~> 1.0.0" に変更

# 2. .tf ファイルの属性名を更新
#    trigger_description → trigger
#    parent_folder_id    → folder_id

# 3. 既存 State からリソースを除去（実リソースは削除しない）
terraform state rm devin_knowledge.example

# 4. v3 の note_id で再インポート
terraform import devin_knowledge.example note-xxxx
```

### import ブロックを使う場合（Terraform 1.5.0+）

リソースが多い場合は import ブロックでまとめて処理できます：

```hcl
import {
  to = devin_knowledge.rule_1
  id = "note-aaaa"
}

import {
  to = devin_knowledge.rule_2
  id = "note-bbbb"
}
```

```bash
terraform plan -generate-config-out=generated.tf  # 設定ファイルの自動生成
terraform apply
```

### 一括移行スクリプト例

```bash
#!/bin/bash
# 全 devin_knowledge リソースを state から除去
terraform state list | grep '^devin_knowledge\.' | while read resource; do
  terraform state rm "$resource"
done

# import ブロックを含む .tf を用意した上で
terraform apply
```

---

## 8. 自動生成スクリプトの更新

Knowledge リソースを自動生成するスクリプトがある場合、以下の変更が必要です：

| 変更箇所 | Before | After |
|----------|--------|-------|
| 属性名 | `trigger_description` | `trigger` |
| 属性名 | `parent_folder_id` | `folder_id` |
| Provider 設定 | `api_key` のみ | `api_key` + `org_id` |
| ID 形式 | 旧形式 | `note-xxxx` 形式 |

---

## FAQ

### Q: v0.0.7 と v1.0.0 を同時に使えますか？

いいえ。API v1 と v3 は認証方式が異なるため、同一 Provider インスタンスでの共存はできません。

### Q: API Key はどう変わりますか？

`apk_*` (Personal/Service API Key) から `cog_*` (Service User credential) に変更されます。Devin の Settings → Service Users で新しいトークンを発行してください。

### Q: 既存の Knowledge は消えますか？

いいえ。Devin 上のリソースには影響ありません。Terraform State のみを再構築します。`terraform state rm` は Terraform の管理からリソースを外すだけで、実リソースは削除しません。

### Q: folder の作成・削除はできますか？

v1.0.0 でもフォルダの作成・更新・削除 API は提供されていません。`devin_folder` は引き続きデータソース（読み取り専用）です。フォルダは Devin Web UI で事前に作成してください。
