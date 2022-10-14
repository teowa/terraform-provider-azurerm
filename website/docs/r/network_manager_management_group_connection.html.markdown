---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_management_group_network_manager_connection"
description: |-
  Manages a Network Management Group Network Manager Connections.
---

# azurerm_network_management_group_network_manager_connection

Manages a Network Management Group Network Manager Connections.

## Example Usage

```hcl
resource "azurerm_network_management_group_network_manager_connection" "example" {
  name                = "example-nmgnmc"
  management_group_id = ""
  connection_state    = ""
  description         = ""
  network_manager_id  = ""

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Management Group Network Manager Connections. Changing this forces a new Network Management Group Network Manager Connections to be created.

* `management_group_id` - (Required) Specifies the Management Group Id. Changing this forces a new Network Management Group Network Manager Connections to be created.

* `connection_state` - (Optional) .

* `description` - (Optional) A description of the network manager connection.

* `network_manager_id` - (Optional) Network Manager Id.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Management Group Network Manager Connections.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Management Group Network Manager Connections.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Management Group Network Manager Connections.
* `update` - (Defaults to 30 minutes) Used when updating the Network Management Group Network Manager Connections.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Management Group Network Manager Connections.

## Import

Network Management Group Network Manager Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_management_group_network_manager_connection.example /providers/Microsoft.Management/managementGroups/{managementGroupId}/providers/Microsoft.Network/networkManagerConnections/networkManagerConnection1
```
