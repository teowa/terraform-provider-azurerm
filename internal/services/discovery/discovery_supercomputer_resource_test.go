package discovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/supercomputers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SupercomputerResource struct{}

func TestAccDiscoverySupercomputer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_supercomputer", "test")
	r := SupercomputerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoverySupercomputer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_supercomputer", "test")
	r := SupercomputerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.RequiresImportErrorStep(r.requiresImport)})
}
func TestAccDiscoverySupercomputer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_supercomputer", "test")
	r := SupercomputerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoverySupercomputer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_supercomputer", "test")
	r := SupercomputerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep(), {Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func (SupercomputerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := supercomputers.ParseSupercomputerID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Discovery.SupercomputersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}
func (SupercomputerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-discovery-%d"
  location = "West Europe"
}

resource "azurerm_discovery_supercomputer" "test" {
  name                = "acctest-supercomputer-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.RandomInteger)
}
func (r SupercomputerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_discovery_supercomputer" "import" {
  name                = azurerm_discovery_supercomputer.test.name
  resource_group_name = azurerm_discovery_supercomputer.test.resource_group_name
  location            = azurerm_discovery_supercomputer.test.location
}
`, r.basic(data))
}
