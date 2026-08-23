package discovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storagecontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageContainerResource struct{}

func TestAccDiscoveryStorageContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_storage_container", "test")
	r := StorageContainerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryStorageContainer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_storage_container", "test")
	r := StorageContainerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.RequiresImportErrorStep(r.requiresImport)})
}
func TestAccDiscoveryStorageContainer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_storage_container", "test")
	r := StorageContainerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryStorageContainer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_storage_container", "test")
	r := StorageContainerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep(), {Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func (StorageContainerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := storagecontainers.ParseStorageContainerID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Discovery.StorageContainersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}
func (StorageContainerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-discovery-%d"
  location = "West Europe"
}

resource "azurerm_discovery_storage_container" "test" {
  name                = "acctest-storage_container-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.RandomInteger)
}
func (r StorageContainerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_discovery_storage_container" "import" {
  name                = azurerm_discovery_storage_container.test.name
  resource_group_name = azurerm_discovery_storage_container.test.resource_group_name
  location            = azurerm_discovery_storage_container.test.location
}
`, r.basic(data))
}
