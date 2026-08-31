// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func GeoCatalogName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(3, 24),
		validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9-]*$"), "can contain only alphanumeric characters and hyphens"),
	)(i, k)
}
