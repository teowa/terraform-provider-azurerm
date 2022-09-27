package network

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	virtualNetworkManager "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2022-01-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagerConnectivityConfigurationModel struct {
	Name                    string                                     `tfschema:"name"`
	NetworkNetworkManagerId string                                     `tfschema:"network_network_manager_id"`
	AppliesToGroups         []ConnectivityGroupItemModel               `tfschema:"applies_to_groups"`
	ConnectivityTopology    virtualNetworkManager.ConnectivityTopology `tfschema:"connectivity_topology"`
	DeleteExistingPeering   bool                                       `tfschema:"delete_existing_peering"`
	Description             string                                     `tfschema:"description"`
	Hubs                    []HubModel                                 `tfschema:"hubs"`
	IsGlobal                bool                                       `tfschema:"is_global"`
}

type ConnectivityGroupItemModel struct {
	GroupConnectivity virtualNetworkManager.GroupConnectivity `tfschema:"group_connectivity"`
	IsGlobal          bool                                    `tfschema:"is_global"`
	NetworkGroupId    string                                  `tfschema:"network_group_id"`
	UseHubGateway     bool                                    `tfschema:"use_hub_gateway"`
}

type HubModel struct {
	ResourceId   string `tfschema:"resource_id"`
	ResourceType string `tfschema:"resource_type"`
}

type ManagerConnectivityConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = ManagerConnectivityConfigurationResource{}

func (r ManagerConnectivityConfigurationResource) ResourceType() string {
	return "azurerm_network_manager_connectivity_configuration"
}

func (r ManagerConnectivityConfigurationResource) ModelObject() interface{} {
	return &ManagerConnectivityConfigurationModel{}
}

func (r ManagerConnectivityConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerConnectivityConfigurationID
}

func (r ManagerConnectivityConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validate.NetworkManagerID,
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
							string(virtualNetworkManager.GroupConnectivityNone),
							string(virtualNetworkManager.GroupConnectivityDirectlyConnected),
						}, false),
					},

					"is_global": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(virtualNetworkManager.IsGlobalTrue),
							string(virtualNetworkManager.IsGlobalFalse),
						}, false),
					},

					"network_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"use_hub_gateway": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(virtualNetworkManager.UseHubGatewayFalse),
							string(virtualNetworkManager.UseHubGatewayTrue),
						}, false),
					},
				},
			},
		},

		"connectivity_topology": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualNetworkManager.ConnectivityTopologyHubAndSpoke),
				string(virtualNetworkManager.ConnectivityTopologyMesh),
			}, false),
		},

		"delete_existing_peering": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualNetworkManager.DeleteExistingPeeringFalse),
				string(virtualNetworkManager.DeleteExistingPeeringTrue),
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
				string(virtualNetworkManager.IsGlobalFalse),
				string(virtualNetworkManager.IsGlobalTrue),
			}, false),
		},
	}
}

