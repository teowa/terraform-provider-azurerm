---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_discovery_supercomputer"
description: |-
  Manages a Discovery Supercomputer.
---

# azurerm_discovery_supercomputer

Manages a Discovery Supercomputer.

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
```

## Arguments Reference

* `location` - (Required) The Azure Region where the Discovery Supercomputer should exist. Changing this forces a new resource to be created.
* `name` - (Required) The name of the Discovery Supercomputer. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the Resource Group where the Discovery Supercomputer should exist. Changing this forces a new resource to be created.
* `tags` - (Optional) A mapping of tags which should be assigned to the Discovery Supercomputer.

## Attributes Reference

* `id` - The ID of the Discovery Supercomputer.

## Timeouts

The `timeouts` block supports `create`, `read`, `update`, and `delete`.

## Import

Discovery Supercomputers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_discovery_supercomputer.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Discovery/supercomputers/example
```
