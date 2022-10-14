---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_admin_rule"
description: |-
  Manages a Network Admin Rules.
---

# azurerm_network_admin_rule

Manages a Network Admin Rules.

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
}

resource "azurerm_network_admin_rule" "example" {
  name                             = "example-nar"
  network_admin_rule_collection_id = azurerm_network_admin_rule_collection.test.id
  kind                             = ""

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Admin Rules. Changing this forces a new Network Admin Rules to be created.

* `network_rule_collection_id` - (Required) Specifies the ID of the Network Admin Rules. Changing this forces a new Network Admin Rules to be created.

* `kind` - (Required) Whether the rule is custom or default.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Admin Rules.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Admin Rules.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Admin Rules.
* `update` - (Defaults to 30 minutes) Used when updating the Network Admin Rules.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Admin Rules.

## Import

Network Admin Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_admin_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/securityAdminConfigurations/configuration1/ruleCollections/ruleCollection1/rules/rule1
```
