---
page_title: "Usage Guide - terraform-provider-devin"
subcategory: ""
description: |-
  Explains various usage scenarios and best practices for the Devin Terraform provider.
---

# Devin Provider Usage Guide

This guide explains various usage scenarios and best practices for the Devin Terraform provider.

## Basic Workflow

### 1. Provider Setup

First, set up the Devin provider in your Terraform configuration file (typically `main.tf`):

```terraform
terraform {
  required_providers {
    devin = {
      source  = "hirosi1900day/devin"
      version = "~> 0.0.1"
    }
  }
}

provider "devin" {
  api_key = "your_api_key" # Can also use the DEVIN_API_KEY environment variable
}
```

### 2. Creating Knowledge Resources

Define your knowledge resources:

```terraform
resource "devin_knowledge" "example" {
  name                = "User Authentication Best Practices"
  body                = "When implementing user authentication, follow these best practices: ..."
  trigger_description = "When questions about user authentication, security, or password management are asked"
}
```

### 3. Running Terraform Commands

Use standard Terraform commands to apply your configuration:

```bash
terraform init
terraform plan
terraform apply
```

## Advanced Use Cases

### Managing Knowledge with Folder Hierarchy

You can organize multiple knowledge resources in folders:

```terraform
resource "devin_knowledge" "security_guideline1" {
  name                = "Security Guideline: Passwords"
  body                = "Detailed information about password policies..."
  trigger_description = "When questions about password security are asked"
  parent_folder_id    = "security_folder_id"
}

resource "devin_knowledge" "security_guideline2" {
  name                = "Security Guideline: API Authentication"
  body                = "Best practices for API authentication..."
  trigger_description = "When questions about API authentication are asked"
  parent_folder_id    = "security_folder_id"
}
```

### Importing Existing Resources

To manage an existing knowledge resource with Terraform, write the configuration like this:

```terraform
resource "devin_knowledge" "imported" {
  # Other attributes will be set after import
}
```

Then import it with this command:

```bash
terraform import devin_knowledge.imported existing-knowledge-id
```

After importing, run `terraform plan` to see the differences between the actual resource state and your configuration, then update your Terraform configuration accordingly.

## Best Practices

### 1. Naming Conventions

It's recommended to use clear and consistent naming conventions for your knowledge resources:

```terraform
resource "devin_knowledge" "dev_python_style_guide" {
  name                = "[Development] Python Style Guide"
  body                = "# Python Coding Style Guide\n\n..."
  trigger_description = "When questions about Python coding conventions or style are asked"
}

resource "devin_knowledge" "dev_javascript_style_guide" {
  name                = "[Development] JavaScript Style Guide"
  body                = "# JavaScript Coding Style Guide\n\n..."
  trigger_description = "When questions about JavaScript coding conventions or style are asked"
}
```

### 2. Optimizing Trigger Descriptions

Make trigger descriptions specific and clear to ensure knowledge is used in appropriate situations:

```terraform
resource "devin_knowledge" "database_backup" {
  name                = "Database Backup Procedures"
  body                = "# Database Backup Procedures\n\n..."
  trigger_description = "When questions or procedures about PostgreSQL, MySQL, database backups, restoration, or disaster recovery are asked"
}
```

### 3. Loading Content from External Files

For large knowledge content, it's recommended to load from external files:

```terraform
resource "devin_knowledge" "architecture_guide" {
  name                = "System Architecture Guide"
  body                = file("${path.module}/knowledge/architecture_guide.md")
  trigger_description = "When questions about system architecture, design patterns, or component structure are asked"
}
```

### 4. Using Variables and Modules

Leverage Terraform variables and modules for reusable patterns:

```terraform
variable "common_prefix" {
  type    = string
  default = "[Company] "
}

resource "devin_knowledge" "company_policy" {
  name                = "${var.common_prefix}Information Security Policy"
  body                = file("${path.module}/policies/security_policy.md")
  trigger_description = "When questions about company information security policies are asked"
}
```

### 5. Using Environment Variables

Use environment variables for sensitive information instead of hardcoding:

```bash
export DEVIN_API_KEY="your_api_key"
terraform apply
```
