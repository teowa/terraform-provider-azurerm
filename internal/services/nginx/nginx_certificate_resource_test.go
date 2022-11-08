package nginx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxcertificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CertificateResource struct{}

func (a CertificateResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := nginxcertificate.ParseCertificateID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Nginx.NginxCertificate.CertificatesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Certificate %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func TestAccCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.CertificateResource{}.ResourceType(), "test")
	r := CertificateResource{}
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

func TestAccCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.CertificateResource{}.ResourceType(), "test")
	r := CertificateResource{}
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

func (a CertificateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_certificate" "test" {
  name                     = "acctest%[2]d"
  nginx_deployment_id      = azurerm_nginx_deployment.test.id
  key_virtual_path         = "/src/cert/soservermekey.key"
  certificate_virtual_path = "/src/cert/server.cert"
  key_vault_secret_id      = azurerm_key_vault_secret.test.id
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a CertificateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nginx_certificate" "import" {
  name                     = azurerm_nginx_certificate.test.name
  nginx_deployment_id      = azurerm_nginx_certificate.test.nginx_deployment_id
  key_virtual_path         = azurerm_nginx_certificate.test.key_virtual_path
  certificate_virtual_path = azurerm_nginx_certificate.test.certificate_virtual_path
  key_vault_secret_id      = azurerm_nginx_certificate.test.key_vault_secret_id
}
`, a.basic(data))
}

func (a CertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "subnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "delegation"

    service_delegation {
      name = "NGINX.NGINXPLUS/nginxDeployments"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_nginx_deployment" "test" {
  name                     = "acctest-%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "publicpreview_Monthly_gmz7xq9ge3py"
  location                 = azurerm_resource_group.test.location
  diagnose_support_enabled = true

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "List",
      "Purge",
      "Recover",
      "Set",
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
