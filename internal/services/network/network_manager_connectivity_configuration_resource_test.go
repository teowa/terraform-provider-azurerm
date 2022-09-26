package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/connectivityconfigurations"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkConnectivityConfigurationResource struct{}

func TestAccNetworkConnectivityConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connectivity_configuration", "test")
	r := NetworkConnectivityConfigurationResource{}
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

func TestAccNetworkConnectivityConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connectivity_configuration", "test")
	r := NetworkConnectivityConfigurationResource{}
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

func TestAccNetworkConnectivityConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connectivity_configuration", "test")
	r := NetworkConnectivityConfigurationResource{}
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

func TestAccNetworkConnectivityConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connectivity_configuration", "test")
	r := NetworkConnectivityConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r NetworkConnectivityConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := connectivityconfigurations.ParseConnectivityConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.ConnectivityConfigurationsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkConnectivityConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_network_network_manager" "test" {
  name                = "acctest-nnm-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r NetworkConnectivityConfigurationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_connectivity_configuration" "test" {
  name                       = "acctest-ncc-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
  connectivity_topology      = ""
  applies_to_groups {
    group_connectivity = ""
    is_global          = ""
    network_group_id   = ""
    use_hub_gateway    = ""
  }
}
`, template, data.RandomInteger)
}

func (r NetworkConnectivityConfigurationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_connectivity_configuration" "import" {
  name                       = azurerm_network_connectivity_configuration.test.name
  network_network_manager_id = azurerm_network_network_manager.test.id
  connectivity_topology      = ""
  applies_to_groups {
    group_connectivity = ""
    is_global          = ""
    network_group_id   = ""
    use_hub_gateway    = ""
  }
}
`, config)
}

func (r NetworkConnectivityConfigurationResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_connectivity_configuration" "test" {
  name                       = "acctest-ncc-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
  connectivity_topology      = ""
  delete_existing_peering    = ""
  description                = ""
  is_global                  = ""
  applies_to_groups {
    group_connectivity = ""
    is_global          = ""
    network_group_id   = ""
    use_hub_gateway    = ""
  }
  hubs {
    resource_id   = ""
    resource_type = ""
  }

}
`, template, data.RandomInteger)
}

func (r NetworkConnectivityConfigurationResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_connectivity_configuration" "test" {
  name                       = "acctest-ncc-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
  connectivity_topology      = ""
  delete_existing_peering    = ""
  description                = ""
  is_global                  = ""
  applies_to_groups {
    group_connectivity = ""
    is_global          = ""
    network_group_id   = ""
    use_hub_gateway    = ""
  }
  hubs {
    resource_id   = ""
    resource_type = ""
  }

}
`, template, data.RandomInteger)
}
