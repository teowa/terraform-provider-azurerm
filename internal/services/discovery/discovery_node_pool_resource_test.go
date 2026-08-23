package discovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/nodepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NodePoolResource struct{}

func TestAccDiscoveryNodePool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_node_pool", "test")
	r := NodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryNodePool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_node_pool", "test")
	r := NodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.RequiresImportErrorStep(r.requiresImport)})
}
func TestAccDiscoveryNodePool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_node_pool", "test")
	r := NodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryNodePool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_node_pool", "test")
	r := NodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep(), {Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func (NodePoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := nodepools.ParseNodePoolID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Discovery.NodePoolsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}
func (NodePoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-discovery-%d"
  location = "West Europe"
}

resource "azurerm_discovery_supercomputer" "test" {
  name                = "acctest-super-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_discovery_node_pool" "test" {
  name                        = "acctest-nodepool-%d"
  discovery_supercomputer_id  = azurerm_discovery_supercomputer.test.id
  location                    = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
func (r NodePoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_discovery_node_pool" "import" {
  name                       = azurerm_discovery_node_pool.test.name
  discovery_supercomputer_id = azurerm_discovery_node_pool.test.discovery_supercomputer_id
  location                   = azurerm_discovery_node_pool.test.location
}
`, r.basic(data))
}
