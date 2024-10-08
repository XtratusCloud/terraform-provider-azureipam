---
page_title: "azureipam_reservation Data Source - azureipam"
subcategory: ""
description: |-
  The reservation data source allows you to retrieve a specific reservation by id in the specified space and block.
---

# azureipam_reservation (Data Source)

The reservation data source allows you to retrieve a specific reservation by id in the specified space and block.

## Example Usage

```terraform
# Returns one reservation in au/AustraliaEast block
data "azureipam_reservation" "example" {
  space = "au"
  block = "AustraliaSoutheast"
  id    = "dVUApz478E6mwjhmp6tr44"
}
output "reservation" {
  value = data.azureipam_reservation.example
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `block` (String) Name of the  block where the reservation is allocated.
- `id` (String) The unique identifier of the reservation.
- `space` (String) Name of the space where the reservation is allocated.

### Read-Only

- `cidr` (String) The assigned and reserved range, in cidr notation.
- `created_by` (String) Email or identification of user that created the reservation.
- `created_on` (String) The date and time that the reservacion was created.
- `description` (String) Description text that describe the reservation.
- `settled_by` (String) Email or identification of user that have settled the reservation.
- `settled_on` (String) The date and time when the reservation was settled.
- `status` (String) Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation.
- `tags` (Map of String) Auto-generated tags for the reservation. Particular relevance the 'X-IPAM-RES-ID' tag, since it must be included in the vnet creation in order that the IPAM solution automatically considers the reservation as completed.