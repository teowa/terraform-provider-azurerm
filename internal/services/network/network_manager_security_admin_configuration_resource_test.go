package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/securityadminconfigurations"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkSecurityAdminConfigurationResource struct{}

func TestAccNetworkSecurityAdminConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_admin_configuration", "test")
	r := NetworkSecurityAdminConfigurationResource{}
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

func TestAccNetworkSecurityAdminConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_admin_configuration", "test")
	r := NetworkSecurityAdminConfigurationResource{}
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

func TestAccNetworkSecurityAdminConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_admin_configuration", "test")
	r := NetworkSecurityAdminConfigurationResource{}
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

func TestAccNetworkSecurityAdminConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_admin_configuration", "test")
	r := NetworkSecurityAdminConfigurationResource{}
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

func (r NetworkSecurityAdminConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.SecurityAdminConfigurationsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkSecurityAdminConfigurationResource) template(data acceptance.TestData) string {
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

func (r NetworkSecurityAdminConfigurationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_security_admin_configuration" "test" {
  name                       = "acctest-nsac-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
}
`, template, data.RandomInteger)
}

func (r NetworkSecurityAdminConfigurationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_security_admin_configuration" "import" {
  name                       = azurerm_network_security_admin_configuration.test.name
  network_network_manager_id = azurerm_network_network_manager.test.id
}
`, config)
}

func (r NetworkSecurityAdminConfigurationResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_security_admin_configuration" "test" {
  name                       = "acctest-nsac-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
  description                = ""
  apply_on_network_intent_policy_based_services {

  }

}
`, template, data.RandomInteger)
}

func (r NetworkSecurityAdminConfigurationResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_security_admin_configuration" "test" {
  name                       = "acctest-nsac-%d"
  network_network_manager_id = azurerm_network_network_manager.test.id
  description                = ""
  apply_on_network_intent_policy_based_services {

  }

}
`, template, data.RandomInteger)
}
