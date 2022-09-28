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

type ManagerStaticMemberResource struct{}

func TestAccNetworkStaticMember_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_static_member", "test")
	r := ManagerStaticMemberResource{}
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

func TestAccNetworkStaticMember_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_static_member", "test")
	r := ManagerStaticMemberResource{}
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

func TestAccNetworkStaticMember_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_static_member", "test")
	r := ManagerStaticMemberResource{}
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

func TestAccNetworkStaticMember_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_static_member", "test")
	r := ManagerStaticMemberResource{}
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

func (r ManagerStaticMemberResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkManagerStaticMemberID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.ManagerStaticMembersClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.StaticMemberProperties != nil), nil
}

func (r ManagerStaticMemberResource) template(data acceptance.TestData) string {
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
resource "azurerm_network_manager_network_group" "test" {
  name                       = "acctest-nng-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ManagerStaticMemberResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_manager_static_member" "test" {
  name                     = "acctest-nsm-%d"
  network_network_group_id = azurerm_network_manager_network_group.test.id
}
`, template, data.RandomInteger)
}

func (r ManagerStaticMemberResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_static_member" "import" {
  name                     = azurerm_network_manager_static_member.test.name
  network_network_group_id = azurerm_network_manager_network_group.test.id
}
`, config)
}

func (r ManagerStaticMemberResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_static_member" "test" {
  name                     = "acctest-nsm-%d"
  network_network_group_id = azurerm_network_manager_network_group.test.id
  resource_id              = ""

}
`, template, data.RandomInteger)
}

func (r ManagerStaticMemberResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_static_member" "test" {
  name                     = "acctest-nsm-%d"
  network_network_group_id = azurerm_network_manager_network_group.test.id
  resource_id              = ""

}
`, template, data.RandomInteger)
}
