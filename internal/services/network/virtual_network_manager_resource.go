package network

import (
	"context"
	"fmt"
	virtualNetworkManager "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2022-01-01/network"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	managementGroupValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"
)

type VirtualNetworkManagerModel struct {
	CrossTenantScopes []CrossTenantScope     `tfschema:"cross_tenant_scopes"`
	Scope             []Scope                `tfschema:"scope"`
	ScopeAccess       []string               `tfschema:"scope_access"`
	Description       string                 `tfschema:"description"`
	Name              string                 `tfschema:"name"`
	Location          string                 `tfschema:"location"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

type Scope struct {
	ManagementGroups []string `tfschema:"management_group_ids"`
	Subscriptions    []string `tfschema:"subscription_ids"`
}

type CrossTenantScope struct {
	TenantId         string   `tfschema:"tenant_id"`
	ManagementGroups []string `tfschema:"management_groups"`
	Subscriptions    []string `tfschema:"subscriptions"`
}

type VirtualNetworkManagerResource struct{}

func (r VirtualNetworkManagerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"scope": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"management_group_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: managementGroupValidate.ManagementGroupID,
						},
						AtLeastOneOf: []string{"scope.0.management_group_ids", "scope.0.subscription_ids"},
					},
					"subscription_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: commonids.ValidateSubscriptionID,
						},
						AtLeastOneOf: []string{"scope.0.management_group_ids", "scope.0.subscription_ids"},
					},
				},
			},
		},

		"scope_access": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualNetworkManager.ConfigurationTypeConnectivity),
					string(virtualNetworkManager.ConfigurationTypeSecurityAdmin),
				}, false),
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r VirtualNetworkManagerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cross_tenant_scopes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"subscriptions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
					"management_groups": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r VirtualNetworkManagerResource) ResourceType() string {
	return "azurerm_virtual_network_manager"
}

func (r VirtualNetworkManagerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.VirtualNetworkManagerID
}

func (r VirtualNetworkManagerResource) ModelObject() interface{} {
	return &VirtualNetworkManagerModel{}
}

func (r VirtualNetworkManagerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state VirtualNetworkManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Network.VirtualNetworkManagersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewVirtualNetworkManagerID(subscriptionId, state.ResourceGroupName, state.Name)
			metadata.Logger.Infof("creating %s", id)

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := virtualNetworkManager.Manager{
				Location: utils.String(azure.NormalizeLocation(state.Location)),
				Name:     utils.String(state.Name),
				ManagerProperties: &virtualNetworkManager.ManagerProperties{
					Description:                 utils.String(state.Description),
					NetworkManagerScopes:        expandVirtualNetworkManagerScope(state.Scope),
					NetworkManagerScopeAccesses: expandVirtualNetworkManagerScopeAccess(state.ScopeAccess),
				},
				Tags: tags.Expand(state.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, input, id.ResourceGroup, id.NetworkManagerName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r VirtualNetworkManagerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualNetworkManagersClient
			id, err := parse.VirtualNetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)
			resp, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var description string
			var scope []Scope
			var scopeAccess []string
			if prop := resp.ManagerProperties; prop != nil {
				if prop.Description != nil {
					description = *resp.Description
				}
				scope = flattenVirtualNetworkManagerScope(resp.NetworkManagerScopes)
				scopeAccess = flattenVirtualNetworkManagerScopeAccess(resp.NetworkManagerScopeAccesses)
			}

			return metadata.Encode(&VirtualNetworkManagerModel{
				Description:       description,
				Location:          location.NormalizeNilable(resp.Location),
				Name:              id.NetworkManagerName,
				ResourceGroupName: id.ResourceGroup,
				ScopeAccess:       scopeAccess,
				Scope:             scope,
				Tags:              tags.Flatten(resp.Tags),
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r VirtualNetworkManagerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.VirtualNetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Network.VirtualNetworkManagersClient
			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.ManagerProperties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}
			var state VirtualNetworkManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("description") {
				existing.ManagerProperties.Description = utils.String(state.Description)
			}

			if metadata.ResourceData.HasChange("scope") {
				existing.ManagerProperties.NetworkManagerScopes = expandVirtualNetworkManagerScope(state.Scope)
			}

			if metadata.ResourceData.HasChange("scope_access") {
				existing.ManagerProperties.NetworkManagerScopeAccesses = expandVirtualNetworkManagerScopeAccess(state.ScopeAccess)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.Expand(state.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, existing, id.ResourceGroup, id.NetworkManagerName); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil

		},
		Timeout: 30 * time.Minute,
	}
}

func (r VirtualNetworkManagerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualNetworkManagersClient
			id, err := parse.VirtualNetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			future, err := client.Delete(ctx, id.ResourceGroup, id.NetworkManagerName, utils.Bool(true))
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func stringSlice(input []string) *[]string {
	return &input
}

func expandVirtualNetworkManagerScope(input []Scope) *virtualNetworkManager.ManagerPropertiesNetworkManagerScopes {
	if len(input) == 0 {
		return nil
	}

	return &virtualNetworkManager.ManagerPropertiesNetworkManagerScopes{
		ManagementGroups: stringSlice(input[0].ManagementGroups),
		Subscriptions:    stringSlice(input[0].Subscriptions),
	}
}

func expandVirtualNetworkManagerScopeAccess(input []string) *[]virtualNetworkManager.ConfigurationType {
	if len(input) == 0 {
		return nil
	}

	result := make([]virtualNetworkManager.ConfigurationType, 0)
	for _, v := range input {
		result = append(result, virtualNetworkManager.ConfigurationType(v))
	}
	return &result
}

func flattenStringSlicePtr(input *[]string) []string {
	if input == nil {
		return make([]string, 0)
	}
	return *input
}

func flattenVirtualNetworkManagerScope(input *virtualNetworkManager.ManagerPropertiesNetworkManagerScopes) []Scope {
	if input == nil {
		return make([]Scope, 0)
	}

	return []Scope{{
		ManagementGroups: flattenStringSlicePtr(input.ManagementGroups),
		Subscriptions:    flattenStringSlicePtr(input.Subscriptions),
	}}
}

func flattenVirtualNetworkManagerScopeAccess(input *[]virtualNetworkManager.ConfigurationType) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}
