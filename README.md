# Terraform Provider AzureIPAM

This provider is intended to manage the reservation of network ranges in the [Azure IPAM](https://github.com/Azure/ipam) solution. IPAM solution is a simple, straightforward way to manage IP address spaces in Azure, and it's's required to have a previous implementation of this solution.

The provider makes use of the IPAM REST API to manage CIDR range reservations in a space and block from those configured in the application.

> **NOTE** that this provider makes use of a functionality not implemented in the open-source version of the Azure IPAM solution, which allows to include additional tags in the creation of new reservations (used for description). 

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-azureipam
```
or

```shell
$ make build
```

## Local release build

```shell
$ go install github.com/goreleaser/goreleaser/v2@latest
```

```shell
$ make release
```

You will find the releases in the `/dist` directory. You will need to rename the provider binary to `terraform-provider-azureipam` before use it.
To run locally you can proceed in the following ways:  

- Create a [Terraform CLI Configuration File with Development Overrides](https://developer.hashicorp.com/terraform/plugin/debugging#terraform-cli-development-overrides) that includes a `provider_installation` block with a `dev_overrides` block, specifiyng the path where your local binary is created.
- Copy the binary file into one of the [implied configuration `filesystem_mirror` folder](https://developer.hashicorp.com/terraform/cli/config/config-file#implied-local-mirror-directories) after each build.


## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory. 

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```