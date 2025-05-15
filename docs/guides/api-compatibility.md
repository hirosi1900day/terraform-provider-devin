---
page_title: "Devin API Compatibility - terraform-provider-devin"
subcategory: ""
description: |-
  Explains the compatibility between the Devin Terraform provider and its underlying API.
---

# Devin API Compatibility Guide

This document describes the compatibility between the Devin Terraform provider and the Devin API.

## API Structure and Compatibility

This provider works directly with the Devin API specification:

### Field Structure

- The Terraform provider uses the same field names as the API: `body` for content and `trigger_description` for trigger conditions.
- The `parent_folder_id` field is optional and can be used to organize knowledge resources in folders.

### API Endpoints

- The provider supports the `/knowledge` endpoint for listing all knowledge resources.
- The response includes both knowledge items and folders as described in the [Devin API documentation](https://docs.devin.ai/api-reference/knowledge/list-knowledge).

### Field Requirements

- Both `body` and `trigger_description` are required fields when creating knowledge resources.
- The API requires these fields, and the provider enforces this requirement.

## Authentication

Authentication with the Devin API is done using an API key. The API key can be specified directly in the Terraform configuration or provided through the `DEVIN_API_KEY` environment variable.

```terraform
provider "devin" {
  api_key = "your_api_key"
}
```

Or

```bash
export DEVIN_API_KEY="your_api_key"
terraform apply
```

## Error Handling

The provider properly captures and displays error messages returned from the API. Errors from the API include a detailed description along with a message type.

Common errors:

- **Authentication errors**: When the API key is invalid or expired
- **Resource not found errors**: When referencing a resource ID that doesn't exist
- **Validation errors**: When required fields are missing or in an invalid format

## Versioning and API Changes

As the Devin API evolves, this provider will be updated accordingly. The provider versioning follows Terraform's semantic versioning convention:

- **Patch version increases**: Bug fixes and documentation updates
- **Minor version increases**: Backward-compatible feature additions
- **Major version increases**: Breaking changes

When there are breaking changes in the API, those changes will be implemented through a major version update.
