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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/networkgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-01-01/staticmembers"
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

type NetworkStaticMemberModel struct {
	Name                  string `tfschema:"name"`
	NetworkNetworkGroupId string `tfschema:"network_network_group_id"`
	ResourceId            string `tfschema:"resource_id"`
	Region                string `tfschema:"region"`
}

type NetworkStaticMemberResource struct{}

var _ sdk.ResourceWithUpdate = NetworkStaticMemberResource{}

func (r NetworkStaticMemberResource) ResourceType() string {
	return "azurerm_network_static_member"
}

func (r NetworkStaticMemberResource) ModelObject() interface{} {
	return &NetworkStaticMemberModel{}
}

func (r NetworkStaticMemberResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticmembers.ValidateStaticMemberID
}

func (r NetworkStaticMemberResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_network_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkgroups.ValidateNetworkGroupID,
		},

		"resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r NetworkStaticMemberResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"region": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r NetworkStaticMemberResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkStaticMemberModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.StaticMembersClient
			networkGroupId, err := networkgroups.ParseNetworkGroupID(model.NetworkNetworkGroupId)
			if err != nil {
				return err
			}

			id := staticmembers.NewStaticMemberID(networkGroupId.SubscriptionId, networkGroupId.ResourceGroupName, networkGroupId.NetworkManagerName, networkGroupId.NetworkGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &staticmembers.StaticMember{
				Properties: &staticmembers.StaticMemberProperties{},
			}

			if model.ResourceId != "" {
				properties.Properties.ResourceId = &model.ResourceId
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkStaticMemberResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.StaticMembersClient

			id, err := staticmembers.ParseStaticMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkStaticMemberModel
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

			if metadata.ResourceData.HasChange("resource_id") {
				if model.ResourceId != "" {
					properties.Properties.ResourceId = &model.ResourceId
				} else {
					properties.Properties.ResourceId = nil
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

func (r NetworkStaticMemberResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.StaticMembersClient

			id, err := staticmembers.ParseStaticMemberID(metadata.ResourceData.Id())
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

			state := NetworkStaticMemberModel{
				Name:                  id.StaticMemberName,
				NetworkNetworkGroupId: networkgroups.NewNetworkGroupID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.NetworkGroupName).ID(),
			}

			if properties := model.Properties; properties != nil {
				if properties.Region != nil {
					state.Region = *properties.Region
				}

				if properties.ResourceId != nil {
					state.ResourceId = *properties.ResourceId
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkStaticMemberResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.StaticMembersClient

			id, err := staticmembers.ParseStaticMemberID(metadata.ResourceData.Id())
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
