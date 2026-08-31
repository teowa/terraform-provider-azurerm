// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package orbital_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbitalplanetarycomputer/2026-04-15/geocatalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type OrbitalGeoCatalogResource struct{}

func TestAccOrbitalGeoCatalog_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orbital_geocatalog", "test")
	r := OrbitalGeoCatalogResource{}

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

func TestAccOrbitalGeoCatalog_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orbital_geocatalog", "test")
	r := OrbitalGeoCatalogResource{}

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

func TestAccOrbitalGeoCatalog_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orbital_geocatalog", "test")
	r := OrbitalGeoCatalogResource{}

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

func TestAccOrbitalGeoCatalog_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orbital_geocatalog", "test")
	r := OrbitalGeoCatalogResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (OrbitalGeoCatalogResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := geocatalogs.ParseGeoCatalogID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Orbital.GeoCatalogsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (OrbitalGeoCatalogResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-orbital-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r OrbitalGeoCatalogResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_orbital_geocatalog" "test" {
  name                = "acctest-gc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func (r OrbitalGeoCatalogResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_orbital_geocatalog" "import" {
  name                = azurerm_orbital_geocatalog.test.name
  resource_group_name = azurerm_orbital_geocatalog.test.resource_group_name
  location            = azurerm_orbital_geocatalog.test.location
}
`, template)
}

func (r OrbitalGeoCatalogResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-gc-uai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_orbital_geocatalog" "test" {
  name                                    = "acctest-gc-%d"
  resource_group_name                     = azurerm_resource_group.test.name
  location                                = azurerm_resource_group.test.location
  auto_generated_domain_name_label_scope  = "NoReuse"
  tier                                    = "Basic"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  tags = {
    environment = "Production"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}
