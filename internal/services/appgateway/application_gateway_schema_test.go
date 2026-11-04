// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package appgateway

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestApplicationGatewaySchemaSupportedValues(t *testing.T) {
	resource := resourceApplicationGateway()
	tests := []struct {
		block    string
		field    string
		accepted []string
		rejected []string
	}{
		{
			block:    "request_routing_rule",
			field:    "rule_type",
			accepted: []string{"Basic", "PathBasedRouting"},
			rejected: []string{"AdvancedRouting"},
		},
		{
			block:    "sku",
			field:    "name",
			accepted: []string{"Basic", "Standard_Large", "Standard_Medium", "Standard_Small", "Standard_v2", "WAF_Large", "WAF_Medium", "WAF_v2"},
			rejected: []string{"Basic_v2", "Basic_WAF_v2"},
		},
		{
			block:    "sku",
			field:    "tier",
			accepted: []string{"Basic", "Standard", "Standard_v2", "WAF", "WAF_v2"},
			rejected: []string{"Basic_v2", "Basic_WAF_v2"},
		},
	}

	for _, test := range tests {
		t.Run(test.block+"."+test.field, func(t *testing.T) {
			field := resource.Schema[test.block].Elem.(*pluginsdk.Resource).Schema[test.field]
			for _, value := range test.accepted {
				if _, errors := field.ValidateFunc(value, test.field); len(errors) != 0 {
					t.Errorf("expected %q to be accepted, got %+v", value, errors)
				}
			}
			for _, value := range test.rejected {
				if _, errors := field.ValidateFunc(value, test.field); len(errors) == 0 {
					t.Errorf("expected %q to be rejected", value)
				}
			}
		})
	}
}
