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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/networkmanagerconnections"
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

type NetworkSubscriptionNetworkManagerConnectionModel struct {
	Name             string                                         `tfschema:"name"`
	ConnectionState  networkmanagerconnections.ScopeConnectionState `tfschema:"connection_state"`
	Description      string                                         `tfschema:"description"`
	NetworkManagerId string                                         `tfschema:"network_manager_id"`
}

type NetworkSubscriptionNetworkManagerConnectionResource struct{}

var _ sdk.ResourceWithUpdate = NetworkSubscriptionNetworkManagerConnectionResource{}

func (r NetworkSubscriptionNetworkManagerConnectionResource) ResourceType() string {
	return "azurerm_network_subscription_network_manager_connection"
}

func (r NetworkSubscriptionNetworkManagerConnectionResource) ModelObject() interface{} {
	return &NetworkSubscriptionNetworkManagerConnectionModel{}
}

func (r NetworkSubscriptionNetworkManagerConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkmanagerconnections.ValidateSubscriptionNetworkManagerConnectionID
}

func (r NetworkSubscriptionNetworkManagerConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"connection_state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(networkmanagerconnections.ScopeConnectionStateConnected),
				string(networkmanagerconnections.ScopeConnectionStatePending),
				string(networkmanagerconnections.ScopeConnectionStateConflict),
				string(networkmanagerconnections.ScopeConnectionStateRevoked),
				string(networkmanagerconnections.ScopeConnectionStateRejected),
			}, false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r NetworkSubscriptionNetworkManagerConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkSubscriptionNetworkManagerConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkSubscriptionNetworkManagerConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.NetworkManagerConnectionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := networkmanagerconnections.NewNetworkManagerConnectionID(subscriptionId, model.Name)
			existing, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &networkmanagerconnections.NetworkManagerConnection{
				Properties: &networkmanagerconnections.NetworkManagerConnectionProperties{
					ConnectionState: &model.ConnectionState,
				},
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if model.NetworkManagerId != "" {
				properties.Properties.NetworkManagerId = &model.NetworkManagerId
			}

			if _, err := client.SubscriptionNetworkManagerConnectionsCreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkSubscriptionNetworkManagerConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnectionsClient

			id, err := networkmanagerconnections.ParseNetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkSubscriptionNetworkManagerConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, *id)
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

			if metadata.ResourceData.HasChange("network_manager_id") {
				if model.NetworkManagerId != "" {
					properties.Properties.NetworkManagerId = &model.NetworkManagerId
				} else {
					properties.Properties.NetworkManagerId = nil
				}
			}

			properties.SystemData = nil

			if _, err := client.SubscriptionNetworkManagerConnectionsCreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NetworkSubscriptionNetworkManagerConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnectionsClient

			id, err := networkmanagerconnections.ParseNetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, *id)
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

			state := NetworkSubscriptionNetworkManagerConnectionModel{
				Name: id.NetworkManagerConnectionName,
			}

			if properties := model.Properties; properties != nil {
				if properties.ConnectionState != nil {
					state.ConnectionState = *properties.ConnectionState
				}

				if properties.Description != nil {
					state.Description = *properties.Description
				}

				if properties.NetworkManagerId != nil {
					state.NetworkManagerId = *properties.NetworkManagerId
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkSubscriptionNetworkManagerConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnectionsClient

			id, err := networkmanagerconnections.ParseNetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.SubscriptionNetworkManagerConnectionsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
