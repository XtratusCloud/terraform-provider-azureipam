# Returns all blocks in au space with all allowed information
data "azureipam_blocks" "full" {
  space              = "au"
  expand             = true
  append_utilization = true
}
output "all_au_spaces_full" {
  value = data.azureipam_blocks.full.spaces
}

# Returns all blocks in au space without usage and child resources information
data "azureipam_blocks" "minimal" {
  space              = "au"
  expand             = false
  append_utilization = false
}
output "all_au_blocks_minimal" {
  value = data.azureipam_blocks.minimal.spaces
}
