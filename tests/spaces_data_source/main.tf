# Returns all spaces with network and usage information
data "azureipam_spaces" "expanded" {
  expand             = true
  append_utilization = true
}
output "expanded_spaces" {
  value = data.azureipam_spaces.expanded
}


# Returns all spaces without network and usage information
data "azureipam_spaces" "not_expanded" {
  expand             = false
  append_utilization = false
}
output "not_expanded_spaces" {
  value = data.azureipam_spaces.not_expanded
}
