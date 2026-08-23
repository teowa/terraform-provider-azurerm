package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/nodepools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/supercomputers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = NodePoolResource{}

type NodePoolResource struct{}

type NodePoolModel struct {
	Name            string            `tfschema:"name"`
	SupercomputerId string            `tfschema:"discovery_supercomputer_id"`
	Location        string            `tfschema:"location"`
	Tags            map[string]string `tfschema:"tags"`
}

func (NodePoolResource) ResourceType() string     { return "azurerm_discovery_node_pool" }
func (NodePoolResource) ModelObject() interface{} { return &NodePoolModel{} }
func (NodePoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nodepools.ValidateNodePoolID
}
func (NodePoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name":                       {Type: pluginsdk.TypeString, Required: true, ForceNew: true, ValidateFunc: validation.StringIsNotEmpty},
		"discovery_supercomputer_id": {Type: pluginsdk.TypeString, Required: true, ForceNew: true, ValidateFunc: supercomputers.ValidateSupercomputerID},
		"location":                   commonschema.Location(),
		"tags":                       commonschema.Tags(),
	}
}
func (NodePoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (NodePoolResource) Timeouts() *pluginsdk.ResourceTimeout {
	return &pluginsdk.ResourceTimeout{Create: pluginsdk.DefaultTimeout(30 * time.Minute), Read: pluginsdk.DefaultTimeout(5 * time.Minute), Update: pluginsdk.DefaultTimeout(30 * time.Minute), Delete: pluginsdk.DefaultTimeout(30 * time.Minute)}
}
func (r NodePoolResource) Create() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.create} }
func (r NodePoolResource) Read() sdk.ResourceFunc   { return sdk.ResourceFunc{Func: r.read} }
func (r NodePoolResource) Update() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.update} }
func (r NodePoolResource) Delete() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.delete} }

func (r NodePoolResource) create(ctx context.Context, metadata sdk.ResourceMetaData) error {
	var model NodePoolModel
	if err := metadata.Decode(&model); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}
	supercomputerId, err := supercomputers.ParseSupercomputerID(model.SupercomputerId)
	if err != nil {
		return err
	}
	id := nodepools.NewNodePoolID(supercomputerId.SubscriptionId, supercomputerId.ResourceGroupName, supercomputerId.SupercomputerName, model.Name)
	if metadata.ResourceData.IsNewResource() {
		existing, err := metadata.Client.Discovery.NodePoolsClient.Get(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing node pool `%s`: %+v", id, err)
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return metadata.ResourceRequiresImport(r.ResourceType(), id)
		}
	}
	payload := nodepools.NodePool{Location: model.Location, Tags: &model.Tags}
	if err := metadata.Client.Discovery.NodePoolsClient.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating node pool `%s`: %+v", id, err)
	}
	metadata.SetID(id)
	return nil
}

func (r NodePoolResource) read(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := nodepools.ParseNodePoolID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	resp, err := metadata.Client.Discovery.NodePoolsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			metadata.MarkAsGone(id)
			return nil
		}
		return fmt.Errorf("reading resource `%s`: %+v", id, err)
	}
	output := NodePoolModel{
		Name:     *resp.Model.Name,
		Location: resp.Model.Location,
	}
	if resp.Model.Tags != nil {
		output.Tags = *resp.Model.Tags
	}
	output.SupercomputerId = supercomputers.NewSupercomputerID(id.SubscriptionId, id.ResourceGroupName, id.SupercomputerName).ID()
	return metadata.Encode(&output)
}

func (r NodePoolResource) update(ctx context.Context, metadata sdk.ResourceMetaData) error {
	return r.read(ctx, metadata)
}

func (r NodePoolResource) delete(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := nodepools.ParseNodePoolID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	if err := metadata.Client.Discovery.NodePoolsClient.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting node pool `%s`: %+v", id, err)
	}
	return nil
}
