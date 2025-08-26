// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CognitiveAccountConnectionResource struct{}

func TestAccCognitiveAccountConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("None"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccountConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionResource{}

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

func TestAccCognitiveAccountConnection_apiKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("ApiKey"),
				check.That(data.ResourceName).Key("api_key").HasValue("test-key"),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccCognitiveAccountConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("None"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("None"),
				check.That(data.ResourceName).Key("target").HasValue("updated-target"),
			),
		},
		data.ImportStep(),
	})
}

func (r CognitiveAccountConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accountconnectionresource.ParseConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cognitive.AccountConnectionsClient
	resp, err := client.AccountConnectionsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r CognitiveAccountConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CognitiveAccountConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                  = "acctest-conn-%d"
  cognitive_account_id  = azurerm_cognitive_account.test.id
  auth_type            = "None"
  category             = "OpenAI"
  target               = "test-target"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "import" {
  name                  = azurerm_cognitive_account_connection.test.name
  cognitive_account_id  = azurerm_cognitive_account_connection.test.cognitive_account_id
  auth_type            = azurerm_cognitive_account_connection.test.auth_type
  category             = azurerm_cognitive_account_connection.test.category
  target               = azurerm_cognitive_account_connection.test.target
}
`, r.basic(data))
}

func (r CognitiveAccountConnectionResource) basicUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                  = "acctest-conn-%d"
  cognitive_account_id  = azurerm_cognitive_account.test.id
  auth_type            = "None"
  category             = "OpenAI"
  target               = "updated-target"
  is_shared_to_all     = true
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionResource) apiKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                  = "acctest-conn-%d"
  cognitive_account_id  = azurerm_cognitive_account.test.id
  auth_type            = "ApiKey"
  category             = "OpenAI"
  target               = "test-target"
  api_key              = "test-key"
}
`, r.template(data), data.RandomInteger)
}
