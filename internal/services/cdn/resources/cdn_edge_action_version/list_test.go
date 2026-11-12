// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_version_test

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

func TestAccCdnEdgeActionVersion_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action_version", "test")
	r := CdnEdgeActionVersionResource{}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_14_0)},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.basic(data) + `
resource "azurerm_cdn_edge_action_version" "additional" {
  count = 2

  name               = "v${count.index + 2}"
  cdn_edge_action_id = azurerm_cdn_edge_action.test.id
  location           = "global"
  deployment_type    = "file"
}
`},
			{
				Query: true,
				Config: `
list "azurerm_cdn_edge_action_version" "test" {
  provider         = azurerm
  include_resource = true
  config {
    cdn_edge_action_id = azurerm_cdn_edge_action.test.id
  }
}
`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_cdn_edge_action_version.test", 3),
					querycheck.ExpectIdentity("azurerm_cdn_edge_action_version.test", map[string]knownvalue.Check{
						"name":                knownvalue.StringRegexp(regexp.MustCompile(`^v[1-3]$`)),
						"edge_action_name":    knownvalue.StringExact("acctest" + strconv.Itoa(data.RandomInteger)),
						"resource_group_name": knownvalue.StringExact("acctestRG-" + strconv.Itoa(data.RandomInteger)),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}
