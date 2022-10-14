---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_subscription_network_manager_connection"
description: |-
  Manages a Network Subscription Network Manager Connections.
---

# azurerm_network_subscription_network_manager_connection

Manages a Network Subscription Network Manager Connections.

## Example Usage

```hcl
resource "azurerm_network_subscription_network_manager_connection" "example" {
  name               = "example-nsnmc"
  connection_state   = ""
  description        = ""
  network_manager_id = ""

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Subscription Network Manager Connections. Changing this forces a new Network Subscription Network Manager Connections to be created.

* `connection_state` - (Optional) .

* `description` - (Optional) A description of the network manager connection.

* `network_manager_id` - (Optional) Network Manager Id.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Subscription Network Manager Connections.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Subscription Network Manager Connections.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Subscription Network Manager Connections.
* `update` - (Defaults to 30 minutes) Used when updating the Network Subscription Network Manager Connections.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Subscription Network Manager Connections.

## Import

Network Subscription Network Manager Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_subscription_network_manager_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/networkManagerConnections/networkManagerConnection1
```
