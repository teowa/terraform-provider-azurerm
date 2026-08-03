---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_cloud_hosted_search"
description: |-
  Lists Elastic Cloud hosted search deployments.
---

# List resource: azurerm_elastic_cloud_hosted_search

Lists Elastic Cloud hosted search deployments.

## Example Usage

### List all Elastic Cloud hosted search deployments in the subscription

```hcl
list "azurerm_elastic_cloud_hosted_search" "example" {
  provider = azurerm
  config {}
}
```

### List all Elastic Cloud hosted search deployments in a specific resource group

```hcl
list "azurerm_elastic_cloud_hosted_search" "example" {
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
