# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
