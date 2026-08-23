package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ProjectResource{}

type ProjectResource struct{}

type ProjectModel struct {
	Name        string            `tfschema:"name"`
	WorkspaceId string            `tfschema:"discovery_workspace_id"`
	Location    string            `tfschema:"location"`
	Tags        map[string]string `tfschema:"tags"`
}

func (ProjectResource) ResourceType() string     { return "azurerm_discovery_project" }
func (ProjectResource) ModelObject() interface{} { return &ProjectModel{} }
func (ProjectResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return projects.ValidateProjectID
}
func (ProjectResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name":                   {Type: pluginsdk.TypeString, Required: true, ForceNew: true, ValidateFunc: validation.StringIsNotEmpty},
		"discovery_workspace_id": {Type: pluginsdk.TypeString, Required: true, ForceNew: true, ValidateFunc: workspaces.ValidateWorkspaceID},
		"location":               commonschema.Location(),
		"tags":                   commonschema.Tags(),
	}
}
func (ProjectResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (ProjectResource) Timeouts() *pluginsdk.ResourceTimeout {
	return &pluginsdk.ResourceTimeout{Create: pluginsdk.DefaultTimeout(30 * time.Minute), Read: pluginsdk.DefaultTimeout(5 * time.Minute), Update: pluginsdk.DefaultTimeout(30 * time.Minute), Delete: pluginsdk.DefaultTimeout(30 * time.Minute)}
}
func (r ProjectResource) Create() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.create} }
func (r ProjectResource) Read() sdk.ResourceFunc   { return sdk.ResourceFunc{Func: r.read} }
func (r ProjectResource) Update() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.update} }
func (r ProjectResource) Delete() sdk.ResourceFunc { return sdk.ResourceFunc{Func: r.delete} }

func (r ProjectResource) create(ctx context.Context, metadata sdk.ResourceMetaData) error {
	var model ProjectModel
	if err := metadata.Decode(&model); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}
	workspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
	if err != nil {
		return err
	}
	id := projects.NewProjectID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, model.Name)
	if metadata.ResourceData.IsNewResource() {
		existing, err := metadata.Client.Discovery.ProjectsClient.Get(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing project `%s`: %+v", id, err)
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return metadata.ResourceRequiresImport(r.ResourceType(), id)
		}
	}
	payload := projects.Project{Location: model.Location, Tags: &model.Tags}
	if err := metadata.Client.Discovery.ProjectsClient.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating project `%s`: %+v", id, err)
	}
	metadata.SetID(id)
	return nil
}

func (r ProjectResource) read(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := projects.ParseProjectID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	resp, err := metadata.Client.Discovery.ProjectsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			metadata.MarkAsGone(id)
			return nil
		}
		return fmt.Errorf("reading resource `%s`: %+v", id, err)
	}
	output := ProjectModel{
		Name:     *resp.Model.Name,
		Location: resp.Model.Location,
	}
	if resp.Model.Tags != nil {
		output.Tags = *resp.Model.Tags
	}
	output.WorkspaceId = workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID()
	return metadata.Encode(&output)
}

func (r ProjectResource) update(ctx context.Context, metadata sdk.ResourceMetaData) error {
	return r.read(ctx, metadata)
}

func (r ProjectResource) delete(ctx context.Context, metadata sdk.ResourceMetaData) error {
	id, err := projects.ParseProjectID(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
	if err := metadata.Client.Discovery.ProjectsClient.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting project `%s`: %+v", id, err)
	}
	return nil
}
