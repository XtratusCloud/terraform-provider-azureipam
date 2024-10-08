# Changelog

## 2023.01.21 - v0.1.0

### Added
+ Provider azureipam
+ Resource reservation
+ Data reservation

### Fixed


## 2023.02.20 - v0.1.1

### Added

### Fixed
+ Provider azureipam `api_url` param can also be sourced from the `AZUREIPAM_API_URL` Environment Variable.
+ Avoid error when reservation is not found in redeployments (see [Special Considerations](https://registry.terraform.io/providers/XtratusCloud/azureipam/latest/docs#special-considerations))



## 2023.02.20 - v1.0.0

### Added
+ Added `description` in resource Reservation
+ Added `reverse_search` and `smallest_cidr` in resource Reservation
+ Added new fields in data Reservations: `description`, `created_by`, `settled_on`, `settled_by`S

### Fixed
+ Reservation is not longer deleted, removed error prevention when not found in redeployments. Also removed `Special Considerations` section in documentation.


## 2024.09.15 - v2.0.0
### Added
+ migration from [SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2) to [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
+ templates for [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs), to allow Terraform plugin doc generation
+ acceptance tests.
+ data resource `azureipam_spaces` to get a list of all spaces with related information.
+ resource `azureipam_space` to allow operations with spaces.
+ resource `azureipam_block` to allow operations with blocks.
+ data resource `azureipam_blocks` to get a list of all blocks in the specified space with related information.
+ resource `azureipam_external` to allow to associate an external network to the target space and block.
+ data resource `azureipam_blocks` to get a list of all external networks associated with a space and block. 
+ resource `azureipam_reservation_cidr` to allow to create reservation specifying the cidr, space and block.
+ data resource `azureipam_external` to get a specific external network associated with a space and block. 
+ data resource `azureipam_block` to get a specific block by space and name with related information.
+ data resource `azureipam_space` to get a specific space by name with related information.
+ data resource `azureipam_reservation` to get a specific reservation by id in the specified space and block.
+ resource `azureipam_block_network` that allow to associate an existing azure network to the target block.
+ data resource `azureipam_block_networks_availables` to get the list of azure virtual networks availables to be associated with a space and block.
+ data resource `azureipam_block_networks` to get the list of networks actively associated to a space block.

### Modified (Breaking Change)
+ resource `azureipam_reservation` now allow to specify a block list. The list is evaluated in the order provider
+ New provider attribute `skip_cert_verification` allow to specify if tls certificate verification must be skipped in API calls. Can be helpful in development environments. In previous versions it was always omitted by default. 

### Fixed
+ [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) implementation and generation