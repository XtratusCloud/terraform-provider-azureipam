# Create a CIDR reservation specifying multiple blocks
resource "azureipam_reservation" "new" {
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
output "reservation" {
  value = azureipam_reservation.new
}
