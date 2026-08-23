package discovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/tools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ToolResource struct{}

func TestAccDiscoveryTool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_tool", "test")
	r := ToolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryTool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_tool", "test")
	r := ToolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.RequiresImportErrorStep(r.requiresImport)})
}
func TestAccDiscoveryTool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_tool", "test")
	r := ToolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryTool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_tool", "test")
	r := ToolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep(), {Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func (ToolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tools.ParseToolID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Discovery.ToolsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}
func (ToolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-discovery-%d"
  location = "West Europe"
}

resource "azurerm_discovery_tool" "test" {
  name                = "acctest-tool-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.RandomInteger)
}
func (r ToolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_discovery_tool" "import" {
  name                = azurerm_discovery_tool.test.name
  resource_group_name = azurerm_discovery_tool.test.resource_group_name
  location            = azurerm_discovery_tool.test.location
}
`, r.basic(data))
}
