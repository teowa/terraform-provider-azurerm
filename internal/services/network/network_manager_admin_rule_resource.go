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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/adminrules"
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

type NetworkAdminRuleModel struct {
	Name                    string                   `tfschema:"name"`
	NetworkRuleCollectionId string                   `tfschema:"network_rule_collection_id"`
	Kind                    adminrules.AdminRuleKind `tfschema:"kind"`
}

type NetworkAdminRuleResource struct{}

var _ sdk.ResourceWithUpdate = NetworkAdminRuleResource{}

func (r NetworkAdminRuleResource) ResourceType() string {
	return "azurerm_network_admin_rule"
}

func (r NetworkAdminRuleResource) ModelObject() interface{} {
	return &NetworkAdminRuleModel{}
}

func (r NetworkAdminRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return adminrules.ValidateAdminRuleID
}

func (r NetworkAdminRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_rule_collection_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: adminrulecollections.ValidateRuleCollectionID,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(adminrules.AdminRuleKindCustom),
				string(adminrules.AdminRuleKindDefault),
			}, false),
		},
	}
}

func (r NetworkAdminRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkAdminRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkAdminRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.AdminRulesClient
			ruleCollectionId, err := adminrulecollections.ParseRuleCollectionID(model.NetworkRuleCollectionId)
			if err != nil {
				return err
			}

			id := adminrules.NewRuleID(ruleCollectionId.SubscriptionId, ruleCollectionId.ResourceGroupName, ruleCollectionId.NetworkManagerName, ruleCollectionId.ConfigurationName, ruleCollectionId.RuleCollectionName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &adminrules.BaseAdminRule{
				Kind: model.Kind,
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkAdminRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.AdminRulesClient

			id, err := adminrules.ParseRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkAdminRuleModel
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

			if metadata.ResourceData.HasChange("kind") {
				properties.Kind = model.Kind
			}

			properties.SystemData = nil

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NetworkAdminRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.AdminRulesClient

			id, err := adminrules.ParseRuleID(metadata.ResourceData.Id())
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

			state := NetworkAdminRuleModel{
				Name:                    id.RuleName,
				NetworkRuleCollectionId: adminrulecollections.NewRuleCollectionID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.ConfigurationName, id.RuleCollectionName).ID(),
				Kind:                    model.Kind,
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkAdminRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.AdminRulesClient

			id, err := adminrules.ParseRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, adminrules.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
