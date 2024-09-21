# Returns the list of networks actively associated to a space block
data "azureipam_block_networks" "example" {
  space = "au"
  block = "AustraliaEast"
}
output "associated" {
  value = data.azureipam_block_networks.example
}
