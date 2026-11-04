// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/expressroute/resources/express_route_circuit"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestExpressRouteCircuitSkuTier(t *testing.T) {
	schema := express_route_circuit.RegisterResource().Schema["sku"].Elem.(*pluginsdk.Resource).Schema["tier"]
	tests := []struct {
		tier  string
		valid bool
	}{
		{tier: "Basic", valid: true},
		{tier: "Local", valid: true},
		{tier: "Premium", valid: true},
		{tier: "Standard", valid: true},
		{tier: "MultiCloud", valid: false},
		{tier: "standard", valid: false},
		{tier: "", valid: false},
	}

	for _, test := range tests {
		t.Run(test.tier, func(t *testing.T) {
			_, errors := schema.ValidateFunc(test.tier, "sku.0.tier")
			if valid := len(errors) == 0; valid != test.valid {
				t.Fatalf("expected valid=%t, got errors: %+v", test.valid, errors)
			}
		})
	}
}
