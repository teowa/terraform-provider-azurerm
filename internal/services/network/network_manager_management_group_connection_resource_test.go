package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerManagementGroupConnectionResource struct{}

func TestAccNetworkManagementGroupNetworkManagerConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_management_group_connection", "test")
	r := ManagerManagementGroupConnectionResource{}
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
	data := acceptance.BuildTestData(t, "azurerm_network_manager_management_group_connection", "test")
	r := ManagerManagementGroupConnectionResource{}
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
	data := acceptance.BuildTestData(t, "azurerm_network_manager_management_group_connection", "test")
	r := ManagerManagementGroupConnectionResource{}
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
	data := acceptance.BuildTestData(t, "azurerm_network_manager_management_group_connection", "test")
	r := ManagerManagementGroupConnectionResource{}
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

func (r ManagerManagementGroupConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkManagerManagementGroupConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.ManagerManagementGrpConnectionsClient
	resp, err := client.Get(ctx, id.ManagementGroupName, id.NetworkManagerConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.ManagerConnectionProperties != nil), nil
}

func (r ManagerManagementGroupConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
			provider "azurerm" {
				features {}
			}

			
`)
}

func (r ManagerManagementGroupConnectionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_manager_management_group_connection" "test" {
  name                = "acctest-nmgnmc-%d"
  management_group_id = ""
}
`, template, data.RandomInteger)
}

func (r ManagerManagementGroupConnectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_management_group_connection" "import" {
  name                = azurerm_network_manager_management_group_connection.test.name
  management_group_id = ""
}
`, config)
}

func (r ManagerManagementGroupConnectionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_management_group_connection" "test" {
  name                = "acctest-nmgnmc-%d"
  management_group_id = ""
  connection_state    = ""
  description         = ""
  network_manager_id  = ""

}
`, template, data.RandomInteger)
}

func (r ManagerManagementGroupConnectionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_management_group_connection" "test" {
  name                = "acctest-nmgnmc-%d"
  management_group_id = ""
  connection_state    = ""
  description         = ""
  network_manager_id  = ""

}
`, template, data.RandomInteger)
}
