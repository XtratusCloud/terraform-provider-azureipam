# Create a CIDR reservation
resource "azureipam_reservation" "created" {
  space          = "au"
  block          = "AustraliaEast"
  size           = 24
  description    = "this is a test"
  reverse_search = true
  smallest_cidr  = true
}
output "created" {
  value = azureipam_reservation.created
}


# Deploy the azurerm vnet
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "Australia East"
}
resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group

  address_space = [azureipam_reservation.example.cidr]
  tags          = azureipam_reservation.example.tags ##Don't forget to add the auto-generated `X-IPAM-RES-ID` tag to the vnet.
}