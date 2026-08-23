---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_discovery_workspace"
description: |-
  Manages a Discovery Workspace.
---

# azurerm_discovery_workspace

Manages a Discovery Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-discovery-rg"
  location = "West Europe"
}

resource "azurerm_discovery_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

* `location` - (Required) The Azure Region where the Discovery Workspace should exist. Changing this forces a new resource to be created.
* `name` - (Required) The name of the Discovery Workspace. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the Resource Group where the Discovery Workspace should exist. Changing this forces a new resource to be created.
* `tags` - (Optional) A mapping of tags which should be assigned to the Discovery Workspace.

## Attributes Reference

* `id` - The ID of the Discovery Workspace.

## Timeouts

The `timeouts` block supports `create`, `read`, `update`, and `delete`.

## Import

Discovery Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_discovery_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Discovery/workspaces/example
```
