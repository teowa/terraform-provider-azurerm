---
subcategory: "Orbital"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_orbital_geocatalog"
description: |-
    Lists Orbital GeoCatalog resources.
---

# List resource: azurerm_orbital_geocatalog

Lists Orbital GeoCatalog resources.

## Example Usage

### List all Orbital GeoCatalogs in the subscription

```hcl
list "azurerm_orbital_geocatalog" "example" {
  provider = azurerm
  config {}
}
```

### List all Orbital GeoCatalogs in a specific resource group

```hcl
list "azurerm_orbital_geocatalog" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
