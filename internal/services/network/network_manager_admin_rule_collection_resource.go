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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/adminrulecollections"
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

type NetworkAdminRuleCollectionModel struct {
	Name                                string                                 `tfschema:"name"`
	NetworkSecurityAdminConfigurationId string                                 `tfschema:"network_security_admin_configuration_id"`
	AppliesToGroups                     []NetworkManagerSecurityGroupItemModel `tfschema:"applies_to_groups"`
	Description                         string                                 `tfschema:"description"`
}

type NetworkManagerSecurityGroupItemModel struct {
	NetworkGroupId string `tfschema:"network_group_id"`
}

type NetworkAdminRuleCollectionResource struct{}

var _ sdk.ResourceWithUpdate = NetworkAdminRuleCollectionResource{}

func (r NetworkAdminRuleCollectionResource) ResourceType() string {
	return "azurerm_network_admin_rule_collection"
}

func (r NetworkAdminRuleCollectionResource) ModelObject() interface{} {
	return &NetworkAdminRuleCollectionModel{}
}

func (r NetworkAdminRuleCollectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return adminrulecollections.ValidateAdminRuleCollectionID
}

func (r NetworkAdminRuleCollectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_security_admin_configuration_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: securityadminconfigurations.ValidateSecurityAdminConfigurationID,
		},

		"applies_to_groups": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r NetworkAdminRuleCollectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkAdminRuleCollectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkAdminRuleCollectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.AdminRuleCollectionsClient
			securityAdminConfigurationId, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(model.NetworkSecurityAdminConfigurationId)
			if err != nil {
				return err
			}

			id := adminrulecollections.NewRuleCollectionID(securityAdminConfigurationId.SubscriptionId, securityAdminConfigurationId.ResourceGroupName, securityAdminConfigurationId.NetworkManagerName, securityAdminConfigurationId.ConfigurationName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &adminrulecollections.AdminRuleCollection{
				Properties: &adminrulecollections.AdminRuleCollectionPropertiesFormat{},
			}

			appliesToGroupsValue, err := expandNetworkManagerSecurityGroupItemModel(model.AppliesToGroups)
			if err != nil {
				return err
			}

			if appliesToGroupsValue != nil {
				properties.Properties.AppliesToGroups = *appliesToGroupsValue
			}

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

func (r NetworkAdminRuleCollectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.AdminRuleCollectionsClient

			id, err := adminrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkAdminRuleCollectionModel
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
				appliesToGroupsValue, err := expandNetworkManagerSecurityGroupItemModel(model.AppliesToGroups)
				if err != nil {
					return err
				}

				if appliesToGroupsValue != nil {
					properties.Properties.AppliesToGroups = *appliesToGroupsValue
				}
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

func (r NetworkAdminRuleCollectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.AdminRuleCollectionsClient

			id, err := adminrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
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

			state := NetworkAdminRuleCollectionModel{
				Name:                                id.RuleCollectionName,
				NetworkSecurityAdminConfigurationId: securityadminconfigurations.NewSecurityAdminConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.ConfigurationName).ID(),
			}

			if properties := model.Properties; properties != nil {
				appliesToGroupsValue, err := flattenNetworkManagerSecurityGroupItemModel(&properties.AppliesToGroups)
				if err != nil {
					return err
				}

				state.AppliesToGroups = appliesToGroupsValue

				if properties.Description != nil {
					state.Description = *properties.Description
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkAdminRuleCollectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.AdminRuleCollectionsClient

			id, err := adminrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, adminrulecollections.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandNetworkManagerSecurityGroupItemModel(inputList []NetworkManagerSecurityGroupItemModel) (*[]adminrulecollections.NetworkManagerSecurityGroupItem, error) {
	var outputList []adminrulecollections.NetworkManagerSecurityGroupItem
	for _, v := range inputList {
		input := v
		output := adminrulecollections.NetworkManagerSecurityGroupItem{
			NetworkGroupId: input.NetworkGroupId,
		}

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func flattenNetworkManagerSecurityGroupItemModel(inputList *[]adminrulecollections.NetworkManagerSecurityGroupItem) ([]NetworkManagerSecurityGroupItemModel, error) {
	var outputList []NetworkManagerSecurityGroupItemModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := NetworkManagerSecurityGroupItemModel{
			NetworkGroupId: input.NetworkGroupId,
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}
