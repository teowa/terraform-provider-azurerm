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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/connectivityconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/networkmanagers"
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

type NetworkConnectivityConfigurationModel struct {
	Name                    string                                           `tfschema:"name"`
	NetworkNetworkManagerId string                                           `tfschema:"network_network_manager_id"`
	AppliesToGroups         []ConnectivityGroupItemModel                     `tfschema:"applies_to_groups"`
	ConnectivityTopology    connectivityconfigurations.ConnectivityTopology  `tfschema:"connectivity_topology"`
	DeleteExistingPeering   connectivityconfigurations.DeleteExistingPeering `tfschema:"delete_existing_peering"`
	Description             string                                           `tfschema:"description"`
	Hubs                    []HubModel                                       `tfschema:"hubs"`
	IsGlobal                connectivityconfigurations.IsGlobal              `tfschema:"is_global"`
}

type ConnectivityGroupItemModel struct {
	GroupConnectivity connectivityconfigurations.GroupConnectivity `tfschema:"group_connectivity"`
	IsGlobal          connectivityconfigurations.IsGlobal          `tfschema:"is_global"`
	NetworkGroupId    string                                       `tfschema:"network_group_id"`
	UseHubGateway     connectivityconfigurations.UseHubGateway     `tfschema:"use_hub_gateway"`
}

type HubModel struct {
	ResourceId   string `tfschema:"resource_id"`
	ResourceType string `tfschema:"resource_type"`
}

type NetworkConnectivityConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = NetworkConnectivityConfigurationResource{}

func (r NetworkConnectivityConfigurationResource) ResourceType() string {
	return "azurerm_network_connectivity_configuration"
}

func (r NetworkConnectivityConfigurationResource) ModelObject() interface{} {
	return &NetworkConnectivityConfigurationModel{}
}

func (r NetworkConnectivityConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connectivityconfigurations.ValidateConnectivityConfigurationID
}

func (r NetworkConnectivityConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
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

		"applies_to_groups": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"group_connectivity": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(connectivityconfigurations.GroupConnectivityNone),
							string(connectivityconfigurations.GroupConnectivityDirectlyConnected),
						}, false),
					},

					"is_global": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(connectivityconfigurations.IsGlobalTrue),
							string(connectivityconfigurations.IsGlobalFalse),
						}, false),
					},

					"network_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"use_hub_gateway": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(connectivityconfigurations.UseHubGatewayFalse),
							string(connectivityconfigurations.UseHubGatewayTrue),
						}, false),
					},
				},
			},
		},

		"connectivity_topology": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(connectivityconfigurations.ConnectivityTopologyHubAndSpoke),
				string(connectivityconfigurations.ConnectivityTopologyMesh),
			}, false),
		},

		"delete_existing_peering": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(connectivityconfigurations.DeleteExistingPeeringFalse),
				string(connectivityconfigurations.DeleteExistingPeeringTrue),
			}, false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hubs": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"resource_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"is_global": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(connectivityconfigurations.IsGlobalFalse),
				string(connectivityconfigurations.IsGlobalTrue),
			}, false),
		},
	}
}

