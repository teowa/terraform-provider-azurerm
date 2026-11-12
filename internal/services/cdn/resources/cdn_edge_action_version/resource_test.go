// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_version_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnEdgeActionVersionResource struct{}

func TestAccCdnEdgeActionVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_version", "test")
	r := CdnEdgeActionVersionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep(),
	})
}

func TestAccCdnEdgeActionVersion_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_version", "test")
	r := CdnEdgeActionVersionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnEdgeActionVersion_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_version", "test")
	r := CdnEdgeActionVersionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.complete(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep("code"),
	})
}

func TestAccCdnEdgeActionVersion_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_version", "test")
	r := CdnEdgeActionVersionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep(),
		{Config: r.complete(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep("code"),
		{Config: r.updatedCode(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep("code"),
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep(),
	})
}

func TestAccCdnEdgeActionVersion_swapDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_version", "test")
	r := CdnEdgeActionVersionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.swapDefault(data, false), Check: check.That(data.ResourceName).Key("default_version_enabled").HasValue("true")},
		data.ImportStep(),
		{
			Config: r.swapDefault(data, true),
			Check:  check.That("azurerm_cdn_edge_action_version.other").Key("default_version_enabled").HasValue("true"),
		},
		{
			Config: r.swapDefault(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("default_version_enabled").HasValue("false"),
				check.That("azurerm_cdn_edge_action_version.other").Key("default_version_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func (CdnEdgeActionVersionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := edgeactionversions.ParseVersionID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Cdn.EdgeActionVersionsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (CdnEdgeActionVersionResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnEdgeActionVersionResource) basic(data acceptance.TestData) string {
	return r.template(data) + `
resource "azurerm_cdn_edge_action_version" "test" {
  name               = "v1"
  cdn_edge_action_id = azurerm_cdn_edge_action.test.id
  location           = "global"
  deployment_type    = "file"
}
`
}

func (r CdnEdgeActionVersionResource) complete(data acceptance.TestData) string {
	return r.withCode(data, `addEventListener("fetch", event => event.respondWith(new Response("hello")));`)
}

func (r CdnEdgeActionVersionResource) updatedCode(data acceptance.TestData) string {
	return r.withCode(data, `addEventListener("fetch", event => event.respondWith(new Response("updated")));`)
}

func (r CdnEdgeActionVersionResource) withCode(data acceptance.TestData, code string) string {
	return r.template(data) + fmt.Sprintf(`
resource "azurerm_cdn_edge_action_version" "test" {
  name                    = "v1"
  cdn_edge_action_id      = azurerm_cdn_edge_action.test.id
  location                = "global"
  deployment_type         = "file"
  default_version_enabled = true

  code {
    file_name      = "index.js"
    content_base64 = base64encode(%q)
  }

  tags = {
    environment = "test"
  }
}
`, code)
}

func (r CdnEdgeActionVersionResource) swapDefault(data acceptance.TestData, promote bool) string {
	promotion := ""
	if promote {
		promotion = "  default_version_enabled = true\n"
	}
	return r.basic(data) + fmt.Sprintf(`
resource "azurerm_cdn_edge_action_version" "other" {
  name               = "v2"
  cdn_edge_action_id = azurerm_cdn_edge_action.test.id
  location           = "global"
  deployment_type    = "file"
%s
  depends_on = [azurerm_cdn_edge_action_version.test]
}
`, promotion)
}

func (r CdnEdgeActionVersionResource) requiresImport(data acceptance.TestData) string {
	return r.basic(data) + `
resource "azurerm_cdn_edge_action_version" "import" {
  name               = azurerm_cdn_edge_action_version.test.name
  cdn_edge_action_id = azurerm_cdn_edge_action_version.test.cdn_edge_action_id
  location           = azurerm_cdn_edge_action_version.test.location
  deployment_type    = azurerm_cdn_edge_action_version.test.deployment_type
}
`
}
