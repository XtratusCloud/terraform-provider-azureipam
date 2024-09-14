# Create a new external in ua/AustraliaSoutheast block
resource "azureipam_external" "new" {
  space       = "au"
  block       = "AustraliaSoutheast"
  name        = "acctest"
  description = "External Network for Acceptance Tests"
  cidr        = "10.83.6.0/24"
}
output "external" {
  value = azureipam_external.new
}
