# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-04-11

### Breaking Changes
- **Provider**: `org_id` パラメータが必須に。環境変数 `DEVIN_ORG_ID` でも設定可能
- **Provider**: API Key が `cog_*` (Service User credential) 形式に変更
- **Provider**: Base URL が `https://api.devin.ai/v3/organizations/{org_id}` に変更
- **devin_knowledge**: `trigger_description` → `trigger` にリネーム
- **devin_knowledge**: `parent_folder_id` → `folder_id` にリネーム
- **devin_knowledge**: `id` の値が `note-xxxx` 形式に変更
- **devin_knowledge**: `created_at` の型が ISO 8601 文字列 → UNIX timestamp (float64) に変更
- **data.devin_folder**: `description` 属性が廃止
- **data.devin_folder**: `id` の値が `folder-xxxx` 形式に変更

### Added
- Devin API v3 対応
- `devin_knowledge`: `is_enabled`, `pinned_repo`, `folder_path`, `macro`, `access_type`, `updated_at` 属性を追加
- `data.devin_folder`: `path`, `note_count`, `parent_folder_id` 属性を追加
- 新リソース `devin_playbook`: Playbook の CRUD 管理
- 新リソース `devin_secret`: Secret の作成・削除管理（Update 非対応、ForceNew）
- 新リソース `devin_schedule`: Schedule の CRUD 管理（更新は PATCH）
- Knowledge の個別取得 API 対応（`GET /knowledge/notes/{note_id}`）
- カーソルベースページネーション対応

### Removed
- API v1 のサポート
- Knowledge 一覧のキャッシュ機構（v3 では個別取得 API により不要）

### Migration
- [UPGRADE_GUIDE.md](UPGRADE_GUIDE.md) を参照してください

## [0.0.7] - 2025-11-30

### Added
- Added caching mechanism for ListKnowledge API responses to avoid rate limiting
- Cache TTL defaults to 5 minutes, suitable for terraform plan/apply duration
- Thread-safe cache implementation using sync.RWMutex

### Changed
- GetKnowledge, GetFolderByID, GetFolderByName now use cached ListKnowledge response
- Cache is automatically invalidated after Create, Update, or Delete operations

### Fixed
- Fixed rate limiting issues during terraform plan with many knowledge resources
- Fixed CI workflow to use Go version from go.mod instead of hardcoded version

## [0.0.6] - 2025-05-23

### Changed
- Updated version to 0.0.6

## [0.0.5] - 2025-05-23

### Changed
- Updated example configurations


## [0.0.4] - 2025-05-16

### Added
- Enhanced documentation for knowledge resources
- Additional examples for import operations
- Support for import blocks (Terraform 1.5.0+)

### Changed
- Improved error handling for API responses
- Optimized resource state management
- Updated example configurations

## [0.0.3] - 2025-05-16

### Added
- Folder data source support for improved organization
- Enhanced error handling and validation
- Performance optimizations for API requests

### Changed
- Refactored provider configuration for better usability
- Improved documentation with more detailed examples

## [0.0.2] - 2025-05-16

### Added
- Comprehensive documentation explaining how to use the Devin Terraform provider
- Troubleshooting guides and API compatibility information
- Examples for Devin knowledge data source and resource usage

### Fixed
- Updated GetKnowledge method to correctly match note ID format
- Improved ParentFolderID update logic
- Various code optimizations and stability improvements

## [0.0.1] - 2025-05-14

### Added
- Initial release
- Create, read, update, and delete functionality for knowledge resources
- Implementation of knowledge data source
- Implementation using Terraform Framework API
- Development and testing support with mock data functionality

[0.0.7]: https://github.com/hirosi1900day/terraform-provider-devin-knowledge/releases/tag/v0.0.7
[0.0.6]: https://github.com/hirosi1900day/terraform-provider-devin-knowledge/releases/tag/v0.0.6
[0.0.5]: https://github.com/hirosi1900day/terraform-provider-devin-knowledge/releases/tag/v0.0.5
[0.0.4]: https://github.com/hirosi1900day/terraform-provider-devin-knowledge/releases/tag/v0.0.4
[0.0.3]: https://github.com/hirosi1900day/terraform-provider-devin-knowledge/releases/tag/v0.0.3
[0.0.2]: https://github.com/hirosi1900day/terraform-provider-devin-knowledge/releases/tag/v0.0.2
[0.0.1]: https://github.com/hirosi1900day/terraform-provider-devin-knowledge/releases/tag/v0.0.1
