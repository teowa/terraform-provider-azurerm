---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_connectivity_configuration"
description: |-
  Manages a Network Connectivity Configurations.
---

# azurerm_network_connectivity_configuration

Manages a Network Connectivity Configurations.

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

resource "azurerm_network_connectivity_configuration" "example" {
  name                       = "example-ncc"
  network_manager_id = azurerm_network_manager.test.id
  connectivity_topology      = ""
  delete_existing_peering    = ""
  description                = ""
  is_global                  = ""
  applies_to_groups {
    group_connectivity = ""
    is_global          = ""
    network_group_id   = ""
    use_hub_gateway    = ""
  }
  hubs {
    resource_id   = ""
    resource_type = ""
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Connectivity Configurations. Changing this forces a new Network Connectivity Configurations to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Connectivity Configurations. Changing this forces a new Network Connectivity Configurations to be created.

* `applies_to_groups` - (Required) An `applies_to_groups` block as defined below.

* `connectivity_topology` - (Required) Connectivity topology type.

* `delete_existing_peering` - (Optional) Flag if need to remove current existing peerings.

* `description` - (Optional) A description of the connectivity configuration.

* `hubs` - (Optional) A `hubs` block as defined below.

* `is_global` - (Optional) Flag if global mesh is supported.

---

An `applies_to_groups` block supports the following:

* `group_connectivity` - (Required) Group connectivity type.

* `is_global` - (Optional) Flag if global is supported.

* `network_group_id` - (Required) Network group Id.

* `use_hub_gateway` - (Optional) Flag if need to use hub gateway.

---

A `hubs` block supports the following:

* `resource_id` - (Optional) Resource Id.

* `resource_type` - (Optional) Resource Type.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Connectivity Configurations.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Connectivity Configurations.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Connectivity Configurations.
* `update` - (Defaults to 30 minutes) Used when updating the Network Connectivity Configurations.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Connectivity Configurations.

## Import

Network Connectivity Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_connectivity_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/connectivityConfigurations/configuration1
```
