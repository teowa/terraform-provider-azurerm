// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CognitiveAccountConnectionTestResource struct{}

func TestAccCognitiveAccountConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctest-conn-basic"),
				check.That(data.ResourceName).Key("auth_type").HasValue("None"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccountConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

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
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("ApiKey"),
				check.That(data.ResourceName).Key("api_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("api_key.0.key").HasValue("test-api-key-123"),
			),
		},
		data.ImportStep("api_key.0.key"),
	})
}

func TestAccCognitiveAccountConnection_accessKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.accessKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("AccessKey"),
				check.That(data.ResourceName).Key("access_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("access_key.0.access_key").HasValue("test-access-key-456"),
			),
		},
		data.ImportStep("access_key.0.access_key"),
	})
}

func TestAccCognitiveAccountConnection_accountKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.accountKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("AccountKey"),
				check.That(data.ResourceName).Key("account_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("account_key.0.account_key").HasValue("test-account-key-789"),
			),
		},
		data.ImportStep("account_key.0.account_key"),
	})
}

func TestAccCognitiveAccountConnection_managedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("ManagedIdentity"),
				check.That(data.ResourceName).Key("managed_identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("managed_identity.0.client_id").Exists(),
				check.That(data.ResourceName).Key("managed_identity.0.resource_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAccountConnection_oauth2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oauth2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("OAuth2"),
				check.That(data.ResourceName).Key("oauth2.#").HasValue("1"),
				check.That(data.ResourceName).Key("oauth2.0.client_id").Exists(),
				check.That(data.ResourceName).Key("oauth2.0.client_secret").HasValue("test-client-secret"),
			),
		},
		data.ImportStep("oauth2.0.client_secret"),
	})
}

func TestAccCognitiveAccountConnection_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.servicePrincipal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("ServicePrincipal"),
				check.That(data.ResourceName).Key("service_principal.#").HasValue("1"),
				check.That(data.ResourceName).Key("service_principal.0.client_id").Exists(),
				check.That(data.ResourceName).Key("service_principal.0.client_secret").HasValue("test-sp-secret"),
				check.That(data.ResourceName).Key("service_principal.0.tenant_id").Exists(),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccCognitiveAccountConnection_usernamePassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.usernamePassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("UsernamePassword"),
				check.That(data.ResourceName).Key("username_password.#").HasValue("1"),
				check.That(data.ResourceName).Key("username_password.0.username").HasValue("testuser"),
				check.That(data.ResourceName).Key("username_password.0.password").HasValue("TestPassword123!"),
			),
		},
		data.ImportStep("username_password.0.password"),
	})
}

func TestAccCognitiveAccountConnection_customKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("CustomKeys"),
				check.That(data.ResourceName).Key("custom_keys.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_keys.0.keys.key1").HasValue("value1"),
				check.That(data.ResourceName).Key("custom_keys.0.keys.key2").HasValue("value2"),
			),
		},
		data.ImportStep("custom_keys.0.keys"),
	})
}

func TestAccCognitiveAccountConnection_pat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pat(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("PAT"),
				check.That(data.ResourceName).Key("pat.#").HasValue("1"),
				check.That(data.ResourceName).Key("pat.0.pat").HasValue("test-personal-access-token"),
			),
		},
		data.ImportStep("pat.0.pat"),
	})
}

func TestAccCognitiveAccountConnection_sas(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sas(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("SAS"),
				check.That(data.ResourceName).Key("sas.#").HasValue("1"),
				check.That(data.ResourceName).Key("sas.0.sas").HasValue("test-shared-access-signature"),
			),
		},
		data.ImportStep("sas.0.sas"),
	})
}

func TestAccCognitiveAccountConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_type").HasValue("ApiKey"),
				check.That(data.ResourceName).Key("category").HasValue("AzureOpenAI"),
				check.That(data.ResourceName).Key("group").HasValue("DataSets"),
				check.That(data.ResourceName).Key("target").HasValue("https://example.com"),
				check.That(data.ResourceName).Key("shared_to_all").HasValue("true"),
				check.That(data.ResourceName).Key("shared_user_list.#").HasValue("2"),
				check.That(data.ResourceName).Key("metadata.environment").HasValue("test"),
				check.That(data.ResourceName).Key("metadata.version").HasValue("1.0"),
			),
		},
		data.ImportStep("api_key.0.key"),
	})
}

func TestAccCognitiveAccountConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("shared_to_all").HasValue("true"),
			),
		},
		data.ImportStep("api_key.0.key"),
		{
			Config: r.updateSharing(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("shared_to_all").HasValue("false"),
				check.That(data.ResourceName).Key("shared_user_list.#").HasValue("1"),
			),
		},
		data.ImportStep("api_key.0.key"),
	})
}

func (r CognitiveAccountConnectionTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accountconnectionresource.ParseConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cognitive.AccountConnectionResourceClient
	resp, err := client.AccountConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r CognitiveAccountConnectionTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CognitiveAccountConnectionTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-basic"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "None"
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "import" {
  name                 = azurerm_cognitive_account_connection.test.name
  cognitive_account_id = azurerm_cognitive_account_connection.test.cognitive_account_id
  auth_type           = azurerm_cognitive_account_connection.test.auth_type
}
`, config)
}

func (r CognitiveAccountConnectionTestResource) apiKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-apikey"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "ApiKey"

  api_key {
    key = "test-api-key-123"
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) accessKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-accesskey"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "AccessKey"

  access_key {
    access_key = "test-access-key-456"
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) accountKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-accountkey"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "AccountKey"

  account_key {
    account_key = "test-account-key-789"
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) managedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-mi"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "ManagedIdentity"

  managed_identity {
    client_id   = azurerm_user_assigned_identity.test.client_id
    resource_id = azurerm_user_assigned_identity.test.id
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) oauth2(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-oauth2"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "OAuth2"

  oauth2 {
    client_id     = azurerm_user_assigned_identity.test.client_id
    client_secret = "test-client-secret"
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) servicePrincipal(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "test" {}

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-sp"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "ServicePrincipal"

  service_principal {
    client_id     = azurerm_user_assigned_identity.test.client_id
    client_secret = "test-sp-secret"
    tenant_id     = data.azurerm_client_config.test.tenant_id
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) usernamePassword(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-userpass"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "UsernamePassword"

  username_password {
    username = "testuser"
    password = "TestPassword123!"
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) customKeys(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-customkeys"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "CustomKeys"

  custom_keys {
    keys = {
      key1 = "value1"
      key2 = "value2"
    }
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) pat(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-pat"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "PAT"

  pat {
    pat = "test-personal-access-token"
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) sas(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-sas"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "SAS"

  sas {
    sas = "test-shared-access-signature"
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-complete"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "ApiKey"
  category            = "AzureOpenAI"
  group               = "DataSets"
  target              = "https://example.com"
  shared_to_all       = true
  shared_user_list    = ["user1@example.com", "user2@example.com"]

  metadata = {
    environment = "test"
    version     = "1.0"
  }

  api_key {
    key = "test-api-key-complete"
  }
}
`, template)
}

func (r CognitiveAccountConnectionTestResource) updateSharing(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-complete"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type           = "ApiKey"
  category            = "AzureOpenAI"
  group               = "DataSets"
  target              = "https://example.com"
  shared_to_all       = false
  shared_user_list    = ["user1@example.com"]

  metadata = {
    environment = "test"
    version     = "1.0"
  }

  api_key {
    key = "test-api-key-complete"
  }
}
`, template)
}