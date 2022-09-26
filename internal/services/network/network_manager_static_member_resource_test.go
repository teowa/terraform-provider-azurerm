package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/staticmembers"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkStaticMemberResource struct{}

func TestAccNetworkStaticMember_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_static_member", "test")
	r := NetworkStaticMemberResource{}
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
	data := acceptance.BuildTestData(t, "azurerm_network_static_member", "test")
	r := NetworkStaticMemberResource{}
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
	data := acceptance.BuildTestData(t, "azurerm_network_static_member", "test")
	r := NetworkStaticMemberResource{}
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
	data := acceptance.BuildTestData(t, "azurerm_network_static_member", "test")
	r := NetworkStaticMemberResource{}
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

func (r NetworkStaticMemberResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := staticmembers.ParseStaticMemberID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.StaticMembersClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkStaticMemberResource) template(data acceptance.TestData) string {
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
resource "azurerm_network_network_group" "test" {
  name                       = "acctest-nng-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r NetworkStaticMemberResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_static_member" "test" {
  name                     = "acctest-nsm-%d"
  network_network_group_id = azurerm_network_network_group.test.id
}
`, template, data.RandomInteger)
}

func (r NetworkStaticMemberResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_static_member" "import" {
  name                     = azurerm_network_static_member.test.name
  network_network_group_id = azurerm_network_network_group.test.id
}
`, config)
}

func (r NetworkStaticMemberResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_static_member" "test" {
  name                     = "acctest-nsm-%d"
  network_network_group_id = azurerm_network_network_group.test.id
  resource_id              = ""

}
`, template, data.RandomInteger)
}

func (r NetworkStaticMemberResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_static_member" "test" {
  name                     = "acctest-nsm-%d"
  network_network_group_id = azurerm_network_network_group.test.id
  resource_id              = ""

}
`, template, data.RandomInteger)
}
