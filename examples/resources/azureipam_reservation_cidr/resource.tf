# Deploy the azurerm vnet
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "Australia East"
}
resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  address_space = [azureipam_reservation_cidr.new.cidr]
  tags          = azureipam_reservation_cidr.new.tags ##Don't forget to add the auto-generated `X-IPAM-RES-ID` tag to the vnet.
}

# Create a CIDR reservation specifying a custom cidr
resource "azureipam_reservation_cidr" "new" {
  space         = "au"
  block         = "AustraliaEast"
  specific_cidr = "10.82.4.0/24"
  description   = "this is a test"
}
output "created" {
  value = azureipam_reservation_cidr.new
}
