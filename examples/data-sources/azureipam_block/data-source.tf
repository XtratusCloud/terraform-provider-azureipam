# Returns one block with network and usage information
data "azureipam_block" "expanded" {
  space              = "au"
  name               = "AustraliaEast"
  expand             = true
  append_utilization = true
}
output "expanded_block" {
  value = data.azureipam_block.expanded
}

# Returns one block without network and usage information
data "azureipam_block" "not_expanded" {
  space              = "au"
  name               = "AustraliaEast"
  expand             = false
  append_utilization = false
}
output "not_expanded_block" {
  value = data.azureipam_block.not_expanded
}
