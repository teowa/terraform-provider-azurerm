package azurestackhci_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StackHCILogicalNetworkResource struct{}

func TestAccStackHCILogicalNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_logical_network", "test")
	r := StackHCILogicalNetworkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCILogicalNetwork_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_logical_network", "test")
	r := StackHCILogicalNetworkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStackHCILogicalNetwork_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_logical_network", "test")
	r := StackHCILogicalNetworkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StackHCILogicalNetworkResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.LogicalNetworks
	id, err := logicalnetworks.ParseLogicalNetworkID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r StackHCILogicalNetworkResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
  
}
`, template)
}

func (r StackHCILogicalNetworkResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_logical_network" "import" {
  stack_hci_cluster_id = azurerm_stack_hci_cluster.test.id
  arc_resource_ids     = azurerm_stack_hci_logical_network.test.arc_resource_ids
  version              = azurerm_stack_hci_logical_network.test.version

  scale_unit {
    adou_path                     = azurerm_stack_hci_logical_network.test.scale_unit.0.adou_path
    domain_fqdn                   = azurerm_stack_hci_logical_network.test.scale_unit.0.domain_fqdn
    secrets_location              = azurerm_stack_hci_logical_network.test.scale_unit.0.secrets_location
    naming_prefix                 = azurerm_stack_hci_logical_network.test.scale_unit.0.naming_prefix
    streaming_data_client_enabled = azurerm_stack_hci_logical_network.test.scale_unit.0.streaming_data_client_enabled
    eu_location_enabled           = azurerm_stack_hci_logical_network.test.scale_unit.0.eu_location_enabled
    episodic_data_upload_enabled  = azurerm_stack_hci_logical_network.test.scale_unit.0.episodic_data_upload_enabled

    bitlocker_boot_volume_enabled   = azurerm_stack_hci_logical_network.test.scale_unit.0.bitlocker_boot_volume_enabled
    bitlocker_data_volume_enabled   = azurerm_stack_hci_logical_network.test.scale_unit.0.bitlocker_data_volume_enabled
    credential_guard_enabled        = azurerm_stack_hci_logical_network.test.scale_unit.0.credential_guard_enabled
    drift_control_enabled           = azurerm_stack_hci_logical_network.test.scale_unit.0.drift_control_enabled
    drtm_protection_enabled         = azurerm_stack_hci_logical_network.test.scale_unit.0.drtm_protection_enabled
    hvci_protection_enabled         = azurerm_stack_hci_logical_network.test.scale_unit.0.hvci_protection_enabled
    side_channel_mitigation_enabled = azurerm_stack_hci_logical_network.test.scale_unit.0.side_channel_mitigation_enabled
    smb_cluster_encryption_enabled  = azurerm_stack_hci_logical_network.test.scale_unit.0.smb_cluster_encryption_enabled
    smb_signing_enabled             = azurerm_stack_hci_logical_network.test.scale_unit.0.smb_signing_enabled
    wdac_enabled                    = azurerm_stack_hci_logical_network.test.scale_unit.0.wdac_enabled


    cluster {
      azure_service_endpoint = azurerm_stack_hci_logical_network.test.scale_unit.0.cluster.0.azure_service_endpoint
      cloud_account_name     = azurerm_stack_hci_logical_network.test.scale_unit.0.cluster.0.cloud_account_name
      name                   = azurerm_stack_hci_logical_network.test.scale_unit.0.cluster.0.name
      witness_type           = azurerm_stack_hci_logical_network.test.scale_unit.0.cluster.0.witness_type
      witness_path           = azurerm_stack_hci_logical_network.test.scale_unit.0.cluster.0.witness_path
    }

    host_network {
      intent {
        name         = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.intent.0.name
        adapter      = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.intent.0.adapter
        traffic_type = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.intent.0.traffic_type
      }

      intent {
        name         = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.intent.1.name
        adapter      = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.intent.1.adapter
        traffic_type = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.intent.1.traffic_type
      }

      storage_network {
        name                 = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.storage_network.0.name
        network_adapter_name = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.storage_network.0.network_adapter_name
        vlan_id              = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.storage_network.0.vlan_id
      }

      storage_network {
        name                 = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.storage_network.1.name
        network_adapter_name = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.storage_network.1.network_adapter_name
        vlan_id              = azurerm_stack_hci_logical_network.test.scale_unit.0.host_network.0.storage_network.1.vlan_id
      }
    }

    infrastructure_network {
      gateway     = azurerm_stack_hci_logical_network.test.scale_unit.0.infrastructure_network.0.gateway
      subnet_mask = azurerm_stack_hci_logical_network.test.scale_unit.0.infrastructure_network.0.subnet_mask
      dns_server  = azurerm_stack_hci_logical_network.test.scale_unit.0.infrastructure_network.0.dns_server
      ip_pool {
        ending_address   = azurerm_stack_hci_logical_network.test.scale_unit.0.infrastructure_network.0.ip_pool.0.ending_address
        starting_address = azurerm_stack_hci_logical_network.test.scale_unit.0.infrastructure_network.0.ip_pool.0.starting_address
      }
    }

    optional_service {
      custom_location = azurerm_stack_hci_logical_network.test.scale_unit.0.optional_service.0.custom_location
    }

    physical_node {
      ipv4_address = azurerm_stack_hci_logical_network.test.scale_unit.0.physical_node.0.ipv4_address
      name         = azurerm_stack_hci_logical_network.test.scale_unit.0.physical_node.0.name
    }

    physical_node {
      ipv4_address = azurerm_stack_hci_logical_network.test.scale_unit.0.physical_node.1.ipv4_address
      name         = azurerm_stack_hci_logical_network.test.scale_unit.0.physical_node.1.name
    }

    storage {
      configuration_mode = azurerm_stack_hci_logical_network.test.scale_unit.0.storage.0.configuration_mode
    }
  }
}
`, config)
}

func (r StackHCILogicalNetworkResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_logical_network" "test" {
}
`, template)
}

func (r StackHCILogicalNetworkResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}

variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-vm-${var.random_string}"
  location = var.primary_location
}

`, data.Locations.Primary, data.RandomString)
}
