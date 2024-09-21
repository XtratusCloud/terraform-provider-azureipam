# Returns all spaces with all allowed information
data "azureipam_spaces" "full" {
  expand             = true
  append_utilization = true
}
output "all_spaces_full" {
  value = data.azureipam_spaces.full.spaces
}

# Returns all spaces without usage and child resources information
data "azureipam_spaces" "minimal" {
  expand             = false
  append_utilization = false
}
output "all_spaces_minimal" {
  value = data.azureipam_spaces.minimal.spaces
}
