---
page_title: "Troubleshooting Guide - terraform-provider-devin"
subcategory: ""
description: |-
  Provides common issues and solutions when using the Devin Terraform provider.
---

# Troubleshooting Guide

This guide provides common issues and solutions that you might encounter when using the Devin Terraform provider.

## Authentication Errors

### Problem: API Authentication Failure

```
Error: Failed to configure provider
...
Error during Devin API request: API error: Invalid API key
```

### Solution:

1. Verify that your API key is correctly set in your configuration:
   ```terraform
   provider "devin" {
     api_key = "your_correct_api_key"
   }
   ```

2. Verify that your environment variable is correctly set:
   ```bash
   export DEVIN_API_KEY="your_correct_api_key"
   echo $DEVIN_API_KEY  # Verify the setting
   terraform apply
   ```

3. Verify that your API key is valid and not expired. Check the status of your API key in the Devin admin dashboard.

## Resource Operation Errors

### Problem: Knowledge Resource Creation Failure

```
Error: Failed to create knowledge
...
Error during Devin API request: API error: Required field missing
```

### Solution:

1. Verify that all required fields are set:
   ```terraform
   resource "devin_knowledge" "example" {
     name                = "Sample Knowledge"  # Required
     body                = "Knowledge content"  # Required
     trigger_description = "Trigger conditions"  # Required
     # parent_folder_id is optional
   }
   ```

2. Check that the length or content of your fields does not exceed any limits. Be especially careful when dealing with large knowledge content.

### Problem: Resource Not Found

```
Error: Failed to retrieve knowledge
...
Error during Devin API request: knowledge resource with ID 'xyz123' not found
```

### Solution:

1. Verify that the specified ID exists.
2. Check if the knowledge resource exists in the Devin admin dashboard.
3. Verify that your API key has permission to access the resource.

## Import Issues

### Problem: Import Command Fails

```
Error: resource address "devin_knowledge.imported" does not exist
```

### Solution:

1. Before importing, make sure the resource configuration exists in your Terraform file:
   ```terraform
   resource "devin_knowledge" "imported" {
     # Attributes will be set after import
   }
   ```

2. Verify that you are using the correct resource name and ID:
   ```bash
   terraform import devin_knowledge.imported correct-knowledge-id
   ```

## Development Mode and Debugging

### Enabling Debug Logging

You can get detailed debug information by setting the Terraform log level:

```bash
export TF_LOG=DEBUG
terraform apply
```

### Using Test Mode

For development or testing purposes, you can use test mode that returns mock data:

```bash
export DEVIN_API_KEY="test_api_key"
terraform apply
```

## General Troubleshooting Steps

1. **Use the Latest Version**: Make sure you're using the latest version of the provider.
   ```terraform
   terraform {
     required_providers {
       devin = {
         source  = "hirosi1900day/devin"
         version = "~> 0.0.1"  # Specify the latest version
       }
     }
   }
   ```

2. **Reinitialize Terraform**: After changing configurations, reinitialize Terraform.
   ```bash
   terraform init -upgrade
   ```

3. **Refresh State**: Update your current state.
   ```bash
   terraform refresh
   ```

4. **Clean Run**: In some cases, removing the `.terraform` directory and reinitializing might solve issues.
   ```bash
   rm -rf .terraform
   terraform init
   ```

5. **Check Support Resources**: If your issue persists, check these resources:
   - [Devin API Documentation](https://docs.devin.ai/api-reference)
   - [GitHub Issues Page](https://github.com/hirosi1900day/terraform-provider-devin-knowledge/issues)
   - Provider's [Release Notes](https://github.com/hirosi1900day/terraform-provider-devin-knowledge/releases)
