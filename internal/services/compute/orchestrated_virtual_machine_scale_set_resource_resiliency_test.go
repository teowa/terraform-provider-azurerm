// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccOrchestratedVirtualMachineScaleSet_automaticZoneRebalancingPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2"
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.automaticZoneRebalancingPolicy(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("automatic_zone_rebalancing_policy.#").DoesNotExist(),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.automaticZoneRebalancingPolicy(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("automatic_zone_rebalancing_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("automatic_zone_rebalancing_policy.0.rebalance_behavior").HasValue("CreateBeforeDelete"),
				check.That(data.ResourceName).Key("automatic_zone_rebalancing_policy.0.rebalance_strategy").HasValue("Recreate"),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.automaticZoneRebalancingPolicy(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("automatic_zone_rebalancing_policy.#").DoesNotExist(),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r OrchestratedVirtualMachineScaleSetResource) automaticZoneRebalancingPolicy(data acceptance.TestData, enabled bool) string {
	policy := ""
	if enabled {
		policy = `
  automatic_zone_rebalancing_policy {
    enabled            = true
    rebalance_behavior = "CreateBeforeDelete"
    rebalance_strategy = "Recreate"
  }
`
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  zones = ["1", "2"]

  sku_name  = "Standard_D1_v2"
  instances = 2

  platform_fault_domain_count = 2

%[4]s

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data), policy)
}
