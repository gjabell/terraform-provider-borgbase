# Terraform BorgBase Provider

This is a Terraform provider for [BorgBase](https://www.borgbase.com). It can currently be used to manage SSH keys and Borg repositories. Managing Restic repositories will be supported in a future version.

This provider is available on the [Terraform registry](https://registry.terraform.io/providers/gjabell/borgbase).

## Getting started

You will need an API key for your BorgBase account. You can create one on the [Account page](https://www.borgbase.com/account) under the "API" tab.

Once you have an API key, create a `main.tf` with the following settings:

```hcl
terraform {
  required_providers {
    borgbase = {
      source  = "gjabell/borgbase"
      version = "0.1.0"
    }
  }
}

provider "borgbase" {
	api_token = "your token here"
}
```

You can also set the token via the `BORGBASE_API_TOKEN` environment variable.

Now run `terraform init` to initialize the Terraform project and provider.

### Creating an SSH key

Create a new SSH key resource:

```hcl
resource "borgbase_ssh_key" "test" {
	name       = "test"
	public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBAt/X37WDQ3cNPEVHQBsW3lH7XPeea5rUoeXuhoTkzR user@hostname"
}
```

And add it your account by running `terraform apply`.

You can also import an existing key using:

```shell
$ terraform import borgbase_ssh_key.test "key name in borgbase"
```

### Creating a Borg repository

Create a new SSH key data source and Borg repository resource:

```hcl
data "borgbase_ssh_key" "test" {
	name = "test"
}

resource "borgbase_borg_repo" "test" {
	name             = "test"
	full_access_keys = [data.borgbase_ssh_key.test.id]
	region           = "eu"
}
```

And add it your account by running `terraform apply`.

You can also import an existing repository using:

```shell
$ terraform import borgbase_borg_repo.test "repository name in borgbase"
```

See the [documentation](https://registry.terraform.io/providers/gjabell/borgbase/latest/docs/resources/borg_repo#schema) for a full list of available repository options.

## Contributing

Clone the project and build:

```shell
$ git clone https://github.com/gjabell/terraform-provider-borgbase
$ cd terraform-provider-borgbase
$ make
```

Running the acceptance tests requires a BorgBase API key, as the tests interact with the actual BorgBase backend. **Please double check the test resource names to ensure they don't conflict with resources in your account!** Otherwise you may end up losing data.

```shell
$ export BORGBASE_API_TOKEN="your token here"
$ make testacc
```

To build the docs after changing one of the examples, run `make docs`.