func (r NetworkConnectivityConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkConnectivityConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkConnectivityConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ConnectivityConfigurationsClient
			networkManagerId, err := networkmanagers.ParseNetworkManagerID(model.NetworkNetworkManagerId)
			if err != nil {
				return err
			}

			id := connectivityconfigurations.NewConnectivityConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &connectivityconfigurations.ConnectivityConfiguration{
				Properties: &connectivityconfigurations.ConnectivityConfigurationProperties{
					ConnectivityTopology:  model.ConnectivityTopology,
					DeleteExistingPeering: &model.DeleteExistingPeering,
					IsGlobal:              &model.IsGlobal,
				},
			}

			appliesToGroupsValue, err := expandConnectivityGroupItemModel(model.AppliesToGroups)
			if err != nil {
				return err
			}

			if appliesToGroupsValue != nil {
				properties.Properties.AppliesToGroups = *appliesToGroupsValue
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			hubsValue, err := expandHubModel(model.Hubs)
			if err != nil {
				return err
			}

			properties.Properties.Hubs = hubsValue

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkConnectivityConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ConnectivityConfigurationsClient

			id, err := connectivityconfigurations.ParseConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkConnectivityConfigurationModel
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

			if metadata.ResourceData.HasChange("applies_to_groups") {
				appliesToGroupsValue, err := expandConnectivityGroupItemModel(model.AppliesToGroups)
				if err != nil {
					return err
				}

				if appliesToGroupsValue != nil {
					properties.Properties.AppliesToGroups = *appliesToGroupsValue
				}
			}

			if metadata.ResourceData.HasChange("connectivity_topology") {
				properties.Properties.ConnectivityTopology = model.ConnectivityTopology
			}

			if metadata.ResourceData.HasChange("delete_existing_peering") {
				properties.Properties.DeleteExistingPeering = &model.DeleteExistingPeering
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Properties.Description = &model.Description
				} else {
					properties.Properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("hubs") {
				hubsValue, err := expandHubModel(model.Hubs)
				if err != nil {
					return err
				}

				properties.Properties.Hubs = hubsValue
			}

			if metadata.ResourceData.HasChange("is_global") {
				properties.Properties.IsGlobal = &model.IsGlobal
			}

			properties.SystemData = nil

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NetworkConnectivityConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ConnectivityConfigurationsClient

			id, err := connectivityconfigurations.ParseConnectivityConfigurationID(metadata.ResourceData.Id())
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

			state := NetworkConnectivityConfigurationModel{
				Name:                    id.ConfigurationName,
				NetworkNetworkManagerId: networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
			}

			if properties := model.Properties; properties != nil {
				appliesToGroupsValue, err := flattenConnectivityGroupItemModel(&properties.AppliesToGroups)
				if err != nil {
					return err
				}

				state.AppliesToGroups = appliesToGroupsValue

				state.ConnectivityTopology = properties.ConnectivityTopology

				if properties.DeleteExistingPeering != nil {
					state.DeleteExistingPeering = *properties.DeleteExistingPeering
				}

				if properties.Description != nil {
					state.Description = *properties.Description
				}

				hubsValue, err := flattenHubModel(properties.Hubs)
				if err != nil {
					return err
				}

				state.Hubs = hubsValue

				if properties.IsGlobal != nil {
					state.IsGlobal = *properties.IsGlobal
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkConnectivityConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ConnectivityConfigurationsClient

			id, err := connectivityconfigurations.ParseConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, connectivityconfigurations.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandConnectivityGroupItemModel(inputList []ConnectivityGroupItemModel) (*[]connectivityconfigurations.ConnectivityGroupItem, error) {
	var outputList []connectivityconfigurations.ConnectivityGroupItem
	for _, v := range inputList {
		input := v
		output := connectivityconfigurations.ConnectivityGroupItem{
			GroupConnectivity: input.GroupConnectivity,
			IsGlobal:          &input.IsGlobal,
			NetworkGroupId:    input.NetworkGroupId,
			UseHubGateway:     &input.UseHubGateway,
		}

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func expandHubModel(inputList []HubModel) (*[]connectivityconfigurations.Hub, error) {
	var outputList []connectivityconfigurations.Hub
	for _, v := range inputList {
		input := v
		output := connectivityconfigurations.Hub{}

		if input.ResourceId != "" {
			output.ResourceId = &input.ResourceId
		}

		if input.ResourceType != "" {
			output.ResourceType = &input.ResourceType
		}

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func flattenConnectivityGroupItemModel(inputList *[]connectivityconfigurations.ConnectivityGroupItem) ([]ConnectivityGroupItemModel, error) {
	var outputList []ConnectivityGroupItemModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := ConnectivityGroupItemModel{
			GroupConnectivity: input.GroupConnectivity,
			NetworkGroupId:    input.NetworkGroupId,
		}

		if input.IsGlobal != nil {
			output.IsGlobal = *input.IsGlobal
		}

		if input.UseHubGateway != nil {
			output.UseHubGateway = *input.UseHubGateway
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func flattenHubModel(inputList *[]connectivityconfigurations.Hub) ([]HubModel, error) {
	var outputList []HubModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := HubModel{}

		if input.ResourceId != nil {
			output.ResourceId = *input.ResourceId
		}

		if input.ResourceType != nil {
			output.ResourceType = *input.ResourceType
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}
