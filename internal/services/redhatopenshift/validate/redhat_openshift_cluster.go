package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/parse"
)

func ValidateClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := parse.RedhatOpenShiftClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

func ValidateDiskSizeGB(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 128 {
		errors = append(errors, fmt.Errorf(
			"The `disk_size_gb` must be 128 or greater"))
	}
	return warnings, errors
}
