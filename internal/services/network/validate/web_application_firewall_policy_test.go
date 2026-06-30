// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestValidateWebApplicationFirewallPolicyRuleGroupName(t *testing.T) {
	cases := []struct {
		Value     string
		HasErrors bool
	}{
		{
			Value:     "General",
			HasErrors: false,
		},
		{
			Value:     "NotARealRuleGroup",
			HasErrors: true,
		},
		{
			Value:     "ExcessiveRequests",
			HasErrors: true,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateWebApplicationFirewallPolicyRuleGroupName(tc.Value, "rule_group_name")
		hasErrors := len(errors) > 0
		if hasErrors != tc.HasErrors {
			t.Fatalf("expected validation errors for rule group `%s` to be %t", tc.Value, tc.HasErrors)
		}
	}
}

func TestValidateWebApplicationFirewallPolicyManagedRuleOverrideGroupName(t *testing.T) {
	cases := []struct {
		Value     string
		HasErrors bool
	}{
		{
			Value:     "ExcessiveRequests",
			HasErrors: false,
		},
		{
			Value:     "NotARealRuleGroup",
			HasErrors: true,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateWebApplicationFirewallPolicyManagedRuleOverrideGroupName(tc.Value, "rule_group_name")
		hasErrors := len(errors) > 0
		if hasErrors != tc.HasErrors {
			t.Fatalf("expected validation errors for managed rule override group `%s` to be %t", tc.Value, tc.HasErrors)
		}
	}
}

func TestValidateWebApplicationFirewallPolicyRuleSetType(t *testing.T) {
	cases := []struct {
		Value     string
		HasErrors bool
	}{
		{
			Value:     "Microsoft_HTTPDDoSRuleSet",
			HasErrors: true,
		},
		{
			Value:     "NotARealRuleSetType",
			HasErrors: true,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateWebApplicationFirewallPolicyRuleSetType(tc.Value, "type")
		hasErrors := len(errors) > 0
		if hasErrors != tc.HasErrors {
			t.Fatalf("expected validation errors for rule set type `%s` to be %t", tc.Value, tc.HasErrors)
		}
	}
}

func TestValidateWebApplicationFirewallPolicyManagedRuleSetType(t *testing.T) {
	cases := []struct {
		Value     string
		HasErrors bool
	}{
		{
			Value:     "Microsoft_HTTPDDoSRuleSet",
			HasErrors: false,
		},
		{
			Value:     "NotARealRuleSetType",
			HasErrors: true,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateWebApplicationFirewallPolicyManagedRuleSetType(tc.Value, "type")
		hasErrors := len(errors) > 0
		if hasErrors != tc.HasErrors {
			t.Fatalf("expected validation errors for managed rule set type `%s` to be %t", tc.Value, tc.HasErrors)
		}
	}
}
