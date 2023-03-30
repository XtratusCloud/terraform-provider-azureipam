provider "azurerm" {
  subscription_id = "5c5ffb7c-8ed5-4fd1-9eb5-da5ee72a6490" #IA-CORP-TFMODULES-NONPROD
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}