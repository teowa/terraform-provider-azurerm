---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_group"
description: |-
  Manages a Network Groups.
---

# azurerm_network_group

Manages a Network Groups.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_manager" "example" {
  name                = "example-nnm"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_group" "example" {
  name                       = "example-nng"
  network_manager_id = azurerm_network_manager.test.id
  description                = ""

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Groups. Changing this forces a new Network Groups to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Groups. Changing this forces a new Network Groups to be created.

* `description` - (Optional) A description of the network group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Groups.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Groups.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Groups.
* `update` - (Defaults to 30 minutes) Used when updating the Network Groups.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Groups.

## Import

Network Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/networkGroups/networkGroup1
```
