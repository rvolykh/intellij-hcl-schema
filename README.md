# intellij-hcl-schema **(UNOFFICIAL)**

> **DEPRECATED**:
>
> This approach does not work anymore (at least with lastest PyCharm).
> The right approach would be to run terraform init, so that .terraform.d folder is created
> inside your terraform root/module and then IDEA can obtain schemas on its own.

----

Helper to build terraform provider autocompletion for Intellij HCL plugin.
It starts the provider, gets it schema using GRPC, saves to $HOME/.terraform.d/metadata-repo/terraform/model/providers
(path from which Intellij HCL plugin loads custom providers schemas)

## Installation

Build version can be downloaded from releases page
https://github.com/rvolykh/intellij-hcl-schema/releases.

As an alternative, build yourself
```shell
make build
```
\* Requirements: Go

## Usage

```text
Example:
        intellij-hcl-schema -path ./terraform-provider-hashicups -name hashicups

Flags:
  -name string
        Name to use for provider
  -path string
        Path to already build terraform provider
  -ver string
        Version to use for provider (default "0.0.0")
```

## Development Notes

1. Terraform Provider Protocol updates

    In future new protocol version might be incompatible with tfplugin5.3, which is currently
    used for obtaining provider schema. In order to get latest schema support:

   1. Replace [tfplugin.proto](./proto/tfplugin.proto) with new protocol (https://github.com/hashicorp/terraform/tree/v1.3.9/docs/plugin-protocol).
   2. Download requirements: `make requirements`, potentially with changing version of protoc-gen-go-grpc (in Makefile).
   3. Change option go_package to `option go_package = "github.com/rvolykh/intellij-hcl-schema/proto";`.
   4. Generate code: `make generate`, might require changing source code.

## References

- https://github.com/VladRassokhin/terraform-metadata
- https://github.com/hashicorp/terraform/tree/main/docs/plugin-protocol
- https://gist.github.com/ArthurHlt/45d4207c557dd614735a4aabc3d12976
