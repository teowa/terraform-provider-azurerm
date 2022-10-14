---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_admin_rule_collection"
description: |-
  Manages a Network Admin Rule Collections.
---

# azurerm_network_admin_rule_collection

Manages a Network Admin Rule Collections.

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
}

resource "azurerm_network_admin_rule_collection" "example" {
  name                                    = "example-narc"
  network_security_admin_configuration_id = azurerm_network_security_admin_configuration.test.id
  description                             = ""
  applies_to_groups {
    network_group_id = ""
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Admin Rule Collections. Changing this forces a new Network Admin Rule Collections to be created.

* `network_security_admin_configuration_id` - (Required) Specifies the ID of the Network Admin Rule Collections. Changing this forces a new Network Admin Rule Collections to be created.

* `applies_to_groups` - (Required) An `applies_to_groups` block as defined below.

* `description` - (Optional) A description of the admin rule collection.

---

An `applies_to_groups` block supports the following:

* `network_group_id` - (Required) Network manager group Id.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Admin Rule Collections.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Admin Rule Collections.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Admin Rule Collections.
* `update` - (Defaults to 30 minutes) Used when updating the Network Admin Rule Collections.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Admin Rule Collections.

## Import

Network Admin Rule Collections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_admin_rule_collection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/securityAdminConfigurations/configuration1/ruleCollections/ruleCollection1
```
