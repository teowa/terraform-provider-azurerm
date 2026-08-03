// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ElasticCloudHostedSearchResource struct{}

func TestAccElasticCloudHostedSearch_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_hosted_search", "test")
	r := ElasticCloudHostedSearchResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("elastic_cloud_deployment_id").Exists(),
				check.That(data.ResourceName).Key("elastic_cloud_sso_default_url").Exists(),
				check.That(data.ResourceName).Key("elastic_cloud_user_id").Exists(),
				check.That(data.ResourceName).Key("elasticsearch_service_url").Exists(),
				check.That(data.ResourceName).Key("kibana_service_url").Exists(),
				check.That(data.ResourceName).Key("kibana_sso_uri").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticCloudHostedSearch_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_hosted_search", "test")
	r := ElasticCloudHostedSearchResource{}
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

func TestAccElasticCloudHostedSearch_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_hosted_search", "test")
	r := ElasticCloudHostedSearchResource{}
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

func TestAccElasticCloudHostedSearch_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_hosted_search", "test")
	r := ElasticCloudHostedSearchResource{}
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

func (ElasticCloudHostedSearchResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := elasticmonitorresources.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Elastic.HostedSearchMonitorClient.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil || resp.Model.Kind == nil || *resp.Model.Kind != "elastic-hosted-deployment" {
		return pointer.To(false), nil
	}

	return pointer.To(true), nil
}

func (ElasticCloudHostedSearchResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_cloud_hosted_search" "test" {
  name                        = "acctesths-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ElasticCloudHostedSearchResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_cloud_hosted_search" "import" {
  name                        = azurerm_elastic_cloud_hosted_search.test.name
  resource_group_name         = azurerm_elastic_cloud_hosted_search.test.resource_group_name
  location                    = azurerm_elastic_cloud_hosted_search.test.location
  sku_name                    = azurerm_elastic_cloud_hosted_search.test.sku_name
  elastic_cloud_email_address = azurerm_elastic_cloud_hosted_search.test.elastic_cloud_email_address
}
`, r.basic(data))
}

func (ElasticCloudHostedSearchResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_cloud_hosted_search" "test" {
  name                        = "acctesths-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"
  monitoring_enabled          = false

  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ElasticCloudHostedSearchResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_cloud_hosted_search" "test" {
  name                        = "acctesths-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"

  tags = {
    environment = "updated"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
