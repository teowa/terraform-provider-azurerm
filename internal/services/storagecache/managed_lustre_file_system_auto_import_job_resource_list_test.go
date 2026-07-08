// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storagecache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccManagedLustreFileSystemAutoImportJob_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "testlist")
	r := ManagedLustreFileSystemAutoImportJobResource{}
	resourceName := fmt.Sprintf("acctest-autoimportjob-%d", data.RandomInteger)
	amlFilesystemName := fmt.Sprintf("acctest-amlfs-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-amlfs-%d", data.RandomInteger)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.basic(data)},
			{
				Query:  true,
				Config: r.listQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_managed_lustre_file_system_auto_import_job.list", 1),
					querycheck.ExpectIdentity("azurerm_managed_lustre_file_system_auto_import_job.list", map[string]knownvalue.Check{
						"aml_filesystem_name": knownvalue.StringExact(amlFilesystemName),
						"name":                knownvalue.StringExact(resourceName),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (r ManagedLustreFileSystemAutoImportJobResource) listQuery() string {
	return `
list "azurerm_managed_lustre_file_system_auto_import_job" "list" {
  provider = azurerm
  config {
    managed_lustre_file_system_id = azurerm_managed_lustre_file_system.test.id
  }
}
`
}
