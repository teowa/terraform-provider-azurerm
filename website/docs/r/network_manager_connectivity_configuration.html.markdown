---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_connectivity_configuration"
description: |-
  Manages a Network Manager Connectivity Configuration.
---

# azurerm_network_manager_connectivity_configuration

Manages a Network Manager Connectivity Configuration.

-> **Note:** The `azurerm_network_manager_connectivity_configuration` may delete or modify existing Network Peering.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin"]
  description    = "example network manager"
}

resource "azurerm_network_manager_network_group" "example" {
  name               = "example-group"
  network_manager_id = azurerm_network_manager.example.id
}

resource "azurerm_virtual_network" "example" {
  name                    = "example-net"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  address_space           = ["10.0.0.0/16"]
  flow_timeout_in_minutes = 10
}

resource "azurerm_network_manager_connectivity_configuration" "example" {
  name                  = "example-connectivity-conf"
  network_manager_id    = azurerm_network_manager.example.id
  connectivity_topology = "HubAndSpoke"
  applies_to_group {
    group_connectivity = "None"
    network_group_id   = azurerm_network_manager_network_group.example.id
  }
  hub {
    resource_id   = azurerm_virtual_network.example.id
    resource_type = "Microsoft.Network/virtualNetworks"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager Connectivity Configuration. Changing this forces a new Network Manager Connectivity Configuration to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Manager. Changing this forces a new Network Manager Connectivity Configuration to be created.

* `applies_to_group` - (Required) An `applies_to_group` block as defined below.

* `connectivity_topology` - (Required) Connectivity topology type. Possible values are `HubAndSpoke` and `Mesh`.

* `delete_existing_peering` - (Optional) Boolean flag whether to remove current existing Virtual Network Peerings.

* `description` - (Optional) A description of the Connectivity Configuration.

* `hubs` - (Optional) A `hubs` block as defined below.

* `is_global` - (Optional) Boolean flag if global mesh is supported.

---

An `applies_to_group` block supports the following:

* `group_connectivity` - (Required) Group connectivity type. Possible values are `None` and `DirectlyConnected`.

* `is_global` - (Optional) Boolean flag if global is supported.

* `network_group_id` - (Required) Network Group ID.

* `use_hub_gateway` - (Optional) Boolean flag if hub gateway is used.

---

A `hubs` block supports the following:

* `resource_id` - (Optional) Resource ID used as hub in Hub And Spoke topology.

* `resource_type` - (Optional) Resource Type used as hub in Hub And Spoke topology.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Connectivity Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Connectivity Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Connectivity Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Connectivity Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Connectivity Configuration.

## Import

Network Manager Connectivity Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_connectivity_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/connectivityConfigurations/configuration1
```
