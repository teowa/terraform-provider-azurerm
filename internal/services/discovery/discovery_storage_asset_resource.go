package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storageassets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storagecontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = StorageAssetResource{}

type StorageAssetResource struct{}

type StorageAssetModel struct {
	Name               string            `tfschema:"name"`
	StorageContainerId string            `tfschema:"discovery_storage_container_id"`
	Location           string            `tfschema:"location"`
	Tags               map[string]string `tfschema:"tags"`
}

func (StorageAssetResource) ResourceType() string     { return "azurerm_discovery_storage_asset" }
func (StorageAssetResource) ModelObject() interface{} { return &StorageAssetModel{} }
func (StorageAssetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return storageassets.ValidateStorageAssetID
}
func (StorageAssetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name":                           {Type: pluginsdk.TypeString, Required: true, ForceNew: true, ValidateFunc: validation.StringIsNotEmpty},
		"discovery_storage_container_id": {Type: pluginsdk.TypeString, Required: true, ForceNew: true, ValidateFunc: storagecontainers.ValidateStorageContainerID},
		"location":                       commonschema.Location(),
		"tags":                           commonschema.Tags(),
	}
}
func (StorageAssetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (StorageAssetResource) Timeouts() *pluginsdk.ResourceTimeout {
	return &pluginsdk.ResourceTimeout{Create: pluginsdk.DefaultTimeout(30 * time.Minute), Read: pluginsdk.DefaultTimeout(5 * time.Minute), Update: pluginsdk.DefaultTimeout(30 * time.Minute), Delete: pluginsdk.DefaultTimeout(30 * time.Minute)}
}
func (r StorageAssetResource) Create() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.create} }
func (r StorageAssetResource) Read() sdk.ResourceFunc   { return sdk.ResourceFunc{Func: r.read} }
func (r StorageAssetResource) Update() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.update} }
func (r StorageAssetResource) Delete() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.delete} }

func (r StorageAssetResource) create(ctx context.Context, metadata sdk.ResourceMetaData) error {
	var model StorageAssetModel
	if err := metadata.Decode(&model); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}
	storageContainerId, err := storagecontainers.ParseStorageContainerID(model.StorageContainerId)
	if err != nil {
		return err
	}
	id := storageassets.NewStorageAssetID(storageContainerId.SubscriptionId, storageContainerId.ResourceGroupName, storageContainerId.StorageContainerName, model.Name)
	if metadata.ResourceData.IsNewResource() {
		existing, err := metadata.Client.Discovery.StorageAssetsClient.Get(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing storage asset `%s`: %+v", id, err)
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return metadata.ResourceRequiresImport(r.ResourceType(), id)
		}
	}
	payload := storageassets.StorageAsset{Location: model.Location, Tags: &model.Tags}
	if err := metadata.Client.Discovery.StorageAssetsClient.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating storage asset `%s`: %+v", id, err)
	}
	metadata.SetID(id)
	return nil
}

func (r StorageAssetResource) read(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := storageassets.ParseStorageAssetID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	resp, err := metadata.Client.Discovery.StorageAssetsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			metadata.MarkAsGone(id)
			return nil
		}
		return fmt.Errorf("reading resource `%s`: %+v", id, err)
	}
	output := StorageAssetModel{
		Name:     *resp.Model.Name,
		Location: resp.Model.Location,
	}
	if resp.Model.Tags != nil {
		output.Tags = *resp.Model.Tags
	}
	output.StorageContainerId = storagecontainers.NewStorageContainerID(id.SubscriptionId, id.ResourceGroupName, id.StorageContainerName).ID()
	return metadata.Encode(&output)
}

func (r StorageAssetResource) update(ctx context.Context, metadata sdk.ResourceMetaData) error {
	return r.read(ctx, metadata)
}

func (r StorageAssetResource) delete(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := storageassets.ParseStorageAssetID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	if err := metadata.Client.Discovery.StorageAssetsClient.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting storage asset `%s`: %+v", id, err)
	}
	return nil
}
