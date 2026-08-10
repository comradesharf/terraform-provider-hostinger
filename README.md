# Hostinger Terraform Provider

The [Hostinger Terraform provider](https://registry.terraform.io/providers/hostinger/hostinger/latest/docs) provides convenient access to the [Hostinger API](https://developers.hostinger.com/) from Terraform.

## Requirements

This provider requires Terraform CLI 1.0 or later. You can [install Terraform for your system](https://developer.hashicorp.com/terraform/install) from HashiCorp.

## Usage

Add the following to your `main.tf` file:

```hcl
terraform {
  required_providers {
    hostinger = {
      source  = "hostinger/hostinger"
      version = "~> 0.1"
    }
  }
}

provider "hostinger" {
  # The provider reads HOSTINGER_API_TOKEN from the environment.
}

resource "hostinger_vps_public_key" "example" {
  name = "example-key"
  key  = file("~/.ssh/id_ed25519.pub")
}
```

The API token can be created in the [Hostinger hPanel](https://hpanel.hostinger.com/). For security, use the `HOSTINGER_API_TOKEN` environment variable instead of storing the token in Terraform configuration:

```shell
export HOSTINGER_API_TOKEN="your-api-token"
```

Initialize your project by running `terraform init` in the directory.

Additional examples can be found in the [examples](./examples) directory. Full resource, data source, list resource, and action documentation is available in the [Terraform Registry](https://registry.terraform.io/providers/hostinger/hostinger/latest/docs) and the generated [provider documentation](./docs/index.md).

### Provider Options

The following provider options are supported. It is recommended to use environment variables for sensitive values such as API tokens. When an environment variable is provided, the corresponding option does not need to be set in Terraform configuration.

| Property | Environment variable | Required | Default value |
| --- | --- | --- | --- |
| `api_token` | `HOSTINGER_API_TOKEN` | true | — |
| `host` | `HOSTINGER_HOST` | false | Production API server |

The `host` option is intended for custom API endpoints and testing. It should generally be left unset for normal use.

## Development

To build the provider locally, make sure [Go](https://go.dev/doc/install) 1.25 or later is installed, then run:

```shell
go install
```

To generate or update the provider documentation, run:

```shell
make generate
```

To run the unit test suite, run:

```shell
make test
```

Acceptance tests use a real or test Hostinger API endpoint and may create billable resources. Set `HOSTINGER_HOST` and `HOSTINGER_API_TOKEN` before running them:

```shell
make testacc
```

## Contributing

Bug reports, feature requests, and pull requests are welcome. Please open an [issue](https://github.com/hostinger/terraform-provider-hostinger/issues) or submit a pull request on GitHub.
