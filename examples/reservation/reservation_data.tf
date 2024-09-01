# Returns all reservations in au/AustraliaEast block
data "azureipam_reservations" "all" {
  space           = "au"
  block           = "AustraliaEast"
  include_settled = true
}
output "all_reservations" {
  value = data.azureipam_reservations.all.reservations
}


# Returns not settled reservations in au/AustraliaEast block
data "azureipam_reservations" "not_settled" {
  space = "au"
  block = "AustraliaEast"
}
output "not_settled_reservations" {
  value = data.azureipam_reservations.not_settled.reservations
}
