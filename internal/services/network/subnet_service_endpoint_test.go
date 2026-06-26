// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/subnets"
)

func TestExpandSubnetServiceEndpoints(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"service":               "Microsoft.Sql",
			"locations":             []interface{}{"eastus2"},
			"network_identifier_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/virtualNetworks/virtualNetwork1",
		},
	}

	result := expandSubnetServiceEndpoints([]interface{}{"Microsoft.Storage"}, input)
	if result == nil {
		t.Fatalf("expected service endpoints to be expanded")
	}

	if len(*result) != 2 {
		t.Fatalf("expected 2 service endpoints, got %d", len(*result))
	}

	if got := *(*result)[0].Service; got != "Microsoft.Storage" {
		t.Fatalf("expected first service endpoint to be Microsoft.Storage, got %s", got)
	}

	detailed := (*result)[1]
	if detailed.Service == nil || *detailed.Service != "Microsoft.Sql" {
		t.Fatalf("expected second service endpoint to be Microsoft.Sql, got %+v", detailed.Service)
	}

	if detailed.Locations == nil || !reflect.DeepEqual(*detailed.Locations, []string{"eastus2"}) {
		t.Fatalf("expected locations to be [eastus2], got %+v", detailed.Locations)
	}

	if detailed.NetworkIdentifier == nil || detailed.NetworkIdentifier.Id == nil || *detailed.NetworkIdentifier.Id == "" {
		t.Fatalf("expected network identifier id to be set")
	}
}

func TestFlattenSubnetServiceEndpoints(t *testing.T) {
	input := []subnets.ServiceEndpointPropertiesFormat{
		{
			Service: stringPointer("Microsoft.Storage"),
		},
		{
			Service:           stringPointer("Microsoft.Sql"),
			Locations:         &[]string{"eastus2"},
			NetworkIdentifier: &subnets.SubResource{Id: stringPointer("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/virtualNetworks/virtualNetwork1")},
		},
	}

	simple, detailed := flattenSubnetServiceEndpoints(&input)

	if !reflect.DeepEqual(simple, []interface{}{"Microsoft.Storage"}) {
		t.Fatalf("expected simple service endpoints to contain Microsoft.Storage, got %+v", simple)
	}

	expectedDetailed := []interface{}{
		map[string]interface{}{
			"service":               "Microsoft.Sql",
			"locations":             []string{"eastus2"},
			"network_identifier_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/virtualNetworks/virtualNetwork1",
		},
	}
	if !reflect.DeepEqual(detailed, expectedDetailed) {
		t.Fatalf("unexpected detailed service endpoint output: %+v", detailed)
	}
}

func TestValidateSubnetServiceEndpointConfiguration(t *testing.T) {
	testData := []struct {
		name                         string
		serviceEndpoints             []interface{}
		serviceEndpointConfiguration []interface{}
		expectedError                string
	}{
		{
			name:             "requires additional configuration",
			serviceEndpoints: []interface{}{},
			serviceEndpointConfiguration: []interface{}{
				map[string]interface{}{
					"service":               "Microsoft.Sql",
					"locations":             []interface{}{},
					"network_identifier_id": "",
				},
			},
			expectedError: "service endpoint `Microsoft.Sql` must specify at least one of `locations` or `network_identifier_id`",
		},
		{
			name:             "disallows duplicate service across schemas",
			serviceEndpoints: []interface{}{"Microsoft.Sql"},
			serviceEndpointConfiguration: []interface{}{
				map[string]interface{}{
					"service":               "Microsoft.Sql",
					"locations":             []interface{}{"eastus2"},
					"network_identifier_id": "",
				},
			},
			expectedError: "service endpoint `Microsoft.Sql` cannot be configured in both `service_endpoints` and `service_endpoint`",
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			err := validateSubnetServiceEndpointConfiguration(test.serviceEndpoints, test.serviceEndpointConfiguration)
			if err == nil {
				t.Fatalf("expected an error")
			}

			if err.Error() != test.expectedError {
				t.Fatalf("expected error %q, got %q", test.expectedError, err.Error())
			}
		})
	}
}

func stringPointer(input string) *string {
	return &input
}
