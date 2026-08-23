---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_discovery_storage_asset"
description: |-
  Manages a Discovery Storage Asset.
---

# azurerm_discovery_storage_asset

Manages a Discovery Storage Asset.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-discovery-rg"
  location = "West Europe"
}

resource "azurerm_discovery_storage_container" "example" {
  name                = "example-storage-container"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_discovery_storage_asset" "example" {
  name                = "example-storage-asset"
  discovery_storage_container_id = azurerm_discovery_storage_container.example.id
  location                        = azurerm_resource_group.example.location
}
```

## Arguments Reference

* `location` - (Required) The Azure Region where the Discovery Storage Asset should exist. Changing this forces a new resource to be created.
* `name` - (Required) The name of the Discovery Storage Asset. Changing this forces a new resource to be created.
* `discovery_storage_container_id` - (Required) The ID of the Discovery Storage Container where this Storage Asset should exist. Changing this forces a new resource to be created.
* `tags` - (Optional) A mapping of tags which should be assigned to the Discovery Storage Asset.

## Attributes Reference

* `id` - The ID of the Discovery Storage Asset.

## Timeouts

The `timeouts` block supports `create`, `read`, `update`, and `delete`.

## Import

Discovery Storage Assets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_discovery_storage_asset.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Discovery/storageContainers/example/storageAssets/example
```
