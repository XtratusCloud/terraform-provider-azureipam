# Terraform Provider AzureIPAM

This provider is intended to manage the reservation of network ranges in the [Azure IPAM](https://github.com/Azure/ipam) solution. IPAM solution is a simple, straightforward way to manage IP address spaces in Azure, and it's's required to have a previous implementation of this solution.

The provider makes use of the IPAM REST API to manage CIDR range reservations in a space and block from those configured in the application.

## Build provider

Run the following command to build the provider
```shell
$ make build
```

## Acceptance tests

To locally validate the implemented acceptance tests, simply run

```shell
$ make testacc
```

## Local release build
For the release creation process [goreleaser](https://goreleaser.com/) v2 or later is used, so it has to be previously installed.

```shell
$ go install github.com/goreleaser/goreleaser/v2@latest
```

And to run the release process locally, simply run
```shell
$ make release
```

You will find the releases in the `/dist` directory. Probably you will need to rename the provider binary to `terraform-provider-azureipam` before use it.

To run locally you can proceed in one of the following ways:

- Create a [Terraform CLI Configuration File with Development Overrides](https://developer.hashicorp.com/terraform/plugin/debugging#terraform-cli-development-overrides) that includes a `provider_installation` block with a `dev_overrides` block, specifiyng the path where your local binary is created.

- Copy the binary file into one of the [implied configuration `filesystem_mirror` folder](https://developer.hashicorp.com/terraform/cli/config/config-file#implied-local-mirror-directories) after each build.


## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to a specific folder inside `tests` directory. 

```shell
$ cd tests/reservation_resource
```

Remember to configure the provider with your environment information
```shell
export AZUREIPAM_TOKEN="eyJ0eXAi......"
export AZUREIPAM_API_URL="https://myazureipam.azurewebsites.net"
```

And initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```