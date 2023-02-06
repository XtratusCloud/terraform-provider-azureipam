# Terraform Provider AzureIPAM

This provider is intended to manage the reservation of network ranges in the [Azure IPAM](https://github.com/Azure/ipam) solution. IPAM solution is a simple, straightforward way to manage IP address spaces in Azure, and it's's required to have a previous implementation of this solution.

The provider makes use of the IPAM REST API to manage CIDR range reservations in a space and block from those configured in the application.

> **NOTE** that this provider makes use of a functionality not implemented in the open-source version of the Azure IPAM solution, which allows to include additional tags in the creation of new reservations (used for description). 

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-azureipam
```

## Local release build

```shell
$ go install github.com/goreleaser/goreleaser@latest
```

```shell
$ make release
```

You will find the releases in the `/dist` directory. You will need to rename the provider binary to `terraform-provider-azureipam` and move the binary into [the appropriate subdirectory within the user plugins directory](https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers#install-hashicups-provider).


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