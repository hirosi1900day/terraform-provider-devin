# Import Block Usage Example

This directory contains an example of importing existing Devin API knowledge resources using the "import block" feature introduced in Terraform 1.5 and later.

## Preparation

1. First, initialize this directory:

```bash
terraform init
```

## How Import Blocks Work

Unlike traditional command-line imports (`terraform import`), import blocks are declared in the configuration file. This provides:

- Import operations can be managed as code
- Easier automation in CI pipelines
- Reproducible import processes

## Usage

1. The `main.tf` file includes an import block like this:

```hcl
import {
  to = devin_knowledge.imported_block
  id = "mock-knowledge-1"
}

resource "devin_knowledge" "imported_block" {
  name                = "Mock Knowledge 1"
  body                = "This is a mock knowledge for testing"
  trigger_description = "Test trigger description"
}
```

2. Run the following commands to apply the import block:

```bash
terraform plan
terraform apply
```

3. After completion, you can manage it as a regular resource.

## Re-running Import

To redo the import, follow these steps:

1. Remove the resource from the state file:

```bash
terraform state rm devin_knowledge.imported_block
```

2. Run `terraform plan` and `terraform apply` again.

## Importing Multiple Resources

The `main.tf` file also includes commented-out examples of importing multiple resources. To import multiple resources, uncomment and modify as needed.

## Notes

- Import blocks fulfill their purpose once applied. You can remove import blocks from the file after applying them.
- Specifying appropriate resource attributes in advance minimizes changes needed after import.
- When using a test API Key, use `test_api_key` to utilize the mock server.
- In production environments, set the API key as an environment variable (`DEVIN_API_KEY`) or specify it directly in the `api_key` parameter.

## More Information

For more details, refer to the [official Terraform documentation](https://developer.hashicorp.com/terraform/language/import). 