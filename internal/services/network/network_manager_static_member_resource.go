package network

import (
	"context"
	"fmt"
	"time"

	virtualNetworkManager "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2022-01-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerStaticMemberModel struct {
	Name                  string `tfschema:"name"`
	NetworkNetworkGroupId string `tfschema:"network_network_group_id"`
	ResourceId            string `tfschema:"resource_id"`
	Region                string `tfschema:"region"`
}

type ManagerStaticMemberResource struct{}

var _ sdk.ResourceWithUpdate = ManagerStaticMemberResource{}

func (r ManagerStaticMemberResource) ResourceType() string {
	return "azurerm_network_manager_static_member"
}

func (r ManagerStaticMemberResource) ModelObject() interface{} {
	return &ManagerStaticMemberModel{}
}

func (r ManagerStaticMemberResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerStaticMemberID
}

func (r ManagerStaticMemberResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validate.NetworkManagerNetworkGroupID,
		},

		"resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ManagerStaticMemberResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"region": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagerStaticMemberResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerStaticMemberModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ManagerStaticMembersClient
			networkGroupId, err := parse.NetworkManagerNetworkGroupID(model.NetworkNetworkGroupId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerStaticMemberID(networkGroupId.SubscriptionId, networkGroupId.ResourceGroup, networkGroupId.NetworkManagerName, networkGroupId.NetworkGroupName, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			staticMember := &virtualNetworkManager.StaticMember{
				StaticMemberProperties: &virtualNetworkManager.StaticMemberProperties{},
			}

			if model.ResourceId != "" {
				staticMember.StaticMemberProperties.ResourceID = &model.ResourceId
			}

			if _, err := client.CreateOrUpdate(ctx, *staticMember, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerStaticMemberResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerStaticMembersClient

			id, err := parse.NetworkManagerStaticMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerStaticMemberModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.StaticMemberProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("resource_id") {
				if model.ResourceId != "" {
					properties.ResourceID = &model.ResourceId
				} else {
					properties.ResourceID = nil
				}
			}

			existing.SystemData = nil

			if _, err := client.CreateOrUpdate(ctx, existing, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerStaticMemberResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerStaticMembersClient

			id, err := parse.NetworkManagerStaticMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.StaticMemberProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ManagerStaticMemberModel{
				Name:                  id.StaticMemberName,
				NetworkNetworkGroupId: parse.NewNetworkManagerNetworkGroupID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName).ID(),
			}

			if properties.Region != nil {
				state.Region = *properties.Region
			}

			if properties.ResourceID != nil {
				state.ResourceId = *properties.ResourceID
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerStaticMemberResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerStaticMembersClient

			id, err := parse.NetworkManagerStaticMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
