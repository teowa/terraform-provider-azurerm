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

type ManagerScopeConnectionResource struct{}

func TestAccNetworkScopeConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_scope_connection", "test")
	r := ManagerScopeConnectionResource{}
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

func TestAccNetworkScopeConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_scope_connection", "test")
	r := ManagerScopeConnectionResource{}
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

func TestAccNetworkScopeConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_scope_connection", "test")
	r := ManagerScopeConnectionResource{}
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

func TestAccNetworkScopeConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_scope_connection", "test")
	r := ManagerScopeConnectionResource{}
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

func (r ManagerScopeConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkManagerScopeConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.ManagerScopeConnectionsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ScopeConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.ScopeConnectionProperties != nil), nil
}

func (r ManagerScopeConnectionResource) template(data acceptance.TestData) string {
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

func (r ManagerScopeConnectionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_manager_scope_connection" "test" {
  name                       = "acctest-nsc-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
}
`, template, data.RandomInteger)
}

func (r ManagerScopeConnectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_scope_connection" "import" {
  name                       = azurerm_network_manager_scope_connection.test.name
  network_network_manager_id = azurerm_network_network_manager.test.id
}
`, config)
}

func (r ManagerScopeConnectionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_scope_connection" "test" {
  name                       = "acctest-nsc-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
  connection_state           = ""
  description                = ""
  resource_id                = ""
  tenant_id                  = ""

}
`, template, data.RandomInteger)
}

func (r ManagerScopeConnectionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_scope_connection" "test" {
  name                       = "acctest-nsc-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
  connection_state           = ""
  description                = ""
  resource_id                = ""
  tenant_id                  = ""

}
`, template, data.RandomInteger)
}
