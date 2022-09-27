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

type NetworkAdminRuleCollectionResource struct{}

func TestAccNetworkAdminRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule_collection", "test")
	r := NetworkAdminRuleCollectionResource{}
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

func TestAccNetworkAdminRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule_collection", "test")
	r := NetworkAdminRuleCollectionResource{}
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

func TestAccNetworkAdminRuleCollection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule_collection", "test")
	r := NetworkAdminRuleCollectionResource{}
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

func TestAccNetworkAdminRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule_collection", "test")
	r := NetworkAdminRuleCollectionResource{}
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

func (r NetworkAdminRuleCollectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkManagerAdminRuleCollectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.ManagerAdminRuleCollectionsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.RuleCollectionName, id.NetworkManagerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.AdminRuleCollectionPropertiesFormat != nil), nil
}

func (r NetworkAdminRuleCollectionResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r NetworkAdminRuleCollectionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_manager_admin_rule_collection" "test" {
  name                                    = "acctest-narc-%d"
  network_security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.test.id
  applies_to_groups {
    network_group_id = ""
  }
}
`, template, data.RandomInteger)
}

func (r NetworkAdminRuleCollectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_admin_rule_collection" "import" {
  name                                    = azurerm_network_manager_admin_rule_collection.test.name
  network_security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.test.id
  applies_to_groups {
    network_group_id = ""
  }
}
`, config)
}

func (r NetworkAdminRuleCollectionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_admin_rule_collection" "test" {
  name                                    = "acctest-narc-%d"
  network_security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.test.id
  description                             = ""
  applies_to_groups {
    network_group_id = ""
  }

}
`, template, data.RandomInteger)
}

func (r NetworkAdminRuleCollectionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_manager_admin_rule_collection" "test" {
  name                                    = "acctest-narc-%d"
  network_security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.test.id
  description                             = ""
  applies_to_groups {
    network_group_id = ""
  }

}
`, template, data.RandomInteger)
}
