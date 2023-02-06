

variable "reservation_id" {
  type    = string
  default = "nVDVdrB7YCrtGzwxuXjzyh"
}

# Returns all reservations in au/AustraliaEast
data "azureipam_reservations" "all" {
  space = "au"
  block = "AustraliaEast"
}
output "all_reservations" {
  value = data.azureipam_reservations.all.reservations
}

resource "azureipam_reservation" "created" {
  space       = "au"
  block       = "AustraliaEast"
  size        = 24
  description = "prueba fnieto"
}
output "created" {
  value = azureipam_reservation.created
}