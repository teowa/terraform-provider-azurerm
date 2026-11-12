// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_batch_rule_set

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
)

func TestFrontDoorBatchEdgeAction(t *testing.T) {
	for _, invocationPoint := range rulesets.PossibleValuesForInvocationPoint() {
		input := []helpers.CdnFrontDoorRuleActionsModel{{
			EdgeAction: []helpers.CdnFrontDoorRuleEdgeActionModel{{
				CdnEdgeActionID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Cdn/edgeActions/test",
				InvocationPoint: invocationPoint,
			}},
		}}
		expanded, err := expandCdnFrontDoorBatchRuleSetActions(input)
		if err != nil {
			t.Fatal(err)
		}
		result, err := flattenCdnFrontDoorBatchRuleSetActions(&expanded)
		if err != nil || !reflect.DeepEqual(input, result) {
			t.Fatalf("EdgeAction round trip failed: %#v, %v", result, err)
		}
		input[0].EdgeAction = append(input[0].EdgeAction, input[0].EdgeAction[0], input[0].EdgeAction[0], input[0].EdgeAction[0], input[0].EdgeAction[0], input[0].EdgeAction[0])
		if _, err := expandCdnFrontDoorBatchRuleSetActions(input); err == nil {
			t.Fatal("EdgeActions must count towards the five-action limit")
		}
	}
}
