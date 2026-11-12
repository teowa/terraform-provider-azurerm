---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_edge_action"
description: |-
  Manages a CDN Edge Action.
---

# azurerm_cdn_edge_action

Manages a CDN Edge Action.

## Example Usage

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this CDN Edge Action. Must contain between 1 and 50 alphanumeric characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where this CDN Edge Action should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where this CDN Edge Action should exist. Changing this forces a new resource to be created.

* `sku_name` - (Required) The name of the SKU for this CDN Edge Action. Changing this forces a new resource to be created.

* `sku_tier` - (Required) The tier of the SKU for this CDN Edge Action. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to this CDN Edge Action.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CDN Edge Action.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CDN Edge Action.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Edge Action.
* `update` - (Defaults to 30 minutes) Used when updating the CDN Edge Action.
* `delete` - (Defaults to 30 minutes) Used when deleting the CDN Edge Action.

## Import

A CDN Edge Action can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_edge_action.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/edgeActions/edgeAction1
```
