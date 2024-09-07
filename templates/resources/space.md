---
page_title: "{{.Name}} {{.Type}} - {{.ProviderName}}"
subcategory: ""
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.Name}} ({{.Type}})

{{ .Description | trimspace }}

## Example Usage

{{ tffile (printf "examples/resources/%s/resource.tf" .Name)}}

{{ .SchemaMarkdown | trimspace }}

## Import

Reservations can be imported using the ID of the IPAM space, e.g.

```shell
terraform import azureipam_space.new asia
```
