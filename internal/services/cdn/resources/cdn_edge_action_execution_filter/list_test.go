// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_execution_filter_test

import (
	"context"
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

func TestAccCdnEdgeActionExecutionFilter_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_execution_filter", "test")
	r := CdnEdgeActionExecutionFilterResource{}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_14_0)},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.basic(data) + `
resource "azurerm_cdn_edge_action_execution_filter" "additional" {
  count = 2

  name                       = "filter${count.index}"
  cdn_edge_action_id         = azurerm_cdn_edge_action.test.id
  cdn_edge_action_version_id = azurerm_cdn_edge_action_version.test[0].id
  location                   = "global"
  header_name                = "X-Edge-Filter-${count.index}"
  header_value               = "test"
}
`},
			{
				Query: true,
				Config: `
list "azurerm_cdn_edge_action_execution_filter" "test" {
  provider         = azurerm
  include_resource = true
  config {
    cdn_edge_action_id = azurerm_cdn_edge_action.test.id
  }
}
`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_cdn_edge_action_execution_filter.test", 3),
					querycheck.ExpectIdentity("azurerm_cdn_edge_action_execution_filter.test", map[string]knownvalue.Check{
						"name":                knownvalue.StringRegexp(regexp.MustCompile(`^filter[01]?$`)),
						"edge_action_name":    knownvalue.StringExact("acctest" + strconv.Itoa(data.RandomInteger)),
						"resource_group_name": knownvalue.StringExact("acctestRG-" + strconv.Itoa(data.RandomInteger)),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}
