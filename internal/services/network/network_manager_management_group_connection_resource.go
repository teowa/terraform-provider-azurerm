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

type NetworkManagementGroupNetworkManagerConnectionModel struct {
	Name              string                                         `tfschema:"name"`
	ManagementGroupId string                                         `tfschema:"management_group_id"`
	ConnectionState   networkmanagerconnections.ScopeConnectionState `tfschema:"connection_state"`
	Description       string                                         `tfschema:"description"`
	NetworkManagerId  string                                         `tfschema:"network_manager_id"`
}

type NetworkManagementGroupNetworkManagerConnectionResource struct{}

var _ sdk.ResourceWithUpdate = NetworkManagementGroupNetworkManagerConnectionResource{}

func (r NetworkManagementGroupNetworkManagerConnectionResource) ResourceType() string {
	return "azurerm_network_management_group_network_manager_connection"
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) ModelObject() interface{} {
	return &NetworkManagementGroupNetworkManagerConnectionModel{}
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkmanagerconnections.ValidateManagementGroupNetworkManagerConnectionID
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"management_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"connection_state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(networkmanagerconnections.ScopeConnectionStatePending),
				string(networkmanagerconnections.ScopeConnectionStateConflict),
				string(networkmanagerconnections.ScopeConnectionStateRevoked),
				string(networkmanagerconnections.ScopeConnectionStateRejected),
				string(networkmanagerconnections.ScopeConnectionStateConnected),
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

func (r NetworkManagementGroupNetworkManagerConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkManagementGroupNetworkManagerConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.NetworkManagerConnectionsClient
			id := networkmanagerconnections.NewProviders2NetworkManagerConnectionID(model.ManagementGroupId, model.Name)
			existing, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, id)
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

			if _, err := client.ManagementGroupNetworkManagerConnectionsCreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnectionsClient

			id, err := networkmanagerconnections.ParseProviders2NetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkManagementGroupNetworkManagerConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, *id)
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

			if _, err := client.ManagementGroupNetworkManagerConnectionsCreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NetworkManagementGroupNetworkManagerConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnectionsClient

			id, err := networkmanagerconnections.ParseProviders2NetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, *id)
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

			state := NetworkManagementGroupNetworkManagerConnectionModel{
				Name:              id.NetworkManagerConnectionName,
				ManagementGroupId: id.ManagementGroupId,
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

func (r NetworkManagementGroupNetworkManagerConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkManagerConnectionsClient

			id, err := networkmanagerconnections.ParseProviders2NetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.ManagementGroupNetworkManagerConnectionsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
