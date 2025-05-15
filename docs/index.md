---
page_title: "Devin Provider"
subcategory: ""
description: |-
  Manage Devin API resources with the Terraform provider.
---

# Devin Provider

This Terraform provider allows you to manage [Devin API](https://api.devin.ai) resources as infrastructure as code. The provider supports creating, retrieving, updating, and deleting knowledge resources, enabling you to maintain knowledge resources as version-controlled configuration files.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.0+
- Access token for the [Devin API](https://docs.devin.ai/api-reference)

## Installation

```terraform
terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "~> 0.0.1"
    }
  }
}
```

## Usage

### Provider Configuration

```terraform
terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "~> 0.0.1"
    }
  }
}

# Configure the provider
provider "devin" {
  # Set API Key (or use the DEVIN_API_KEY environment variable)
  api_key = "your_api_key"
}
```

### Creating a Knowledge Resource

```terraform
resource "devin_knowledge" "example" {
  name                = "Sample Knowledge"
  body                = "This is a sample knowledge created with Terraform"
  trigger_description = "This knowledge is triggered under specific conditions"
  # parent_folder_id  = "optional-folder-id" # Optional parameter
}
```

### Using a Knowledge Data Source

```terraform
data "devin_knowledge" "example" {
  id = "knowledge_id"
}

output "knowledge_name" {
  value = data.devin_knowledge.example.name
}

output "knowledge_body" {
  value = data.devin_knowledge.example.body
}

output "knowledge_trigger_description" {
  value = data.devin_knowledge.example.trigger_description
}
```

## Authentication and Provider Configuration

To configure the provider, specify the API key required for authenticating with the Devin API:

```terraform
provider "devin" {
  api_key = "your_api_key"
}
```

You can also provide the API key using an environment variable:

```shell
export DEVIN_API_KEY="your_api_key"
```

And then specify the provider in your Terraform configuration file without the API key:

```terraform
provider "devin" {}
```

### Provider Arguments

The following arguments are supported:

* `api_key` - (Optional) The API key for the Devin API. Can also be set via the `DEVIN_API_KEY` environment variable.

## Development Mode

For development or testing purposes, you can use a test API key that returns mock data instead of connecting to the actual Devin API:

```shell
DEVIN_API_KEY="test_api_key" terraform apply
```
