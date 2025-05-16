# Example usage of devin_folder data source
# By ID
data "devin_folder" "example_by_id" {
  id = "mock-folder-1"
}

# By Name
data "devin_folder" "example_by_name" {
  name = "モックフォルダ1"
}

# Output example
output "folder_by_id" {
  value = data.devin_folder.example_by_id
}

output "folder_by_name" {
  value = data.devin_folder.example_by_name
}
