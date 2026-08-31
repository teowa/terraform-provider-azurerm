// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateGeoCatalogName(t *testing.T) {
	const errLen = "to be in the range (3 - 24)"
	const errAllowList = "can contain only alphanumeric characters and hyphens"

	cases := []struct {
		Name           string
		Input          string
		ExpectedErrors []string
	}{
		// Happy paths:
		{
			Name:  "Entire character allow-list",
			Input: "aZ09-",
		},
		{
			Name:  "Minimum character length",
			Input: "aaa",
		},
		{
			Name:  "Maximum character length",
			Input: "012345678901234567890123", // 24 chars
		},

		// Simple negative cases:
		{
			Name:           "Introduce a non-allowed character",
			Input:          "aZ09_", // underscore
			ExpectedErrors: []string{errAllowList},
		},
		{
			Name:           "Above maximum character length",
			Input:          "0123456789012345678901234", // 25 chars
			ExpectedErrors: []string{errLen},
		},
		{
			Name:           "Below minimum character length",
			Input:          "aa",
			ExpectedErrors: []string{errLen},
		},
		{
			Name:           "Specifically test for emptiness",
			Input:          "",
			ExpectedErrors: []string{errLen},
		},
	}

	errsContain := func(errors []error, text string) bool {
		for _, err := range errors {
			if strings.Contains(err.Error(), text) {
				return true
			}
		}
		return false
	}

	t.Parallel()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errors := GeoCatalogName(tc.Input, "azurerm_orbital_geocatalog.test.name")

			if len(errors) != len(tc.ExpectedErrors) {
				t.Fatalf("Expected %d errors but got %d for %q: %v", len(tc.ExpectedErrors), len(errors), tc.Input, errors)
			}

			for _, expectedError := range tc.ExpectedErrors {
				if !errsContain(errors, expectedError) {
					t.Fatalf("Errors did not contain expected error: %s", expectedError)
				}
			}
		})
	}
}
