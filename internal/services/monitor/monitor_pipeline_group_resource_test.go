// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MonitorPipelineGroupResource struct{}

func (r MonitorPipelineGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := pipelinegroups.ParsePipelineGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Monitor.PipelineGroupsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccMonitorPipelineGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_pipeline_group", "test")
	r := MonitorPipelineGroupResource{}

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

func TestAccMonitorPipelineGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_pipeline_group", "test")
	r := MonitorPipelineGroupResource{}

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

func TestAccMonitorPipelineGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_pipeline_group", "test")
	r := MonitorPipelineGroupResource{}

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

func TestAccMonitorPipelineGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_pipeline_group", "test")
	r := MonitorPipelineGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MonitorPipelineGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_pipeline_group" "test" {
  name                = "acctest-mpg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  exporter {
    name = "workspace-logs-exporter"
    type = "AzureMonitorWorkspaceLogs"

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url = azurerm_monitor_data_collection_endpoint.test.logs_ingestion_endpoint
        data_collection_rule         = azurerm_monitor_data_collection_rule.test.immutable_id
        stream                       = "Custom-Table_CL"

        schema {
          record_map {
            from = "body"
            to   = "Body"
          }

          record_map {
            from = "severity_text"
            to   = "SeverityText"
          }

          record_map {
            from = "time_unix_nano"
            to   = "TimeGenerated"
          }
        }
      }
    }
  }

  processor {
    name = "batch-processor"
    type = "Batch"
  }

  receiver {
    name = "syslog-receiver"
    type = "Syslog"

    syslog {
      endpoint = "0.0.0.0:514"
    }
  }

  service {
    pipeline {
      name      = "logs-pipeline"
      type      = "Logs"
      receiver  = ["syslog-receiver"]
      processor = ["batch-processor"]
      exporter  = ["workspace-logs-exporter"]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorPipelineGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_pipeline_group" "test" {
  name                = "acctest-mpg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  replicas            = 2

  exporter {
    name = "workspace-logs-exporter"
    type = "AzureMonitorWorkspaceLogs"

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url = azurerm_monitor_data_collection_endpoint.test.logs_ingestion_endpoint
        data_collection_rule         = azurerm_monitor_data_collection_rule.test.immutable_id
        stream                       = "Custom-Table_CL"

        schema {
          record_map {
            from = "body"
            to   = "Body"
          }

          record_map {
            from = "severity_text"
            to   = "SeverityText"
          }

          record_map {
            from = "time_unix_nano"
            to   = "TimeGenerated"
          }
        }
      }

      persistence {
        maximum_storage_usage = 50
        retention_period  = 7
      }
    }
  }

  processor {
    name = "batch-processor"
    type = "Batch"

    batch {
      batch_size = 10
      timeout    = 30
    }
  }

  receiver {
    name = "syslog-receiver"
    type = "Syslog"

    syslog {
      allow_skip_pri_header = true
      allowed_formats       = ["all"]
      endpoint              = "0.0.0.0:514"
      transport_protocol    = "tcp"
    }
  }

  service {
    persistent_volume_name = "pipeline-storage"

    pipeline {
      name      = "logs-pipeline"
      type      = "Logs"
      receiver  = ["syslog-receiver"]
      processor = ["batch-processor"]
      exporter  = ["workspace-logs-exporter"]
    }
  }

  tags = {
    environment = "acctest"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorPipelineGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_pipeline_group" "import" {
  name                = azurerm_monitor_pipeline_group.test.name
  resource_group_name = azurerm_monitor_pipeline_group.test.resource_group_name
  location            = azurerm_monitor_pipeline_group.test.location

  exporter {
    name = "workspace-logs-exporter"
    type = "AzureMonitorWorkspaceLogs"

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url = azurerm_monitor_data_collection_endpoint.test.logs_ingestion_endpoint
        data_collection_rule         = azurerm_monitor_data_collection_rule.test.immutable_id
        stream                       = "Custom-Table_CL"

        schema {
          record_map {
            from = "body"
            to   = "Body"
          }

          record_map {
            from = "severity_text"
            to   = "SeverityText"
          }

          record_map {
            from = "time_unix_nano"
            to   = "TimeGenerated"
          }
        }
      }
    }
  }

  processor {
    name = "batch-processor"
    type = "Batch"
  }

  receiver {
    name = "syslog-receiver"
    type = "Syslog"

    syslog {
      endpoint = "0.0.0.0:514"
    }
  }

  service {
    pipeline {
      name      = "logs-pipeline"
      type      = "Logs"
      receiver  = ["syslog-receiver"]
      processor = ["batch-processor"]
      exporter  = ["workspace-logs-exporter"]
    }
  }
}
`, r.basic(data))
}

func (r MonitorPipelineGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-mpg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestlaw%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_monitor_data_collection_endpoint" "test" {
  name                = "acctest-dce-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_monitor_data_collection_rule" "test" {
  name                        = "acctest-dcr-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  data_collection_endpoint_id = azurerm_monitor_data_collection_endpoint.test.id

  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.test.id
      name                  = "test-destination-log"
    }
  }

  data_flow {
    streams      = ["Custom-Table_CL"]
    destinations = ["test-destination-log"]
  }

  stream_declaration {
    stream_name = "Custom-Table_CL"

    column {
      name = "TimeGenerated"
      type = "datetime"
    }

    column {
      name = "Body"
      type = "string"
    }

    column {
      name = "SeverityText"
      type = "string"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
