// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_elasticsearch

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceElasticsearchImporter(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
	id, err := elasticmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return []*pluginsdk.ResourceData{d}, err
	}

	client := meta.(*clients.Client).Elastic.ServerlessMonitorClient
	resp, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := validate.ElasticCloudHostedMonitor(id, resp.Model); err != nil {
		return []*pluginsdk.ResourceData{d}, err
	}

	return []*pluginsdk.ResourceData{d}, nil
}
