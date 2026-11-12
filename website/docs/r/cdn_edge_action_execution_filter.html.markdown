---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_edge_action_execution_filter"
description: |-
  Manages a CDN Edge Action Execution Filter.
---

# azurerm_cdn_edge_action_execution_filter

Manages a CDN Edge Action Execution Filter.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_cdn_edge_action" "example" {
  name                = "examplecdnedgeaction"
  resource_group_name = azurerm_resource_group.example.name
  location            = "global"
  sku_name            = "Standard"
  sku_tier            = "Standard"
}

resource "azurerm_cdn_edge_action_version" "example" {
  name               = "examplecdnedgeactionversion"
  cdn_edge_action_id = azurerm_cdn_edge_action.example.id
  location           = "global"
  deployment_type    = "file"
}

resource "azurerm_cdn_edge_action_execution_filter" "example" {
  name                       = "examplecdnedgeactionexecutionfilter"
  cdn_edge_action_id         = azurerm_cdn_edge_action.example.id
  cdn_edge_action_version_id = azurerm_cdn_edge_action_version.example.id
  location                   = "global"
  header_name                = "X-Edge-Action"
  header_value               = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this CDN Edge Action Execution Filter. Must contain between 1 and 50 alphanumeric characters. Changing this forces a new resource to be created.

* `cdn_edge_action_id` - (Required) The ID of the CDN Edge Action where this Execution Filter should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where this CDN Edge Action Execution Filter should exist. Changing this forces a new resource to be created.

* `cdn_edge_action_version_id` - (Required) The ID of the CDN Edge Action Version referenced by this Execution Filter.

* `header_name` - (Required) The name of the HTTP header identifying this Execution Filter.

* `header_value` - (Required) The value of the HTTP header identifying this Execution Filter.

* `tags` - (Optional) A mapping of tags which should be assigned to this CDN Edge Action Execution Filter.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CDN Edge Action Execution Filter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CDN Edge Action Execution Filter.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Edge Action Execution Filter.
* `update` - (Defaults to 30 minutes) Used when updating the CDN Edge Action Execution Filter.
* `delete` - (Defaults to 30 minutes) Used when deleting the CDN Edge Action Execution Filter.

## Import

A CDN Edge Action Execution Filter can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_edge_action_execution_filter.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/edgeActions/edgeAction1/executionFilters/executionFilter1
```
