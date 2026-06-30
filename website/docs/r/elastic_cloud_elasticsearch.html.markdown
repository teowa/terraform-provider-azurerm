---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_cloud_elasticsearch"
description: |-
  Manages an Elasticsearch in Elastic Cloud.
---

# azurerm_elastic_cloud_elasticsearch

Manages an Elasticsearch in Elastic Cloud.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_elastic_cloud_elasticsearch" "test" {
  name                        = "example-elasticsearch"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "user@example.com"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Elasticsearch resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Elasticsearch resource should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Elasticsearch resource should exist. Changing this forces a new resource to be created.

* `elastic_cloud_email_address` - (Required) The email address to associate with this Elasticsearch account. Changing this forces a new resource to be created.

* `sku_name` - (Required) The SKU name for this Elasticsearch. Changing this forces a new resource to be created.

-> **Note:** The SKU depends on the Elasticsearch Plans available for your account and is a combination of PlanID_Term.
Ex: If the plan ID is "planXYZ" and term is "Yearly", the SKU will be "planXYZ_Yearly".
You may find your eligible plans [here](https://portal.azure.com/#view/Microsoft_Azure_Marketplace/GalleryItemDetailsBladeNopdl/id/elastic.ec-azure-pp) or in the online documentation [here](https://azuremarketplace.microsoft.com/marketplace/apps/elastic.ec-azure-pp?tab=PlansAndPrice) for more details or in case of any issues with the SKU.

---

* `logs` - (Optional) A `logs` block as defined below.

* `monitoring_enabled` - (Optional) Whether `monitoring` is enabled. Defaults to `true`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Elasticsearch resource.

---

A `filtering_tag` block supports the following:

* `action` - (Required) Specifies the type of action which should be taken when the Tag matches the `name` and `value`. Possible values are `Exclude` and `Include`.

* `name` - (Required) Specifies the name (key) of the Tag which should be filtered.

* `value` - (Required) Specifies the value of the Tag which should be filtered.

---

A `logs` block supports the following:

* `filtering_tag` - (Optional) A list of `filtering_tag` blocks as defined above.

* `send_activity_logs` - (Optional) Whether Azure Activity Logs are sent to the Elasticsearch cluster. Defaults to `false`.

* `send_azuread_logs` - (Optional) Whether Azure AD logs are sent to the Elasticsearch cluster. Defaults to `false`.

* `send_subscription_logs` - (Optional) Whether Azure Subscription Logs are sent to the Elasticsearch cluster. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Elasticsearch.

* `elastic_cloud_deployment_id` - The ID of the Deployment within Elastic Cloud.

* `elastic_cloud_sso_default_url` - The Default URL used for Single Sign On (SSO) to Elastic Cloud.

* `elastic_cloud_user_id` - The ID of the User Account within Elastic Cloud.

* `elasticsearch_service_url` - The URL to the Elasticsearch Service associated with this Elasticsearch.

* `kibana_service_url` - The URL to the Kibana Dashboard associated with this Elasticsearch.

* `kibana_sso_uri` - The URI used for SSO to the Kibana Dashboard associated with this Elasticsearch.

* `monitor_properties` - A `monitor_properties` block as defined below.

---

A `monitor_properties` block exports the following:

* `generate_api_key` - Whether an API key is generated for this Elasticsearch.

* `hosting_type` - The hosting type for this Elasticsearch.

* `liftr_resource_category` - The Liftr resource category for this Elasticsearch.

* `liftr_resource_preference` - The Liftr resource preference for this Elasticsearch.

* `monitoring_status` - The monitoring status for this Elasticsearch.

* `plan_details` - A `plan_details` block as defined below.

* `project_details` - A `project_details` block as defined below.

* `provisioning_state` - The provisioning state of this Elasticsearch.

* `saas_azure_subscription_status` - The SaaS Azure subscription status for this Elasticsearch.

* `source_campaign_id` - The source campaign ID for this Elasticsearch.

* `source_campaign_name` - The source campaign name for this Elasticsearch.

* `subscription_state` - The subscription state for this Elasticsearch.

* `user_info` - A `user_info` block as defined below.

* `version` - The Elastic version for this Elasticsearch.

---

A `plan_details` block exports the following:

* `offer_id` - The Marketplace offer ID for this Elasticsearch.

* `plan_id` - The Marketplace plan ID for this Elasticsearch.

* `plan_name` - The Marketplace plan name for this Elasticsearch.

* `publisher_id` - The Marketplace publisher ID for this Elasticsearch.

* `term_id` - The Marketplace term ID for this Elasticsearch.

---

A `project_details` block exports the following:

* `configuration_type` - The project configuration type for this Elasticsearch.

* `project_type` - The project type for this Elasticsearch.

---

A `user_info` block exports the following:

* `company_info` - A `company_info` block as defined below.

* `company_name` - The company name associated with this Elasticsearch.

* `email_address` - The email address associated with this Elasticsearch.

* `first_name` - The first name associated with this Elasticsearch.

* `last_name` - The last name associated with this Elasticsearch.

---

A `company_info` block exports the following:

* `business` - The business associated with this Elasticsearch.

* `country` - The country associated with this Elasticsearch.

* `domain` - The domain associated with this Elasticsearch.

* `employees_number` - The employee count associated with this Elasticsearch.

* `state` - The state associated with this Elasticsearch.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Elasticsearch.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elasticsearch.
* `update` - (Defaults to 1 hour) Used when updating the Elasticsearch.
* `delete` - (Defaults to 1 hour) Used when deleting the Elasticsearch.

## Import

An Elasticsearch can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_cloud_elasticsearch.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Elastic/monitors/monitor1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Elastic` - 2025-06-01
