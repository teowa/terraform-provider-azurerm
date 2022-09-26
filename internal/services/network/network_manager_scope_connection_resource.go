package network

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/networkmanagers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/scopeconnections"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkScopeConnectionModel struct {
	Name                    string                                `tfschema:"name"`
	NetworkNetworkManagerId string                                `tfschema:"network_network_manager_id"`
	ConnectionState         scopeconnections.ScopeConnectionState `tfschema:"connection_state"`
	Description             string                                `tfschema:"description"`
	ResourceId              string                                `tfschema:"resource_id"`
	TenantId                string                                `tfschema:"tenant_id"`
}

type NetworkScopeConnectionResource struct{}

var _ sdk.ResourceWithUpdate = NetworkScopeConnectionResource{}

func (r NetworkScopeConnectionResource) ResourceType() string {
	return "azurerm_network_scope_connection"
}

func (r NetworkScopeConnectionResource) ModelObject() interface{} {
	return &NetworkScopeConnectionModel{}
}

func (r NetworkScopeConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scopeconnections.ValidateScopeConnectionID
}

func (r NetworkScopeConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkmanagers.ValidateNetworkManagerID,
		},

		"connection_state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(scopeconnections.ScopeConnectionStateConnected),
				string(scopeconnections.ScopeConnectionStatePending),
				string(scopeconnections.ScopeConnectionStateConflict),
				string(scopeconnections.ScopeConnectionStateRevoked),
				string(scopeconnections.ScopeConnectionStateRejected),
			}, false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r NetworkScopeConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkScopeConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkScopeConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ScopeConnectionsClient
			networkManagerId, err := networkmanagers.ParseNetworkManagerID(model.NetworkNetworkManagerId)
			if err != nil {
				return err
			}

			id := scopeconnections.NewScopeConnectionID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &scopeconnections.ScopeConnection{
				Properties: &scopeconnections.ScopeConnectionProperties{
					ConnectionState: &model.ConnectionState,
				},
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if model.ResourceId != "" {
				properties.Properties.ResourceId = &model.ResourceId
			}

			if model.TenantId != "" {
				properties.Properties.TenantId = &model.TenantId
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkScopeConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ScopeConnectionsClient

			id, err := scopeconnections.ParseScopeConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkScopeConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("connection_state") {
				properties.Properties.ConnectionState = &model.ConnectionState
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Properties.Description = &model.Description
				} else {
					properties.Properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("resource_id") {
				if model.ResourceId != "" {
					properties.Properties.ResourceId = &model.ResourceId
				} else {
					properties.Properties.ResourceId = nil
				}
			}

			if metadata.ResourceData.HasChange("tenant_id") {
				if model.TenantId != "" {
					properties.Properties.TenantId = &model.TenantId
				} else {
					properties.Properties.TenantId = nil
				}
			}

			properties.SystemData = nil

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NetworkScopeConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ScopeConnectionsClient

			id, err := scopeconnections.ParseScopeConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := NetworkScopeConnectionModel{
				Name:                    id.ScopeConnectionName,
				NetworkNetworkManagerId: networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
			}

			if properties := model.Properties; properties != nil {
				if properties.ConnectionState != nil {
					state.ConnectionState = *properties.ConnectionState
				}

				if properties.Description != nil {
					state.Description = *properties.Description
				}

				if properties.ResourceId != nil {
					state.ResourceId = *properties.ResourceId
				}

				if properties.TenantId != nil {
					state.TenantId = *properties.TenantId
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkScopeConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ScopeConnectionsClient

			id, err := scopeconnections.ParseScopeConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
