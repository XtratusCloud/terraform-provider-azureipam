---
page_title: "azureipam_space Resource - azureipam"
subcategory: ""
description: |-
  The space resource allows you to create a IPAM space.
---

# azureipam_space (Resource)

The space resource allows you to create a IPAM space.

## Example Usage

```terraform
# Create a new space
resource "azureipam_space" "new" {
  name        = "asia"
  description = "Asia Description"
}
output "space" {
  value = azureipam_space.new
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) Description text that describe the space.
- `name` (String) Name of the space.

## Import

Spaces can be imported using the name of the IPAM space, e.g.

```shell
terraform import azureipam_space.new asia
```
