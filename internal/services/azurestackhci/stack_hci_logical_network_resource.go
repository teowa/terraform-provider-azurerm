package azurestackhci

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = StackHCILogicalNetworkResource{}
	_ sdk.ResourceWithUpdate = StackHCILogicalNetworkResource{}
)

type StackHCILogicalNetworkResource struct{}

func (StackHCILogicalNetworkResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return logicalnetworks.ValidateLogicalNetworkID
}

func (StackHCILogicalNetworkResource) ResourceType() string {
	return "azurerm_stack_hci_logical_network"
}

func (StackHCILogicalNetworkResource) ModelObject() interface{} {
	return &StackHCIDeploymentSettingModel{}
}

type StackHCILogicalNetworkResourceModel struct {
	Name string `tfschema:"name"`
  ResourceGroupName string `tfschema:"ressource_group_name"`
  Location  string `tfschema:"location"`
  CustomLocationId string `tfschema:"customLocationId"`
  DHCPOption []StackHCIDHCPOptionModel `tfschema:"dhcp_option"`
  Subnets []StackHCISubnetModel `tfschema:"subnet"`


}

type StackHCIDHCPOptionModel struct {
  DNSServers []string `tfschema:"dns_servers"`
}

type StackHCISubnetModel struct {
}

func (StackHCILogicalNetworkResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (StackHCILogicalNetworkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCILogicalNetworkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			var config StackStackHCILogicalNetworkResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			stackHCILogicalNetworkID, err := logicalnetworks.ParseLogicalNetworkID(config.)
			if err != nil {
				return err
			}
		},
	}
}

func (r StackHCILogicalNetworkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func:    func(ctx context.Context, metadata sdk.ResourceMetaData) error {},
	}
}

func (r StackHCILogicalNetworkResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func:    func(ctx context.Context, metadata sdk.ResourceMetaData) error {},
	}
}

func (r StackHCILogicalNetworkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func:    func(ctx context.Context, metadata sdk.ResourceMetaData) error {},
	}
}
