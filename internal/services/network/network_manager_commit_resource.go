package network

import (
	"context"
	"fmt"
	"time"

	networkManager "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2022-01-01/network"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagerCommitModel struct {
	NetworkManagerId string   `tfschema:"network_manager_id"`
	ScopeAccess      string   `tfschema:"scope_access"`
	Location         string   `tfschema:"location"`
	ConfigurationIds []string `tfschema:"configuration_ids"`
	DeploymentStatus string   `tfschema:"deployment_status"`
}

type ManagerCommitResource struct{}

func (r ManagerCommitResource) ResourceType() string {
	return "azurerm_network_manager_commit"
}

func (r ManagerCommitResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerCommitID
}

func (r ManagerCommitResource) ModelObject() interface{} {
	return &ManagerCommitModel{}
}

func (r ManagerCommitResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.NetworkManagerID,
		},

		"location": commonschema.Location(),

		"scope_access": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(networkManager.ConfigurationTypeConnectivity),
				string(networkManager.ConfigurationTypeSecurityAdmin),
			}, false),
		},

		"configuration_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r ManagerCommitResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"deployment_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagerCommitResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state ManagerCommitModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Network.ManagerCommitsClient
			statusClient := metadata.Client.Network.ManagerDeploymentStatusClient
			networkManagerId, err := parse.NetworkManagerID(state.NetworkManagerId)
			if err != nil {
				return err
			}
			id := parse.NewNetworkManagerCommitID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroup, networkManagerId.Name, state.Location, state.ScopeAccess)

			metadata.Logger.Infof("creating %s", id)

			if !metadata.Client.Features.Network.ManagerOverwriteCommitted {
				listParam := networkManager.ManagerDeploymentStatusParameter{
					Regions:         &[]string{azure.NormalizeLocation(state.Location)},
					DeploymentTypes: &[]networkManager.ConfigurationType{networkManager.ConfigurationType(state.ScopeAccess)},
				}
				existing, err := statusClient.List(ctx, listParam, id.ResourceGroup, id.NetworkManagerName)
				if err != nil {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
				if existing.Value != nil && len(*existing.Value) > 0 {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := networkManager.ManagerCommit{
				ConfigurationIds: &state.ConfigurationIds,
				TargetLocations:  &[]string{state.Location},
				CommitType:       networkManager.ConfigurationType(state.ScopeAccess),
			}

			if _, err := client.Post(ctx, input, id.ResourceGroup, id.NetworkManagerName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ManagerCommitResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerDeploymentStatusClient
			id, err := parse.NetworkManagerCommitID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)

			listParam := networkManager.ManagerDeploymentStatusParameter{
				Regions:         &[]string{azure.NormalizeLocation(id.Location)},
				DeploymentTypes: &[]networkManager.ConfigurationType{networkManager.ConfigurationType(id.ScopeAccess)},
			}
			resp, err := client.List(ctx, listParam, id.ResourceGroup, id.NetworkManagerName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Value == nil || len(*resp.Value) == 0 {
				metadata.Logger.Infof("%s was not found - removing from state!", *id)
				return metadata.MarkAsGone(id)
			}

			commit := (*resp.Value)[0]
			if commit.ConfigurationIds == nil {
				return fmt.Errorf("retrieving %s error null configuration ID of commit", *id)
			}
			return metadata.Encode(&ManagerCommitModel{
				NetworkManagerId: parse.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName).ID(),
				Location:         location.NormalizeNilable(commit.Region),
				ScopeAccess:      string(commit.DeploymentType),
				ConfigurationIds: *commit.ConfigurationIds,
				DeploymentStatus: string(commit.DeploymentStatus),
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ManagerCommitResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.NetworkManagerCommitID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Network.ManagerCommitsClient
			statusClient := metadata.Client.Network.ManagerDeploymentStatusClient

			listParam := networkManager.ManagerDeploymentStatusParameter{
				Regions:         &[]string{azure.NormalizeLocation(id.Location)},
				DeploymentTypes: &[]networkManager.ConfigurationType{networkManager.ConfigurationType(id.ScopeAccess)},
			}
			resp, err := statusClient.List(ctx, listParam, id.ResourceGroup, id.NetworkManagerName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Value == nil || len(*resp.Value) == 0 {
				metadata.Logger.Infof("%s was not found - removing from state!", *id)
				return metadata.MarkAsGone(id)
			}

			commit := (*resp.Value)[0]
			if commit.ConfigurationIds == nil {
				return fmt.Errorf("unexpected null configuration ID of %s", *id)
			}

			var state ManagerCommitModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("configuration_ids") {
				commit.ConfigurationIds = &state.ConfigurationIds
			}

			input := networkManager.ManagerCommit{
				ConfigurationIds: commit.ConfigurationIds,
				TargetLocations:  &[]string{state.Location},
				CommitType:       networkManager.ConfigurationType(state.ScopeAccess),
			}

			if _, err := client.Post(ctx, input, id.ResourceGroup, id.NetworkManagerName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			return nil

		},
		Timeout: 30 * time.Minute,
	}
}

func (r ManagerCommitResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerCommitsClient
			id, err := parse.NetworkManagerCommitID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)

			if !metadata.Client.Features.Network.ManagerKeepCommittedOnDestroy {
				input := networkManager.ManagerCommit{
					ConfigurationIds: &[]string{},
					TargetLocations:  &[]string{id.Location},
					CommitType:       networkManager.ConfigurationType(id.ScopeAccess),
				}

				future, err := client.Post(ctx, input, id.ResourceGroup, id.NetworkManagerName)
				if err != nil {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}

				if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
					return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
				}
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
