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
}
```

## Argument Reference

- `space` - (Required) name of the existing space in the IPAM application.
- `block` - (Required) name of the existing block, related to the specified space, in which the reservation is to be made.
- `size` - (Required) integer value to indicate the subnet mask bits, which defines the size of the vnet to reserve (example 24 for a /24 subnet).
- `description` - (Optional) description text that describe the reservation, that will be added as an additional tag.

## Attributes Reference

The reservation item also contains the following attributes. 

### Reservation

- `id` - The unique identifier of the generated reservation.
- `cidr` - The assigned and reserved range, in cidr notation.
- `created_on` - The date and time that the reservacion was created.
- `status` - Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation.
- `user_id` - Email or identification of user that created the reservation.
- `tags` - All tags specified for the reservation, including the description (if specified). Particular relevance has the tag with id 'X-IPAM-RES-ID', since it must be included in the vnet creation in order that the IPAM solution automatically considers the reservation as completed.