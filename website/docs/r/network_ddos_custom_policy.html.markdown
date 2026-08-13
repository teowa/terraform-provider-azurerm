---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_ddos_custom_policy"
description: |-
  Manages an Azure Network DDoS Custom Policy.
---

# azurerm_network_ddos_custom_policy

Manages an Azure Network DDoS Custom Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_ddos_custom_policy" "example" {
  name                = "example-ddos-custom-policy"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  detection_rule {
    name           = "detectionRuleTcp"
    detection_mode = "TrafficThreshold"

    traffic_detection_rule {
      packets_per_second = 1000000
      traffic_type       = "Tcp"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Network DDoS Custom Policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the Network DDoS Custom Policy. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Network DDoS Custom Policy should exist. Changing this forces a new resource to be created.

* `detection_rule` - (Optional) One or more `detection_rule` blocks as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `detection_rule` block supports the following:

* `name` - (Required) The name of the DDoS detection rule.

* `detection_mode` - (Required) The detection mode for the DDoS detection rule. Possible values are `TrafficThreshold`.

* `traffic_detection_rule` - (Required) A `traffic_detection_rule` block as documented below.

---

A `traffic_detection_rule` block supports the following:

* `packets_per_second` - (Required) The packets per second threshold for the traffic detection rule.

* `traffic_type` - (Required) The type of traffic to match. Possible values are `Tcp`, `TcpSyn`, and `Udp`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network DDoS Custom Policy.

* `frontend_ip_configuration_ids` - A list of frontend IP configuration IDs associated with the Network DDoS Custom Policy.

* `resource_guid` - The resource GUID of the Network DDoS Custom Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network DDoS Custom Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network DDoS Custom Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Network DDoS Custom Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network DDoS Custom Policy.

## Import

A Network DDoS Custom Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_ddos_custom_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/ddosCustomPolicies/ddosCustomPolicy1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-05-01
