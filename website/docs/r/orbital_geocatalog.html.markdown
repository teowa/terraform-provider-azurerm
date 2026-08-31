---
subcategory: "Orbital"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_orbital_geocatalog"
description: |-
  Manages an Orbital GeoCatalog.
---

# azurerm_orbital_geocatalog

Manages an Orbital GeoCatalog.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_orbital_geocatalog" "example" {
  name                = "example-geocatalog"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Orbital GeoCatalog. Changing this forces a new Orbital GeoCatalog to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where this Orbital GeoCatalog should exist. Changing this forces a new Orbital GeoCatalog to be created.

* `location` - (Required) Specifies the Azure Region where the Orbital GeoCatalog should exist. Changing this forces a new Orbital GeoCatalog to be created.

* `auto_generated_domain_name_label_scope` - (Optional) Specifies the scope for how the auto-generated domain name for this Orbital GeoCatalog can be reused. Possible values are `TenantReuse`, `SubscriptionReuse`, `ResourceGroupReuse` and `NoReuse`. Defaults to `TenantReuse`. Changing this forces a new Orbital GeoCatalog to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `tier` - (Optional) Specifies the tier of this Orbital GeoCatalog. The only possible value is `Basic`. Defaults to `Basic`. Changing this forces a new Orbital GeoCatalog to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Orbital GeoCatalog.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Orbital GeoCatalog. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Orbital GeoCatalog.

~> **Note:** `identity_ids` is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Orbital GeoCatalog.

* `catalog_uri` - The URI of this Orbital GeoCatalog.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the System-Assigned Managed Identity assigned to this Orbital GeoCatalog.

* `tenant_id` - The Tenant ID for the System-Assigned Managed Identity assigned to this Orbital GeoCatalog.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Orbital GeoCatalog.
* `read` - (Defaults to 5 minutes) Used when retrieving the Orbital GeoCatalog.
* `update` - (Defaults to 30 minutes) Used when updating the Orbital GeoCatalog.
* `delete` - (Defaults to 30 minutes) Used when deleting the Orbital GeoCatalog.

## Import

Orbital GeoCatalogs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_orbital_geocatalog.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Orbital/geoCatalogs/geocatalog1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Orbital` - 2026-04-15
