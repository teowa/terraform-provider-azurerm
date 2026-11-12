// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_execution_filter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnEdgeActionExecutionFilterResource struct{}

func TestAccCdnEdgeActionExecutionFilter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_execution_filter", "test")
	r := CdnEdgeActionExecutionFilterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep(),
	})
}

func TestAccCdnEdgeActionExecutionFilter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_execution_filter", "test")
	r := CdnEdgeActionExecutionFilterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnEdgeActionExecutionFilter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_execution_filter", "test")
	r := CdnEdgeActionExecutionFilterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.complete(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep(),
	})
}

func TestAccCdnEdgeActionExecutionFilter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_execution_filter", "test")
	r := CdnEdgeActionExecutionFilterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep(),
		{Config: r.complete(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep(),
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep(),
	})
}

func (CdnEdgeActionExecutionFilterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := edgeactionexecutionfilters.ParseExecutionFilterID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Cdn.EdgeActionExecutionFiltersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (CdnEdgeActionExecutionFilterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_edge_action" "test" {
  name                = "acctest%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "global"
  sku_name            = "Standard"
  sku_tier            = "Standard"
}

resource "azurerm_cdn_edge_action_version" "test" {
  count = 2

  name               = "v${count.index}"
  cdn_edge_action_id = azurerm_cdn_edge_action.test.id
  location           = "global"
  deployment_type    = "file"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnEdgeActionExecutionFilterResource) basic(data acceptance.TestData) string {
	return r.config(data, 0, "X-Edge-Action", "one", "")
}

func (r CdnEdgeActionExecutionFilterResource) complete(data acceptance.TestData) string {
	return r.config(data, 1, "X-Edge-Version", "two", `
  tags = {
    environment = "test"
  }
`)
}

func (r CdnEdgeActionExecutionFilterResource) config(data acceptance.TestData, version int, header, value, extra string) string {
	return r.template(data) + fmt.Sprintf(`
resource "azurerm_cdn_edge_action_execution_filter" "test" {
  name                       = "filter"
  cdn_edge_action_id         = azurerm_cdn_edge_action.test.id
  cdn_edge_action_version_id = azurerm_cdn_edge_action_version.test[%d].id
  location                   = "global"
  header_name                = %q
  header_value               = %q
%s
}
`, version, header, value, extra)
}

func (r CdnEdgeActionExecutionFilterResource) requiresImport(data acceptance.TestData) string {
	return r.basic(data) + `
resource "azurerm_cdn_edge_action_execution_filter" "import" {
  name                       = azurerm_cdn_edge_action_execution_filter.test.name
  cdn_edge_action_id         = azurerm_cdn_edge_action_execution_filter.test.cdn_edge_action_id
  cdn_edge_action_version_id = azurerm_cdn_edge_action_execution_filter.test.cdn_edge_action_version_id
  location                   = azurerm_cdn_edge_action_execution_filter.test.location
  header_name                = azurerm_cdn_edge_action_execution_filter.test.header_name
  header_value               = azurerm_cdn_edge_action_execution_filter.test.header_value
}
`
}
