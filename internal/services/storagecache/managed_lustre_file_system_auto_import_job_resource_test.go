// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storagecache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2025-07-01/autoimportjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedLustreFileSystemAutoImportJobResource struct{}

func TestAccManagedLustreFileSystemAutoImportJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "test")
	r := ManagedLustreFileSystemAutoImportJobResource{}

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

func TestAccManagedLustreFileSystemAutoImportJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "test")
	r := ManagedLustreFileSystemAutoImportJobResource{}

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

func TestAccManagedLustreFileSystemAutoImportJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "test")
	r := ManagedLustreFileSystemAutoImportJobResource{}

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

func TestAccManagedLustreFileSystemAutoImportJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "test")
	r := ManagedLustreFileSystemAutoImportJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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

func (r ManagedLustreFileSystemAutoImportJobResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autoimportjobs.ParseAutoImportJobID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageCache_2025_07_01.AutoImportJobs
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagedLustreFileSystemAutoImportJobResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_import_job" "test" {
  name                           = "acctest-autoimportjob-%d"
  managed_lustre_file_system_id  = azurerm_managed_lustre_file_system.test.id
  location                       = azurerm_resource_group.test.location
  admin_status                   = "Enable"
  auto_import_prefixes           = ["/"]
  conflict_resolution_mode       = "Skip"
}
`, ManagedLustreFileSystemResource{}.complete(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemAutoImportJobResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_import_job" "test" {
  name                          = "acctest-autoimportjob-%d"
  managed_lustre_file_system_id = azurerm_managed_lustre_file_system.test.id
  location                      = azurerm_resource_group.test.location
  admin_status                  = "Enable"
  auto_import_prefixes          = ["/data", "/archive"]
  conflict_resolution_mode      = "OverwriteAlways"
  enable_deletions              = true
  maximum_errors                = 5

  tags = {
    Env = "Test"
  }
}
`, ManagedLustreFileSystemResource{}.complete(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemAutoImportJobResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_import_job" "test" {
  name                          = "acctest-autoimportjob-%d"
  managed_lustre_file_system_id = azurerm_managed_lustre_file_system.test.id
  location                      = azurerm_resource_group.test.location
  admin_status                  = "Disable"
  auto_import_prefixes          = ["/data", "/archive"]
  conflict_resolution_mode      = "OverwriteAlways"
  enable_deletions              = true
  maximum_errors                = 5

  tags = {
    Env = "Test2"
  }
}
`, ManagedLustreFileSystemResource{}.complete(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemAutoImportJobResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_import_job" "import" {
  name                          = azurerm_managed_lustre_file_system_auto_import_job.test.name
  managed_lustre_file_system_id = azurerm_managed_lustre_file_system_auto_import_job.test.managed_lustre_file_system_id
  location                      = azurerm_managed_lustre_file_system_auto_import_job.test.location
  admin_status                  = azurerm_managed_lustre_file_system_auto_import_job.test.admin_status
  auto_import_prefixes          = azurerm_managed_lustre_file_system_auto_import_job.test.auto_import_prefixes
  conflict_resolution_mode      = azurerm_managed_lustre_file_system_auto_import_job.test.conflict_resolution_mode
  enable_deletions              = azurerm_managed_lustre_file_system_auto_import_job.test.enable_deletions
  maximum_errors                = azurerm_managed_lustre_file_system_auto_import_job.test.maximum_errors
  tags                          = azurerm_managed_lustre_file_system_auto_import_job.test.tags
}
`, r.complete(data))
}
