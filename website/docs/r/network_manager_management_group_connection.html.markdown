---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_management_group_connection"
description: |-
  Manages a Network Manager Management Group Connections.
---

# azurerm_network_manager_management_group_connection

Manages a Network Manager Management Group Connections.

## Example Usage

```hcl
resource "azurerm_network_manager_management_group_connection" "example" {
  name                = "example-nmgnmc"
  management_group_id = ""
  connection_state    = ""
  description         = ""
  network_manager_id  = ""

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager Management Group Connections. Changing this forces a new Network Manager Management Group Connections to be created.

* `management_group_id` - (Required) Specifies the Management Group ID. Changing this forces a new Network Manager Management Group Connections to be created.

* `connection_state` - (Optional) .

* `description` - (Optional) A description of the network manager connection.

* `network_manager_id` - (Optional) Network Manager ID.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Management Group Connections.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Management Group Connections.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Management Group Connections.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Management Group Connections.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Management Group Connections.

## Import

Network Manager Management Group Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_management_group_connection.example /providers/Microsoft.Management/managementGroups/{managementGroupId}/providers/Microsoft.Network/networkManagerConnections/networkManagerConnection1
```
