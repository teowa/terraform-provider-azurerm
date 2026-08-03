---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_cloud_hosted_search"
description: |-
  Manages an Elastic Cloud hosted search deployment.
---

# azurerm_elastic_cloud_hosted_search

Manages an Elastic Cloud hosted search deployment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_elastic_cloud_hosted_search" "example" {
  name                        = "example-hosted-search"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "user@example.com"
}
```

## Arguments Reference

The following arguments are supported:

* `elastic_cloud_email_address` - (Required) The email address associated with this Elastic Cloud hosted search deployment. Changing this forces a new Elastic Cloud hosted search deployment to be created.

* `location` - (Required) The Azure Region where the Elastic Cloud hosted search deployment should exist. Changing this forces a new Elastic Cloud hosted search deployment to be created.

* `name` - (Required) The name of the Elastic Cloud hosted search deployment. Changing this forces a new Elastic Cloud hosted search deployment to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Elastic Cloud hosted search deployment should exist. Changing this forces a new Elastic Cloud hosted search deployment to be created.

* `sku_name` - (Required) The SKU name for this Elastic Cloud hosted search deployment. Changing this forces a new Elastic Cloud hosted search deployment to be created.

-> **Note:** The SKU depends on the Elastic plans available for your account and is a combination of Plan ID and term. For example, if the plan ID is `planXYZ` and the term is `Yearly`, the SKU is `planXYZ_Yearly`. You can review eligible plans in the [Azure portal](https://portal.azure.com/#view/Microsoft_Azure_Marketplace/GalleryItemDetailsBladeNopdl/id/elastic.ec-azure-pp) or the [Azure Marketplace listing](https://azuremarketplace.microsoft.com/en-us/marketplace/apps/elastic.ec-azure-pp?tab=PlansAndPrice).

---

* `monitoring_enabled` - (Optional) Whether monitoring should be enabled for this Elastic Cloud hosted search deployment. Defaults to `true`. Changing this forces a new Elastic Cloud hosted search deployment to be created.

* `tags` - (Optional) A mapping of tags assigned to the Elastic Cloud hosted search deployment.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Elastic Cloud hosted search deployment.

* `elastic_cloud_deployment_id` - The Elastic Cloud deployment ID.

* `elastic_cloud_sso_default_url` - The default single sign-on URL for Elastic Cloud.

* `elastic_cloud_user_id` - The Elastic Cloud user ID.

* `elasticsearch_service_url` - The Elasticsearch service URL for the deployment.

* `kibana_service_url` - The Kibana service URL for the deployment.

* `kibana_sso_uri` - The Kibana single sign-on URI for the deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Elastic Cloud hosted search deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic Cloud hosted search deployment.
* `update` - (Defaults to 1 hour) Used when updating the Elastic Cloud hosted search deployment.
* `delete` - (Defaults to 1 hour) Used when deleting the Elastic Cloud hosted search deployment.

## Import

Elastic Cloud hosted search deployments can be imported using the `resource id`, for example:

```shell
terraform import azurerm_elastic_cloud_hosted_search.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Elastic/monitors/monitor1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Elastic` - 2025-06-01
