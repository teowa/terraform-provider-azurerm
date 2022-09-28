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

type ManagerSubscriptionConnectionResource struct{}

func TestAccNetworkSubscriptionNetworkManagerConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_subscription_connection", "test")
	r := ManagerSubscriptionConnectionResource{}
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

func TestAccNetworkSubscriptionNetworkManagerConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_subscription_connection", "test")
	r := ManagerSubscriptionConnectionResource{}
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

func TestAccNetworkSubscriptionNetworkManagerConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_subscription_connection", "test")
	r := ManagerSubscriptionConnectionResource{}
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

func TestAccNetworkSubscriptionNetworkManagerConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_subscription_connection", "test")
	r := ManagerSubscriptionConnectionResource{}
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

func (r ManagerSubscriptionConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkManagerSubscriptionConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.ManagerSubscriptionConnectionsClient
	resp, err := client.Get(ctx, id.SubscriptionId, id.NetworkManagerConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.ManagerConnectionProperties != nil), nil
}

func (r ManagerSubscriptionConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
			provider "azurerm" {
				features {}
			}

			
`)
}

func (r ManagerSubscriptionConnectionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_manager_subscription_connection" "test" {
  name = "acctest-nsnmc-%d"
}
`, template, data.RandomInteger)
}

func (r ManagerSubscriptionConnectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_subscription_connection" "import" {
  name = azurerm_network_manager_subscription_connection.test.name
}
`, config)
}

func (r ManagerSubscriptionConnectionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_subscription_connection" "test" {
  name               = "acctest-nsnmc-%d"
  connection_state   = ""
  description        = ""
  network_manager_id = ""

}
`, template, data.RandomInteger)
}

func (r ManagerSubscriptionConnectionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_subscription_connection" "test" {
  name               = "acctest-nsnmc-%d"
  connection_state   = ""
  description        = ""
  network_manager_id = ""

}
`, template, data.RandomInteger)
}
