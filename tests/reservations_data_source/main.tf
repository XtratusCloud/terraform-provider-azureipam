# We strongly recommend using the required_providers block to set the
# azureipam provider source and version being used
terraform {
  required_providers {
    azureipam = {
      #version = "~>1.1"
      source = "xtratuscloud/azureipam"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>3.116"
    }
  }
}

# Replace with appropriate values for your AZURE IPAM implementation. 
locals {
  # ipam_url   = "https://xtratusipam.dev.azure.apps-ferrovial.int" #"https://myazureipam.azurewebsites.net/"
  # ipam_apiId = "91db5c72-7c77-4ed7-81ab-7692754f7b64"             # "d47d5cd9-b599-4a6a-9d54-254565ff08de" #ApplicationId of the Engine Azure AD Application, see also the [IPAM deployment documentation](https://github.com/Azure/ipam/tree/main/docs/deployment)
}

## Get an access token for ipam engine application
# data "external" "get_access_token" {
#   program = ["az", "account", "get-access-token", "--resource", "api://${local.ipam_apiId}", "--query", "{accessToken:accessToken}"]
# }

provider "azurerm" {
  features {}
}

# Configure the Azure IPAM provider
provider "azureipam" {
  # api_url = local.ipam_url
  # token   = data.external.get_access_token.result.accessToken
}

