// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagecache_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccManagedLustreFileSystemAutoImportJob_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "test")
	r := ManagedLustreFileSystemAutoImportJobResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"aml_filesystem_name": {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_managed_lustre_file_system_auto_import_job.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_managed_lustre_file_system_auto_import_job.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_managed_lustre_file_system_auto_import_job.test", tfjsonpath.New("aml_filesystem_name"), tfjsonpath.New("managed_lustre_file_system_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_managed_lustre_file_system_auto_import_job.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("managed_lustre_file_system_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_managed_lustre_file_system_auto_import_job.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("managed_lustre_file_system_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
