---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_edge_action_version"
description: |-
  Manages a CDN Edge Action Version.
---

# azurerm_cdn_edge_action_version

Manages a CDN Edge Action Version.

## Example Usage

### Deploy code and promote a version to default

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_cdn_edge_action" "example" {
  name                = "examplecdnedgeaction"
  resource_group_name = azurerm_resource_group.example.name
  location            = "global"
  sku_name            = "Standard"
  sku_tier            = "Standard"
}

resource "azurerm_cdn_edge_action_version" "example" {
  name                    = "examplecdnedgeactionversion"
  cdn_edge_action_id      = azurerm_cdn_edge_action.example.id
  location                = "global"
  deployment_type         = "file"
  default_version_enabled = true

  code {
    file_name      = "index.js"
    content_base64 = base64encode("addEventListener(\"fetch\", event => event.respondWith(new Response(\"hello\")));")
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this CDN Edge Action Version. Must contain between 1 and 50 alphanumeric characters. Changing this forces a new resource to be created.

* `cdn_edge_action_id` - (Required) The ID of the CDN Edge Action where this Version should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where this CDN Edge Action Version should exist. Changing this forces a new resource to be created.

* `deployment_type` - (Required) The deployment type for this CDN Edge Action Version. Possible values are `file`, `others`, and `zip`. Changing this forces a new resource to be created.

* `code` - (Optional) A `code` block as defined below.

~> **Note:** Changing a `code` block deploys the new content to the existing Version. Removing the block forces a new CDN Edge Action Version to be created because the API cannot clear deployed code. Code is not read back during refresh or import, so out-of-band code changes cannot be detected.

* `default_version_enabled` - (Optional) Whether to promote this Version to the default. When configured, this must be `true`. If omitted, the current default status is read from Azure.

~> **Note:** Azure automatically makes the first Version the default. Configure `default_version_enabled = true` on only the Version that should be promoted, leaving it unset on other Versions. To change the default, remove this argument from the previous default's configuration and set it to `true` on the target Version. Promotion uses the default-version swap operation after any configured code deployment.

* `tags` - (Optional) A mapping of tags which should be assigned to this CDN Edge Action Version.

---

A `code` block supports the following:

* `content_base64` - (Required) The Base64-encoded content to deploy.

~> **Note:** This value is marked sensitive but is stored in Terraform state. Secure access to the state file.

* `file_name` - (Required) The name of the code file or deployment package.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CDN Edge Action Version.

* `validation_status` - The validation status of the deployed code.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CDN Edge Action Version, deploying code, and promoting the default Version.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Edge Action Version.
* `update` - (Defaults to 30 minutes) Used when updating the CDN Edge Action Version, deploying code, and promoting the default Version.
* `delete` - (Defaults to 30 minutes) Used when deleting the CDN Edge Action Version.

## Import

A CDN Edge Action Version can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_edge_action_version.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/edgeActions/edgeAction1/versions/version1
```

-> **Note:** The `code` block is not populated during import. Adding it to an imported Version deploys the configured content.
