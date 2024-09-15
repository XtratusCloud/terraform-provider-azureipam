# We strongly recommend using the required_providers block to set the
# azureipam provider source and version being used
terraform {
  required_providers {
    azureipam = {
      source = "xtratuscloud/azureipam"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>3.116"
    }
  }
}

provider "azurerm" {
  features {}
}

# REMEMBER to set AZUREIPAM_API_URL and AZUREIPAM_TOKEN env variables
provider "azureipam" {
  skip_cert_verification = true
}