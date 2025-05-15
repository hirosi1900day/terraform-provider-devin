---
page_title: "devin_knowledge Data Source - terraform-provider-devin"
subcategory: ""
description: |-
  Retrieves knowledge resource information from the Devin API.
---

# devin_knowledge (Data Source)

This data source retrieves knowledge resource information from the Devin API. It allows you to access details of existing knowledge resources and use that information in other Terraform resources or outputs.

## Usage

### Basic Usage

```terraform
data "devin_knowledge" "example" {
  id = "knowledge_id_123"
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

### Combined with Resources

```terraform
resource "devin_knowledge" "new_example" {
  name                = "New Knowledge"
  body                = "This is a newly created knowledge resource"
  trigger_description = "Trigger description for the new knowledge"
}

data "devin_knowledge" "existing" {
  id = devin_knowledge.new_example.id
}

output "created_knowledge" {
  value = "${data.devin_knowledge.existing.name} has been created. ID: ${data.devin_knowledge.existing.id}"
}
```

### Conditional Resource Creation

```terraform
data "devin_knowledge" "lookup" {
  id = "knowledge_id_123"
}

resource "devin_knowledge" "conditional" {
  count = data.devin_knowledge.lookup.name == "Specific Name" ? 1 : 0
  
  name                = "Conditional Knowledge"
  body                = "This knowledge is created based on a condition"
  trigger_description = "Trigger description for conditional knowledge"
  parent_folder_id    = data.devin_knowledge.lookup.parent_folder_id
}
```

## Schema

### Required

- `id` (String) - The ID of the knowledge resource. This is the unique identifier for the knowledge resource you want to retrieve.

### Read-Only

- `name` (String) - The name of the knowledge resource.
- `body` (String) - The content of the knowledge resource. This contains the full content of the knowledge.
- `trigger_description` (String) - The trigger description for the knowledge resource. This describes under what conditions the knowledge should be triggered.
- `parent_folder_id` (String) - The ID of the parent folder. If the knowledge is placed within a specific folder, this will contain the folder's ID. If not specified, it will be empty.

## Data Source Usage

This data source is particularly useful in the following scenarios:

1. **Accessing Existing Knowledge Information**: Retrieve detailed information about knowledge resources that have already been created.

2. **Conditional Resource Creation**: Create resources conditionally based on attributes of existing knowledge.

3. **Outputs and Dependencies**: Output knowledge details or use them in other resources.

4. **Configuration Validation**: Compare the actual API data with your configuration files to ensure consistency.
