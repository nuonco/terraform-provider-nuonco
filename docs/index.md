---
page_title: "Provider: Nuon"
description: |-
    The Nuon provider allows you to interact with the Nuon API.
---

# Nuon Provider

Use the Nuon provider to configure resources in the [Nuon platform](https://www.nuon.co/). It must be configure with the correct credentials before you can use it.

To learn the basics of using this provider, follow the [Nuon quickstart](https://github.com/nuonco/quickstart).

Use the navigation to the left to read about available resources.

## Example Usage

```terraform
terraform {
  required_version = ">= 1.3.7"

  required_providers {
    nuon = {
      source  = "nuonco/nuon"
      version = "~> 0.2"
    }
  }
}

variable "api_token" {
  type = string
}
variable "org_id" {
  type = string
}

provider "nuon" {
  api_token = var.api_token
  org_id    = var.org_id
}

resource "nuon_app" "my_app" {
  name = "My App"
}
```

## Authentication and Configuration

Configuration can comes from 3 places. These are, in order of precedence:

1. Parameters in the provider configuration
1. Environment variables
1. Shared config file

### Provider Configuration

!> **Warning:** Hard-coded credentials are not recommended in any Terraform configuration and risks secret leakage should this file ever be committed to a public version control system.

Credentials can be provided by adding the attributes `org_id` and `api_token` to the provider block.

```terraform
provider "nuon" {
  org_id    = "my-org-id"
  api_token = "my-api-token"
}
```

## Environment Variables

Credentials can be provided via environment variables by setting `NUON_ORG_ID` and `NUON_API_TOKEN`.

```console
export NUON_ORG_ID="my-org-id"
export NUON_API_TOKEN="my-api-token"
```

## Shared Configuration File

The provider will also read credentials from the same config file used by the Nuon CLI, found at `~/.nuon`.

```yaml
org_id: "my-org-id"
api_token: "my-api-token"
```

If you would like to use a custom config file, you can set the environment variable `$NUON_CONFIG_FILE=<path>` and the provider will automatically use it.
