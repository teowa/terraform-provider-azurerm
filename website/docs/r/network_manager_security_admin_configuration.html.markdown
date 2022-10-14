---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_admin_configuration"
description: |-
  Manages a Network Security Admin Configurations.
---

# azurerm_network_security_admin_configuration

Manages a Network Security Admin Configurations.

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

resource "azurerm_network_security_admin_configuration" "example" {
  name                       = "example-nsac"
  network_manager_id = azurerm_network_manager.test.id
  description                = ""
  apply_on_network_intent_policy_based_services {

  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Security Admin Configurations. Changing this forces a new Network Security Admin Configurations to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Security Admin Configurations. Changing this forces a new Network Security Admin Configurations to be created.

* `description` - (Optional) A description of the security configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Security Admin Configurations.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Security Admin Configurations.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Security Admin Configurations.
* `update` - (Defaults to 30 minutes) Used when updating the Network Security Admin Configurations.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Security Admin Configurations.

## Import

Network Security Admin Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_security_admin_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/securityAdminConfigurations/configuration1
```
