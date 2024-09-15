# Returns all blocks with network and usage information
data "azureipam_blocks" "expanded" {
  space              = "au"
  expand             = true
  append_utilization = true
}
output "expanded_blocks" {
  value = data.azureipam_blocks.expanded
}


# Returns all blocks without network and usage information
data "azureipam_blocks" "not_expanded" {
  space              = "au"
  expand             = false
  append_utilization = false
}
output "not_expanded_blocks" {
  value = data.azureipam_blocks.not_expanded
}
