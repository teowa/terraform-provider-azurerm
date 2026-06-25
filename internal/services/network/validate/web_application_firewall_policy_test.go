// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestValidateWebApplicationFirewallPolicyManagedRuleSetType(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "default rule set",
			input:    "Microsoft_DefaultRuleSet",
			expected: true,
		},
		{
			name:     "http ddos rule set",
			input:    "Microsoft_HTTPDDoSRuleSet",
			expected: true,
		},
		{
			name:     "invalid rule set",
			input:    "Microsoft_UnknownRuleSet",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errors := ValidateWebApplicationFirewallPolicyManagedRuleSetType(testCase.input, "type")
			if ok := len(errors) == 0; ok != testCase.expected {
				t.Fatalf("input %q: expected %t got %t", testCase.input, testCase.expected, ok)
			}
		})
	}
}

func TestValidateWebApplicationFirewallPolicyManagedRuleGroupName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "existing rule group",
			input:    "REQUEST-920-PROTOCOL-ENFORCEMENT",
			expected: true,
		},
		{
			name:     "http ddos rule group",
			input:    "ExcessiveRequests",
			expected: true,
		},
		{
			name:     "invalid rule group",
			input:    "DoesNotExist",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errors := ValidateWebApplicationFirewallPolicyManagedRuleGroupName(testCase.input, "rule_group_name")
			if ok := len(errors) == 0; ok != testCase.expected {
				t.Fatalf("input %q: expected %t got %t", testCase.input, testCase.expected, ok)
			}
		})
	}
}
