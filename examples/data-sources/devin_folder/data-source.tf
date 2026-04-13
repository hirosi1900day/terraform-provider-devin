# Example usage of devin_folder data source
# By ID
data "devin_folder" "example_by_id" {
  id = "folder-xxxx"  # folder_id
}

# By Name
data "devin_folder" "example_by_name" {
  name = "My Folder"
}

# Output example
output "folder_by_id" {
  value = data.devin_folder.example_by_id
}

output "folder_by_name" {
  value = data.devin_folder.example_by_name
}
