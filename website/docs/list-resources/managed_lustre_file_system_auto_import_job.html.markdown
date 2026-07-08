---
subcategory: "Azure Managed Lustre File System"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_lustre_file_system_auto_import_job"
description: |-
  Lists Azure Managed Lustre File System Auto Import Job resources.
---

# List resource: azurerm_managed_lustre_file_system_auto_import_job

Lists Azure Managed Lustre File System Auto Import Job resources.

## Example Usage

### List Auto Import Jobs in an Azure Managed Lustre File System

```hcl
list "azurerm_managed_lustre_file_system_auto_import_job" "example" {
  provider = azurerm
  config {
    managed_lustre_file_system_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resource-group/providers/Microsoft.StorageCache/amlFilesystems/example-managed-lustre-file-system"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `managed_lustre_file_system_id` - (Required) The ID of the Azure Managed Lustre File System to query.
