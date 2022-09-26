package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/networkgroups"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkNetworkGroupResource struct{}

func TestAccNetworkNetworkGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_network_group", "test")
	r := NetworkNetworkGroupResource{}
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

func TestAccNetworkNetworkGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_network_group", "test")
	r := NetworkNetworkGroupResource{}
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

func TestAccNetworkNetworkGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_network_group", "test")
	r := NetworkNetworkGroupResource{}
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

func TestAccNetworkNetworkGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_network_group", "test")
	r := NetworkNetworkGroupResource{}
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

func (r NetworkNetworkGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkgroups.ParseNetworkGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.NetworkGroupsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkNetworkGroupResource) template(data acceptance.TestData) string {
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

func (r NetworkNetworkGroupResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_network_group" "test" {
  name                       = "acctest-nng-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
}
`, template, data.RandomInteger)
}

func (r NetworkNetworkGroupResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_network_group" "import" {
  name                       = azurerm_network_network_group.test.name
  network_network_manager_id = azurerm_network_network_manager.test.id
}
`, config)
}

func (r NetworkNetworkGroupResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_network_group" "test" {
  name                       = "acctest-nng-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
  description                = ""

}
`, template, data.RandomInteger)
}

func (r NetworkNetworkGroupResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_network_group" "test" {
  name                       = "acctest-nng-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
  description                = ""

}
`, template, data.RandomInteger)
}
