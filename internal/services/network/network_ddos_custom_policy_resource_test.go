// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-05-01/ddoscustompolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkDdosCustomPolicyResource struct{}

func TestAccNetworkDDoSCustomPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_custom_policy", "test")
	r := NetworkDdosCustomPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("detection_rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("detection_rule.0.traffic_detection_rule.0.traffic_type").HasValue("Tcp"),
				check.That(data.ResourceName).Key("resource_guid").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkDDoSCustomPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_custom_policy", "test")
	r := NetworkDdosCustomPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_ddos_custom_policy"),
		},
	})
}

func TestAccNetworkDDoSCustomPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_custom_policy", "test")
	r := NetworkDdosCustomPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updatedConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("detection_rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("detection_rule.0.traffic_detection_rule.0.packets_per_second").HasValue("2000000"),
				check.That(data.ResourceName).Key("detection_rule.0.traffic_detection_rule.0.traffic_type").HasValue("Udp"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkDDoSCustomPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_custom_policy", "test")
	r := NetworkDdosCustomPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("detection_rule.0.traffic_detection_rule.0.packets_per_second").HasValue("1000000"),
			),
		},
		{
			Config: r.updatedConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("detection_rule.0.traffic_detection_rule.0.packets_per_second").HasValue("2000000"),
				check.That(data.ResourceName).Key("detection_rule.0.traffic_detection_rule.0.traffic_type").HasValue("Udp"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
	})
}

func (NetworkDdosCustomPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := ddoscustompolicies.ParseDdosCustomPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.DdosCustomPoliciesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (NetworkDdosCustomPolicyResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := ddoscustompolicies.ParseDdosCustomPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	if err := client.Network.DdosCustomPoliciesClient.DeleteThenPoll(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (NetworkDdosCustomPolicyResource) basicConfigIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_custom_policy" "test" {
  name                = "acctestddoscp-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func (NetworkDdosCustomPolicyResource) basicConfigList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_custom_policy" "test" {
  name                = "acctestddoscp-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Ternary, data.RandomInteger)
}

func (NetworkDdosCustomPolicyResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_custom_policy" "test" {
  name                = "acctestddoscp-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  detection_rule {
    name           = "detectionRuleTcp"
    detection_mode = "TrafficThreshold"

    traffic_detection_rule {
      packets_per_second = 1000000
      traffic_type       = "Tcp"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r NetworkDdosCustomPolicyResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_ddos_custom_policy" "import" {
  name                = azurerm_network_ddos_custom_policy.test.name
  location            = azurerm_network_ddos_custom_policy.test.location
  resource_group_name = azurerm_network_ddos_custom_policy.test.resource_group_name

  detection_rule {
    name           = azurerm_network_ddos_custom_policy.test.detection_rule.0.name
    detection_mode = azurerm_network_ddos_custom_policy.test.detection_rule.0.detection_mode

    traffic_detection_rule {
      packets_per_second = azurerm_network_ddos_custom_policy.test.detection_rule.0.traffic_detection_rule.0.packets_per_second
      traffic_type       = azurerm_network_ddos_custom_policy.test.detection_rule.0.traffic_detection_rule.0.traffic_type
    }
  }
}
`, r.basicConfig(data))
}

func (NetworkDdosCustomPolicyResource) updatedConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_custom_policy" "test" {
  name                = "acctestddoscp-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  detection_rule {
    name           = "detectionRuleUdp"
    detection_mode = "TrafficThreshold"

    traffic_detection_rule {
      packets_per_second = 2000000
      traffic_type       = "Udp"
    }
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
