# Return all external networks in a space/block
data "azureipam_external" "example" {
  space = "au"
  block = "AustraliaSoutheast"
  name  = "acctest"
}
output "external" {
  value = data.azureipam_external.example
}
