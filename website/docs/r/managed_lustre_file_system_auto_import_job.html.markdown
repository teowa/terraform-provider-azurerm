---
subcategory: "Azure Managed Lustre File System"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_lustre_file_system_auto_import_job"
description: |-
  Manages an Azure Managed Lustre File System Auto Import Job.
---

# azurerm_managed_lustre_file_system_auto_import_job

Manages an Azure Managed Lustre File System Auto Import Job.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-virtual-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_storage_account" "example" {
  name                            = "examplestorageaccount"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "example" {
  name                  = "example-storage-container"
  storage_account_id    = azurerm_storage_account.example.id
  container_access_type = "private"
}

resource "azurerm_storage_container" "example_logging" {
  name                  = "example-logging-container"
  storage_account_id    = azurerm_storage_account.example.id
  container_access_type = "private"
}

data "azuread_service_principal" "example" {
  display_name = "HPC Cache Resource Provider"
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_role_assignment" "example_blob" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_managed_lustre_file_system" "example" {
  name                   = "example-managed-lustre-file-system"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  sku_name               = "AMLFS-Durable-Premium-250"
  subnet_id              = azurerm_subnet.example.id
  storage_capacity_in_tb = 8
  zones                  = ["2"]

  maintenance_window {
    day_of_week        = "Friday"
    time_of_day_in_utc = "22:00"
  }

  hsm_setting {
    container_id         = azurerm_storage_container.example.id
    logging_container_id = azurerm_storage_container.example_logging.id
    import_prefix        = "/"
  }

  depends_on = [
    azurerm_role_assignment.example,
    azurerm_role_assignment.example_blob,
  ]
}

resource "azurerm_managed_lustre_file_system_auto_import_job" "example" {
  name                          = "example-managed-lustre-auto-import-job"
  managed_lustre_file_system_id = azurerm_managed_lustre_file_system.example.id
  location                      = azurerm_resource_group.example.location
  admin_status                  = "Enable"
  auto_import_prefixes          = ["/data", "/archive"]
  conflict_resolution_mode      = "OverwriteAlways"
  enable_deletions              = true
  maximum_errors                = 5

  tags = {
    environment = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Managed Lustre File System Auto Import Job. Changing this forces a new resource to be created.

-> **Note:** The name must be between 2 and 80 characters in length, start and end with an alphanumeric character, and contain only alphanumeric characters, hyphens, and underscores.

* `managed_lustre_file_system_id` - (Required) The ID of the Azure Managed Lustre File System where this Auto Import Job should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Managed Lustre File System Auto Import Job should exist. Changing this forces a new resource to be created.

* `admin_status` - (Optional) The administrative status of this Azure Managed Lustre File System Auto Import Job. Possible values are `Disable` and `Enable`. Defaults to `Enable`.

* `auto_import_prefixes` - (Optional) A list of namespace prefixes to import into the Azure Managed Lustre File System. Changing this forces a new resource to be created.

-> **Note:** Defaults to `["/"]`.

* `conflict_resolution_mode` - (Optional) The conflict resolution mode for imported content. Possible values are `Fail`, `OverwriteAlways`, `OverwriteIfDirty`, and `Skip`. Defaults to `Skip`. Changing this forces a new resource to be created.

* `enable_deletions` - (Optional) Whether delete events from the backing storage should be propagated into the Azure Managed Lustre File System namespace. Defaults to `false`. Changing this forces a new resource to be created.

* `maximum_errors` - (Optional) The maximum number of import errors to tolerate before the job stops. Use `-1` for unlimited errors. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Managed Lustre File System Auto Import Job.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Managed Lustre File System Auto Import Job.

* `provisioning_state` - The provisioning state of the Azure Managed Lustre File System Auto Import Job.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Managed Lustre File System Auto Import Job.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Managed Lustre File System Auto Import Job.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Managed Lustre File System Auto Import Job.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Managed Lustre File System Auto Import Job.

## Import

Azure Managed Lustre File System Auto Import Jobs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_lustre_file_system_auto_import_job.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageCache/amlFilesystems/amlFilesystem1/autoImportJobs/autoImportJob1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageCache` - 2025-07-01
