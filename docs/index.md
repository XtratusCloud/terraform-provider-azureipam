---
page_title: "Provider: AzureIpam"
subcategory: ""
description: |-
  Terraform provider to manage reservations in Azure IPAM solution through REST API
---

# AzureIpam Provider
This provider is intended to manage the reservation of network ranges in the [Azure IPAM](https://github.com/Azure/ipam) solution. IPAM solution is a simple, straightforward way to manage IP address spaces in Azure, and it's's required to have a previous implementation of this solution.

The provider makes use of the IPAM REST API to manage CIDR range reservations in a space and block from those configured in the application.

> **NOTE** that this provider makes use of a functionality not implemented in the open-source version of the Azure IPAM solution, which allows to include additional tags in the creation of new reservations (used for description). 

## Example Usage

Do not keep your authentication token in HCL, use Terraform environment variables or generate as part of the deploymenet process.

```terraform
# We strongly recommend using the required_providers block to set the
# azureipam provider source and version being used
terraform {
  required_providers {
    azureipam = {
      version = "0.1.1"
      source  = "xtratuscloud/azureipam"
    }
  }
}

## get an access token for ipam engine application
data "external" "get_access_token" {
  program = ["az", "account", "get-access-token", "--resource", "api://d47d5cd9-b599-4a6a-9d54-254565ff08de"]
}

# Configure the Azure IPAM provider
provider "azureipam" {
  api_url = local.ipam_url_dev
  token   = data.external.get_access_token.result.accessToken
}

# Create a CIDR reservation
resource "azureipam_reservation" "example" {
  space       = "au"
  block       = "AustraliaEast"
  size        = 24
  description = "this is a test"
}
```

## Argument Reference

* `api_url` - (Optional) The root url of the APIM REST API solution to be used, without the /api url suffix. This can also be sourced from the `AZUREIPAM_API_URL` Environment Variable.
* `token` - (Optional) The bearer token to be used when authenticating to the API. This can also be sourced from the `AZUREIPAM_TOKEN` Environment Variable.


## Special Considerations
Due to the current behaviour of the IPAM application, as the reservation is deleted once the vnet is deployed, an error avoidance mechanism has been implemented, which takes the current values when trying to update the state. This mechanism assumes that the reservation search is only performed when the element is already in the tfstate, to refresh the state information if needed, and it's not performed in the initial plan.