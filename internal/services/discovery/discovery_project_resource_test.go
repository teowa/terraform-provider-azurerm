package discovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ProjectResource struct{}

func TestAccDiscoveryProject_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_project", "test")
	r := ProjectResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryProject_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_project", "test")
	r := ProjectResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.RequiresImportErrorStep(r.requiresImport)})
}
func TestAccDiscoveryProject_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_project", "test")
	r := ProjectResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func TestAccDiscoveryProject_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_project", "test")
	r := ProjectResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep(), {Config: r.basic(data), Check: acceptance.ComposeTestCheckFunc(check.That(data.ResourceName).ExistsInAzure(r))}, data.ImportStep()})
}
func (ProjectResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := projects.ParseProjectID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Discovery.ProjectsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}
func (ProjectResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-discovery-%d"
  location = "West Europe"
}

resource "azurerm_discovery_workspace" "test" {
  name                = "acctest-workspace-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_discovery_project" "test" {
  name                    = "acctest-project-%d"
  discovery_workspace_id  = azurerm_discovery_workspace.test.id
  location                = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
func (r ProjectResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_discovery_project" "import" {
  name                   = azurerm_discovery_project.test.name
  discovery_workspace_id = azurerm_discovery_project.test.discovery_workspace_id
  location               = azurerm_discovery_project.test.location
}
`, r.basic(data))
}
