---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_connection_account_managed_identity"
description: |-
  Lists Cognitive Services Account Connection with Account Managed Identity authentication resources.
---

# List resource: azurerm_cognitive_account_connection_account_managed_identity

Lists Cognitive Services Account Connection with Account Managed Identity authentication resources.

## Example Usage

### List all Cognitive Services Account Connection with Account Managed Identity authentication resources in a Cognitive Services Account

```hcl
list "azurerm_cognitive_account_connection_account_managed_identity" "example" {
  provider = azurerm
  config {
    cognitive_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.CognitiveServices/accounts/example-cognitive-account"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cognitive_account_id` - (Required) The full ID of an existing Cognitive Services Account.
