package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/chatmodeldeployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ChatModelDeploymentResource{}

type ChatModelDeploymentResource struct{}

type ChatModelDeploymentModel struct {
	Name        string            `tfschema:"name"`
	WorkspaceId string            `tfschema:"discovery_workspace_id"`
	Location    string            `tfschema:"location"`
	Tags        map[string]string `tfschema:"tags"`
}

func (ChatModelDeploymentResource) ResourceType() string {
	return "azurerm_discovery_chat_model_deployment"
}
func (ChatModelDeploymentResource) ModelObject() interface{} { return &ChatModelDeploymentModel{} }
func (ChatModelDeploymentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return chatmodeldeployments.ValidateChatModelDeploymentID
}
func (ChatModelDeploymentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name":                   {Type: pluginsdk.TypeString, Required: true, ForceNew: true, ValidateFunc: validation.StringIsNotEmpty},
		"discovery_workspace_id": {Type: pluginsdk.TypeString, Required: true, ForceNew: true, ValidateFunc: workspaces.ValidateWorkspaceID},
		"location":               commonschema.Location(),
		"tags":                   commonschema.Tags(),
	}
}
func (ChatModelDeploymentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (ChatModelDeploymentResource) Timeouts() *pluginsdk.ResourceTimeout {
	return &pluginsdk.ResourceTimeout{Create: pluginsdk.DefaultTimeout(30 * time.Minute), Read: pluginsdk.DefaultTimeout(5 * time.Minute), Update: pluginsdk.DefaultTimeout(30 * time.Minute), Delete: pluginsdk.DefaultTimeout(30 * time.Minute)}
}
func (r ChatModelDeploymentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{Func: r.create}
}
func (r ChatModelDeploymentResource) Read() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.read} }
func (r ChatModelDeploymentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{Func: r.update}
}
func (r ChatModelDeploymentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{Func: r.delete}
}

func (r ChatModelDeploymentResource) create(ctx context.Context, metadata sdk.ResourceMetaData) error {
	var model ChatModelDeploymentModel
	if err := metadata.Decode(&model); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}
	workspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
	if err != nil {
		return err
	}
	id := chatmodeldeployments.NewChatModelDeploymentID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, model.Name)
	if metadata.ResourceData.IsNewResource() {
		existing, err := metadata.Client.Discovery.ChatModelDeploymentsClient.Get(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing chat model deployment `%s`: %+v", id, err)
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return metadata.ResourceRequiresImport(r.ResourceType(), id)
		}
	}
	payload := chatmodeldeployments.ChatModelDeployment{Location: model.Location, Tags: &model.Tags}
	if err := metadata.Client.Discovery.ChatModelDeploymentsClient.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating chat model deployment `%s`: %+v", id, err)
	}
	metadata.SetID(id)
	return nil
}

func (r ChatModelDeploymentResource) read(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := chatmodeldeployments.ParseChatModelDeploymentID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	resp, err := metadata.Client.Discovery.ChatModelDeploymentsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			metadata.MarkAsGone(id)
			return nil
		}
		return fmt.Errorf("reading resource `%s`: %+v", id, err)
	}
	output := ChatModelDeploymentModel{
		Name:     *resp.Model.Name,
		Location: resp.Model.Location,
	}
	if resp.Model.Tags != nil {
		output.Tags = *resp.Model.Tags
	}
	output.WorkspaceId = workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID()
	return metadata.Encode(&output)
}

func (r ChatModelDeploymentResource) update(ctx context.Context, metadata sdk.ResourceMetaData) error {
	return r.read(ctx, metadata)
}

func (r ChatModelDeploymentResource) delete(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := chatmodeldeployments.ParseChatModelDeploymentID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	if err := metadata.Client.Discovery.ChatModelDeploymentsClient.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting chat model deployment `%s`: %+v", id, err)
	}
	return nil
}
