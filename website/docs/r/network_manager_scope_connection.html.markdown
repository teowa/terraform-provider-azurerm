---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_scope_connection"
description: |-
  Manages a Network Scope Connections.
---

# azurerm_network_scope_connection

Manages a Network Scope Connections.

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

resource "azurerm_network_scope_connection" "example" {
  name               = "example-nsc"
  network_manager_id = azurerm_network_manager.test.id
  connection_state   = ""
  description        = ""
  resource_id        = ""
  tenant_id          = ""

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Scope Connections. Changing this forces a new Network Scope Connections to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Scope Connections. Changing this forces a new Network Scope Connections to be created.

* `connection_state` - (Optional) Connection State.

* `description` - (Optional) A description of the scope connection.

* `resource_id` - (Optional) Resource ID.

* `tenant_id` - (Optional) Tenant ID.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Scope Connections.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Scope Connections.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Scope Connections.
* `update` - (Defaults to 30 minutes) Used when updating the Network Scope Connections.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Scope Connections.

## Import

Network Scope Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_scope_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/scopeConnections/scopeConnection1
```
