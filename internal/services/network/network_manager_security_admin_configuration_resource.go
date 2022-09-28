package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/utils"

	virtualNetworkManager "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2022-01-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagerSecurityAdminConfigurationModel struct {
	Name                                    string   `tfschema:"name"`
	NetworkNetworkManagerId                 string   `tfschema:"network_network_manager_id"`
	ApplyOnNetworkIntentPolicyBasedServices []string `tfschema:"apply_on_network_intent_policy_based_services"`
	Description                             string   `tfschema:"description"`
}

type ManagerSecurityAdminConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = ManagerSecurityAdminConfigurationResource{}

func (r ManagerSecurityAdminConfigurationResource) ResourceType() string {
	return "azurerm_network_manager_security_admin_configuration"
}

func (r ManagerSecurityAdminConfigurationResource) ModelObject() interface{} {
	return &ManagerSecurityAdminConfigurationModel{}
}

func (r ManagerSecurityAdminConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerSecurityAdminConfigurationID
}

func (r ManagerSecurityAdminConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
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

		"apply_on_network_intent_policy_based_services": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualNetworkManager.IntentPolicyBasedServiceNone),
					string(virtualNetworkManager.IntentPolicyBasedServiceAll),
				}, false),
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ManagerSecurityAdminConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerSecurityAdminConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerSecurityAdminConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ManagerSecurityAdminConfClient
			networkManagerId, err := parse.NetworkManagerID(model.NetworkNetworkManagerId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerSecurityAdminConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroup, networkManagerId.Name, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			conf := &virtualNetworkManager.SecurityAdminConfiguration{
				SecurityAdminConfigurationPropertiesFormat: &virtualNetworkManager.SecurityAdminConfigurationPropertiesFormat{},
			}

			applyOnNetworkIntentPolicyBasedServicesValue, err := expandNetworkIntentPolicyBasedServiceModel(model.ApplyOnNetworkIntentPolicyBasedServices)
			if err != nil {
				return err
			}

			conf.SecurityAdminConfigurationPropertiesFormat.ApplyOnNetworkIntentPolicyBasedServices = applyOnNetworkIntentPolicyBasedServicesValue

			if model.Description != "" {
				conf.SecurityAdminConfigurationPropertiesFormat.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, *conf, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerSecurityAdminConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerSecurityAdminConfClient

			id, err := parse.NetworkManagerSecurityAdminConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerSecurityAdminConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.SecurityAdminConfigurationPropertiesFormat
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("apply_on_network_intent_policy_based_services") {
				applyOnNetworkIntentPolicyBasedServicesValue, err := expandNetworkIntentPolicyBasedServiceModel(model.ApplyOnNetworkIntentPolicyBasedServices)
				if err != nil {
					return err
				}

				properties.ApplyOnNetworkIntentPolicyBasedServices = applyOnNetworkIntentPolicyBasedServicesValue
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Description = &model.Description
				} else {
					properties.Description = nil
				}
			}

			existing.SystemData = nil

			if _, err := client.CreateOrUpdate(ctx, existing, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerSecurityAdminConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerSecurityAdminConfClient

			id, err := parse.NetworkManagerSecurityAdminConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.SecurityAdminConfigurationPropertiesFormat
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ManagerSecurityAdminConfigurationModel{
				Name:                    id.SecurityAdminConfigurationName,
				NetworkNetworkManagerId: parse.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName).ID(),
			}

			applyOnNetworkIntentPolicyBasedServicesValue, err := flattenNetworkIntentPolicyBasedServiceModel(properties.ApplyOnNetworkIntentPolicyBasedServices)
			if err != nil {
				return err
			}

			state.ApplyOnNetworkIntentPolicyBasedServices = applyOnNetworkIntentPolicyBasedServicesValue

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerSecurityAdminConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerSecurityAdminConfClient

			id, err := parse.NetworkManagerSecurityAdminConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, utils.Bool(true))
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

func expandNetworkIntentPolicyBasedServiceModel(inputList []string) (*[]virtualNetworkManager.IntentPolicyBasedService, error) {
	var outputList []virtualNetworkManager.IntentPolicyBasedService
	for _, input := range inputList {
		output := virtualNetworkManager.IntentPolicyBasedService(input)

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func flattenNetworkIntentPolicyBasedServiceModel(inputList *[]virtualNetworkManager.IntentPolicyBasedService) ([]string, error) {
	var outputList []string
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		outputList = append(outputList, string(input))
	}

	return outputList, nil
}
