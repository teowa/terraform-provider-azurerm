---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_subnet"
description: |-
  Gets information about an existing Subnet located within a Virtual Network.
---

# Data Source: azurerm_subnet

Gets information about an existing Subnet within a Virtual Network.

## Example Usage

```hcl
data "azurerm_subnet" "example" {
  name                 = "backend"
  virtual_network_name = "production"
  resource_group_name  = "networking"
}

output "subnet_id" {
  value = data.azurerm_subnet.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - Specifies the name of the Subnet.

* `resource_group_name` - Specifies the name of the Resource Group the Virtual Network is located in.

* `virtual_network_name` - Specifies the name of the Virtual Network this Subnet is located within.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Subnet.
* `address_prefixes` - The address prefixes for the subnet.
* `default_outbound_access_enabled` - Whether `default outbound access` is enabled for the subnet.
* `network_security_group_id` - The ID of the Network Security Group associated with the subnet.
* `private_endpoint_network_policies` - The network policies for the private endpoint on the subnet.
* `private_link_service_network_policies_enabled` - Whether `private link service network policies` is enabled for the subnet.
* `route_table_id` - The ID of the Route Table associated with this subnet.
* `service_endpoint` - One or more `service_endpoint` blocks as defined below.
* `service_endpoints` - A list of Service Endpoints within this subnet.

---

A `service_endpoint` block exports the following:

* `locations` - The list of locations scoped to the service endpoint.

* `network_identifier_id` - The resource ID used as the network identifier for the service endpoint.

* `service` - The service endpoint associated with the subnet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Subnet located within a Virtual Network.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Network` - 2025-01-01
