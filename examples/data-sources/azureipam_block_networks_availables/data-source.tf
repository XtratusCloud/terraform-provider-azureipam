# Returns one block with network and usage information
data "azureipam_block_networks_availables" "example" {
  space = "au"
  block = "AustraliaEast"
}
output "availables" {
  value = data.azureipam_block_networks_availables.example
}