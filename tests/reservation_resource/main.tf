# Create a CIDR reservation specifying only a block
resource "azureipam_reservation" "new" {
  space = "au"
  blocks = [
    "AustraliaEast"
  ]
  size           = 24
  description    = "Reservation test in specified block"
  reverse_search = true
  smallest_cidr  = true
}
output "reservation" {
  value = azureipam_reservation.new
}

# Create a CIDR reservation specifying multiple blocks
resource "azureipam_reservation" "block_list" {
  space = "au"
  blocks = [
    "AustraliaSoutheast",
    "AustraliaEast"
  ]
  size           = 25
  description    = "Reservation test with block list"
  reverse_search = false
  smallest_cidr  = false
}
output "reservation_block_list" {
  value = azureipam_reservation.block_list
}
