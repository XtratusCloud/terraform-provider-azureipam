# Returns one reservation in au/AustraliaEast block
data "azureipam_reservation" "example" {
  space = "au"
  block = "AustraliaSoutheast"
  id    = "dVUApz478E6mwjhmp6tr44"
}
output "reservation" {
  value = data.azureipam_reservation.example
}
