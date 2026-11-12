---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_edge_action_execution_filter"
description: |-
  Lists CDN Edge Action Execution Filter resources.
---

# List resource: azurerm_cdn_edge_action_execution_filter

Lists CDN Edge Action Execution Filter resources.

## Example Usage

```hcl
list "azurerm_cdn_edge_action_execution_filter" "example" {
  provider = azurerm
  config {
    cdn_edge_action_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Cdn/edgeActions/examplecdnedgeaction"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cdn_edge_action_id` - (Required) The ID of the CDN Edge Action to query.
