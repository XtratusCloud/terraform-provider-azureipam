# Returns the list of resource Ids of available networks to be associated to a space block
data "azureipam_block_networks_availables" "example" {
  space = "au"
  block = "AustraliaEast"
}
output "availables" {
  value = data.azureipam_block_networks_availables.example
}
