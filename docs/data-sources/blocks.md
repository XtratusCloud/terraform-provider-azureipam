---
page_title: "azureipam_blocks Data Source - azureipam"
subcategory: ""
description: |-
  The blocks data source allows you to retrieve information about all blocks in the specified space, and their related information.
---

# azureipam_blocks (Data Source)

The blocks data source allows you to retrieve information about all blocks in the specified space, and their related information.

## Example Usage

```terraform
# Returns all blocks in au space with all allowed information
data "azureipam_blocks" "full" {
  space              = "au"
  expand             = true
  append_utilization = true
}
output "all_au_spaces_full" {
  value = data.azureipam_blocks.full.spaces
}

# Returns all blocks in au space without usage and child resources information
data "azureipam_blocks" "minimal" {
  space              = "au"
  expand             = false
  append_utilization = false
}
output "all_au_blocks_minimal" {
  value = data.azureipam_blocks.minimal.spaces
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `space` (String) Name of the `space` for which to read the related `blocks`.

### Optional

- `append_utilization` (Boolean) Indicates if utilization information for each network must be included.
- `expand` (Boolean) Indicates if network references to full network objects must be included.

### Read-Only

- `blocks` (Attributes List) List containing the `blocks` included in the specified `space`. (see [below for nested schema](#nestedatt--blocks))

<a id="nestedatt--blocks"></a>
### Nested Schema for `blocks`

Read-Only:

- `cidr` (String) The IPV4 range assigned to this block, in cidr notation.
- `externals` (Attributes List) List containing the `external networks` included in this `block`. (see [below for nested schema](#nestedatt--blocks--externals))
- `name` (String) Name of the block.
- `reservations` (Attributes List) List containing the `reservations` included in this `block`. (see [below for nested schema](#nestedatt--blocks--reservations))
- `size` (Number) Total IP's allowed in the `block` by its size.
- `used` (Number) Assigned IP's in the `block`.
- `vnets` (Attributes List) List containing the `vnet` included in this `block`. (see [below for nested schema](#nestedatt--blocks--vnets))

<a id="nestedatt--blocks--externals"></a>
### Nested Schema for `blocks.externals`

Read-Only:

- `cidr` (String) The IPV4 range reserved for the external network, in cidr notation.
- `description` (String) Text that describes the external network.
- `name` (String) Name of the external network.


<a id="nestedatt--blocks--reservations"></a>
### Nested Schema for `blocks.reservations`

Read-Only:

- `cidr` (String) The IPv4 range assigned to this reservation, in cidr notation.
- `created_by` (String) Email or identification of user that created the reservation.
- `created_on` (String) The date and time that the reservacion was created.
- `description` (String) Text that describes the reservation.
- `id` (String) The unique identifier of the reservation.
- `settled_by` (String) Email or identification of user that have settled the reservation.
- `settled_on` (String) The date and time when the reservation was settled.
- `status` (String) Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation


<a id="nestedatt--blocks--vnets"></a>
### Nested Schema for `blocks.vnets`

Read-Only:

- `id` (String) Resourece Id of the virtual network.
- `name` (String) Name of the virtual network.
- `prefixes` (List of String) The list of IPV4 prefixes assigned to this vnet, in cidr notation.
- `resource_group` (String) Name of the resource group where the `vnet` is deployed.
- `size` (Number) Total IP's allowed in the `vnet` by its size.
- `subnets` (Attributes List) List containing the `subnets` included in this `vnet`. (see [below for nested schema](#nestedatt--blocks--vnets--subnets))
- `subscription_id` (String) Id of the Azure subscription where the `vnet` is deployed.
- `tenant_id` (String) Id of the Azure tenant where the `vnet` is deployed.
- `used` (Number) Assigned IP's in the `vnet`.

<a id="nestedatt--blocks--vnets--subnets"></a>
### Nested Schema for `blocks.vnets.subnets`

Read-Only:

- `name` (String) Name of the subnet.
- `prefix` (String) The IPV4 prefix assigned to this block, in cidr notation.
- `size` (Number) Total IP's allowed in the `subnet` by its size.
- `used` (Number) Assigned IP's in the `subnet`.