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


## 2024.09.01 - v1.1.0
### Added
+ data resource `azureipam_spaces` to get a list of all Spaces.

### Fixed
+ None