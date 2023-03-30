

variable "reservation_id" {
  type    = string
  default = "G7Ngtt3Ek2GY5MHoQtmoTh"
}

# Returns all reservations in au/AustraliaEast
data "azureipam_reservations" "all" {
  space = "au"
  block = "AustraliaEast"
}
output "all_reservations" {
  value = data.azureipam_reservations.all.reservations
}

###
## IPAM Reservation
###
resource "azureipam_reservation" "created" {
  space          = "au"
  block          = "AustraliaEast"
  size           = 24
  description    = "test xtratusipam terraform provider"
  reverse_search = true
  smallest_cidr  = true
}
output "created" {
  value = azureipam_reservation.created
}


###
## VNET
###
resource "azurerm_resource_group" "vnet" {
  name     = "RG-WE-D-TERRAFORM-PROVIDER-01"
  location = "westeurope"
}
resource "azurerm_virtual_network" "vnet" {
  name                = "vnet-we-t-terraform-provider-01"
  resource_group_name = azurerm_resource_group.vnet.name
  location            = azurerm_resource_group.vnet.location

  address_space = [azureipam_reservation.created.cidr]
  tags          = azureipam_reservation.created.tags
}
