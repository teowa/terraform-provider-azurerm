---
subcategory: "Datadog"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_datadog_monitor"
description: |-
  Manages a Datadog Monitor.
---

# azurerm_datadog_monitor

Manages a Datadog Monitor.

## Example Usage

### Monitor creation with linking to Datadog organization

```hcl
variable "datadog_api_key" {
  type      = string
  sensitive = true
}

variable "datadog_application_key" {
  type      = string
  sensitive = true
}

resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West US 2"
}

resource "azurerm_datadog_monitor" "example" {
  name                = "example-datadog-monitor"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  datadog_organization {
    api_key         = var.datadog_api_key
    application_key = var.datadog_application_key
  }

  user {
    name  = "Example"
    email = "abc@xyz.com"
  }

  sku_name = "Linked"

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Datadog Monitor. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Datadog Monitor should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Datadog Monitor should exist. Changing this forces a new resource to be created.

* `datadog_organization` - (Required) A `datadog_organization` block as defined below.

* `sku_name` - (Required) The name of the SKU for the Datadog Monitor.

* `user` - (Required) A `user` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `monitoring_enabled` - (Optional) Whether monitoring is enabled. Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Datadog Monitor.

---

A `datadog_organization` block supports the following:

* `api_key` - (Required) The API key associated with the Datadog organization. Changing this forces a new resource to be created.

* `application_key` - (Required) The application key associated with the Datadog organization. Changing this forces a new resource to be created.

* `cspm` - (Optional) Whether CSPM is enabled for the Datadog organization.

* `enterprise_app_id` - (Optional) The ID of the enterprise application. Changing this forces a new resource to be created.

* `linking_auth_code` - (Optional) The authorization code used to link an existing Datadog organization. Changing this forces a new resource to be created.

* `linking_client_id` - (Optional) The client ID used to link an existing Datadog organization. Changing this forces a new resource to be created.

* `name` - (Optional) The name of the Datadog organization. Changing this forces a new resource to be created.

* `redirect_uri` - (Optional) The redirect URI for linking. Changing this forces a new resource to be created.

* `resource_collection` - (Optional) Whether resource collection is enabled for the Datadog organization.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Datadog Monitor. The only possible value is `SystemAssigned`.

-> **Note:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and the Datadog Monitor has been created. More details are available below.

---

A `user` block supports the following:

* `email` - (Required) The email address of the user that Datadog can contact if needed. Changing this forces a new resource to be created.

* `name` - (Required) The name of the user that Datadog can contact if needed. Changing this forces a new resource to be created.

* `phone_number` - (Optional) The phone number of the user that Datadog can contact if needed. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Datadog Monitor.

* `identity` - An `identity` block as defined below.

* `marketplace_subscription_status` - The Marketplace subscription status of the Datadog Monitor.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Datadog Monitor.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Datadog Monitor.

-> **Note:** You can access the Principal ID via `${azurerm_datadog_monitor.example.identity[0].principal_id}` and the Tenant ID via `${azurerm_datadog_monitor.example.identity[0].tenant_id}`

## Role Assignment

To enable metrics flow, perform a role assignment on the identity created above. The `Monitoring Reader` role (`43d0d8ad-25c7-4714-9337-8ba259a9fe05`) is required.

### Role assignment on the monitor created

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "monitoring_reader" {
  name = "Monitoring Reader"
}

resource "azurerm_role_assignment" "example" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = data.azurerm_role_definition.monitoring_reader.role_definition_id
  principal_id       = azurerm_datadog_monitor.example.identity[0].principal_id
}
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Datadog Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Datadog Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the Datadog Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Datadog Monitor.

## Import

A Datadog Monitor can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_datadog_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Datadog` - 2025-06-11
