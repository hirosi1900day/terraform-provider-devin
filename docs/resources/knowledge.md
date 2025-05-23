---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "devin_knowledge Resource - devin"
subcategory: ""
description: |-
  Manages knowledge resources in the Devin API. Knowledge resources provide a way to supply information to the Devin AI system, which is triggered under specific conditions to generate appropriate responses.
---

# devin_knowledge (Resource)

This resource manages knowledge resources in the Devin API. A knowledge resource includes a name, body content, trigger description, and an optional parent folder ID.

Knowledge resources provide a way to supply information to the Devin AI system, which is triggered under specific conditions to generate appropriate responses.

## Example Usage

```terraform
resource "devin_knowledge" "example" {
  name                = "Example Knowledge"
  body                = "This is an example knowledge resource."
  trigger_description = "Use this knowledge when talking about examples."
  parent_folder_id    = "optional-folder-id"
}
```

## Import

Knowledge resources can be imported using the ID, which can be obtained from the Devin API:

### Command-line Import

```terraform
terraform import devin_knowledge.example knowledge_id
```

The ID should be provided as-is, without any additional prefixes. For example, if the knowledge ID is `note-123abc`, use it directly:

```terraform
terraform import devin_knowledge.example note-123abc
```

### Import Block (Terraform 1.5.0 and later)

Alternatively, you can use an import block in your configuration file:

```terraform
import {
  to = devin_knowledge.example
  id = "note-123abc"
}

resource "devin_knowledge" "example" {
  # Resource configuration will be filled in by Terraform
  # after importing the resource
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `body` (String) The content of the knowledge resource. Include any necessary information, such as text, markdown, or code snippets.
- `name` (String) The name of the knowledge resource. Set a clear and unique name for easy identification.
- `trigger_description` (String) The trigger description for the knowledge resource. This describes the scenarios in which this knowledge should be triggered.

### Optional

- `parent_folder_id` (String) The ID of the parent folder. Used to organize knowledge in folders.

### Read-Only

- `id` (String) The unique ID of the knowledge resource generated by the Devin API.