func (r ManagerConnectivityConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerConnectivityConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerConnectivityConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ManagerConnectivityConfClient
			networkManagerId, err := parse.NetworkManagerID(model.NetworkNetworkManagerId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerConnectivityConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroup, networkManagerId.Name, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			conf := &virtualNetworkManager.ConnectivityConfiguration{
				ConnectivityConfigurationProperties: &virtualNetworkManager.ConnectivityConfigurationProperties{
					ConnectivityTopology:  model.ConnectivityTopology,
					DeleteExistingPeering: expandDeleteExistingPeering(model.DeleteExistingPeering),
					IsGlobal:              expandConnectivityConfIsGlobal(model.IsGlobal),
				},
			}

			appliesToGroupsValue, err := expandConnectivityGroupItemModel(model.AppliesToGroups)
			if err != nil {
				return err
			}

			conf.ConnectivityConfigurationProperties.AppliesToGroups = appliesToGroupsValue

			if model.Description != "" {
				conf.ConnectivityConfigurationProperties.Description = &model.Description
			}

			hubsValue, err := expandHubModel(model.Hubs)
			if err != nil {
				return err
			}

			conf.ConnectivityConfigurationProperties.Hubs = hubsValue

			if _, err := client.CreateOrUpdate(ctx, *conf, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerConnectivityConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerConnectivityConfClient

			id, err := parse.NetworkManagerConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerConnectivityConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.ConnectivityConfigurationProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("applies_to_groups") {
				appliesToGroupsValue, err := expandConnectivityGroupItemModel(model.AppliesToGroups)
				if err != nil {
					return err
				}

				properties.AppliesToGroups = appliesToGroupsValue
			}

			if metadata.ResourceData.HasChange("connectivity_topology") {
				properties.ConnectivityTopology = model.ConnectivityTopology
			}

			if metadata.ResourceData.HasChange("delete_existing_peering") {
				properties.DeleteExistingPeering = expandDeleteExistingPeering(model.DeleteExistingPeering)
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Description = &model.Description
				} else {
					properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("hubs") {
				hubsValue, err := expandHubModel(model.Hubs)
				if err != nil {
					return err
				}

				properties.Hubs = hubsValue
			}

			if metadata.ResourceData.HasChange("is_global") {
				properties.IsGlobal = expandConnectivityConfIsGlobal(model.IsGlobal)
			}

			existing.SystemData = nil

			if _, err := client.CreateOrUpdate(ctx, existing, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerConnectivityConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerConnectivityConfClient

			id, err := parse.NetworkManagerConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.ConnectivityConfigurationProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ManagerConnectivityConfigurationModel{
				Name:                    id.ConnectivityConfigurationName,
				NetworkNetworkManagerId: parse.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName).ID(),
			}

			appliesToGroupsValue, err := flattenConnectivityGroupItemModel(properties.AppliesToGroups)
			if err != nil {
				return err
			}

			state.AppliesToGroups = appliesToGroupsValue
			state.ConnectivityTopology = properties.ConnectivityTopology
			state.DeleteExistingPeering = flattenDeleteExistingPeering(properties.DeleteExistingPeering)
			state.IsGlobal = flattenConnectivityConfIsGlobal(properties.IsGlobal)

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			hubsValue, err := flattenHubModel(properties.Hubs)
			if err != nil {
				return err
			}

			state.Hubs = hubsValue

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerConnectivityConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerConnectivityConfClient

			id, err := parse.NetworkManagerConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName, utils.Bool(true))
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDeleteExistingPeering(input bool) virtualNetworkManager.DeleteExistingPeering {
	if input {
		return virtualNetworkManager.DeleteExistingPeeringTrue
	}
	return virtualNetworkManager.DeleteExistingPeeringFalse
}

func expandConnectivityConfIsGlobal(input bool) virtualNetworkManager.IsGlobal {
	if input {
		return virtualNetworkManager.IsGlobalTrue
	}
	return virtualNetworkManager.IsGlobalFalse
}

func expandConnectivityGroupItemModel(inputList []ConnectivityGroupItemModel) (*[]virtualNetworkManager.ConnectivityGroupItem, error) {
	var outputList []virtualNetworkManager.ConnectivityGroupItem
	for _, v := range inputList {
		input := v
		output := virtualNetworkManager.ConnectivityGroupItem{
			GroupConnectivity: input.GroupConnectivity,
			IsGlobal:          expandConnectivityConfIsGlobal(input.IsGlobal),
			NetworkGroupID:    utils.String(input.NetworkGroupId),
			UseHubGateway:     expandUseHubGateWay(input.UseHubGateway),
		}

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func expandUseHubGateWay(input bool) virtualNetworkManager.UseHubGateway {
	if input {
		return virtualNetworkManager.UseHubGatewayTrue
	}
	return virtualNetworkManager.UseHubGatewayFalse
}

func expandHubModel(inputList []HubModel) (*[]virtualNetworkManager.Hub, error) {
	var outputList []virtualNetworkManager.Hub
	for _, v := range inputList {
		input := v
		output := virtualNetworkManager.Hub{}

		if input.ResourceId != "" {
			output.ResourceID = &input.ResourceId
		}

		if input.ResourceType != "" {
			output.ResourceType = &input.ResourceType
		}

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func flattenDeleteExistingPeering(input virtualNetworkManager.DeleteExistingPeering) bool {
	if input == virtualNetworkManager.DeleteExistingPeeringTrue {
		return true
	}
	return false
}

func flattenConnectivityConfIsGlobal(input virtualNetworkManager.IsGlobal) bool {
	if input == virtualNetworkManager.IsGlobalTrue {
		return true
	}
	return false
}

func flattenConnectivityGroupItemModel(inputList *[]virtualNetworkManager.ConnectivityGroupItem) ([]ConnectivityGroupItemModel, error) {
	var outputList []ConnectivityGroupItemModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := ConnectivityGroupItemModel{
			GroupConnectivity: input.GroupConnectivity,
			UseHubGateway:     flattenUseHubGateWay(input.UseHubGateway),
			IsGlobal:          flattenConnectivityConfIsGlobal(input.IsGlobal),
		}

		if input.NetworkGroupID != nil {
			output.NetworkGroupId = *input.NetworkGroupID
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func flattenUseHubGateWay(input virtualNetworkManager.UseHubGateway) bool {
	if input == virtualNetworkManager.UseHubGatewayTrue {
		return true
	}
	return false
}

func flattenHubModel(inputList *[]virtualNetworkManager.Hub) ([]HubModel, error) {
	var outputList []HubModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := HubModel{}

		if input.ResourceID != nil {
			output.ResourceId = *input.ResourceID
		}

		if input.ResourceType != nil {
			output.ResourceType = *input.ResourceType
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}
