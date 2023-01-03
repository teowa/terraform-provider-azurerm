package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type RedhatOpenShiftClusterId struct {
	SubscriptionId       string
	ResourceGroup        string
	OpenShiftClusterName string
}

func NewRedhatOpenShiftClusterID(subscriptionId, resourceGroup, openShiftClusterName string) RedhatOpenShiftClusterId {
	return RedhatOpenShiftClusterId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		OpenShiftClusterName: openShiftClusterName,
	}
}

func (id RedhatOpenShiftClusterId) String() string {
	segments := []string{
		fmt.Sprintf("Open Shift Cluster Name %q", id.OpenShiftClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Redhat Open Shift Cluster", segmentsStr)
}

func (id RedhatOpenShiftClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RedHatOpenShift/openShiftClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.OpenShiftClusterName)
}

// RedhatOpenShiftClusterID parses a RedhatOpenShiftCluster ID into an RedhatOpenShiftClusterId struct
func RedhatOpenShiftClusterID(input string) (*RedhatOpenShiftClusterId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RedhatOpenShiftClusterId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.OpenShiftClusterName, err = id.PopSegment("openShiftClusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
