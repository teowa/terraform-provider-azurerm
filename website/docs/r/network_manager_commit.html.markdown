---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_commit"
description: |-
  Manages a Network Manager Commit.
---

# azurerm_network_manager_commit

Manages a Network Manager Commit.

-> **Note:** The Azure Provider include a Feature Toggle `manager_replace_committed` to control whether to enable the replacement mode of the commit, the default is `false`. If `manager_replace_committed` is set to `true`, the deployed resource will not be cleaned when the `azurerm_network_manager_commit` is removed from the config, and provisioning a new resource will not check if the remote existence. This is designed to avoid downtime when use [`replace_triggered_by`](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#replace_triggered_by). If `manager_replace_committed` is set to `false`, the deployed resource will be cleaned when the config is removed. See [the Features block documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#features) for more information on Feature Toggles within Terraform.

## Example Usage

```hcl
provider "azurerm" {
  network {
    manager_replace_committed = false
  }
}

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

resource "azurerm_network_manager_commit" "example" {
  network_manager_id = azurerm_network_manager.example.id
  location           = "eastus"
  scope_access       = "Connectivity"
  configuration_ids  = [azurerm_network_manager_connectivity_configuration.example.id]
  lifecycle {
    replace_triggered_by = [
      # Replace `azurerm_network_manager_commit` each time 
      # the `azurerm_network_manager_connectivity_configuration` is updated.
      azurerm_network_manager_connectivity_configuration.example
    ]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `network_manager_id` - (Required) Specifies the ID of the Network Manager. Changing this forces a new Network Manager Commit to be created.

* `location` - (Required) Specifies the location which the configurations will be deployed to. Changing this forces a new Network Manager Commit to be created.

* `scope_access` - (Required) Specifies the configuration deployment type. Possible values are `Connectivity` and `SecurityAdmin`.

* `configuration_ids` - (Required) A list of Network Manager Configuration IDs. If an empty list is passed, it means to clean all the commit.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Admin Rule Collection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Commit.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Commit.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Commit.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Commit.

## Import

Network Manager Commit can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_commit.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/commit|eastus|Connectivity
```
