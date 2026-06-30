---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_storage_defender"
description: |-
    Manages the Defender for Storage. 
---

# azurerm_security_center_storage_defender

Manages the Defender for Storage.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "westus2"
}

resource "azurerm_storage_account" "example" {
  name                = "examplestorageaccount"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_security_center_storage_defender" "example" {
  storage_account_id = azurerm_storage_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the Storage Account where Defender for Storage is applied. Changing this forces a new resource to be created.

* `malware_scanning_on_upload_cap_gb_per_month` - (Optional) The max GB to be scanned per month. Defaults to `-1`. Possible values are `-1` and any positive integer.

* `malware_scanning_on_upload_enabled` - (Optional) Whether `malware scanning on upload` is enabled. Defaults to `false`.

* `malware_scanning_on_upload_exclude_blobs_larger_than` - (Optional) The maximum blob size in bytes to scan. Possible values are any positive integer.

* `malware_scanning_on_upload_exclude_blobs_with_prefix` - (Optional) A list of blob prefixes to exclude from on-upload malware scanning.

-> **Note:** Prefixes use the format `container-name/blob-name`. Use a container name without a trailing `/` to exclude matching container prefixes, or add a trailing `/` to target a single container only.

* `malware_scanning_on_upload_exclude_blobs_with_suffix` - (Optional) A list of blob suffixes to exclude from on-upload malware scanning.

* `override_subscription_settings_enabled` - (Optional) Whether the settings defined for this Storage Account should override the settings defined for the subscription. Defaults to `false`.

* `scan_results_event_grid_topic_id` - (Optional) The ID of the Event Grid Topic where scan results are sent.

~> **Note:** Setting `scan_results_event_grid_topic_id` requires `override_subscription_settings_enabled` to be `true` so the Storage Account can override the subscription-level Defender for Storage settings.

* `sensitive_data_discovery_enabled` - (Optional) Whether Sensitive Data Discovery should be enabled. Defaults to `false`.
 
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Defender for Storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Defender for Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Defender for Storage.
* `update` - (Defaults to 30 minutes) Used when updating the Defender for Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Defender for Storage.

## Import

A Defender for Storage can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_storage_defender.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/storageacc
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Security` - 2025-06-01
