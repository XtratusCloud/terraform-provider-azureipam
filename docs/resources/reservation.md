---
page_title: "reservation Resource - terraform-provider-azureipam"
subcategory: ""
description: |-
  The reservation resource allows you to create a IPAM reservation in the specific space and block.
---

# Resource `azureipam_reservation`

The reservation resource allows you to create a IPAM reservation in the specific space and block.

## Example Usage

```terraform
# Create a CIDR reservation
resource "azureipam_reservation" "example" {
  space       = "au"
  block       = "AustraliaEast"
  size        = 24
  description = "this is a test"
  reverse_search = true
  smallest_cidr  = true
}

# Deploy the azurerm vnet
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "Australia East"
}
resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group

  address_space = [azureipam_reservation.example.cidr]
  tags          = azureipam_reservation.example.tags ##Don't forget to add the auto-generated `X-IPAM-RES-ID` tag to the vnet.
}
```

## Argument Reference

- `space` - (Required) name of the existing space in the IPAM application.
- `block` - (Required) name of the existing block, related to the specified space, in which the reservation is to be made.
- `size` - (Required) integer value to indicate the subnet mask bits, which defines the size of the vnet to reserve (example 24 for a /24 subnet).
- `description` - (Optional) description text that describe the reservation, that will be added as an additional tag.
- `reverse_search` - (Optional) New networks will be created as close to the end of the block as possible?. Defaults to `false`.
- `smallest_cidr` - (Optional) New networks will be created using the smallest possible available block? (e.g. it will not break up large CIDR blocks when possible) .Defaults to `false`.

## Attributes Reference

The reservation item also contains the following attributes. 

### Reservation

- `id` - The unique identifier of the generated reservation.
- `cidr` - The assigned and reserved range, in cidr notation.
- `created_on` - The date and time that the reservacion was created.
- `created_by`- Email or identification of user that created the reservation.
- `status` - Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation.
- `tags` - Auto-generated tags for the reservation. Particular relevance the 'X-IPAM-RES-ID' tag, since it must be included in the vnet creation in order that the IPAM solution automatically considers the reservation as completed.