// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type subnetServiceEndpointConfiguration struct {
	Service             string
	Locations           []string
	NetworkIdentifierID string
}

func subnetServiceEndpointSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"service": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"locations": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"network_identifier_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
		},
		Set: hashSubnetServiceEndpointConfiguration,
	}
}

func subnetServiceEndpointDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"service": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"locations": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"network_identifier_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func hashSubnetServiceEndpointConfiguration(input interface{}) int {
	if values, ok := input.(map[string]interface{}); ok {
		return pluginsdk.HashString(strings.ToLower(values["service"].(string)))
	}

	return 0
}

func expandSubnetServiceEndpointConfigurations(input []interface{}) []subnetServiceEndpointConfiguration {
	configurations := make([]subnetServiceEndpointConfiguration, 0)

	for _, raw := range input {
		configuration := raw.(map[string]interface{})

		locations := make([]string, 0)
		for _, location := range configuration["locations"].([]interface{}) {
			locations = append(locations, location.(string))
		}

		configurations = append(configurations, subnetServiceEndpointConfiguration{
			Service:             configuration["service"].(string),
			Locations:           locations,
			NetworkIdentifierID: configuration["network_identifier_id"].(string),
		})
	}

	return configurations
}

func flattenSubnetServiceEndpointConfiguration(input subnetServiceEndpointConfiguration) map[string]interface{} {
	output := map[string]interface{}{
		"service": input.Service,
	}

	if len(input.Locations) > 0 {
		output["locations"] = input.Locations
	}

	if input.NetworkIdentifierID != "" {
		output["network_identifier_id"] = input.NetworkIdentifierID
	}

	return output
}

func subnetServiceEndpointHasAdditionalConfiguration(input subnetServiceEndpointConfiguration) bool {
	return len(input.Locations) > 0 || input.NetworkIdentifierID != ""
}

func validateSubnetServiceEndpointConfiguration(serviceEndpoints []interface{}, serviceEndpointConfigurations []interface{}) error {
	configuredServices := make(map[string]struct{}, len(serviceEndpoints))

	for _, raw := range serviceEndpoints {
		service := raw.(string)
		configuredServices[strings.ToLower(service)] = struct{}{}
	}

	for _, configuration := range expandSubnetServiceEndpointConfigurations(serviceEndpointConfigurations) {
		if !subnetServiceEndpointHasAdditionalConfiguration(configuration) {
			return fmt.Errorf("service endpoint `%s` must specify at least one of `locations` or `network_identifier_id`", configuration.Service)
		}

		serviceKey := strings.ToLower(configuration.Service)
		if _, exists := configuredServices[serviceKey]; exists {
			return fmt.Errorf("service endpoint `%s` cannot be configured in both `service_endpoints` and `service_endpoint`", configuration.Service)
		}

		configuredServices[serviceKey] = struct{}{}
	}

	return nil
}
