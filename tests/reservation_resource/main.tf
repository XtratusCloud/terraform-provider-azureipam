# Create a CIDR reservation
resource "azureipam_reservation" "new" {
  space          = "au"
  block          = "AustraliaEast"
  size           = 24
  description    = "this is a test"
  reverse_search = true
  smallest_cidr  = true
}
output "reservation" {
  value = azureipam_reservation.new
}