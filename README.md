Terraform Provider for Palo Alto Networks SASE API
==================================================

**NOTE**:  This provider will eventually be entirely auto-generated.


Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) v1+
- [Go](https://go.dev) v1.18+ (to build the provider from source)


Building the Provider
---------------------

1. Install [Go](https://go.dev/dl)

2. Clone the SDK repo:

```sh
git clone https://github.com/paloaltonetworks/sase-go
```

3. Clone this repo:

```sh
git clone https://github.com/paloaltonetworks/terraform-provider-sase
```

4. Build the provider:

```sh
cd terraform-provider-sase
go build
```

5. 4. Specify the `dev_overrides` configuration per the next section below. This tells Terraform where to find the provider you just built. The directory to specify is the full path to the cloned provider repo.

When using the provider, refer to the documentation in the `./docs` directory for all resources and parameters.


Developing the Provider
-----------------------

With Terraform v1 and later, [development overrides for provider developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers) can be leveraged in order to use the provider built from source.

To do this, populate a Terraform CLI configuration file (`~/.terraformrc` for all platforms other than Windows; `terraform.rc` in the `%APPDATA%` directory when using Windows) with at least the following options:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/paloaltonetworks-local/sase" = "/directory/containing/the/provider/binary/here"
  }

  direct {}
}
```

Then when referencing the locally built provider, use the local name in the provider block like so:

```hcl
terraform {
    required_providers {
        sase = {
            source = "paloaltonetworks-local/sase"
            version = "1.0.0"
        }
    }
}
```
