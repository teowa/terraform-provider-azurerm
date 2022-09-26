package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/networkmanagerconnections"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkManagementGroupNetworkManagerConnectionResource struct{}

func TestAccNetworkManagementGroupNetworkManagerConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_management_group_network_manager_connection", "test")
	r := NetworkManagementGroupNetworkManagerConnectionResource{}
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

func TestAccNetworkManagementGroupNetworkManagerConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_management_group_network_manager_connection", "test")
	r := NetworkManagementGroupNetworkManagerConnectionResource{}
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

func TestAccNetworkManagementGroupNetworkManagerConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_management_group_network_manager_connection", "test")
	r := NetworkManagementGroupNetworkManagerConnectionResource{}
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

func TestAccNetworkManagementGroupNetworkManagerConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_management_group_network_manager_connection", "test")
	r := NetworkManagementGroupNetworkManagerConnectionResource{}
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

func (r NetworkManagementGroupNetworkManagerConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkmanagerconnections.ParseProviders2NetworkManagerConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.NetworkManagerConnectionsClient
	resp, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
			provider "azurerm" {
				features {}
			}

			
`)
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_management_group_network_manager_connection" "test" {
  name                = "acctest-nmgnmc-%d"
  management_group_id = ""
}
`, template, data.RandomInteger)
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_management_group_network_manager_connection" "import" {
  name                = azurerm_network_management_group_network_manager_connection.test.name
  management_group_id = ""
}
`, config)
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_management_group_network_manager_connection" "test" {
  name                = "acctest-nmgnmc-%d"
  management_group_id = ""
  connection_state    = ""
  description         = ""
  network_manager_id  = ""

}
`, template, data.RandomInteger)
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_management_group_network_manager_connection" "test" {
  name                = "acctest-nmgnmc-%d"
  management_group_id = ""
  connection_state    = ""
  description         = ""
  network_manager_id  = ""

}
`, template, data.RandomInteger)
}
