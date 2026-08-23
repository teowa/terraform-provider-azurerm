---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_discovery_node_pool"
description: |-
  Manages a Discovery Node Pool.
---

# azurerm_discovery_node_pool

Manages a Discovery Node Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-discovery-rg"
  location = "West Europe"
}

resource "azurerm_discovery_supercomputer" "example" {
  name                = "example-supercomputer"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_discovery_node_pool" "example" {
  name                = "example-node-pool"
  discovery_supercomputer_id = azurerm_discovery_supercomputer.example.id
  location                   = azurerm_resource_group.example.location
}
```

## Arguments Reference

* `location` - (Required) The Azure Region where the Discovery Node Pool should exist. Changing this forces a new resource to be created.
* `name` - (Required) The name of the Discovery Node Pool. Changing this forces a new resource to be created.
* `discovery_supercomputer_id` - (Required) The ID of the Discovery Supercomputer where this Node Pool should exist. Changing this forces a new resource to be created.
* `tags` - (Optional) A mapping of tags which should be assigned to the Discovery Node Pool.

## Attributes Reference

* `id` - The ID of the Discovery Node Pool.

## Timeouts

The `timeouts` block supports `create`, `read`, `update`, and `delete`.

## Import

Discovery Node Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_discovery_node_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Discovery/supercomputers/example/nodePools/example
```
