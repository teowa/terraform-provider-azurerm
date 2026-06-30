// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	metricsobjectfirewallresources "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallMetricsResource struct{}

func TestAccPaloAltoNextGenerationFirewallMetrics_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_metrics", "test")
	r := NextGenerationFirewallMetricsResource{}

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

func TestAccPaloAltoNextGenerationFirewallMetrics_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_metrics", "test")
	r := NextGenerationFirewallMetricsResource{}

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

func TestAccPaloAltoNextGenerationFirewallMetrics_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_metrics", "test")
	r := NextGenerationFirewallMetricsResource{}

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

func TestAccPaloAltoNextGenerationFirewallMetrics_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_metrics", "test")
	r := NextGenerationFirewallMetricsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func (r NextGenerationFirewallMetricsResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := metricsobjectfirewallresources.ParseFirewallID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.MetricsObjectFirewallResources.MetricsObjectFirewallGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving metrics for %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r NextGenerationFirewallMetricsResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "test" {
  name                = "acctest-ai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  application_type    = "web"
}

resource "azurerm_palo_alto_next_generation_firewall_metrics" "test" {
  firewall_id                             = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.id
  application_insights_connection_string = azurerm_application_insights.test.connection_string
  application_insights_resource_id       = azurerm_application_insights.test.id
}
`, NextGenerationFirewallVnetResource{}.basic(data), data.RandomInteger)
}

func (r NextGenerationFirewallMetricsResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_palo_alto_next_generation_firewall_metrics" "import" {
  firewall_id                             = azurerm_palo_alto_next_generation_firewall_metrics.test.firewall_id
  application_insights_connection_string = azurerm_palo_alto_next_generation_firewall_metrics.test.application_insights_connection_string
  application_insights_resource_id       = azurerm_palo_alto_next_generation_firewall_metrics.test.application_insights_resource_id
}
`, r.basic(data))
}

func (r NextGenerationFirewallMetricsResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "test" {
  name                = "acctest-ai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  application_type    = "web"
}

resource "azurerm_palo_alto_next_generation_firewall_metrics" "test" {
  firewall_id                             = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.id
  application_insights_connection_string = azurerm_application_insights.test.connection_string
  application_insights_resource_id       = azurerm_application_insights.test.id
}
`, NextGenerationFirewallVnetResource{}.complete(data), data.RandomInteger)
}

func (r NextGenerationFirewallMetricsResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "test" {
  name                = "acctest-ai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  application_type    = "web"
}

resource "azurerm_application_insights" "alt" {
  name                = "acctest-ai-alt-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  application_type    = "web"
}

resource "azurerm_palo_alto_next_generation_firewall_metrics" "test" {
  firewall_id                             = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.test.id
  application_insights_connection_string = azurerm_application_insights.alt.connection_string
  application_insights_resource_id       = azurerm_application_insights.alt.id
}
`, NextGenerationFirewallVnetResource{}.basic(data), data.RandomInteger)
}
