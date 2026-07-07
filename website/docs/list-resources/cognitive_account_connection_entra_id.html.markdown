---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_connection_entra_id"
description: |-
  Lists Cognitive Services Account Connection with Entra ID authentication resources.
---

# List resource: azurerm_cognitive_account_connection_entra_id

Lists Cognitive Services Account Connection with Entra ID authentication resources.

## Example Usage

### List all Cognitive Services Account Connection with Entra ID authentication resources in a Cognitive Services Account

```hcl
list "azurerm_cognitive_account_connection_entra_id" "example" {
  provider = azurerm
  config {
    cognitive_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.CognitiveServices/accounts/example-cognitive-account"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cognitive_account_id` - (Required) The full ID of an existing Cognitive Services Account.
