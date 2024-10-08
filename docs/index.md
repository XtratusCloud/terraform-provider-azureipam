---
page_title: "Provider: AzureIpam"
subcategory: ""
description: |-
  Terraform provider to manage reservations in Azure IPAM solution through REST API.
---

# AzureIpam Provider

This provider is intended to manage the reservation of network ranges in the [Azure IPAM](https://github.com/Azure/ipam) solution. IPAM solution is a simple, straightforward way to manage IP address spaces in Azure, and it's's required to have a previous implementation of this solution.

The provider makes use of the IPAM REST API to manage CIDR range reservations in a space and block from those configured in the application.

> **NOTE** the provider is aligned with the functionality included in the Azure IPAM solution in the version published on 18 April 2023, in the Pull Request [#113](https://github.com/Azure/ipam/pull/113), so it is necessary that your IPAM implementation have to be based on that version or later.

## Example Usage

Do not keep your authentication token in HCL, use Terraform environment variables or generate as part of the deploymenet process.

```terraform
# We strongly recommend using the required_providers block to set the
# azureipam provider source and version being used
terraform {
  required_providers {
    azureipam = {
      version = "~>2.0"
      source  = "xtratuscloud/azureipam"
    }
  }
}

# Replace with appropriate values for your AZURE IPAM implementation. 
locals {
  ipam_url   = "https://myazureipam.azurewebsites.net"
  ipam_apiId = "d47d5cd9-b599-4a6a-9d54-254565ff08de" #ApplicationId of the Engine Azure AD Application, see also the [IPAM deployment documentation](https://github.com/Azure/ipam/tree/main/docs/deployment)
}

## Get an access token for ipam engine application
data "external" "get_access_token" {
  program = ["az", "account", "get-access-token", "--resource", "api://${local.ipam_apiId}", "--query", "{accessToken:accessToken}"]
}

# Configure the Azure IPAM provider
provider "azureipam" {
  api_url                = local.ipam_url
  token                  = data.external.get_access_token.result.accessToken
  skip_cert_verification = true //ONLY recommended for development environments
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_url` (String) The root url of the APIM REST API solution to be used, without the /api url suffix. Must be also assigned at AZUREIPAM_API_URL environment variable.
- `skip_cert_verification` (Boolean) Specifies it the certificate chain validation must be skipped calling the API endpoint. Default to false.
- `token` (String, Sensitive) The bearer token to be used when authenticating to the API. Must be also assigned at AZUREIPAM_TOKEN environment variable.