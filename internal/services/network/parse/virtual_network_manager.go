package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualNetworkManagerId struct {
	SubscriptionId     string
	ResourceGroup      string
	NetworkManagerName string
}

func NewVirtualNetworkManagerID(subscriptionId, resourceGroup, networkManagerName string) VirtualNetworkManagerId {
	return VirtualNetworkManagerId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		NetworkManagerName: networkManagerName,
	}
}

func (id VirtualNetworkManagerId) String() string {
	segments := []string{
		fmt.Sprintf("Network Manager Name %q", id.NetworkManagerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Manager", segmentsStr)
}

func (id VirtualNetworkManagerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)
}

// VirtualNetworkManagerID parses a VirtualNetworkManager ID into an VirtualNetworkManagerId struct
func VirtualNetworkManagerID(input string) (*VirtualNetworkManagerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualNetworkManagerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NetworkManagerName, err = id.PopSegment("networkManagers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
