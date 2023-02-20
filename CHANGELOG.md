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