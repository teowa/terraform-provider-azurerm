package discovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storageassets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageAssetResource struct{}

func TestAccDiscoveryStorageAsset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_storage_asset", "test")
	r := StorageAssetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryStorageAsset_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_storage_asset", "test")
	r := StorageAssetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.RequiresImportErrorStep(r.requiresImport)})
}
func TestAccDiscoveryStorageAsset_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_storage_asset", "test")
	r := StorageAssetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryStorageAsset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_storage_asset", "test")
	r := StorageAssetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep(), {Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func (StorageAssetResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := storageassets.ParseStorageAssetID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Discovery.StorageAssetsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}
func (StorageAssetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-discovery-%d"
  location = "West Europe"
}

resource "azurerm_discovery_storage_container" "test" {
  name                = "acctest-container-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_discovery_storage_asset" "test" {
  name                            = "acctest-asset-%d"
  discovery_storage_container_id  = azurerm_discovery_storage_container.test.id
  location                        = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
func (r StorageAssetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_discovery_storage_asset" "import" {
  name                           = azurerm_discovery_storage_asset.test.name
  discovery_storage_container_id = azurerm_discovery_storage_asset.test.discovery_storage_container_id
  location                       = azurerm_discovery_storage_asset.test.location
}
`, r.basic(data))
}
