# Return all external networks in a space/block
data "azureipam_externals" "all" {
  space = "au"
  block = "AustraliaSoutheast"
}
output "all_externals" {
  value = data.azureipam_externals.all
}
