package azurestackhci

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = &StackHCIDeploymentSettingResource{}

type StackHCIDeploymentSettingResource struct{}

type StackHCIDeploymentSettingModel struct{}

func (StackHCIDeploymentSettingResource) ID(model StackHCIDeploymentSettingModel) string {
	return ""
}

func (StackHCIDeploymentSettingResource) ModelObject() interface{} {
	return &StackHCIDeploymentSettingModel{}
}

func (StackHCIDeploymentSettingResource) ResourceType() string {
	return "azurerm_stack_hci_deployment_setting"
}

func (StackHCIDeploymentSettingResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    return validate.ResourceGroupID
}

func (StackHCIDeploymentSettingResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": commonschema.Location(),

		"tags": tags.Schema(),
	}
}

func (StackHCIDeploymentSettingResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCIDeploymentSettingResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// the Timeout is how long Terraform should wait for this function to run before returning an error
		// whilst 30 minutes may initially seem excessive, we set this as a default to account for rate
		// limiting - but having this here means that users can override this in their config as necessary
		Timeout: 30 * time.Minute,

		// the Func returns a function which retrieves the current state of the Resource Group into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
      return nil
		},
	}
}

func (r StackHCIDeploymentSettingResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// the Timeout is how long Terraform should wait for this function to run before returning an error
		// whilst 30 minutes may initially seem excessive, we set this as a default to account for rate
		// limiting - but having this here means that users can override this in their config as necessary
		Timeout: 5 * time.Minute,

		// the Func returns a function which retrieves the current state of the Resource Group into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
      return nil
		},
	}
}


func (r StackHCIDeploymentSettingResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// the Timeout is how long Terraform should wait for this function to run before returning an error
		// whilst 30 minutes may initially seem excessive, we set this as a default to account for rate
		// limiting - but having this here means that users can override this in their config as necessary
		Timeout: 30 * time.Minute,

		// the Func returns a function which retrieves the current state of the Resource Group into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
      return nil
		},
	}
}
