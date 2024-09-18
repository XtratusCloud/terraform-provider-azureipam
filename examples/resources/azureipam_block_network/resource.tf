# Associate a virtual network to a block
resource "azureipam_block_network" "new" {
  space = "au"
  block = "AustraliaEast"
  id    = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01"
}
output "block_network" {
  value = azureipam_block_network.new
}
