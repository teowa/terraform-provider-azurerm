// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceApiResource struct{}

func TestAccApiManagementWorkspaceApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiResource{}

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

func TestAccApiManagementWorkspaceApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiResource{}

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

func TestAccApiManagementWorkspaceApi_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiResource{}

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

func TestAccApiManagementWorkspaceApi_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiResource{}

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

func TestAccApiManagementWorkspaceApi_importSwagger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importSwagger(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
	})
}

func TestAccApiManagementWorkspaceApi_oauth2Authorization(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oauth2Authorization(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceApi_openidAuthentication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.openidAuthentication(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementWorkspaceApiResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := api.ParseWorkspaceApiID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.WorkspaceApiClient.WorkspaceApiGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceApiResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  revision                    = "1"
  display_name                = "Test API"
  path                        = "test"
  protocols                   = ["https"]
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_api" "import" {
  name                        = azurerm_api_management_workspace_api.test.name
  api_management_workspace_id = azurerm_api_management_workspace_api.test.api_management_workspace_id
  revision                    = azurerm_api_management_workspace_api.test.revision
  display_name                = azurerm_api_management_workspace_api.test.display_name
  path                        = azurerm_api_management_workspace_api.test.path
  protocols                   = azurerm_api_management_workspace_api.test.protocols
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceApiResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  revision                    = "1"
  display_name                = "Test API Updated"
  path                        = "test-updated"
  protocols                   = ["https", "http"]
  description                 = "This is an updated test API"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  revision                    = "1"
  display_name                = "Complete Test API"
  path                        = "test"
  protocols                   = ["https"]
  description                 = "This is a complete test API"
  service_url                 = "https://example.com/api"
  subscription_required       = true
  terms_of_service_url        = "https://example.com/terms"
  revision_description        = "Test Revision"

  contact {
    name  = "Test Contact"
    email = "test@example.com"
    url   = "https://example.com"
  }

  license {
    name = "MIT"
    url  = "https://opensource.org/licenses/MIT"
  }

  subscription_key_parameter_names {
    header = "Ocp-Apim-Subscription-Key"
    query  = "subscription-key"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiResource) importSwagger(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  revision                    = "1"
  display_name                = "Swagger Test API"
  path                        = "test"
  protocols                   = ["https"]

  import {
    content_format = "openapi+json"
    content_value  = jsonencode({
      openapi = "3.0.1"
      info = {
        title   = "Test API"
        version = "1.0.0"
      }
      paths = {
        "/test" = {
          get = {
            responses = {
              "200" = {
                description = "Success"
              }
            }
          }
        }
      }
    })
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiResource) oauth2Authorization(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauth-%d"
  api_management_name          = azurerm_api_management.test.name
  resource_group_name          = azurerm_resource_group.test.name
  display_name                 = "Test Authorization Server"
  authorization_endpoint       = "https://example.com/oauth/authorize"
  token_endpoint               = "https://example.com/oauth/token"
  client_id                    = "test-client-id"
  grant_types                  = ["authorizationCode"]
  authorization_methods        = ["GET"]
  client_registration_endpoint = "https://example.com/oauth/register"
}

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  revision                    = "1"
  display_name                = "OAuth2 Test API"
  path                        = "test"
  protocols                   = ["https"]

  oauth2_authorization {
    authorization_server_name = azurerm_api_management_authorization_server.test.name
    scope                     = "read write"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceApiResource) openidAuthentication(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctestopenid-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  display_name        = "Test OpenID Connect Provider"
  client_id           = "test-client-id"
  metadata_endpoint   = "https://example.com/.well-known/openid-configuration"
}

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  revision                    = "1"
  display_name                = "OpenID Test API"
  path                        = "test"
  protocols                   = ["https"]

  openid_authentication {
    openid_provider_name          = azurerm_api_management_openid_connect_provider.test.name
    bearer_token_sending_methods  = ["authorizationHeader"]
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (ApiManagementWorkspaceApiResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestws%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
