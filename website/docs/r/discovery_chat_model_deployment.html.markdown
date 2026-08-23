---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_discovery_chat_model_deployment"
description: |-
  Manages a Discovery Chat Model Deployment.
---

# azurerm_discovery_chat_model_deployment

Manages a Discovery Chat Model Deployment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-discovery-rg"
  location = "West Europe"
}

resource "azurerm_discovery_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_discovery_chat_model_deployment" "example" {
  name                = "example-chat-model-deployment"
  discovery_workspace_id = azurerm_discovery_workspace.example.id
  location                = azurerm_resource_group.example.location
}
```

## Arguments Reference

* `location` - (Required) The Azure Region where the Discovery Chat Model Deployment should exist. Changing this forces a new resource to be created.
* `name` - (Required) The name of the Discovery Chat Model Deployment. Changing this forces a new resource to be created.
* `discovery_workspace_id` - (Required) The ID of the Discovery Workspace where this Chat Model Deployment should exist. Changing this forces a new resource to be created.
* `tags` - (Optional) A mapping of tags which should be assigned to the Discovery Chat Model Deployment.

## Attributes Reference

* `id` - The ID of the Discovery Chat Model Deployment.

## Timeouts

The `timeouts` block supports `create`, `read`, `update`, and `delete`.

## Import

Discovery Chat Model Deployments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_discovery_chat_model_deployment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Discovery/workspaces/example/chatModelDeployments/example
```
