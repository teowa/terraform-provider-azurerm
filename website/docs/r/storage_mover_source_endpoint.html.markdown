---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_source_endpoint"
description: |-
  Manages a Storage Mover Source Endpoint.
---

# azurerm_storage_mover_source_endpoint

Manages a Storage Mover Source Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
}

resource "azurerm_storage_mover_source_endpoint" "example" {
  name             = "example-se"
  storage_mover_id = azurerm_storage_mover.example.id
  export           = "/"
  host             = "192.168.0.1"
  nfs_version      = "NFSv3"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover Source Endpoint. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the Storage Mover for this Storage Mover Source Endpoint. Changing this forces a new resource to be created.

-> **Note:** Exactly one of `host`, `azure_multi_cloud_connector`, `azure_storage_nfs_file_share`, `azure_storage_smb_file_share` or `smb_mount` must be specified.

* `host` - (Optional) Specifies the host name or IP address of the server exporting the file system. Changing this forces a new resource to be created.

* `export` - (Optional) Specifies the directory being exported from the server. Changing this forces a new resource to be created.

* `nfs_version` - (Optional) Specifies the NFS protocol version. Possible values are `NFSauto`, `NFSv3` and `NFSv4`. Defaults to `NFSauto`. Changing this forces a new resource to be created.

* `azure_multi_cloud_connector` - (Optional) An `azure_multi_cloud_connector` block as defined below. Changing this forces a new resource to be created.

* `azure_storage_nfs_file_share` - (Optional) An `azure_storage_nfs_file_share` block as defined below. Changing this forces a new resource to be created.

* `azure_storage_smb_file_share` - (Optional) An `azure_storage_smb_file_share` block as defined below. Changing this forces a new resource to be created.

* `smb_mount` - (Optional) A `smb_mount` block as defined below. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `description` - (Optional) Specifies a description for the Storage Mover Source Endpoint.

---

An `azure_multi_cloud_connector` block supports the following:

* `aws_s3_bucket_id` - (Required) Specifies the resource ID of the AWS S3 Bucket connected through the Multi Cloud Connector. Changing this forces a new resource to be created.

* `multi_cloud_connector_id` - (Required) Specifies the resource ID of the Multi Cloud Connector used to access the source. Changing this forces a new resource to be created.

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

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Storage Mover Source Endpoint. The only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover Source Endpoint.

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover Source Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover Source Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover Source Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover Source Endpoint.

## Import

Storage Mover Source Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_source_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/endpoints/endpoint1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageMover` - 2025-07-01
