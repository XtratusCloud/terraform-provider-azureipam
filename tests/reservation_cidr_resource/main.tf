# Create a CIDR reservation specifying only a block
resource "azureipam_reservation_cidr" "new" {
  space         = "au"
  block         = "AustraliaEast"
  specific_cidr = "10.82.4.0/24"
  description   = "Reservation test in specified block"
}
output "reservation" {
  value = azureipam_reservation_cidr.new
}
 
