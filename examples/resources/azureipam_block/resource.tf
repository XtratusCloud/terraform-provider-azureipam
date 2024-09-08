# Create a new block in ua space
resource "azureipam_block" "new" {
  space = "au"
  name  = "AustraliaNorth"
  cidr  = "10.85.0.0/16"
}
output "space" {
  value = azureipam_block.new
}
