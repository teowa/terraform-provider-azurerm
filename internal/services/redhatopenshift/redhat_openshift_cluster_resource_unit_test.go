// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"
)

func TestExpandOpenShiftIdentity(t *testing.T) {
	identityId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/example"

	result, err := expandOpenShiftIdentity([]identity.ModelUserAssigned{
		{
			Type:        identity.TypeUserAssigned,
			IdentityIds: []string{identityId},
		},
	})
	if err != nil {
		t.Fatalf("expanding identity: %+v", err)
	}

	if result == nil {
		t.Fatal("expected identity result, got nil")
	}

	if result.Type != identity.TypeUserAssigned {
		t.Fatalf("expected identity type %q, got %q", identity.TypeUserAssigned, result.Type)
	}

	if _, ok := result.IdentityIds[identityId]; !ok {
		t.Fatalf("expected identity map to contain %q", identityId)
	}
}

func TestFlattenOpenShiftLoadBalancerProfile(t *testing.T) {
	publicIPOne := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-group/providers/Microsoft.Network/publicIPAddresses/ip-b"
	publicIPTwo := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/Example-Group/providers/Microsoft.Network/publicIPAddresses/ip-a"

	count, ids, err := flattenOpenShiftLoadBalancerProfile(&openshiftclusters.LoadBalancerProfile{
		ManagedOutboundIPs: &openshiftclusters.ManagedOutboundIPs{
			Count: pointer.To(int64(2)),
		},
		EffectiveOutboundIPs: &[]openshiftclusters.EffectiveOutboundIP{
			{Id: pointer.To(publicIPOne)},
			{Id: pointer.To(publicIPTwo)},
		},
	})
	if err != nil {
		t.Fatalf("flattening load balancer profile: %+v", err)
	}

	expectedIds := []string{
		commonids.NewPublicIPAddressID("00000000-0000-0000-0000-000000000000", "Example-Group", "ip-a").ID(),
		commonids.NewPublicIPAddressID("00000000-0000-0000-0000-000000000000", "example-group", "ip-b").ID(),
	}

	if count != 2 {
		t.Fatalf("expected managed outbound ip count 2, got %d", count)
	}

	if !reflect.DeepEqual(ids, expectedIds) {
		t.Fatalf("expected ids %#v, got %#v", expectedIds, ids)
	}
}

func TestFlattenOpenShiftWorkloadIdentities(t *testing.T) {
	identities, err := flattenOpenShiftWorkloadIdentities(&openshiftclusters.PlatformWorkloadIdentityProfile{
		PlatformWorkloadIdentities: &map[string]openshiftclusters.PlatformWorkloadIdentity{
			"zeta": {
				ClientId:   pointer.To("client-zeta"),
				ObjectId:   pointer.To("object-zeta"),
				ResourceId: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/zeta"),
			},
			"alpha": {
				ClientId:   pointer.To("client-alpha"),
				ObjectId:   pointer.To("object-alpha"),
				ResourceId: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/Example-Group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/alpha"),
			},
		},
	})
	if err != nil {
		t.Fatalf("flattening workload identities: %+v", err)
	}

	expected := []WorkloadIdentity{
		{
			ClientId:     "client-alpha",
			ObjectId:     "object-alpha",
			OperatorName: "alpha",
			ResourceId:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/Example-Group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/alpha",
		},
		{
			ClientId:     "client-zeta",
			ObjectId:     "object-zeta",
			OperatorName: "zeta",
			ResourceId:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/zeta",
		},
	}

	if !reflect.DeepEqual(identities, expected) {
		t.Fatalf("expected identities %#v, got %#v", expected, identities)
	}
}
