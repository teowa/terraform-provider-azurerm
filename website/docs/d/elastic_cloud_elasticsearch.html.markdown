---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_elastic_cloud_elasticsearch"
description: |-
  Gets information about an existing Elasticsearch resource.

---

# Data Source: azurerm_elastic_cloud_elasticsearch

Gets information about an existing Elasticsearch resource.

## Example Usage

```hcl
data "azurerm_elastic_cloud_elasticsearch" "example" {
  name                = "my-elastic-search"
  resource_group_name = "example-resources"
}

output "elasticsearch_endpoint" {
  value = data.azurerm_elastic_cloud_elasticsearch.example.elasticsearch_service_url
}

output "kibana_endpoint" {
  value = data.azurerm_elastic_cloud_elasticsearch.example.kibana_service_url
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Elasticsearch resource.

* `resource_group_name` - (Required) The name of the Resource Group in which the Elasticsearch exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Elasticsearch.

* `elastic_cloud_deployment_id` - The ID of the Deployment within Elastic Cloud.

* `elastic_cloud_email_address` - The Email Address which is associated with this Elasticsearch account.

* `elastic_cloud_sso_default_url` - The Default URL used for Single Sign On (SSO) to Elastic Cloud.

* `elastic_cloud_user_id` - The ID of the User Account within Elastic Cloud.

* `elasticsearch_service_url` - The URL to the Elasticsearch Service associated with this Elasticsearch.

* `kibana_service_url` - The URL to the Kibana Dashboard associated with this Elasticsearch.

* `kibana_sso_uri` - The URI used for SSO to the Kibana Dashboard associated with this Elasticsearch.

* `location` - The Azure Region in which this Elasticsearch exists.

* `logs` - A `logs` block as defined below.

* `monitor_properties` - A `monitor_properties` block as defined below.

* `monitoring_enabled` - Whether `monitoring` is enabled.

* `sku_name` - The name of the SKU used for this Elasticsearch.

* `tags` - A mapping of tags assigned to the Elasticsearch.

---

A `filtering_tag` block exports the following:

* `action` - The type of action which is taken when the Tag matches the `name` and `value`.

* `name` - The name (key) of the Tag which should be filtered.

* `value` - The value of the Tag which should be filtered.

---

A `logs` block exports the following:

* `filtering_tag` - A list of `filtering_tag` blocks as defined above.

* `send_activity_logs` - Whether Azure Activity Logs are sent to the Elasticsearch cluster.

* `send_azuread_logs` - Whether Azure AD logs are sent to the Elasticsearch cluster.

* `send_subscription_logs` - Whether Azure Subscription Logs are sent to the Elasticsearch cluster.

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

* `read` - (Defaults to 5 minutes) Used when retrieving the Elasticsearch.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Elastic` - 2025-06-01
