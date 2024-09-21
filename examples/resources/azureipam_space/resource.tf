# Create a new space
resource "azureipam_space" "new" {
  name        = "asia"
  description = "Asia Description"
}
output "space" {
  value = azureipam_space.new
}
