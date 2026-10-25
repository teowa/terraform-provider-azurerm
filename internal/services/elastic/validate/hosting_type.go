// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
)

// ElasticCloudServerlessMonitor validates that the given monitor uses `Serverless` hosting, since the same
// underlying monitor can be managed as either `azurerm_elastic_cloud_elasticsearch` (`Hosted`) or
// `azurerm_elastic_cloud_serverless` (`Serverless`).
func ElasticCloudServerlessMonitor(id *elasticmonitorresources.MonitorId, model *elasticmonitorresources.ElasticMonitorResource) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	if model.Properties.HostingType == nil {
		return fmt.Errorf("retrieving %s: `properties.hostingType` was nil", id)
	}
	if *model.Properties.HostingType != elasticmonitorresources.HostingTypeServerless {
		if *model.Properties.HostingType == elasticmonitorresources.HostingTypeHosted {
			return fmt.Errorf("expected %s to use `Serverless` hosting, got `Hosted`; use the `azurerm_elastic_cloud_elasticsearch` resource instead", id)
		}
		return fmt.Errorf("expected %s to use `Serverless` hosting, got `%s`", id, *model.Properties.HostingType)
	}

	return nil
}

// ElasticCloudHostedMonitor validates that the given monitor uses `Hosted` hosting, since the same
// underlying monitor can be managed as either `azurerm_elastic_cloud_elasticsearch` (`Hosted`) or
// `azurerm_elastic_cloud_serverless` (`Serverless`).
func ElasticCloudHostedMonitor(id *elasticmonitorresources.MonitorId, model *elasticmonitorresources.ElasticMonitorResource) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	// Hosted monitors created before hostingType was introduced can omit this field.
	if model.Properties.HostingType == nil {
		return nil
	}
	if *model.Properties.HostingType != elasticmonitorresources.HostingTypeHosted {
		if *model.Properties.HostingType == elasticmonitorresources.HostingTypeServerless {
			return fmt.Errorf("expected %s to use `Hosted` hosting, got `Serverless`; use the `azurerm_elastic_cloud_serverless` resource instead", id)
		}
		return fmt.Errorf("expected %s to use `Hosted` hosting, got `%s`", id, *model.Properties.HostingType)
	}

	return nil
}
