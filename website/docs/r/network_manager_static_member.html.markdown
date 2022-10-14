---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_static_member"
description: |-
  Manages a Network Static Members.
---

# azurerm_network_static_member

Manages a Network Static Members.

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

resource "azurerm_network_group" "example" {
  name                       = "example-nng"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_static_member" "example" {
  name                     = "example-nsm"
  network_group_id = azurerm_network_group.test.id
  resource_id              = ""

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Static Members. Changing this forces a new Network Static Members to be created.

* `network_group_id` - (Required) Specifies the ID of the Network Static Members. Changing this forces a new Network Static Members to be created.

* `resource_id` - (Optional) Resource Id.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Static Members.

* `region` - Resource region.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Static Members.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Static Members.
* `update` - (Defaults to 30 minutes) Used when updating the Network Static Members.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Static Members.

## Import

Network Static Members can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_static_member.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/networkGroups/networkGroup1/staticMembers/staticMember1
```
