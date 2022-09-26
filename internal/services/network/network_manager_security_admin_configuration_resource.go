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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/securityadminconfigurations"
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

type NetworkSecurityAdminConfigurationModel struct {
	Name                                    string                                 `tfschema:"name"`
	NetworkNetworkManagerId                 string                                 `tfschema:"network_network_manager_id"`
	ApplyOnNetworkIntentPolicyBasedServices []NetworkIntentPolicyBasedServiceModel `tfschema:"apply_on_network_intent_policy_based_services"`
	Description                             string                                 `tfschema:"description"`
}

type NetworkSecurityAdminConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = NetworkSecurityAdminConfigurationResource{}

func (r NetworkSecurityAdminConfigurationResource) ResourceType() string {
	return "azurerm_network_security_admin_configuration"
}

func (r NetworkSecurityAdminConfigurationResource) ModelObject() interface{} {
	return &NetworkSecurityAdminConfigurationModel{}
}

func (r NetworkSecurityAdminConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return securityadminconfigurations.ValidateSecurityAdminConfigurationID
}

func (r NetworkSecurityAdminConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
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

		"apply_on_network_intent_policy_based_services": {
			Type:     pluginsdk.TypeList,
			Optional: true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r NetworkSecurityAdminConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkSecurityAdminConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkSecurityAdminConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.SecurityAdminConfigurationsClient
			networkManagerId, err := networkmanagers.ParseNetworkManagerID(model.NetworkNetworkManagerId)
			if err != nil {
				return err
			}

			id := securityadminconfigurations.NewSecurityAdminConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &securityadminconfigurations.SecurityAdminConfiguration{
				Properties: &securityadminconfigurations.SecurityAdminConfigurationPropertiesFormat{},
			}

			applyOnNetworkIntentPolicyBasedServicesValue, err := expandNetworkIntentPolicyBasedServiceModel(model.ApplyOnNetworkIntentPolicyBasedServices)
			if err != nil {
				return err
			}

			properties.Properties.ApplyOnNetworkIntentPolicyBasedServices = applyOnNetworkIntentPolicyBasedServicesValue

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkSecurityAdminConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.SecurityAdminConfigurationsClient

			id, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkSecurityAdminConfigurationModel
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

			if metadata.ResourceData.HasChange("apply_on_network_intent_policy_based_services") {
				applyOnNetworkIntentPolicyBasedServicesValue, err := expandNetworkIntentPolicyBasedServiceModel(model.ApplyOnNetworkIntentPolicyBasedServices)
				if err != nil {
					return err
				}

				properties.Properties.ApplyOnNetworkIntentPolicyBasedServices = applyOnNetworkIntentPolicyBasedServicesValue
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Properties.Description = &model.Description
				} else {
					properties.Properties.Description = nil
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

func (r NetworkSecurityAdminConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.SecurityAdminConfigurationsClient

			id, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(metadata.ResourceData.Id())
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

			state := NetworkSecurityAdminConfigurationModel{
				Name:                    id.ConfigurationName,
				NetworkNetworkManagerId: networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
			}

			if properties := model.Properties; properties != nil {
				applyOnNetworkIntentPolicyBasedServicesValue, err := flattenNetworkIntentPolicyBasedServiceModel(properties.ApplyOnNetworkIntentPolicyBasedServices)
				if err != nil {
					return err
				}

				state.ApplyOnNetworkIntentPolicyBasedServices = applyOnNetworkIntentPolicyBasedServicesValue

				if properties.Description != nil {
					state.Description = *properties.Description
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkSecurityAdminConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.SecurityAdminConfigurationsClient

			id, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, securityadminconfigurations.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandNetworkIntentPolicyBasedServiceModel(inputList []NetworkIntentPolicyBasedServiceModel) (*[]securityadminconfigurations.NetworkIntentPolicyBasedService, error) {
	var outputList []securityadminconfigurations.NetworkIntentPolicyBasedService
	for _, v := range inputList {
		input := v
		output := securityadminconfigurations.NetworkIntentPolicyBasedService{}

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func flattenNetworkIntentPolicyBasedServiceModel(inputList *[]securityadminconfigurations.NetworkIntentPolicyBasedService) ([]NetworkIntentPolicyBasedServiceModel, error) {
	var outputList []NetworkIntentPolicyBasedServiceModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := NetworkIntentPolicyBasedServiceModel{}

		outputList = append(outputList, output)
	}

	return outputList, nil
}
