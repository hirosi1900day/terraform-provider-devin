# Devin Terraform Provider

This repository contains a Terraform Provider for managing [Devin AI](https://devin.ai) knowledge resources.

[![Go Report Card](https://goreportcard.com/badge/github.com/hirosi1900day/terraform-provider-devin-knowledge)](https://goreportcard.com/report/github.com/hirosi1900day/terraform-provider-devin-knowledge)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

The Devin Terraform Provider allows you to manage Devin AI knowledge resources as infrastructure as code. This enables you to maintain knowledge creation, updating, and deletion operations as version-controlled configuration files.

## Features

- Create knowledge resources
- Update knowledge resources
- Delete knowledge resources
- Reference existing knowledge resources
- Import existing resources into Terraform state

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.0+
- [Go](https://golang.org/doc/install) 1.21+ (for development only)
- Access token for the [Devin API](https://docs.devin.ai/api-reference)

## Installation

```hcl
terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "~> 1.0"
    }
  }
}
```

## Usage

### Provider Configuration

```hcl
terraform {
  required_providers {
    devin = {
      source = "hirosi1900day/devin"
    }
  }
}

provider "devin" {
  api_key = "your_api_key" # Or use the DEVIN_API_KEY environment variable
}
```

### Creating a Knowledge Resource

```hcl
resource "devin_knowledge" "example" {
  name        = "Sample Knowledge"
  description = "This is a knowledge resource created with Terraform"
}
```

> **Note**: The `description` field is internally mapped to the `body` field in the Devin API.
> Future versions of this provider will expose additional API fields such as `trigger_description` 
> and `parent_folder_id`.

### Using a Knowledge Data Source

```hcl
data "devin_knowledge" "example" {
  id = "knowledge_id"
}

output "knowledge_name" {
  value = data.devin_knowledge.example.name
}

output "knowledge_description" {
  value = data.devin_knowledge.example.description
}
```

### Importing Knowledge Resources

To import an existing Devin API knowledge resource into Terraform management, use the following command:

```bash
terraform import devin_knowledge.example knowledge_id
```

Where `knowledge_id` is the ID of an existing knowledge resource obtained from the Devin API.

When you run this command, the knowledge resource with the specified ID will be imported into your Terraform state file and can be managed by Terraform from that point forward.

#### Import Example

For example, if you have a Terraform configuration like this:

```hcl
resource "devin_knowledge" "imported" {
  # You only specify the ID initially (other attributes will be set automatically after import)
}
```

To import an existing knowledge resource into this resource:

```bash
terraform import devin_knowledge.imported existing-knowledge-id
```

After importing, running `terraform plan` will show the differences between your current configuration and the actual resource. You can then update your Terraform configuration to match the actual resource state.

## Development

### Clone the Repository

```bash
git clone https://github.com/hirosi1900day/terraform-provider-devin-knowledge.git
cd terraform-provider-devin-knowledge
```

### Install Dependencies

```bash
go mod tidy
```

### Building

```bash
go build -o terraform-provider-devin
# Or
make build
```

### Local Testing

1. Copy the built provider to your local plugin directory:

```bash
# For macOS/Linux
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/hirosi1900day/devin/1.0.0/$(go env GOOS)_$(go env GOARCH)
cp terraform-provider-devin ~/.terraform.d/plugins/registry.terraform.io/hirosi1900day/devin/1.0.0/$(go env GOOS)_$(go env GOARCH)/
```

2. Run the test Terraform code:

```bash
cd examples/simple
terraform init
terraform plan
terraform apply
```

### Test Mode

For development purposes, you can use a test API key that returns mock data instead of connecting to the actual Devin API:

```bash
DEVIN_API_KEY="test_api_key" terraform apply
```

For more details, see [Terraform documentation](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers).

## Release Process

1. Create a version tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

2. Create a GitHub release:
   - Create a release on the GitHub releases page
   - Attach the built binaries to the release

## Contributing

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## API Compatibility

This provider has been updated to work with the latest Devin API specification. There are some important notes about field mappings:

1. **Field Mapping**:
   - The Terraform provider's `description` attribute is mapped to the API's `body` field.
   - The API fields `trigger_description` and `parent_folder_id` are handled internally but not yet exposed as Terraform attributes.

2. **API Structure**:
   - The provider supports the `/knowledge` endpoint for listing all knowledge resources.
   - The response includes both knowledge items and folders as described in the [Devin API documentation](https://docs.devin.ai/api-reference/knowledge/list-knowledge).

3. **Future Enhancements**:
   - Future versions of this provider will expose all API fields directly, using field names that match the API specification.
   - A migration path will be provided for users moving from the current field naming to the updated field naming.

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

## Resources

- [Devin API Documentation](https://docs.devin.ai/api-reference)
- [Terraform Provider Development](https://developer.hashicorp.com/terraform/plugin/framework)
