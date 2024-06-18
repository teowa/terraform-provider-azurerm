package azurestackhci

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"ressource_group_name"`
	Location          string                 `tfschema:"location"`
	CustomLocationId  string                 `tfschema:"custom_location_id"`
	DNSServers        []string               `tfschema:"dns_servers"`
	Subnet            []StackHCISubnetModel  `tfschema:"subnet"`
	VmSwitchName      string                 `tfschema:"vm_switch_name"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

type StackHCISubnetModel struct {
	AddressPrefix      string                `tfschema:"address_prefix"`
	IpAllocationMethod string                `tfschema:"ip_allocation_method"`
	IpPool             []StackHCIIPPoolModel `tfschema:"ip_pool"`
	Route              []StackHCIRouteModel  `tfschema:"route"`
	Vlan               string                `tfschema:"vlan"`
}

type StackHCIIPPoolModel struct {
	Start string `tfschema:"start"`
	End   string `tfschema:"end"`
}

type StackHCIRouteModel struct {
	Name             string `tfschema:"name"`
	AddressPrefix    string `tfschema:"address_prefix"`
	NextHopIpAddress string `tfschema:"next_hop_ip_address"`
}

func (StackHCILogicalNetworkResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},

		"vm_switch_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns_servers": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"subnet": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_allocation_method": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(logicalnetworks.PossibleValuesForIPAllocationMethodEnum(), false),
					},

					"address_prefix": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsCIDR,
					},

					"ip_pool": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
								"end": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
							},
						},
					},

					"route": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"address_prefix": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.IsCIDR,
								},

								"next_hop_ip_address": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.IsCIDR,
								},
							},
						},
					},

					"vlan": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (StackHCILogicalNetworkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCILogicalNetworkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			var config StackHCILogicalNetworkResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := logicalnetworks.NewLogicalNetworkID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := logicalnetworks.LogicalNetworks{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				ExtendedLocation: &logicalnetworks.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(logicalnetworks.ExtendedLocationTypesCustomLocation),
				},
				Properties: &logicalnetworks.LogicalNetworkProperties{
					DhcpOptions: &logicalnetworks.LogicalNetworkPropertiesDhcpOptions{
						DnsServers: pointer.To(config.DNSServers),
					},
					Subnets: ExpandStackHCILogicalNetworkSubnet(config.Subnet),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r StackHCILogicalNetworkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (r StackHCILogicalNetworkResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (r StackHCILogicalNetworkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func ExpandStackHCILogicalNetworkSubnet(input []StackHCISubnetModel) *[]logicalnetworks.Subnet {
	if len(input) == 0 {
		return nil
	}

	results := make([]logicalnetworks.Subnet, 0)
	for _, v := range input {
		results = append(results, logicalnetworks.Subnet{
			Properties: &logicalnetworks.SubnetPropertiesFormat{
				AddressPrefix:      pointer.To(v.AddressPrefix),
				IPAllocationMethod: pointer.To(logicalnetworks.IPAllocationMethodEnum(v.IpAllocationMethod)),
				IPPools:            ExpandStackHCILogicalNetworkIPPool(v.IpPool),
			},
		})
	}

	return &results
}

func ExpandStackHCILogicalNetworkIPPool(input []StackHCIIPPoolModel) *[]logicalnetworks.IPPool {
	if len(input) == 0 {
		return nil
	}

	results := make([]logicalnetworks.IPPool, 0)
	for _, v := range input {
		results = append(results, logicalnetworks.IPPool{
			Start: pointer.To(v.Start),
			End:   pointer.To(v.End),
		})
	}

	return &results
}
