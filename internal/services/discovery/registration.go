package discovery

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

func (Registration) Name() string                  { return "Discovery" }
func (Registration) AssociatedGitHubLabel() string { return "service/discovery" }
func (Registration) WebsiteCategories() []string   { return []string{"Discovery"} }
func (Registration) DataSources() []sdk.DataSource { return []sdk.DataSource{} }
func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		BookshelfResource{},
		ChatModelDeploymentResource{},
		NodePoolResource{},
		ProjectResource{},
		StorageAssetResource{},
		StorageContainerResource{},
		SupercomputerResource{},
		ToolResource{},
		WorkspaceResource{},
	}
}
func (Registration) Actions() []func() action.Action { return []func() action.Action{} }
func (Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}
func (Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}
func (Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}
func (Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
