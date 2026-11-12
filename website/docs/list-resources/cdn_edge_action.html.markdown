---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_edge_action"
description: |-
  Lists CDN Edge Action resources.
---

# List resource: azurerm_cdn_edge_action

Lists CDN Edge Action resources.

## Example Usage

```hcl
list "azurerm_cdn_edge_action" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Required) The name of the Resource Group to query.
