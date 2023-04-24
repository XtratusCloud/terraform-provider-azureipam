---
page_title: "reservations Data Source - terraform-provider-azureipam"
subcategory: ""
description: |-
  The reservations data source allows you to retrieve information about all existing reservations in the specific space and block.
---

# Data Source `azureipam_reservations`

The reservations data source allows you to retrieve information about all existing reservations in the specified space and block.

## Example Usage

```terraform
# Returns all reservations in au/AustraliaEast block
data "azureipam_reservations" "all" {
  space = "au"
  block = "AustraliaEast"
  include_settled = true
}
output "all_reservations" {
  value = data.azureipam_reservations.all.reservations
}


# Returns not settled reservations in au/AustraliaEast block
data "azureipam_reservations" "not_settled" {
  space = "au"
  block = "AustraliaEast"
}
output "not_settled_reservations" {
  value = data.azureipam_reservations.not_settled.reservations
}
```

## Argument Reference

- `space` - (Required) name of the existing space in the IPAM application.
- `block` - (Required) name of the existing block, related to the specified space.
- `include_settled` - (Optional) Settled reservations must be also included? Defaults to `false`.

## Attributes Reference

Each reservation item contains the following attributes.

- `id` - The unique identifier of the reservation.
- `cidr` - The assigned and reserved range, in cidr notation.
- `created_on` - The date and time that the reservacion was created.
- `created_by`- Email or identification of user that created the reservation.
- `settled_on` - The date and time when the reservation was settled.
- `settled_by` - Email or identification of user that have settled the reservation.
- `status` - Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation.
- `tags` - Auto-generated tags for the reservation. Particular relevance the 'X-IPAM-RES-ID' tag, since it must be included in the vnet creation in order that the IPAM solution automatically considers the reservation as completed.