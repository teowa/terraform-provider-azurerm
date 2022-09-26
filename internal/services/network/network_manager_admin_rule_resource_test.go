package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/adminrules"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkAdminRuleResource struct{}

func TestAccNetworkAdminRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_admin_rule", "test")
	r := NetworkAdminRuleResource{}
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
	data := acceptance.BuildTestData(t, "azurerm_network_admin_rule", "test")
	r := NetworkAdminRuleResource{}
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
	data := acceptance.BuildTestData(t, "azurerm_network_admin_rule", "test")
	r := NetworkAdminRuleResource{}
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
	data := acceptance.BuildTestData(t, "azurerm_network_admin_rule", "test")
	r := NetworkAdminRuleResource{}
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

func (r NetworkAdminRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := adminrules.ParseRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.AdminRulesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkAdminRuleResource) template(data acceptance.TestData) string {
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
resource "azurerm_network_security_admin_configuration" "test" {
  name                       = "acctest-nsac-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
}
resource "azurerm_network_admin_rule_collection" "test" {
  name                                    = "acctest-narc-%d"
  network_security_admin_configuration_id = azurerm_network_security_admin_configuration.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r NetworkAdminRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_admin_rule" "test" {
  name                             = "acctest-nar-%d"
  network_admin_rule_collection_id = azurerm_network_admin_rule_collection.test.id
  kind                             = ""
}
`, template, data.RandomInteger)
}

func (r NetworkAdminRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_admin_rule" "import" {
  name                             = azurerm_network_admin_rule.test.name
  network_admin_rule_collection_id = azurerm_network_admin_rule_collection.test.id
  kind                             = ""
}
`, config)
}

func (r NetworkAdminRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_admin_rule" "test" {
  name                             = "acctest-nar-%d"
  network_admin_rule_collection_id = azurerm_network_admin_rule_collection.test.id
  kind                             = ""

}
`, template, data.RandomInteger)
}

func (r NetworkAdminRuleResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_admin_rule" "test" {
  name                             = "acctest-nar-%d"
  network_admin_rule_collection_id = azurerm_network_admin_rule_collection.test.id
  kind                             = ""

}
`, template, data.RandomInteger)
}
