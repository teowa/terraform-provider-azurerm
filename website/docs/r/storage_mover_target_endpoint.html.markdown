---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_target_endpoint"
description: |-
  Manages a Storage Mover Target Endpoint.
---

# azurerm_storage_mover_target_endpoint

Manages a Storage Mover Target Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                            = "examplestr"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "example" {
  name                  = "example-sc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "blob"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
}

resource "azurerm_storage_mover_target_endpoint" "example" {
  name                   = "example-se"
  storage_mover_id       = azurerm_storage_mover.example.id
  storage_account_id     = azurerm_storage_account.example.id
  storage_container_name = azurerm_storage_container.example.name
  description            = "Example Storage Container Endpoint Description"
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover Target Endpoint. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the storage mover for this Storage Mover Target Endpoint. Changing this forces a new resource to be created.

-> **Note:** Exactly one of `storage_account_id`, `azure_multi_cloud_connector`, `azure_storage_nfs_file_share`, `azure_storage_smb_file_share` or `smb_mount` must be specified.

* `storage_account_id` - (Optional) Specifies the ID of the storage account for this Storage Mover Target Endpoint. Changing this forces a new resource to be created.

* `storage_container_name` - (Optional) Specifies the name of the storage blob container for this Storage Mover Target Endpoint. Required when `storage_account_id` is specified. Changing this forces a new resource to be created.

* `azure_multi_cloud_connector` - (Optional) An `azure_multi_cloud_connector` block as defined below. Changing this forces a new resource to be created.

* `azure_storage_nfs_file_share` - (Optional) An `azure_storage_nfs_file_share` block as defined below. Changing this forces a new resource to be created.

* `azure_storage_smb_file_share` - (Optional) An `azure_storage_smb_file_share` block as defined below. Changing this forces a new resource to be created.

* `smb_mount` - (Optional) A `smb_mount` block as defined below. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `description` - (Optional) Specifies a description for the Storage Mover Target Endpoint.

---

An `azure_multi_cloud_connector` block supports the following:

* `aws_s3_bucket_id` - (Required) Specifies the resource ID of the AWS S3 Bucket connected through the Multi Cloud Connector. Changing this forces a new resource to be created.

* `multi_cloud_connector_id` - (Required) Specifies the resource ID of the Multi Cloud Connector used to access the target. Changing this forces a new resource to be created.

---

An `azure_storage_nfs_file_share` block supports the following:

* `file_share_name` - (Required) Specifies the name of the Azure Storage NFS File Share. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) Specifies the ID of the Storage Account containing the NFS File Share. Changing this forces a new resource to be created.

---

An `azure_storage_smb_file_share` block supports the following:

* `file_share_name` - (Required) Specifies the name of the Azure Storage SMB File Share. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) Specifies the ID of the Storage Account containing the SMB File Share. Changing this forces a new resource to be created.

---

A `smb_mount` block supports the following:

* `host` - (Required) Specifies the host name or IP address of the SMB server. Changing this forces a new resource to be created.

* `share_name` - (Required) Specifies the name of the SMB share. Changing this forces a new resource to be created.

* `credentials` - (Optional) A `credentials` block as defined below.

---

A `credentials` block supports the following:

* `username_uri` - (Optional) Specifies the Azure Key Vault URI containing the username used to access the SMB share.

* `password_uri` - (Optional) Specifies the Azure Key Vault URI containing the password used to access the SMB share.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Storage Mover Target Endpoint. The only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover Target Endpoint.

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover Target Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover Target Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover Target Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover Target Endpoint.

## Import

Storage Mover Target Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_target_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/endpoints/endpoint1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageMover` - 2025-07-01
