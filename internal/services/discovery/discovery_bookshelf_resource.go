package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = BookshelfResource{}

type BookshelfResource struct{}

type BookshelfModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

func (BookshelfResource) ResourceType() string     { return "azurerm_discovery_bookshelf" }
func (BookshelfResource) ModelObject() interface{} { return &BookshelfModel{} }
func (BookshelfResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return bookshelves.ValidateBookshelfID
}
func (BookshelfResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name":                {Type: pluginsdk.TypeString, Required: true, ForceNew: true, ValidateFunc: validation.StringIsNotEmpty},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"tags":                commonschema.Tags(),
	}
}
func (BookshelfResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (BookshelfResource) Timeouts() *pluginsdk.ResourceTimeout {
	return &pluginsdk.ResourceTimeout{Create: pluginsdk.DefaultTimeout(30 * time.Minute), Read: pluginsdk.DefaultTimeout(5 * time.Minute), Update: pluginsdk.DefaultTimeout(30 * time.Minute), Delete: pluginsdk.DefaultTimeout(30 * time.Minute)}
}
func (r BookshelfResource) Create() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.create} }
func (r BookshelfResource) Read() sdk.ResourceFunc   { return sdk.ResourceFunc{Func: r.read} }
func (r BookshelfResource) Update() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.update} }
func (r BookshelfResource) Delete() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.delete} }

func (r BookshelfResource) create(ctx context.Context, metadata sdk.ResourceMetaData) error {
	var model BookshelfModel
	if err := metadata.Decode(&model); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}
	id := bookshelves.NewBookshelfID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
	if metadata.ResourceData.IsNewResource() {
		existing, err := metadata.Client.Discovery.BookshelvesClient.Get(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing bookshelf `%s`: %+v", id, err)
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return metadata.ResourceRequiresImport(r.ResourceType(), id)
		}
	}
	payload := bookshelves.Bookshelf{Location: model.Location, Tags: &model.Tags}
	if err := metadata.Client.Discovery.BookshelvesClient.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating bookshelf `%s`: %+v", id, err)
	}
	metadata.SetID(id)
	return nil
}

func (r BookshelfResource) read(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := bookshelves.ParseBookshelfID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	resp, err := metadata.Client.Discovery.BookshelvesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			metadata.MarkAsGone(id)
			return nil
		}
		return fmt.Errorf("reading resource `%s`: %+v", id, err)
	}
	output := BookshelfModel{
		Name:     *resp.Model.Name,
		Location: resp.Model.Location,
	}
	if resp.Model.Tags != nil {
		output.Tags = *resp.Model.Tags
	}
	output.ResourceGroupName = id.ResourceGroupName
	return metadata.Encode(&output)
}

func (r BookshelfResource) update(ctx context.Context, metadata sdk.ResourceMetaData) error {
	return r.read(ctx, metadata)
}

func (r BookshelfResource) delete(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := bookshelves.ParseBookshelfID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	if err := metadata.Client.Discovery.BookshelvesClient.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting bookshelf `%s`: %+v", id, err)
	}
	return nil
}
