---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_pipeline_group"
description: |-
  Manages a Pipeline Group within Azure Monitor.
---

# azurerm_monitor_pipeline_group

Manages a Pipeline Group within Azure Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-law"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_monitor_data_collection_endpoint" "example" {
  name                = "example-dce"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_monitor_data_collection_rule" "example" {
  name                        = "example-dcr"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  data_collection_endpoint_id = azurerm_monitor_data_collection_endpoint.example.id

  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.example.id
      name                  = "example-destination-log"
    }
  }

  data_flow {
    streams      = ["Custom-Table_CL"]
    destinations = ["example-destination-log"]
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

resource "azurerm_monitor_pipeline_group" "example" {
  name                = "example-pipeline-group"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  replicas            = 2

  exporter {
    name = "workspace-logs-exporter"
    type = "AzureMonitorWorkspaceLogs"

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url = azurerm_monitor_data_collection_endpoint.example.logs_ingestion_endpoint
        data_collection_rule         = azurerm_monitor_data_collection_rule.example.immutable_id
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
    environment = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Pipeline Group. Changing this forces a new Pipeline Group to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Pipeline Group should exist. Changing this forces a new Pipeline Group to be created.

* `location` - (Required) The Azure Region where the Pipeline Group should exist.

* `exporter` - (Required) One or more `exporter` blocks as defined below.

* `receiver` - (Required) One or more `receiver` blocks as defined below.

* `service` - (Required) A `service` block as defined below.

* `execution_placement` - (Optional) An `execution_placement` block as defined below.

* `extended_location` - (Optional) An `extended_location` block as defined below.

* `processor` - (Optional) One or more `processor` blocks as defined below.

* `replicas` - (Optional) The number of collector replicas to provision.

* `tls_configuration` - (Optional) One or more `tls_configuration` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Pipeline Group.

---

An `azure_monitor_workspace_logs` block supports the following:

* `api` - (Required) An `api` block as defined below.

* `persistence` - (Optional) A `persistence` block as defined below.

---

An `api` block supports the following:

* `data_collection_endpoint_url` - (Required) The Data Collection Endpoint ingestion URL used by this exporter.

* `data_collection_rule` - (Required) The immutable ID of the Data Collection Rule used by this exporter.

* `schema` - (Required) A `schema` block as defined below.

* `stream` - (Required) The target stream name for this exporter.

---

A `batch` block supports the following:

* `batch_size` - (Optional) The maximum number of records to include in a batch.

* `timeout` - (Optional) The maximum amount of time, in seconds, to wait before flushing a batch.

---

A `certificate` block supports the following:

* `location` - (Required) The Kubernetes object name that stores the certificate.

* `sub_location` - (Required) The key within the Kubernetes object that stores the certificate.

* `type` - (Required) The certificate source type. Possible values are `kubernetesConfigMap` and `kubernetesSecret`.

---

A `client_ca` block supports the following:

* `location` - (Required) The Kubernetes object name that stores the client CA certificate.

* `sub_location` - (Required) The key within the Kubernetes object that stores the client CA certificate.

* `type` - (Required) The certificate source type. Possible values are `kubernetesConfigMap` and `kubernetesSecret`.

---

A `constraint` block supports the following:

* `capability` - (Required) The capability name used by this placement constraint.

* `operator` - (Required) The placement operator. Possible values are `DoesNotExist`, `Exists`, `In`, and `NotIn`.

* `values` - (Optional) The values used by this placement constraint.

---

An `execution_placement` block supports the following:

* `constraint` - (Optional) One or more `constraint` blocks as defined below.

* `maximum_instances_per_host` - (Optional) The maximum number of replicas that can run on the same host.

---

An `exporter` block supports the following:

* `name` - (Required) The name of the exporter.

* `type` - (Required) The exporter type. Possible values are `AzureMonitorWorkspaceLogs`.

* `azure_monitor_workspace_logs` - (Optional) An `azure_monitor_workspace_logs` block as defined above.

-> **Note:** The `azure_monitor_workspace_logs` block must be specified when `type` is `AzureMonitorWorkspaceLogs`.

---

An `extended_location` block supports the following:

* `name` - (Required) The name or resource ID of the extended location.

* `type` - (Required) The extended location type. Possible values are `CustomLocation` and `EdgeZone`.

---

A `persistence` block in an `azure_monitor_workspace_logs` block supports the following:

* `maximum_storage_usage` - (Optional) The maximum storage usage for the exporter persistence buffer.

* `retention_period` - (Optional) The persistence retention period, in days.

---

A `persistence` block in a `service` block supports the following:

* `persistent_volume_name` - (Required) The persistent volume name used by the service.

---

A `pipeline` block supports the following:

* `exporter` - (Required) A list of exporter names defined in `exporter` blocks.

* `name` - (Required) The name of the pipeline.

* `receiver` - (Required) A list of receiver names defined in `receiver` blocks.

* `type` - (Required) The pipeline type. Possible values are `Logs`.

* `processor` - (Optional) A list of processor names defined in `processor` blocks.

---

A `private_key` block supports the following:

* `location` - (Required) The Kubernetes object name that stores the private key.

* `sub_location` - (Required) The key within the Kubernetes object that stores the private key.

* `type` - (Required) The private key source type. Possible values are `kubernetesSecret`.

---

A `processor` block supports the following:

* `name` - (Required) The name of the processor.

* `type` - (Required) The processor type. Possible values are `Batch`, `MicrosoftCommonSecurityLog`, `MicrosoftSyslog`, and `TransformLanguage`.

* `batch` - (Optional) A `batch` block as defined above.

* `transform_statement` - (Optional) The transform statement to run for this processor.

-> **Note:** The `batch` block can only be specified when `type` is `Batch`. The `transform_statement` field must be specified when `type` is `TransformLanguage`.

---

A `receiver` block supports the following:

* `name` - (Required) The name of the receiver.

* `type` - (Required) The receiver type. Possible values are `OTLP` and `Syslog`.

* `otlp_endpoint` - (Optional) The OTLP listener endpoint in `host:port` format.

* `syslog` - (Optional) A `syslog` block as defined below.

* `tls_configuration_name` - (Optional) The name of a `tls_configuration` block to use for this receiver.

-> **Note:** The `otlp_endpoint` field must be specified when `type` is `OTLP`. The `syslog` block must be specified when `type` is `Syslog`.

---

A `record_map` block supports the following:

* `from` - (Required) The source record field name.

* `to` - (Required) The destination column name.

---

A `resource_map` block supports the following:

* `from` - (Required) The source resource field name.

* `to` - (Required) The destination resource attribute name.

---

A `schema` block supports the following:

* `record_map` - (Required) One or more `record_map` blocks as defined above.

* `resource_map` - (Optional) One or more `resource_map` blocks as defined above.

* `scope_map` - (Optional) One or more `scope_map` blocks as defined below.

---

A `scope_map` block supports the following:

* `from` - (Required) The source scope field name.

* `to` - (Required) The destination scope attribute name.

---

A `service` block supports the following:

* `pipeline` - (Required) One or more `pipeline` blocks as defined above.

* `persistent_volume_name` - (Optional) The persistent volume name used by the service.

---

A `syslog` block supports the following:

* `endpoint` - (Required) The syslog listener endpoint in `host:port` format.

* `allow_skip_pri_header` - (Optional) Whether syslog messages may omit the PRI header.

* `allowed_formats` - (Optional) The allowed syslog payload formats. Possible values are `all`, `cefRfc5424`, `cefRfc3164`, `rawCef`, `syslogRfc5424`, and `syslogRfc3164`.

* `transport_protocol` - (Optional) The syslog transport protocol. Possible values are `tcp` and `udp`.

---

A `tls_certificate` block supports the following:

* `certificate` - (Required) A `certificate` block as defined above.

* `private_key` - (Required) A `private_key` block as defined above.

---

A `tls_configuration` block supports the following:

* `name` - (Required) The name of the TLS configuration.

* `client_ca` - (Optional) A `client_ca` block as defined above.

* `mode` - (Optional) The TLS mode. Possible values are `disabled`, `mutualTls`, and `serverOnly`.

* `tls_certificate` - (Optional) A `tls_certificate` block as defined above.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Pipeline Group ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Pipeline Group.

* `read` - (Defaults to 5 minutes) Used when retrieving the Pipeline Group.

* `update` - (Defaults to 30 minutes) Used when updating the Pipeline Group.

* `delete` - (Defaults to 30 minutes) Used when deleting the Pipeline Group.

## Import

A Monitor Pipeline Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_pipeline_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Monitor/pipelineGroups/pipelineGroup1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Monitor` - 2026-04-01
