package elastic_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccElasticCloudHostedSearch_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_hosted_search", "test")
	r := ElasticCloudHostedSearchResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_elastic_cloud_hosted_search.list", 1),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_elastic_cloud_hosted_search.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_elastic_cloud_hosted_search.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (ElasticCloudHostedSearchResource) basicListQuery() string {
	return `
list "azurerm_elastic_cloud_hosted_search" "list" {
  provider = azurerm
  config {}
}
`
}

func (ElasticCloudHostedSearchResource) basicListQueryByResourceGroupName() string {
	return `
list "azurerm_elastic_cloud_hosted_search" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}

func (ElasticCloudHostedSearchResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_cloud_hosted_search" "test" {
  count = 2

  name                        = "acctesths-${count.index}-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"
}
`, data.RandomInteger, data.Locations.Primary)
}
