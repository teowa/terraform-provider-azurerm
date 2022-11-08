package privatednsresolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PrivateDNSResolverForwardingRuleResource struct{}

func TestAccPrivateDNSResolverForwardingRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_forwarding_rule", "test")
	r := PrivateDNSResolverForwardingRuleResource{}
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

func TestAccPrivateDNSResolverForwardingRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_forwarding_rule", "test")
	r := PrivateDNSResolverForwardingRuleResource{}
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

func TestAccPrivateDNSResolverForwardingRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_forwarding_rule", "test")
	r := PrivateDNSResolverForwardingRuleResource{}
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

func TestAccPrivateDNSResolverForwardingRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_forwarding_rule", "test")
	r := PrivateDNSResolverForwardingRuleResource{}
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

func (r PrivateDNSResolverForwardingRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := forwardingrules.ParseForwardingRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PrivateDnsResolver.ForwardingRulesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r PrivateDNSResolverForwardingRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-rg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "outbounddns"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.64/28"]

  delegation {
    name = "Microsoft.Network.dnsResolvers"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.Network/dnsResolvers"
    }
  }
}

resource "azurerm_private_dns_resolver" "test" {
  name                = "acctest-dr-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_network_id  = azurerm_virtual_network.test.id
}

resource "azurerm_private_dns_resolver_outbound_endpoint" "test" {
  name                    = "acctest-droe-%[2]d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location
  subnet_id               = azurerm_subnet.test.id
}

resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
  name                                       = "acctest-drdfr-%[2]d"
  resource_group_name                        = azurerm_resource_group.test.name
  location                                   = azurerm_resource_group.test.location
  private_dns_resolver_outbound_endpoint_ids = [azurerm_private_dns_resolver_outbound_endpoint.test.id]
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r PrivateDNSResolverForwardingRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_private_dns_resolver_forwarding_rule" "test" {
  name                      = "acctest-drfr-%d"
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
  domain_name               = "onprem.local."
  target_dns_servers {
    ip_address = "10.10.0.1"
    port       = 53
  }
}
`, template, data.RandomInteger)
}

func (r PrivateDNSResolverForwardingRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_forwarding_rule" "import" {
  name                      = azurerm_private_dns_resolver_forwarding_rule.test.name
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
  domain_name               = azurerm_private_dns_resolver_forwarding_rule.test.domain_name
  target_dns_servers {
    ip_address = "10.10.0.1"
    port       = 53
  }
}
`, config)
}

func (r PrivateDNSResolverForwardingRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_forwarding_rule" "test" {
  name                      = "acctest-drfr-%d"
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
  domain_name               = "onprem.local."
  enabled                   = true
  target_dns_servers {
    ip_address = "10.10.0.1"
    port       = 53
  }
  metadata = {
    key = "value"
  }
}
`, template, data.RandomInteger)
}

func (r PrivateDNSResolverForwardingRuleResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_forwarding_rule" "test" {
  name                      = "acctest-drfr-%d"
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
  domain_name               = "onprem.local."
  enabled                   = false
  target_dns_servers {
    ip_address = "10.10.0.2"
    port       = 53
  }
  metadata = {
    key = "updated value"
  }
}
`, template, data.RandomInteger)
}
