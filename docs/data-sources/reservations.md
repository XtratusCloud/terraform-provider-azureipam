---
page_title: "reservations Data Source - terraform-provider-azureipam"
subcategory: ""
description: |-
  The reservations data source allows you to retrieve information about all existing reservations in the specific space and block.
---

# Data Source `azureipam_reservations`

The reservations data source allows you to retrieve information about all existing reservations in the specific space and block.

## Example Usage

```terraform
# Returns all reservations in au/AustraliaEast
data "azureipam_reservations" "all" {
  space = "au"
  block = "AustraliaEast"
}
output "all_reservations" {
  value = data.azureipam_reservations.all.reservations
}
```

## Argument Reference

- `space` - (Required) name of the existing space in the IPAM application.
- `block` - (Required) name of the existing block, related to the specified space.

## Attributes Reference

Each reservation item contains the following attributes.

- `id` - The unique identifier of the reservation.
- `cidr` - The assigned and reserved range, in cidr notation.
- `created_on` - The date and time that the reservacion was created.
- `status` - Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation.
- `user_id` - Email or identification of user that created the reservation.
- `tags` - All tags specified for the reservation, including the description (if specified). Particular relevance has the tag with id 'X-IPAM-RES-ID', since it must be included in the vnet creation in order that the IPAM solution automatically considers the reservation as completed.