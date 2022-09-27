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

type ManagerAdminRuleResource struct{}

func TestAccNetworkAdminRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule", "test")
	r := ManagerAdminRuleResource{}
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

func TestAccNetworkAdminRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule", "test")
	r := ManagerAdminRuleResource{}
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

func TestAccNetworkAdminRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule", "test")
	r := ManagerAdminRuleResource{}
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

func TestAccNetworkAdminRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule", "test")
	r := ManagerAdminRuleResource{}
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

func (r ManagerAdminRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkManagerAdminRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.ManagerAdminRulesClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.RuleCollectionName, id.NetworkManagerName, id.RuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	_, ok := resp.Value.AsAdminRule()
	return utils.Bool(ok), nil
}

func (r ManagerAdminRuleResource) template(data acceptance.TestData) string {
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
resource "azurerm_network_manager_security_admin_configuration" "test" {
  name                       = "acctest-nsac-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
}
resource "azurerm_network_manager_admin_rule_collection" "test" {
  name                                    = "acctest-narc-%d"
  network_security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ManagerAdminRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_manager_admin_rule" "test" {
  name                             = "acctest-nar-%d"
  network_admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.test.id
  access                           = ""
  direction                        = ""
  kind                             = ""
  priority                         = 0
  protocol                         = ""
}
`, template, data.RandomInteger)
}

func (r ManagerAdminRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_admin_rule" "import" {
  name                             = azurerm_network_manager_admin_rule.test.name
  network_admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.test.id
  access                           = ""
  direction                        = ""
  kind                             = ""
  priority                         = 0
  protocol                         = ""
}
`, config)
}

func (r ManagerAdminRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_admin_rule" "test" {
  name                             = "acctest-nar-%d"
  network_admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.test.id
  access                           = ""
  description                      = ""
  direction                        = ""
  kind                             = ""
  priority                         = 0
  protocol                         = ""
  destination_port_ranges          = []
  source_port_ranges               = []
  destinations {
    address_prefix      = ""
    address_prefix_type = ""
  }
  sources {
    address_prefix      = ""
    address_prefix_type = ""
  }

}
`, template, data.RandomInteger)
}

func (r ManagerAdminRuleResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_admin_rule" "test" {
  name                             = "acctest-nar-%d"
  network_admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.test.id
  access                           = ""
  description                      = ""
  direction                        = ""
  kind                             = ""
  priority                         = 0
  protocol                         = ""
  destination_port_ranges          = []
  source_port_ranges               = []
  destinations {
    address_prefix      = ""
    address_prefix_type = ""
  }
  sources {
    address_prefix      = ""
    address_prefix_type = ""
  }

}
`, template, data.RandomInteger)
}
